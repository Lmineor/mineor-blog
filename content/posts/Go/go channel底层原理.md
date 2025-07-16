---
title: go channel底层原理介绍
draft: false
date: 2025-04-21
tags:
  - Go
---

Go 语言中的 `chan`（通道）是并发编程的核心组件，其底层实现结合了高效的数据结构和同步机制。以下是其核心原理的总结：

---

### **1. 数据结构**
- **环形队列**：通道内部通过环形缓冲区（`hchan` 结构体）存储数据，包含以下关键字段：
  - `buf`：指向环形缓冲区的指针，用于存储元素。
  - `sendx` 和 `recvx`：记录发送和接收的位置索引（环形队列的头部和尾部）。
  - `qcount`：当前缓冲区中的元素数量。
  - `dataqsiz`：缓冲区容量（带缓冲通道的固定大小，无缓冲通道为 0）。

---

### **2. 同步机制**
- **互斥锁（`lock`）**：保护 `hchan` 的字段，避免并发读写冲突。
- **等待队列**：
  - `sendq`：因发送阻塞的 Goroutine 队列（缓冲区满时）。
  - `recvq`：因接收阻塞的 Goroutine 队列（缓冲区空时）。
- **调度控制**：通过 `runtime.sudog` 表示等待的 Goroutine，由运行时调度器管理。

---

### **3. 操作流程**
#### **发送数据（`ch <- val`）**
1. **加锁**：获取通道的互斥锁。
2. **直接交付**：如果 `recvq` 有等待的接收者，直接将数据拷贝到接收方，解锁并唤醒接收方。
3. **缓冲写入**：若缓冲区未满，数据写入 `buf`，更新索引，解锁。
4. **阻塞等待**：若缓冲区满，当前 Goroutine 加入 `sendq`，解锁并挂起（等待被接收方唤醒）。

#### **接收数据（`val <- ch`）**
1. **加锁**：获取通道的互斥锁。
2. **直接获取**：如果 `sendq` 有等待的发送者，直接从发送方拷贝数据（或从缓冲区读取后唤醒发送方）。
3. **缓冲读取**：若缓冲区非空，从 `buf` 读取数据，更新索引，解锁。
4. **阻塞等待**：若缓冲区空，当前 Goroutine 加入 `recvq`，解锁并挂起（等待被发送方唤醒）。

---

### **4. 关闭通道（`close(ch)`）**
- **设置标志位**：标记通道为关闭状态（`closed = 1`）。
- **唤醒所有等待的 Goroutine**：`sendq` 和 `recvq` 中的 Goroutine 会被唤醒，接收方收到零值，发送方触发 panic。

---

### **5. 无缓冲与带缓冲通道**
- **无缓冲通道**：`dataqsiz = 0`，发送和接收必须同步完成（“同步通信”）。
- **带缓冲通道**：允许异步操作，发送仅在缓冲区满时阻塞，接收仅在空时阻塞。

---

### **6. 底层实现优化**
- **内存复用**：发送/接收数据时尽量复用内存，减少拷贝（如直接拷贝到目标 Goroutine 的栈）。
- **编译器优化**：部分简单操作（如单元素无缓冲通道）可能被内联或优化为原子操作。

---

### **7. 注意事项**
- **并发安全**：通道本身是并发安全的，但需避免在关闭后发送数据。
- **性能权衡**：频繁的小数据传递可能因锁竞争影响性能，需合理选择缓冲大小。

---

