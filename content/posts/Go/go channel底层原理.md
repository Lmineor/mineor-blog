---
title: go channel底层原理介绍
draft: false
date: 2025-04-21
tags:
  - Go
---
# 主要结构定义

```go
type hchan struct {
	qcount   uint           // 当前队列中的元素数量
	dataqsiz uint           // 环形队列的大小
	buf      unsafe.Pointer // 指向环形队列的指针
	elemsize uint16 // 环形队列中的元素的大小
	closed   uint32 // channel是否已关闭
	elemtype *_type // 环形队列中元素的类型
	sendx    uint   // 发送索引
	recvx    uint   // 接收索引
	recvq    waitq  // 等待接收的goroutine队列
	sendq    waitq  // 等待发送的goroutine队列
	lock mutex      // 互斥锁，保护channel的并发访问
}
```

其中waitq结构定义如下
```go
type waitq struct {  
   first *sudog   // 等待队列的头sudog
   last  *sudog  // 等待队列的尾sudog
}
```


# 创建chan

```go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem

	// 元素大小不能超过64KB
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}

	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}

	// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
	// buf points into the same allocation, elemtype is persistent.
	// SudoG's are referenced from their owning thread so they can't be collected.
	// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
	var c *hchan
	switch {
	case mem == 0:
		// Queue or element size is zero.
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
	case elem.ptrdata == 0:
		// Elements do not contain pointers.
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// Elements contain pointers.
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}

	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)
	lockInit(&c.lock, lockRankHchan)

	if debugChan {
		print("makechan: chan=", c, "; elemsize=", elem.size, "; dataqsiz=", size, "\n")
	}
	return c
}
```