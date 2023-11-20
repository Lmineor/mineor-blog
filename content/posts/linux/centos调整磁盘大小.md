---
title: "CentOS调整磁盘大小【磁盘空间移动】"
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
lvreduce -L 100G /dev/mapper/centos-home  # home目录减少100G
lvextend -l +100%FREE /dev/mapper/centos-root # root目录增加所有空闲的
xfs_growfs /dev/mapper/centos-root # 将实际分区落盘
```