通过以上机制，Go 的通道实现了高效的 Goroutine 间通信，平衡了易用性和性能。理解底层原理有助于编写更高效的并发代码。

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

	// 非阻塞，channel未关闭，channel为满的，则发送不成功
	// 怎么理解channel为满的呢，
	// 如果环形队列的大小为0，等待接收的goroutine队列第一个指向nil，那么就可以理解为满的
	// 当前队列中的元素数量==环形队列的大小，那么也可以理解为满的
	if !block && c.closed == 0 && full(c) {
		return false
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}
	// 需要上锁
	lock(&c.lock)

	if c.closed != 0 {
		// 向已关闭的channel中发送数据触发异常
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}

	// 如果有等待接收的队列不为空，那么把值直接发给队列中的sudog
	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	if c.qcount < c.dataqsiz {
		// 环形队列有空间
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			racenotify(c, c.sendx, nil)
		}
		typedmemmove(c.elemtype, qp, ep)
		c.sendx++
		// 环形队列转了一圈了,将发送索引置为0，从头开始
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}
	// 非阻塞场景，队列中没有空间，直接返回false
	if !block {
		unlock(&c.lock)
		return false
	}
	// 以下是阻塞场景，尝试唤醒等待的g
	// Block on the channel. Some receiver will complete our operation for us.
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	// 下面这一堆后面会用
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg //反向标记sudog
	gp.param = nil
	c.sendq.enqueue(mysg) // 将含有当前g指针和send元素的sudog结构体保存到sendq中等待唤醒
	// Signal to anyone trying to shrink our stack that we're about
	// to park on a channel. The window between when this G's status
	// changes and when we set gp.activeStackChans is not safe for
	// stack shrinking.
	gp.parkingOnChan.Store(true)

	// 挂起当前goroutine
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
	// Ensure the value being sent is kept alive until the
	// receiver copies it out. The sudog has a pointer to the
	// stack object, but sudogs aren't considered as roots of the
	// stack tracer.
	KeepAlive(ep)

	// someone woke us up.
	//  解放了，当前g被唤醒了
	// 确保之前当前g的waiting等于sudog
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

在看接收操作
```go
// 实现代码中的 <- c 操作
func chanrecv1(c *hchan, elem unsafe.Pointer) {
	chanrecv(c, elem, true)
}
```

再来看chanrecv的具体实现
```go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	// raceenabled: don't need to check ep, as it is always on the stack
	// or is new memory allocated by reflect.

	if debugChan {
		print("chanrecv: chan=", c, "\n")
	}
	// 当前channel为空,非阻塞场景直接返回，阻塞场景挂起当前goroutine
	if c == nil {
		if !block {
			return
		}
		gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	if !block && empty(c) {
		// 非阻塞场景且发送队列为空
		if atomic.Load(&c.closed) == 0 {
			// 通过原子操作判断如果当前队列没有关闭，则直接返回false，false
			return
		}
		// 再次检查发送队列是否为空（有可能在上面检查完为空之后和检查队列关闭之间发送了数据）
		if empty(c) {
			// The channel is irreversibly closed and empty.
			if raceenabled {
				raceacquire(c.raceaddr())
			}
			if ep != nil {
				typedmemclr(c.elemtype, ep)
			}
			// 通道ok，但没有值
			return true, false
		}
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	// 加锁
	lock(&c.lock)

	// channel已经关闭了
	if c.closed != 0 {
		if c.qcount == 0 {
			if raceenabled {
				raceacquire(c.raceaddr())
			}
			unlock(&c.lock)
			if ep != nil {
				typedmemclr(c.elemtype, ep)
			}
			// 相当于可以接收，但是接收值类型零值
			return true, false
		}
		// The channel has been closed, but the channel's buffer have data.
	} else {
		// 从等待发送的goroutine队列弹出一个sudog，直接传递给接收者  sendq
		if sg := c.sendq.dequeue(); sg != nil {
			// Found a waiting sender. If buffer is size 0, receive value
			// directly from sender. Otherwise, receive from head of queue
			// and add sender's value to the tail of the queue (both map to
			// the same buffer slot because the queue is full).
			recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
			return true, true
		}
	}

	// 环形队列中的元素个数大于0
	if c.qcount > 0 {
		// Receive directly from queue
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			racenotify(c, c.recvx, nil)
		}
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		typedmemclr(c.elemtype, qp)
		c.recvx++
		// 环形队列再次指向队列头部
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.qcount--
		unlock(&c.lock)
		return true, true
	}

	// 非阻塞场景，则直接返回
	if !block {
		unlock(&c.lock)
		return false, false
	}

	// no sender available: block on this channel.
	// 队列中没有元素，开始阻塞，没人发送数据
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
	gp.waiting = mysg
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.param = nil
	// 将包含当前gorouine的sudog放到等待接收的队列中
	c.recvq.enqueue(mysg)
	// Signal to anyone trying to shrink our stack that we're about
	// to park on a channel. The window between when this G's status
	// changes and when we set gp.activeStackChans is not safe for
	// stack shrinking.
	gp.parkingOnChan.Store(true)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanReceive, traceEvGoBlockRecv, 2)

	// 等待接收的队列被唤醒
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	success := mysg.success
	gp.param = nil
	mysg.c = nil
	releaseSudog(mysg)
	return true, success
}
```

