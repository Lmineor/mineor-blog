---
title: "k8s网络模型"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "k8s",
]
---


# 基础原则

- 每个Pod都拥有一个独立的IP地址，而且假定所有Pod都在一个可以直接连通的、扁平的网络空间中，不管是否运行在同一Node上都可以通过Pod的IP来访问。
- 对于支持主机网络的平台,其Pod若采用主机网络方式,则Pod仍然可以不通过NAT的方式访问其余的Pod
- k8s中Pod的IP是最小粒度IP。同一个Pod内所有的容器共享一个网络堆栈，该模型称为IP-per-Pod模型。
- Pod由docker0实际分配的IP，Pod内部看到的IP地址和端口与外部保持一致。同一个Pod内的不同容器共享网络，可以通过localhost来访问对方的端口，类似同一个VM内的不同进程。
- IP-per-Pod模型从端口分配、域名解析、服务发现、负载均衡、应用配置等角度看，Pod可以看作是一台独立的VM或物理机。


# Kubernetes网络地址4大关注点

1. Pod内的容器可以通过loopback口通信
2. 集群网络为Pods间的通信提供条件
3. Service API暴露集群中的Pod的应用给外部,以便外部使用
4. 也可只使用Services在集群内部使用
