---
title: "k8s网络"
date: 2023-05-08
draft: false
tags : [                    # 文章所属标签
    "Linux", "k8s"
]
categories : [              # 文章所属标签
    "技术",
]
---

k8s网络模型设计的一个基础原则是：每个Pod都拥有一个独立的IP地址，而且假定所有Pod都在一个可以直接连通的、扁平的网络空间中。

所以不管他们是否运行在同一个Node（宿主机中），都要求他们可以直接通过对方的IP进行访问。
设计这个原则的原因是，用户不需要额外考虑如何建立Pod之间的连接，也不需要考虑将容器端口映射到主机端口等问题。

# Tips

1. 查看一个网卡是否开启了混杂模式

```bash
[root@centos ~]# ifconfig eth0
eth0: flags=4163<UP,BROADCAST,RUNNING,PROMISC,MULTICAST>  mtu 1500
        inet 172.17.36.2  netmask 255.255.240.0  broadcast 172.17.47.255
        ether 00:16:3e:03:f4:76  txqueuelen 1000  (Ethernet)
```

当输出包含PROMISC时，表明该网络接口处于混杂模式。

```bash
ifconfig eth0 promisc # 启用网卡的混杂模式
ifconfig eth0 -promisc # 关闭网卡的混杂模式
```

将网络设备加入Linux bridge后，会自动进入混杂模式。