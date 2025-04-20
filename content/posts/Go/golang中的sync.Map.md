---
title: "讲明白Golang的sync.map"
date: 2024-12-10
draft: false
tags : [                    # 文章所属标签
    "Go",
]
---

参考：https://www.jianshu.com/p/7c4fb2c3c66e

Go语言版本

```bash
go version go1.20 windows/amd64
```


sync.map 则是一种并发安全的 map，在 Go 1.9 引入


# 1 特点

Map 类型针对两种常见用例进行了优化：

(1) 当给定键的条目只写入一次但读取多次时，如仅增长的缓存，或
(2) 多个 goroutines 读取、写入和覆盖不相交键集的条目。

在这两种情况下，与搭配单独的 Mutex 或 RWMutex 的 Go map 相比，使用 Map 可以显著减少锁争用。

> 需要注意的是:Map的零值为空且可以直接使用，且如果被使用后不能被复制


在 Go 内存模型的术语中，Map 安排“在”任一读取操作观察到写入效果之前“进行同步”，读取和写入操作定义如下。

- Load、LoadAndDelete、LoadOrStore、Swap、CompareAndSwap 和 CompareAndDelete 是读取操作；
- Delete、LoadAndDelete、Store 和 Swap 是写入操作；
- 当 LoadOrStore 返回 loaded 为 false 时，它是写入操作；
- 当 CompareAndSwap 返回 swapped 为 true 时，它是写入操作；
- 当 CompareAndDelete 返回 deleted 为 true 时，它是写入操作。

## 1.1 结构体

### 1.1.1 Map结构体

```go
type Map struct {
	mu Mutex

    // read 包含map内容中可安全进行并发访问的部分（无论是否加锁）
    // read 字段本身可以安全读取，但只有在加锁的情况下才能存储
    // 存储在 read 中的条目可以在不加锁的情况下并发更新，
    // 更新之前删除（expunged）的键值对需要在加锁的情况下将键值对复制到dirty map中，并标记为未删除（unexpunged）
	read atomic.Pointer[readOnly]

    // dirty 包含map内容的部分，需要加锁
    // 为了确保可以快速将 dirty 映射提升为 read 映射，它还包括 read 映射中所有未被删除的条目。
    // 被删除的条目不会存储在 dirty 映射中。在 clean 映射中的被删除条目必须在存储新值之前取消删除并添加到 dirty 映射中。
    // 如果 dirty 映射为 nil，则对地图的下一次写入操作将通过对 clean 映射进行浅拷贝来初始化它，忽略过时的条目。
	dirty map[any]*entry

    // 记录read读取不到数据，需要加锁读取的次数
    // 当misses等于dirty的长度时，会将dirty复制到read
	misses int
}
```

### 1.1.2 readOnly结构体

readOnly 是一个不可变的结构体，以原子方式存储在 Map.read 字段中。

```go
type readOnly struct {
	m       map[any]*entry
	amended bool // dirty 中的key不在m中为true
}
```

### 1.1.3 entry结构体

```go
type entry struct {
	// p指向该条目存储的interface{}的值
	//
	// 如果p == nil，该条目已经被删除，同时m.dirty == nil或者m.dirty[key] 是 e.
	//
	// 如果 p == expunged, 该条目已经被删除，m.dirty != nil, 该条目不在m.dirty中
	//
	// 否则，如果m.dirty != nil, 且在m.dirty[key]中，那么该条目是有效的并记录在m.read.m[key]
	//
	// 一个元素可以用nil通过原子替换的方式进行删除，当下一次创建m.drity，会自动用expunged替换nil，不会将其复制到dirty中
	p atomic.Pointer[any]
}
```

# 2 主要方法

- Load(key any) (value any, ok bool) 
- Store(key, value any)
- Swap(key, value any) (previous any, loaded bool)
- LoadOrStore(key, value any) (actual any, loaded bool)
- LoadAndDelete(key any) (value any, loaded bool)
- Delete(key any)
- Range(f func(key, value any) bool)

## 2.1 Load

Load根据key拿到map中存储的值，如果没有的话返回nil，
ok代表map中是否有结果

