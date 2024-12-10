---
title: "Calico容器网络"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "k8s",
    "CNI"
]
---

# 简介

Calico是一个基于BGP的纯三层的网络方案,与OpenStack,Kubernetes,AWS,GCE等云平台都能够良好地集成.

Calico在每个计算节点利用Linux kernel实现了一个高效的vrouter来负责转发.每个vrouter通过BGP1协议把在本节点上运行的容器的路由信息向整个calico网络广播,并自动设置到达其它节点的路由转发规则.

Calico保证所有容器之间的数据流量都是通过IP路由的方式完成互连互通的.Calico节点组网可以直接利用数据中心的网络结构(L2或者L3),不需要额外的NAT,隧道或者overlay network,没有额外的封包解包,能够节约CPU运算,提高网络效率.
