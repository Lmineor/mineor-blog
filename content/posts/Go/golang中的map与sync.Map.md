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

## 主要方法

#### Load

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
		// Avoid reporting a spurious miss if m.dirty got promoted while we were
		// blocked on m.mu. (If further loads of the same key will not miss, it's
		// not worth copying the dirty map for this key.)
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty map is promoted to the read
			// map.
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	return e.load()
}
```

- **并发安全，且虽然用到了锁，但是显著减少了锁的争用**。 sync.map出现之前，如果想要实现并发安全的map，只能自行构建，使用sync.Mutex或sync.RWMutex，再加上原生的map就可以轻松做到，sync.map也用到了锁，但是在尽可能的避免使用锁，因为使用锁意味着要把一些并行化的东西串行化，会降低程序性能，因此能用原子操作就不要用锁，但是原子操作局限性比较大，只能对一些基本的类型提供支持，在sync.map中将两者做了比较完美的结合。
- **存取删操作的算法复杂度与map一样，都是O(1)**
- **不会做类型检查。**  sync.map只是go语言标准库中的一员，而不是语言层面的东西，也正是因为这一点，go语言的编译器不会对其中的键和值进行特殊的类型检查

作者：xixisuli

链接：https://www.jianshu.com/p/7c4fb2c3c66e

来源：简书

著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。