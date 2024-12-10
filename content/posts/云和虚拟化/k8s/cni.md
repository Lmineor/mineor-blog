---
title: "CNI插件"
date: 2022-12-10
draft: false
tags : [                    # 文章所属标签
    "k8s",
    "CNI"
]
---

> 参考文章：https://developer.aliyun.com/learning/course/572/detail/7866

# 如何开发自己的CNI插件

CNI插件的实现通常包含两个部分：
1. 一个二进制的CNI插件去配置Pod网卡和IP地址。这一步配置完成后相当于给Pod插上了一条网线：有了自己的IP、自己的网卡；
2. 一个Daemon进程去管理Pod之间的网络打通。这一步相当于将Pod真正连上网络，让Pod之间能够互通。

## 给Pod插上网线

1. 给Pod准备虚拟网卡
    - 创建“veth”虚拟网卡对
    - 将一端的网卡挪到Pod中
2. 给Pod分配IP地址
    - 给Pod分配集群中唯一的IP地址
    - 一般把Pod网段按Node分段
    - 每个Pod再从Node段中分配IP
3. 配置Pod的IP和路由
    - 给Pod的虚拟网卡网址分配到的IP
    - 给Pod的网卡上配置集群网段的路由
    - 在宿主机上配置到Pod的IP地址的路由到对端虚拟网卡上

## 给Pod连上网络

刚才是给Pod插上网线，也就是说分配了IP地址和路由表。接下来说明怎么让每一个Pod的IP地址在集群里都能被访问到。

一般是在CNI的daemon进程中去做这些网络打通的事情。

- 首先CNI在每个节点上运行的daemon进程会学习到集群所有Pod的IP地址及其所在节点的信息。学习的方式通过监听K8s APIserver，拿到现有Pod的IP地址以及节点，并且新的节点和新的Pod在创建的时候也能通知到每个daemon；
- 拿到Pod以及Node相关信息后，再去配置网络进行打通。
    - 首先daemon回创建到整个集群所有节点的通道。这里的通道是个抽象的概念， 具体实现一般是通过overlay隧道等。
    - 第二部是将所有Pod的IP地址跟上一步创建的通道关联起来。关联也是个抽象的概念，具体实现通常是通过linux路由、fdb转发表或者ovs流表完成的。