```go
func (m *Map) loadReadOnly() readOnly {
	if p := m.read.Load(); p != nil {
		return *p
	}
	return readOnly{}
}

func (m *Map) Load(key any) (value any, ok bool) {
	read := m.loadReadOnly()
	e, ok := read.m[key]
	if !ok && read.amended { // read中没有取到，且dirty map中的key没有在read中（即该值在dirty中，不在read中），则进行加锁去读
		m.mu.Lock()
		// 再读一次的目的是防止加锁的过程中dirty变成read，防止出现读不到的错误
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// 不管是否存在，都要记录miss: 
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	return e.load()
}

// misss长度小于dirty长度时，只++miss
// 否则把dirty提升为read，然后dirty置为nil。miss归零
func (m *Map) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(&readOnly{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

```

## 2.2 Store&Swap

```go
func (m *Map) Store(key, value any) {
	_, _ = m.Swap(key, value)
}

// Swap 如字面意思，就是新旧交换，存新的，返回旧的
func (m *Map) Swap(key, value any) (previous any, loaded bool) {
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok { //read里面有，则直接进行变更，不加锁，原子操作
		if v, ok := e.trySwap(&value); ok {
			if v == nil {
				return nil, false
			}
			return *v, true
		}
	}

	// read里面没有，加锁进一步操作
	m.mu.Lock()
	read = m.loadReadOnly()
	// 再读一次的目的是防止加锁的过程中dirty变成read
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			//如果读到的值为expunged，对于nil的元素，搞成了expunged，所以意味着dirty不为nil，且这个元素中没有该元素
			m.dirty[key] = e
		}
		// 更新read中的值
		if v := e.swapLocked(&value); v != nil {
			loaded = true
			previous = *v
		}
	} else if e, ok := m.dirty[key]; ok {  // 不在read中，在dirty中
		if v := e.swapLocked(&value); v != nil {
			loaded = true
			previous = *v
		}
	} else {
		if !read.amended {
			// read.amended==false,说明dirty map为空，需要将read map 复制一份到dirty map
			m.dirtyLocked()
			m.read.Store(&readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
	return previous, loaded
}

func (e *entry) trySwap(i *any) (*any, bool) {
	for {
		p := e.p.Load()
		if p == expunged {
			return nil, false
		}
		if e.p.CompareAndSwap(p, i) {
			return p, true
		}
	}
}

func (m *Map) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read := m.loadReadOnly()
	m.dirty = make(map[any]*entry, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}
```

## 2.3 LoadOrStore

有则返回，没有则存储

```go
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool) {
	// Avoid locking if it's a clean hit.
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	// 同Swap
	read = m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			m.dirtyLocked()
			m.read.Store(&readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

func (e *entry) tryLoadOrStore(i any) (actual any, loaded, ok bool) {
	p := e.p.Load()
	if p == expunged {
		return nil, false, false
	}
	if p != nil {
		return *p, true, true
	}

	// Copy the interface after the first load to make this method more amenable
	// to escape analysis: if we hit the "load" path or the entry is expunged, we
	// shouldn't bother heap-allocating.
	ic := i
	for {
		if e.p.CompareAndSwap(nil, &ic) {
			return i, false, true
		}
		p = e.p.Load()
		if p == expunged {
			return nil, false, false
		}
		if p != nil {
			return *p, true, true
		}
	}
}

```

## 2.4 LoadAndDelete

根据key 删除元素，返回已删除元素的值

```go
func (m *Map) LoadAndDelete(key any) (value any, loaded bool) {
	read := m.loadReadOnly()
	// 先来read中找
	e, ok := read.m[key]
	// read中没有，先加锁，再尝试
	// read中有，则直接标记为nil
	if !ok && read.amended {
		m.mu.Lock()
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// 删除dirty中的元素
			delete(m.dirty, key)
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty map is promoted to the read
			// map.
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		// 把read中的元素标记为nil（delete）
		return e.delete()
	}
	return nil, false
}

func (e *entry) delete() (value any, ok bool) {
	for {
		p := e.p.Load()
		if p == nil || p == expunged {
			return nil, false
		}
		// e是read中的entry，删除把p标记为nil
		if e.p.CompareAndSwap(p, nil) {
			return *p, true
		}
	}
}

```

## 2.5 Delete

删除元素

