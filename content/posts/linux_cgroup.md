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


> 参考：https://zhuanlan.zhihu.com/p/434731896

# cgroup

目前我们所提到的容器技术、虚拟化技术（不论何种抽象层次下的虚拟化技术）都能做到资源层面上的隔离和限制。

对于容器技术而言，它实现资源层面上的限制和隔离，依赖于linux内核所提供的cgroup和namespace技术。

两项技术的概括

- cgroup主要作用：管理资源的分配、限制；
- namespace的主要作用：封装抽象，限制，隔离，使命名空间内的进程开起来拥有他们自己的全局资源；

cgroup是Linux内核的一个功能，用来限制、控制与分离一个进程组的资源（如CPU、内存、磁盘、输入输出等）。

cgroup需要限制的资源是：
- CPU
- 内存
- 网络
- 磁盘I/O

# cgroup的组成

cgroup代表“控制组”，并且不会使用大写，是一种分层组织进程的机制，沿层次结构以及受控的方式分配系统资源。

主要组成：
- core 负责分层组织过程
- controller 通常负责沿层次结构分配特性类型的系统资源。每个cgroup都有一个cgroup.controllers文件，其中列出了所有可供cgroup启动的控制器。当在cgroup.subtree_control中指定多个控制器时，要没全部成功，耀目全部失败。在同一个控制器上指定多项操作，那么只有最后一个生效。每个cgroup的控制器销毁是异步的，在引用时同样也有着延迟引用的问题；