---
title: "cgroup"
date: 2022-12-10
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---


> 参考：https://zhuanlan.zhihu.com/p/434731896、https://www.cnblogs.com/zhrx/p/16388175.html

# cgroup

目前我们所提到的容器技术、虚拟化技术（不论何种抽象层次下的虚拟化技术）都能做到资源层面上的隔离和限制。

对于容器技术而言，它实现资源层面上的限制和隔离，依赖于linux内核所提供的cgroup和namespace技术。

两项技术的概括

- cgroup主要作用：管理资源的分配、限制；
- namespace的主要作用：封装抽象，限制，隔离，使命名空间内的进程开起来拥有他们自己的全局资源；

**cgroup是Linux内核的一个功能，用来限制、控制与分离一个进程组的资源（如CPU、内存、磁盘、输入输出等）。**

cgroup需要限制的资源是：
- CPU
- 内存
- 网络
- 磁盘I/O

作用：

- 资源限制：可以配置cgroup，从而限制进程可以对特性资源的使用量
- 优先级：当资源发生冲突时，可以控制一个进程相比另一个cgroup中的进程可以使用的资源量（CPU、磁盘或网络等）
- 记录：在cgroup级别监控和报告资源限制
- 控制：可以使用单个命令更改cgroup中所有进程的状态（冻结、停止或重新启动）

# 依赖的四个核心概念

- 子系统
- 控制组
- 层技树
- 任务

## 控制组（group）

表示一组进程和一组带有参数的子系统的关联关系。例如，一个进程使用了CPU子系统来限制CPU的使用时间，则这个进程和CPU子系统的关联关系称为控制组。

## 层级树

有 一系列的控制组按照树状结构排列组成的。这种排列方式可以使得控制组拥有父子关系，子控制组默认拥有父控制组的属性，也就是子控制组会继承父控制组。

比如，系统中定义了一组控制组c1，限制了CPU可以使用1核，然后另一个控制组c2想实现既限制CPU使用1核，同时限制内存使用2G，那么c2就可以直接继承c1，无需重复定义CPU限制。

## 子系统

一个内核的组件，一个系统代表一类资源调度控制器。例如内存子系统可以限制内存的使用量，CPU子系统可以限制CPU的使用时间。

子系统是真正实现某类资源的限制的基础。

subsystems（子系统）cgroups中的子系统就是一个资源调度控制器（又叫controllers）。
在/sys/fs/cgroup/这个目录下可以看到cgroup子系统

```bash
[root@centos ~]# ll /sys/fs/cgroup/
total 0
dr-xr-xr-x 5 root root  0 Oct 23 06:38 blkio # 为块设备设定输入输出限制，比如物理驱动设备
lrwxrwxrwx 1 root root 11 Oct 23 06:38 cpu -> cpu,cpuacct # 使用调度控制程序控制对CPU的使用
lrwxrwxrwx 1 root root 11 Oct 23 06:38 cpuacct -> cpu,cpuacct # 自动生成cgroup中任务对cpu资源使用情况的报告
dr-xr-xr-x 3 root root  0 Oct 23 06:38 cpuset # 可以为cgroup中的任务分配独立的cpu和内存
dr-xr-xr-x 5 root root  0 Oct 23 06:38 devices # 可以开启或关闭cgroup中任务对设备的访问
dr-xr-xr-x 3 root root  0 Oct 23 06:38 freezer # 可以挂起或恢复cgroup中的任务
dr-xr-xr-x 3 root root  0 Oct 23 06:38 hugetlb
dr-xr-xr-x 6 root root  0 Oct 23 06:38 memory # 可以设定cgroup中任务对内存使用量的限定，并且自动生成这些任务对内存资源使用情况的报告
lrwxrwxrwx 1 root root 16 Oct 23 06:38 net_cls -> net_cls,net_prio # docker没有直接使用它，它通过使用等级识别符标记网络数据包，从而允许linux流量控制程序识别从具体cgroup中生成的数据包
lrwxrwxrwx 1 root root 16 Oct 23 06:38 net_prio -> net_cls,net_prio
dr-xr-xr-x 3 root root  0 Oct 23 06:38 perf_event # 使用后使cgroup中的任务可以进行统一的性能测试
dr-xr-xr-x 5 root root  0 Oct 23 06:38 pids # 限制任务数量
dr-xr-xr-x 3 root root  0 Oct 23 06:38 rdma
dr-xr-xr-x 6 root root  0 Oct 23 06:38 systemd
```


**这篇文章好强，可以抽空研究复现下：https://www.cnblogs.com/zhrx/p/16388175.html**