---
title: "centos检测新的硬盘"
date: 2023-10-23
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---


如下命令向所有的SCSI主机发送一个扫描请求，可以帮助系统重新检测新添加的硬盘。


[root@compute1 ~]# echo "- - -" > /sys/class/scsi_host/host0/scan
[root@compute1 ~]# echo "- - -" > /sys/class/scsi_host/host1/scan
[root@compute1 ~]# echo "- - -" > /sys/class/scsi_host/host2/scan
