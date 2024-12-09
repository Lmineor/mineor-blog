---
title: "ifconfig配置命令"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---


配置IPv6地址:

```bash
ifconfig enp9s0f1 inet6 add 81::d/64 

```

配置IPv4地址:

```bash
ifconfig enp9s0f1 add 10.0.0.6/64 
ifconfig enp9s0f1 delete 10.0.0.6/64 
```