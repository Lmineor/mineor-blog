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
	// 检查内存对齐
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}

	// 计算缓冲区大小
	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}

	// 根据缓冲区大小和元素是否包含指针，选择不同的内存分配策略：
	var c *hchan
	switch {
	case mem == 0:
		// 队列大小或元素大小为零。
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
	case elem.ptrdata == 0:
		// 元素不包含指针
		// 在一次调用中分配hchan和buf（指向环形队列的指针）
		// 将hchan和buf分配到同一块内存中
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		// 将 c.buf 设置为 hchan 内存的末尾。
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// 元素包含指针.
		// 如果元素包含指针，需要分别分配 hchan 和缓冲区。
		// 使用 new 分配 hchan，然后使用 mallocgc 分配缓冲区。
		// 这里是否用到了内存逃逸=>都是分配在堆上，不会出现
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

# chan接收数据的实现

首先，c <- x 的实现是调用了`chansend1(c *hchan, elem unsafe.Pointer)`函数

chansend1又调用了chansend函数

```go
// 实现 c <- x
//
//go:nosplit
func chansend1(c *hchan, elem unsafe.Pointer) {
	chansend(c, elem, true, getcallerpc())
}
```
那么接下来看看chansend的实现

## chansend的实现

```go
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	// 如果channel已经关闭
	if c == nil {
		// channel已经关闭，且为非阻塞的，那么就直接return了，不做动作
		if !block {
			return false
		}
		// 阻塞的话，挂起当前goroutine，将当前goroutine置为waiting状态
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	if debugChan {
		print("chansend: chan=", c, "\n")
	}

	if raceenabled {
		racereadpc(c.raceaddr(), callerpc, abi.FuncPCABIInternal(chansend))
	}

	// Fast path: check for failed non-blocking operation without acquiring the lock.
	//
	// After observing that the channel is not closed, we observe that the channel is
	// not ready for sending. Each of these observations is a single word-sized read
	// (first c.closed and second full()).
	// Because a closed channel cannot transition from 'ready for sending' to
	// 'not ready for sending', even if the channel is closed between the two observations,
	// they imply a moment between the two when the channel was both not yet closed
	// and not ready for sending. We behave as if we observed the channel at that moment,
	// and report that the send cannot proceed.
	//
	// It is okay if the reads are reordered here: if we observe that the channel is not
	// ready for sending and then observe that it is not closed, that implies that the
	// channel wasn't closed during the first observation. However, nothing here
	// guarantees forward progress. We rely on the side effects of lock release in
	// chanrecv() and closechan() to update this thread's view of c.closed and full().
	if !block && c.closed == 0 && full(c) {
		return false
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	lock(&c.lock)

	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}

	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			racenotify(c, c.sendx, nil)
		}
		typedmemmove(c.elemtype, qp, ep)
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}

	if !block {
		unlock(&c.lock)
		return false
	}

	// Block on the channel. Some receiver will complete our operation for us.
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg
	gp.param = nil
	c.sendq.enqueue(mysg)
	// Signal to anyone trying to shrink our stack that we're about
	// to park on a channel. The window between when this G's status
	// changes and when we set gp.activeStackChans is not safe for
	// stack shrinking.
	gp.parkingOnChan.Store(true)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
	// Ensure the value being sent is kept alive until the
	// receiver copies it out. The sudog has a pointer to the
	// stack object, but sudogs aren't considered as roots of the
	// stack tracer.
	KeepAlive(ep)

	// someone woke us up.
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	closed := !mysg.success
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	mysg.c = nil
	releaseSudog(mysg)
	if closed {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
		panic(plainError("send on closed channel"))
	}
	return true
}
```