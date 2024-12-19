---
title: "修改Centos官方云镜像的root密码"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "云与虚拟化", 
]
---


https://developer.aliyun.com/article/799104

还可以看看这个

```bash
yum -y install libguestfs-tools

要设置root密码，请使用以下命令：
virt-customize -a CentOS-7-x86_64-GenericCloud.qcow2 --root-password password:123456

[   0.0] Examining the guest ...

[   1.9] Setting a random seed

[   1.9] Setting passwords

[   6.8] Finishing off

注：

CentOS-7-x86_64-GenericCloud.qcow2是要修改图像的名称。

123456是为root用户设置的密码。
```
