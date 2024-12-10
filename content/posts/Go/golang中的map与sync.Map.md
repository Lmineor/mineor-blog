---
title: "Golang中的map与sync.map"
date: 2024-12-10
draft: true
tags : [                    # 文章所属标签
    "Go",
]
---

# map

todo

# 2 sync.map

## 2.1 特点

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

## 2.2 结构体

### 2.2.1 Map结构体

```go
type Map struct {
	mu Mutex

    // read 包含map内容中可安全进行并发访问的部分（无论是否持有 mu）
    // read 字段本身可以安全读取，但只有在加锁的情况下才能存储
    // 存储在 read 中的条目可以在不加锁的情况下并发更新，
    // 更新先前删除的条目需要将条目复制到dirty map中，并在加锁的情况下取消删除
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

### 2.2.2 readOnly结构体

readOnly 是一个不可变的结构体，以原子方式存储在 Map.read 字段中。

```go
type readOnly struct {
	m       map[any]*entry
	amended bool // dirty 中的key不在m中为true
}
```

### 2.2.3 entry结构体

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

## 2.3 主要方法

#### 2.3.1 Load

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

### 2.3.2 Store&Swap

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

### 2.3.3 LoadOrStore

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

### 2.3.4 LoadAndDelete

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

### 2.3.5 Delete

删除元素

```go
func (m *Map) Delete(key any) {
	m.LoadAndDelete(key)
}
```

### 2.3.6 Range

```go
func (m *Map) Range(f func(key, value any) bool) {
	// We need to be able to iterate over all of the keys that were already
	// present at the start of the call to Range.
	// If read.amended is false, then read.m satisfies that property without
	// requiring us to hold m.mu for a long time.
	read := m.loadReadOnly()
	if read.amended {
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		m.mu.Lock()
		read = m.loadReadOnly()
		if read.amended {
			read = readOnly{m: m.dirty}
			m.read.Store(&read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}
```

- **并发安全，且虽然用到了锁，但是显著减少了锁的争用**。 sync.map出现之前，如果想要实现并发安全的map，只能自行构建，使用sync.Mutex或sync.RWMutex，再加上原生的map就可以轻松做到，sync.map也用到了锁，但是在尽可能的避免使用锁，因为使用锁意味着要把一些并行化的东西串行化，会降低程序性能，因此能用原子操作就不要用锁，但是原子操作局限性比较大，只能对一些基本的类型提供支持，在sync.map中将两者做了比较完美的结合。
- **存取删操作的算法复杂度与map一样，都是O(1)**
- **不会做类型检查。**  sync.map只是go语言标准库中的一员，而不是语言层面的东西，也正是因为这一点，go语言的编译器不会对其中的键和值进行特殊的类型检查

作者：xixisuli

链接：https://www.jianshu.com/p/7c4fb2c3c66e

来源：简书

著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。