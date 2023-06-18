---
title: "hypervisor"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
    "SDN"
]
categories : [              # 文章所属标签
    "技术",
]
---


Hypervisor:一种运行在物理服务器和操作系统之间的中间层软件,可以允许多个操作系统和应用共享一套基础物理硬件.
可以将Hypervisor看做是虚拟环境中的"元"操作系统,可以协调访问服务器上的所有物理设备的虚拟机,所以又称为虚拟机监视器(virtual machine monitor).
Hypervisor是所有虚拟化技术的核心,非中断的支持多工作负载迁移是Hypervisor的基本功能.

当服务器启动并执行Hypervisor时,会给每一台虚拟机分配适量的内存,cpu,网络和磁盘资源,并且加载所有虚拟机的客户操作系统.

Hypervisor之于操作系统类似于操作系统之于进程.他们为执行提供独立的虚拟硬件平台,而虚拟硬件平台反过来有提供对底层机器的虚拟的完整访问.但并不是所有Hypervisor都是一样的.

# 虚拟化和Hypervisor

> **虚拟化**就是通过某种方式隐藏底层物理硬件的过程,从而让多个操作系统可以透明的使用和共享它.

## Hypervisor分类

- 类型1:这种Hypervisor运行在物理硬件之上
- 类型2:Hypervisor运行在另一个操作系统(运行在物理硬件上)中.

类型1例子:基于内核的虚拟机(kvm---它本身是一个基于操作系统的Hypervisor)
类型2例子:包括qemu和wine.
