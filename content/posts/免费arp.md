---
title: "免费arp"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "数通",
]
categories : [              # 文章所属标签
    "技术",
]
---

# 免费(Free ARP)的作用是什么

解决方案
A：主机主动使用自己的IP地址作为目标地址发送ARP请求，此种方式称免费ARP。免费ARP有两个方面的作用：

1. 用于检查重复的IP地址 
正常情况下应当不能收到ARP回应，如果收到，则表明本网络中存在与自身IP地址重复的地址。

2.用于通告一个新的MAC地址
发送方换了块网卡，MAC地址变了，为了能够在ARP表项老化前就通告所有主机，发送方可以发送一个免费ARP。

# arp抓包

执行arping的主机的ip地址为99.0.85.123

1. 免费arp

```bash
[root@lex ~]# arping -I ens160 99.0.85.123
[root@lex ~]# tcpdump  arp -i ens192 host 99.0.85.123 -w arp.pcap
```

抓包结果

![arp2](https://blog.mineor.xyz/images/arp2.png)

2. 一般arp

```bash
[root@lex ~]# arping -I ens160 99.0.85.13
[root@lex ~]# tcpdump  arp -i ens192 host 99.0.85.123 -w arp.pcap
```

抓包结果

![arp1](https://blog.mineor.xyz/images/arp1.png)