代码大概流程就是这样，考虑几种特殊场景
# 几种场景

### 场景1：有缓冲的通道，执行 c <- x这样的操作，如果c的队列已经满了

c的队列已有数据，队列空间为4
```bash
c.队列 = [1,2,3,4]
```

这是goroutine `g1`执行`c <- 5`
`c.队列`满了，`5`进不去

查看chansend函数，可以发现如下代码
```go
...
// mysg已经是包含当前操作数据sudog了
c.sendq.enqueue(mysg)
gp.parkingOnChan.Store(true)
gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
...
```
把mysg加入到发送队列中，然后让gmp调度挂起当前g1，那么什么时候能继续执行这个g1呢

看chanrecv的函数，可以看到如下代码
```go
// sendq也就是说send的时候等待队列中数据
if sg := c.sendq.dequeue(); sg != nil {
	// Found a waiting sender. If buffer is size 0, receive value
	// directly from sender. Otherwise, receive from head of queue
	// and add sender's value to the tail of the queue (both map to
	// the same buffer slot because the queue is full).
	// buffer的大小为0（无缓冲队列）--> 直接sendq中取一个值
	// buffer的大小不为0（有缓冲队列）--> 在sendq中去header的值，然后把sender加到等待队列的末尾）
	recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
	return true, true
}
```

某一个时刻有个`g2`执行了`<-c`的操作，也就是说会走到上面的代码中，这时候`g1`就被唤醒了（肯定是`recv(c, sg, ep, func() { unlock(&c.lock) }, 3)`把g1唤醒的，咱们瞧瞧这个revc函数
```go
// 瞧一瞧recv函数
// sg就是发送者
func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	if c.dataqsiz == 0 {
		// 无缓冲chan
		if raceenabled {
			racesync(c, sg)
		}
		if ep != nil {
			// copy data from sender
			recvDirect(c.elemtype, sg, ep)
		}
	} else {
		// 有缓冲的chan
		// Queue is full. Take the item at the
		// head of the queue. Make the sender enqueue
		// its item at the tail of the queue. Since the
		// queue is full, those are both the same slot.
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			racenotify(c, c.recvx, nil)
			racenotify(c, c.recvx, sg)
		}
		// copy data from queue to receiver
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		// copy data from sender to queue
		typedmemmove(c.elemtype, qp, sg.elem)
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.sendx = c.recvx // c.sendx = (c.sendx+1) % c.dataqsiz
		// 队列[2,3,4,5]
	}
	sg.elem = nil
	// 还记得chansend中的这一堆操作吗
	// mysg.elem = ep
	// mysg.waitlink = nil
	// mysg.g = gp
	// mysg.isSelect = false
	// mysg.c = c
	// gp.waiting = mysg
	// gp.param = nil
	gp := sg.g // 当前阻塞的g1
	unlockf()
	gp.param = unsafe.Pointer(sg)
	sg.success = true
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	// 标记当前的g1可以运行_Grunnable，等待被gmp调度（按照先放到p的本地队列中，如果满的话，这块就不赘述了，可以参考下go的gmp模型
	goready(gp, skip+1)
}
```
这样子就实现了`g1`调用`c <- x`满了以后接受者`g2`出现拯救了`g1`的过程

### 场景2：有缓冲的通道，执行 `<- c`这样的操作，如果c的队列已经空了

