---
title: "CRI"
date: 2022-12-10
draft: false
tags : [                    # 文章所属标签
    "容器", 
]
categories : [              # 文章所属标签
    "技术",
]
---

> 参考文章：https://zhuanlan.zhihu.com/p/102897620

## CRI是什么

CRI： Container Runtime Interface，容器运行时接口；

CRI包括Protocol Buffers、gRPC API、运行库支持及开发中的标准规范和工具，是以容器为中心设计的API，设计CRI的初衷是不希望向容器（比如docker）暴露pod信息或pod的api。

CRI工作在kubelet与container runtime之间，目前常见的runtime有：

- docker
- containerd： 可以通过shim对接不同的low-level runtime
- cri-o：一种轻量级的runtime，支持runc和Clear container作为low-level runtimes。

## CRI是如何工作的

CRI大体包含三部分接口：Sandbox、Container和Image，其中提供了一些操作容器的通用接口，包括Create、Delete、List等。

Sanbox为container提供一定的运行环境，这其中包括pod的网络等。Container包括容器生命周期的具体操作，Image则提供对镜像的操作。

kubelet回通过gRPC调用CRI接口，首先去创建一个环境，也就是所谓的PodSandbox。当Podsandbox可用后，继续调用image或container接口去拉取镜像和创建容器。

## PodSandbox

从虚拟机和容器化两方面看，两者懂事有了cgroups[cgroup](cgroup.md)做资源配额，而且概念上都抽离出一个隔离的运行时环境，只是区别在于资源隔离的实现。
因此sandbox是k8s为兼容不同运行时环境所预留的空间，也就是说k8s允许low-level runtime依据不同的实现去创建不同的podsandbox，对于kata来说podsandbox就是虚拟机，对于docker来说就是linxe namespace。
当pod sandbox建立起来后，kubelet就可以在里面创建用户容器。当删除pod时，kubelet会先移除pod sandbox然后再停止里面的所有容器，对于container来说，当sandbox运行后，只需将新的container的namespace加入到已有的sandbox的 namespace中。

在默认情况下，cri体系里，pod sandbox其实就是pause容器。kubelet代码引用的defaultSandboxImage就是官方提供的gcr.io/google_containers/pause-amd64 镜像