```go
func (m *Map) Delete(key any) {
	m.LoadAndDelete(key)
}
```

## 2.6 Range

```go
func (m *Map) Range(f func(key, value any) bool) {
	read := m.loadReadOnly()
	if read.amended {
		// dirty中的key不在read中
		m.mu.Lock()
		read = m.loadReadOnly()
		if read.amended {
			//升级dirty
			read = readOnly{m: m.dirty}
			m.read.Store(&read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}
	// 此时，dirty和read都一致
	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			// 元素被标记为删除，忽略
			continue
		}
		//函数返回false，终止
		if !f(k, v) {
			break
		}
	}
}

func (e *entry) load() (value any, ok bool) {
	p := e.p.Load()
	if p == nil || p == expunged {
		return nil, false
	}
	return *p, true
}
```


# 3 如何保证键类型和值类型的正确性

map的键类型的不能是哪些类型：即函数类型、切片类型、字典类型，同样的，对sync.map的键类型也是一样的，不能为这3种类型。

而sync.map中涉及到的键和值类型都为any类型，所以必须依赖我们自己来保证键类型和值类型的正确性，那么问题就来了：如何保证并发安全字典中键和值类型的正确性?

## 方案一：让sync.Map只存储某个特定类型的键

```go
// IntStrMap 代表键类型为int、值类型为string的并发安全字典。
type IntStrMap struct {
	m sync.Map
}

func (iMap *IntStrMap) Delete(key int) {
	iMap.m.Delete(key)
}

func (iMap *IntStrMap) Load(key int) (value string, ok bool) {
	v, ok := iMap.m.Load(key)
	if v != nil {
		value = v.(string)
	}
	return
}

func (iMap *IntStrMap) LoadOrStore(key int, value string) (actual string, loaded bool) {
	a, loaded := iMap.m.LoadOrStore(key, value)
	actual = a.(string)
	return
}

func (iMap *IntStrMap) Range(f func(key int, value string) bool) {
	f1 := func(key, value interface{}) bool {
		return f(key.(int), value.(string))
	}
	iMap.m.Range(f1)
}

func (iMap *IntStrMap) Store(key int, value string) {
	iMap.m.Store(key, value)
}
```

方案一的实现很简单，但是缺点也是显而易见的，非常不灵活，不能灵活改变键和值的类型，需求多了之后，会产生很多雷同的代码，因此我们来看方案二

## 方案二：封装的结构体类型的所有方法，与sync.Map类型完全一致，此时需要类型检查

```go
// ConcurrentMap 代表可自定义键类型和值类型的并发安全字典。
type ConcurrentMap struct {
    m         sync.Map
    keyType   reflect.Type  //键类型
    valueType reflect.Type  //值类型
}

func NewConcurrentMap(keyType, valueType reflect.Type) (*ConcurrentMap, error) {
    if keyType == nil {
        return nil, errors.New("nil key type")
    }
    if !keyType.Comparable() {
        return nil, fmt.Errorf("incomparable key type: %s", keyType)
    }
    if valueType == nil {
        return nil, errors.New("nil value type")
    }
    cMap := &ConcurrentMap{
        keyType:   keyType,
        valueType: valueType,
    }
    return cMap, nil
}

func (cMap *ConcurrentMap) Delete(key interface{}) {
    if reflect.TypeOf(key) != cMap.keyType {
        return
    }
    cMap.m.Delete(key)
}

func (cMap *ConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
    if reflect.TypeOf(key) != cMap.keyType {
        return
    }
    return cMap.m.Load(key)
}

func (cMap *ConcurrentMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
    if reflect.TypeOf(key) != cMap.keyType {
        panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
    }
    if reflect.TypeOf(value) != cMap.valueType {
        panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
    }
    actual, loaded = cMap.m.LoadOrStore(key, value)
    return
}

func (cMap *ConcurrentMap) Range(f func(key, value interface{}) bool) {
    cMap.m.Range(f)
}

func (cMap *ConcurrentMap) Store(key, value interface{}) {
    if reflect.TypeOf(key) != cMap.keyType {
        panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
    }
    if reflect.TypeOf(value) != cMap.valueType {
        panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
    }
    cMap.m.Store(key, value)
}
```