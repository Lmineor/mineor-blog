---
title: "linux_lldptool"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

```bash
[root@localhost ~]# lldptool set-lldp -i eno2 adminStatus=rxtx;
[root@localhost ~]# lldptool -T -i eno2 -V sysName enableTx=yes;
[root@localhost ~]# lldptool -T -i eno2 -V portDesc enableTx=yes;
[root@localhost ~]# lldptool -T -i eno2 -V sysDesc enableTx=yes;
[root@localhost ~]# lldptool -T -i eno2 -V sysCap enableTx=yes;

```