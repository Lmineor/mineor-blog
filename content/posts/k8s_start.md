---
title: "《k8s权威指南学习》--入门篇"
date: 2022-11-10
draft: true
tags : [                    # 文章所属标签
    "k8s",
]
categories : [              # 文章所属标签
    "技术",
]
---


# Service

在k8s中，Service（服务）是分布式集群架构的核心，一个Service对象拥有如下关键特征：

- 拥有一个唯一指定的名字（比如mysql-server）
- 拥有一个虚拟IP（Cluster IP、service IP或VIP）和端口号。
- 能够提供某种远程服务能力。
- 被映射到了提供这种服务能力的一组容器应用上。

为了建立Service和Pod之间的关联关系，k8s给每个Pod贴上标签（Label），如运行MySQL的Pod贴上name=mysql标签。然后给相应的Service定义标签选择器（Label Selector）

# Pod

Pod里的容器共享Pause的网络栈和Volume，因此他们之间的通信和数据交换更为高效，在设计时我们可以利用这一特性将一组密切相关的服务进程放到一个Pod中

> 注意：并不是每个Pod和它里面运行的容器都能“映射”到一个Service上，只有那些提供服务（无论是对内还是对外）的一组Pod才会被“映射”成一个服务。

# 节点

## Master

服务如下：
kube-apiserver、kube-controller-manager和kube-scheduler

作用：
实现整个集群的资源管理、Pod调度、弹性伸缩、安全控制、系统监控和纠错等管理功能。

## Node

集群的工作节点，运行真正的应用程序。
在Node上k8s管理的最小运行单元是Pod。

服务如下：
kubelet、kube-proxy服务进程

作用：
负责Pod的创建、启动、监控、重启、销毁及实现软件模式的负载均衡。

# 扩容

在k8s集群中，只需为需要扩容的Service关联的Pod创建一个RC（Replication Controller）。

RC定义文件中包含三个关键信息：
1. 目标Pod的定义。
2. 目标Pod需要运行的副本数量（Replicas）。
3. 要监控的目标Pod的标签（Label）。