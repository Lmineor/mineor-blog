---
title: go垃圾回收
date: 2025-04-06
draft: false
tags:
  - Go
---
原文参考：[Go语言的垃圾回收机制，图文并茂 一篇搞懂！_go垃圾回收-CSDN博客](https://blog.csdn.net/Shoulen/article/details/140456878)
# Go垃圾回收
 
## 标记清除法
> Go 1.5之前使用的垃圾回收策略是标记清除法

> 根对象就是应用程序中可以直接或间接访问的对象
   根对象主要包括：
> - 全局变量：在程序编译期间就能确定，全局变量存在于程序的整个生命周期
> - 执行栈：Go语言中协程是分配在堆上，每个协程都含有自己的执行栈
> - 寄存器：寄存器的值可能表示一个指针，这些指针可能指向某些赋值器分配的堆内存区块

**简单来说，标记清除算法整个过程就是把从根对象遍历不到的对象清除掉，根对象就是垃圾收集器判断哪些对象是活动对象（还在使用）的起点**

比如，声明一个变量 b

 var b int = 1
在垃圾回收过程中，这个变量 b 就是一个根对象

STW(Stop-The-World)，指的是系统在执行特定操作时，暂停所有应用程序线程，直到某个特定的时间发生或处理完成

### 标记清除法过程：

1、STW暂停
2、给所有可达对象做标记
3、回收不可达的对象
4、STW结束

### 标记清除算法实现很简单，缺点如下：
1. STW： Stop The World，让程序暂停
2. 标记需要扫描整个堆
3. 清除碎片会产生堆碎片
##  三色标记法（v1.5)

三色标记法把所有的内存对象分为**黑、灰、白**三类

- 白色对象（可能死亡）：未被回收器访问到的对象，即未被引用的对象。在回收开始阶段，所有对象均为白色，当回收结束后，白色对象均不可达。
- 灰色对象（波面）：已被回收器访问到的对象，但回收器需要对其中的一个或多个指针进行扫描，因为他们可能还指向白色对象。处于活跃状态，是黑色和白色的中间状态，需要判断自身引用的下游对象，所有根对象在标记开始时全部置为灰色
- 黑色对象（确定存活）：已被回收器访问到的对象，其中所有字段都已被扫描，黑色对象中任何一个指针都不可能直接指向白色对象。表示已经判断完成的对象，不会再通过该对象对其下游对象做扫描标记。

### 整个三色标记法流程分为五步：

1、将所有对象标记为白色

![gc1](https://blog.mineor.xyz/images/go/gc1.png)

2、将所有根对象标记为灰色

![gc2](https://blog.mineor.xyz/images/go/gc2.png)

3、遍历整个灰色对象集合，把能从灰色对象遍历到的白色对象标记为灰色，把自己标记为黑色（没有下游对象就直接把自己标为黑色）

![gc3](https://blog.mineor.xyz/images/go/gc3.png)

4、循环第三步，直到灰色集合为空

![gc4](https://blog.mineor.xyz/images/go/gc4.png)

5、回收所有的白色节点

![gc5](https://blog.mineor.xyz/images/go/gc5.png)

那么三色标记法可以实现并发执行回收吗，答案是不行

看一个场景：

![gc6](https://blog.mineor.xyz/images/go/gc6.png)

在执行到第四步时，由于是并发场景，遍历完根节点 A 的所有下游对象后，引用结构发生了变化，A 引用了 D、B 不再引用 D

这样整个流程下来就会导致活跃对象 D 被错误回收，因为灰色对象 B 没有引用 D 对象，引用对象 D 的对象是黑色对象 A
![gc7](https://blog.mineor.xyz/images/go/gc7.png)


可以看到错误原因是：
1. 一个白色对象被黑色对象引用了（白色挂在黑色下）
2. 之前与此白色对象连接的灰色对象连接断开了（灰色同时丢失了该白色）

同时满足这两条，才会出现**错误(丢失对象)**

所以破坏任意一条，就可以避免错误回收活跃对象，由此就有了**强三色不变性**和**弱三色不变性**

- 强三色不变性强制性的不允许黑色对象引用白色对象（破坏条件1）
- 弱三色不变性允许黑色对象引用白色对象，但是白色对象必须存在其他灰色对象对它的引用（破坏条件2）

三色标记法本身不能保证并发场景下垃圾回收的正确性，只有添加强三色不变性或者弱三色不变性限制之后，三色标记法才能保证正确性
而强不变三色性和弱不变三色性通过屏障技术来实现

# 屏障技术
屏障技术可以理解为一种回调机制，在程序的某种执行过程中加一个判断机制，满足判断机制则执行回调函数，类似于钩子函数（Hook）

对于内存操作，可以简单的分为：

- 栈对象的读
- 栈对象的写
- 堆对象的读
- 堆对象的写

实际上，垃圾回收机制只用于回收堆上的内存，栈中的内存如局部变量、函数调用等会在调用结束后自动释放

也就是说，屏障机制只能作用于堆对象

屏障机制分为插入写屏障和删除写屏障

**插入写屏障实现了强三色不变性，给对象添加引用关系时触发**

**删除写屏障实现了弱三色不变性，删除对象引用关系时触发**

## 插入写屏障
Go 1.5 垃圾回收机制由三色标记法+插入写屏障实现

插入写屏障满足了强三色不变式，是对象被引用后触发的机制：

每当一个对象被引用，就会触发判断：如果这次操作是一个白色对象被黑色对象引用，就把这个白色对象标记为灰色

在引入插入写屏障后，再执行上面的流程。在对象 A 添加对 对象 D 的引用时，插入写屏障机制会把对象 D 置灰。保证了活跃对象没有被错误回收
![gc8](https://blog.mineor.xyz/images/go/gc8.png)


上面也提到了，栈对象不启用屏障机制

假如一个栈上的对象扫描完成后又引用了堆上的白色对象，白色对象在堆中也没有其他对象引用，由于栈对象没有屏障机制，在添加引用时不会把引用的白色对象置灰，这样就会造成堆上的活跃对象被错误回收
![gc9](https://blog.mineor.xyz/images/go/gc9.png)


栈对象 A 置黑后引用堆对象 B，现在 B 为活跃对象，但是会因为标记为白色被错误回收

Go 语言的处理方法是：栈在 GC 迭代结束时（使用三色标记法反复遍历到没有灰色节点时），会对栈执行一次标记清除法（STW），重新扫描一遍栈对象，清除掉从栈对象出发访问不到的对象
![gc10](https://blog.mineor.xyz/images/go/gc10.png)


**插入写屏障的缺点**是 GC 迭代结束时，需要一次 STW 来重新扫描栈，虽然比标记清除法扫描栈和堆要好，但是仍然有优化空间

## 删除写屏障
删除写屏障满足了弱三色不变式，是对象引用关系被删除时触发的机制：

每当一个对象被删除时，就会触发判断：如果是一个灰色对象引用的白色对象被删除，那么就把这个白色对象标记为灰色

在引入删除写屏障后，第三步中对象 B 删除对象 D 的引用时，删除写屏障会把对象 D 置灰。保证了活跃对象没有被错误回收
![gc11](https://blog.mineor.xyz/images/go/gc11.png)


缺点：回收精度较低，有些本该删除的对象可能会在下一轮才会被回收。假如对象 B 删除对象 D 的引用后，对象 A 并没有引用对象 D。此时对象 D 已经是应该回收的对象了，但会因为删除连接使对象 D 置灰，导致在这一轮的回收中没有回收对象 D

此外，在引入栈对象不启用屏障机制这一限制条件后，在一些场景下会出现问题：

下图中，对象 C 属于栈对象，在 C 删除对 D 的引用时不会触发删除写屏障，所以对象 D 还是白色。随后对象 A 引用对象 D，此时对象 D 是一个活跃对象，但是会被错误回收
![gc12](https://blog.mineor.xyz/images/go/gc12.png)


为了解决这个问题，采用的方法是：在起始时，STW 扫描整个栈，把在栈上引用的所有对象都置灰

这样处理后，就能保证所有堆上在用的对象都不会错误回收，虽然会有些应该被回收的对象没回收掉，但是一轮一轮的回收机制早晚会把没用的对象回收掉

# 混合写屏障
插入写屏障和删除写屏障都有各自的短板：

插入写屏障结束时需要 STW（标记清除）重新扫描栈

删除写屏障回收精度低，在开始时需要 STW 扫描整个堆栈记录初始快照

Go 语言在 1.8 版本之后引入的混合写屏障机制结合了插入写屏障和删除写屏障，满足了变形的弱三色不变式

混合写屏障的具体操作为：

1. 开始时扫描栈上所有的可达对象全部标记为黑色
2. 在整个扫描期间，在栈上创建的新对象都标记为黑色
3. 将被删除的对象标记为灰色
4. 将被添加的对象标记为灰色

只看步骤很难理解，来看几个例子

1、对象被一个堆对象删除引用，又成为栈对象的下游
![gc13](https://blog.mineor.xyz/images/go/gc13.png)


在第三步中，D删除对E的引用时触发混合写屏障机制，把E置灰，这样就避免了E被错误回收，但是如果A没有引用E，E对象就成了垃圾，但是E对象已经触发屏障置灰，所以在当前轮不会被回收，只可能在后面几轮GC中被回收

2、对象被一个栈对象删除引用，并成为另一个栈对象的下游
![gc14](https://blog.mineor.xyz/images/go/gc14.png)


栈对象不会触发屏障机制，但是GC开始时将所有可达对象置黑就保证了栈上活跃对象不会被错误回收，在第一步，对象C为栈上新创建对象，会直接置黑。第二步中A删除B对象的引用不会触发屏障，但由于B对象已经被置黑，所以不会被回收。如果C没有引用B对象，B对象成为了垃圾，在这一轮的GC中依然不会被清除，只可能在后面几轮GC中被回收

3、对象被一个黑色堆对象添加引用
![gc15](https://blog.mineor.xyz/images/go/gc15.png)


混合写屏障会将被添加的对象标记为灰色，当对象A添加对B的引用时，会触发混合写屏障，把B置灰，这样就避免了对象B被错误回收

经过这三个例子，可以看出来，混合写屏障就是结合了插入写屏障（将被添加的对象标记为灰色）和删除写屏障（将被删除的对象标记为灰色），同时通过将栈上可达对象置黑，栈上新创建对象置黑，通过遍历保存了栈对象的起始状态，不需要 STW

总的来说，混合写屏障结合了插入写屏障和删除对象的触发机制，同时解决了开始和结束的 STW，大大提高了性能（但是回收精度依然不佳）

总结
Go 语言在整个发展过程中，垃圾回收机制的演进为：标记清除法——>三色标记法+插入写屏障->三色标记法+混合写屏障

标记清除法需要一次全局 STW，先标记可达对象，再清除所有的不可达对象，性能很差

三色标记法把对象分为了三种状态，分别用白色、灰色、黑色表示，先把所有可达对象置灰，然后遍历灰色对象，把下游对象置灰，本身置黑；重复这一过程直到没有灰色节点，把白色节点清除。三色标记法本身无法在并发条件下正确执行，单独使用也需要 STW 才能保证正确性

为了优化这一过程，避免全局 SWT，引入了插入写屏障，插入写屏障是一种触发机制，每当一个对象被引用时，如果是黑色对象引用白色对象，就把这个白色对象置灰，但是由于屏障机制无法作用到栈对象，为了避免与栈相关的活跃对象被错误回收，需要在三色标记法结束后，对栈做一次 STW，清除掉垃圾对象

插入写屏障虽然优化了全局 SWT，但是每次GC也需要对栈上的一次 STW

为了解决这一问题，最终的方案是混合写屏障，混合写屏障的触发机制有两个：被删除的对象会被标记为灰色，被添加的对象会被标记为灰色；同时在开始时会把站上所有可达对象置黑，栈上对象被创建时也会置黑。这样就记录了栈上的对象状态，避免了错误回收活跃对象

混合写屏障结合了插入写屏障和删除写屏障的优点，避免了在插入写屏障结束后需要进行全栈扫描的性能问题，实现了更高效的垃圾回收机制。但是在某些场景下仍然存在精度不足的问题