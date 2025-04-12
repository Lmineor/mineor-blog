---
title: Go内存泄露
date: 2025-04-06
draft: 
tags:
  - Go
---
# 内存泄露发生的可能情况

简单归纳一下，还是“临时性”内存泄露和“永久性”内存泄露：

临时性泄露，指的是该释放的内存资源没有及时释放，对应的内存资源仍然有机会在更晚些时候被释放，即便如此在内存资源紧张情况下，也会是个问题。这类主要是string、slice底层buffer的错误共享，导致无用数据对象无法及时释放，或者defer函数导致的资源没有及时释放。
永久性泄露，指的是在进程后续生命周期内，泄露的内存都没有机会回收，如goroutine内部预期之外的 for-loop 或者 chan select-case 导致的无法退出的情况，导致协程栈及引用内存永久泄露问题。

## 暂时性内存泄露

- 获取长字符串中的一段导致长字符串未释放
- 获取长slice中的一段导致长slice未释放
- 在长slice新建slice导致泄漏

string相比于切片少了一个容量的cap字段，可以把string当成一个只读的切片类型。获取长string或切片中的一段内容，由于新生成的对象和老的string或切片共用一个内存空间，  
会导致老的string和切片资源暂时得不到释放，造成短暂的内存泄露。

## 永久性内存泄露

- goroutine泄漏
- time.Ticker未关闭导致泄漏
- Finalizer导致泄漏
- Deferring Function Call导致泄漏

常见的内存泄露场景，go101进行了讨论，总结了如下几种：

[Kind of memory leaking caused by substrings](https://go101.org/article/memory-leaking.html)
[Kind of memory leaking caused by subslices](https://go101.org/article/memory-leaking.html)
[Kind of memory leaking caused by not resetting pointers in lost slice elements](https://go101.org/article/memory-leaking.html)
[Real memory leaking caused by hanging goroutines](https://go101.org/article/memory-leaking.html)
[real memory leadking caused by not stopping time.Ticker values which are not used any more](https://go101.org/article/memory-leaking.html)
[Real memory leaking caused by using finalizers improperly](https://go101.org/article/memory-leaking.html)
[Kind of resource leaking by deferring function calls](https://go101.org/article/memory-leaking.html)




# 下面以实际代码举了几个小例子

## 数组的错误使用
由于数组是Golang的基本数据类型，每个数组占用不同的内存空间，生命周期互不干扰，很难出现内存泄漏的情况。但是数组作为形参传输时，遵循的是值拷贝，如果函数被多次调用且数组过大时，则会导致内存使用激增。

```text
//统计nums中target出现的次数
func countTarget(nums [1000000]int, target int) int {
    num := 0
    for i := 0; i < len(nums) && nums[i] == target; i++ {
        num++
    }
    return num
}
```

例如上面的函数中，每次调用countTarget函数传参时都需要新建一个大小为100万的int数组，大约为8MB内存，如果在短时间内调用100次就需要约800MB的内存空间了。（未达到GC时间或者GC阀值是不会触发GC的）如果是在高并发场景下每个协程都同时调用该函数，内存占用量是非常恐怖的。

对于大数组放在形参场景下，通常使用切片或者指针进行传递，避免短时间的内存使用激增。

## Goroutine引起的内存泄漏，未及时退出
实际开发中更多的还是Goroutine引起的内存泄漏，因为Goroutine的创建非常简单，通过关键字go即可创建，由于开发的进度大部分程序猿只会关心代码的功能是否实现，很少会关心Goroutine何时退出。如果Goroutine在执行时被阻塞而无法退出，就会导致Goroutine的内存泄漏，一个Goroutine的最低栈大小为2KB，在高并发的场景下，对内存的消耗也是非常恐怖的！

**互斥锁未释放**

```text
 1//协程拿到锁未释放，其他协程获取锁会阻塞
 2func mutexTest() {
 3    mutex := sync.Mutex{}
 4    for i := 0; i < 10; i++ {
 5        go func() {
 6            mutex.Lock()
 7            fmt.Printf("%d goroutine get mutex", i)
 8      //模拟实际开发中的操作耗时
 9            time.Sleep(100 * time.Millisecond)
10        }()
11    }
12    time.Sleep(10 * time.Second)
13}
```

**死锁**

```text
 1func mutexTest() {
 2    m1, m2 := sync.Mutex{}, sync.RWMutex{}
 3  //g1得到锁1去获取锁2
 4    go func() {
 5        m1.Lock()
 6        fmt.Println("g1 get m1")
 7        time.Sleep(1 * time.Second)
 8        m2.Lock()
 9        fmt.Println("g1 get m2")
10    }()
11    //g2得到锁2去获取锁1
12    go func() {
13        m2.Lock()
14        fmt.Println("g2 get m2")
15        time.Sleep(1 * time.Second)
16        m1.Lock()
17        fmt.Println("g2 get m1")
18    }()
19  //其余协程获取锁都会失败
20    go func() {
21        m1.Lock()
22        fmt.Println("g3 get m1")
23    }()
24    time.Sleep(10 * time.Second)
25}
```

**空channel**

```text
 1func channelTest() {
 2  //声明未初始化的channel读写都会阻塞
 3    var c chan int
 4  //向channel中写数据
 5    go func() {
 6        c <- 1
 7        fmt.Println("g1 send succeed")
 8        time.Sleep(1 * time.Second)
 9    }()
10  //从channel中读数据
11    go func() {
12        <-c
13        fmt.Println("g2 receive succeed")
14        time.Sleep(1 * time.Second)
15    }()
16    time.Sleep(10 * time.Second)
17}
```

**能出不能进**

```text
 1func channelTest() {
 2    var c = make(chan int)
 3  //10个协程向channel中写数据
 4    for i := 0; i < 10; i++ {
 5        go func() {
 6            c <- 1
 7            fmt.Println("g1 send succeed")
 8            time.Sleep(1 * time.Second)
 9        }()
10    }
11  //1个协程丛channel中读数据
12    go func() {
13        <-c
14        fmt.Println("g2 receive succeed")
15        time.Sleep(1 * time.Second)
16    }()
17  //会有写的9个协程阻塞得不到释放
18    time.Sleep(10 * time.Second)
19}
```

**能进不能出**

```text
 1func channelTest() {
 2    var c = make(chan int)
 3  //10个协程向channel中读数据
 4    for i := 0; i < 10; i++ {
 5        go func() {
 6            <- c
 7            fmt.Println("g1 receive succeed")
 8            time.Sleep(1 * time.Second)
 9        }()
10    }
11  //1个协程丛channel写读数据
12    go func() {
13        c <- 1
14        fmt.Println("g2 send succeed")
15        time.Sleep(1 * time.Second)
16    }()
17  //会有读的9个协程阻塞得不到释放
18    time.Sleep(10 * time.Second)
19}
```
## time.Ticker未及时调用stop导致

time.Ticker是每隔指定的时间就会向通道内写数据。作为循环触发器，必须调用stop方法才会停止，从而被GC掉，否则会一直占用内存空间。

```text
 1func tickerTest() {
 2    //定义一个ticker，每隔500毫秒触发
 3    ticker := time.NewTicker(time.Second * 1)
 4    //Ticker触发
 5    go func() {
 6        for t := range ticker.C {
 7            fmt.Println("ticker被触发", t)
 8        }
 9    }()
10
11    time.Sleep(time.Second * 10)
12    //停止ticker
13    ticker.Stop()
14}
```
// go垃圾回收
// go多线程，多进程怎么用，与go routing的区别
// 工作中与同事争论，怎么去解决争论