`g1`执行`<- c`但c的队列空了
查看chanrecv函数可以发现这段代码
```go
// 将包含当前gorouine的sudog放到等待接收的队列中
// mysg已经是包含当前需要接收数据的对象了
c.recvq.enqueue(mysg)
// Signal to anyone trying to shrink our stack that we're about
// to park on a channel. The window between when this G's status
// changes and when we set gp.activeStackChans is not safe for
// stack shrinking.
gp.parkingOnChan.Store(true)
gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanReceive, traceEvGoBlockRecv, 2)
```
把mysg加入到接收队列中，然后让[gmp](./gmp.md)调度挂起当前g1，那么什么时候能继续执行这个g1呢？
同理，看chansend的函数，可以看到如下代码
```go
// 如果有等待接收的队列不为空，那么把只直接发给队列中的sudog
if sg := c.recvq.dequeue(); sg != nil {
	// Found a waiting receiver. We pass the value we want to send
	// directly to the receiver, bypassing the channel buffer (if any).
	// 从一堆等待的接收队列中取一个（head)
	send(c, sg, ep, func() { unlock(&c.lock) }, 3)
	return true
}
```
某一个时刻有个`g2`执行了`c<-x`的操作，也就是说会走到上面的代码中，这时候`g1`就被唤醒了（肯定是`send(c, sg, ep, func() { unlock(&c.lock) }, 3)`把g1唤醒的，咱们瞧瞧这个send函数
```go
// 仔细瞧瞧send函数
// sg 就是接收者
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	if raceenabled {
		if c.dataqsiz == 0 {
			racesync(c, sg)
		} else {
			// Pretend we go through the buffer, even though
			// we copy directly. Note that we need to increment
			// the head/tail locations only when raceenabled.
			racenotify(c, c.recvx, nil)
			racenotify(c, c.recvx, sg)
			c.recvx++
			if c.recvx == c.dataqsiz {
				c.recvx = 0
			}
			c.sendx = c.recvx // c.sendx = (c.sendx+1) % c.dataqsiz
		}
	}
	if sg.elem != nil {
		// 把需要发送的值发送给sg，也就是等待队列中的goroutine
		sendDirect(c.elemtype, sg, ep)
		sg.elem = nil
	}
	gp := sg.g // 也就是g1
	unlockf()
	gp.param = unsafe.Pointer(sg)
	sg.success = true
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	// 也是通过goready把g1唤醒了
	goready(gp, skip+1)
}
```
这样子就实现了`g1`调用`<- c`但c的队列为空以后发送者`g2`出现拯救了`g1`的过程


### 场景3：无缓冲的通道，执行 `c <- x`这样的操作
和场景1类似，只是在`g2`执行`recv`的时候直接把`g1`的`x`复制给了`g2，后续唤醒流程一样

### 场景4：无缓冲的通道，执行 `<-c`这样的操作
和场景2类似，只是在`g2`执行`send`的时候直接把`g2`的`x`复制给了`g1`，后续唤醒流程一样

# closechan的实现

```go
func closechan(c *hchan) {

	if c == nil {
		panic(plainError("close of nil channel"))
	}

	lock(&c.lock)
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("close of closed channel"))
	}

	if raceenabled {
		callerpc := getcallerpc()
		racewritepc(c.raceaddr(), callerpc, abi.FuncPCABIInternal(closechan))
		racerelease(c.raceaddr())
	}

	c.closed = 1 // 标志位

	var glist gList

	// 释放所有等待接收的goroutine
	for {
		sg := c.recvq.dequeue()
		if sg == nil {
			break
		}
		if sg.elem != nil {
			typedmemclr(c.elemtype, sg.elem)
			sg.elem = nil
		}
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = unsafe.Pointer(sg)
		sg.success = false
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		// 将接受队列中等待的goroutine加入到glist中
		glist.push(gp)
	}


	// 释放所有等待写入的goroutine（panic：不能向已经关闭的channel中写）
	for {
		sg := c.sendq.dequeue()
		if sg == nil {
			break
		}
		sg.elem = nil
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = unsafe.Pointer(sg)
		sg.success = false
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		glist.push(gp)
	}
	unlock(&c.lock)

	// 唤醒所有的g
	for !glist.empty() {
		gp := glist.pop()
		gp.schedlink = 0
		goready(gp, 3)
	}
}
```

综上，就是go语言的chan的底层原理了，谷歌这帮人太6了

 