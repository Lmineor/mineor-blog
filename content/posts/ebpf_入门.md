---
title: "eBPF入门"
date: 2023-06-24
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---


# eBPF特性

eBPF程序都是事件驱动的，他们会在内核或者应用程序经过某个确定的Hook点的时候运行，这些Hook点都是提前定义的，包括系统调用、函数进入/退出、内核tracepoints、网络事件等。

