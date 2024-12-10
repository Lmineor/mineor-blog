---
title: "查看POD的IP地址"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "k8s",
]
---

# 从外面看Pod IP地址之Kubernetes API

1. `kubectl get pod`或者`kubectl describe pod`就可以
2. 如果运行在容器内的进程希望获取该容器的IP,可以通过环境变量的方式来获取IP

```bash
spec:
  containers:
    - name: env-pod
      image: busybox
      command: ["/bin/sh", "-c","env"]
      env:
      - name: POD_IP
        valueFrom:
          fieldRef:
            fieldPath: status.podIP
```

# 从外面看Pod IP之docker命令

假设容器的ID是6e8147cd2f3d, 一般情况下可以通过以下命令查询容器的IP地址:

```bash
docker inspect --format '{{ .NetworkSettings.IPAddress }}' 6e8147cd2f3d
```

但是对于Pod中的容器, 会输出空字符串.原因是Kubernetes调用的是cni插件,而docker的网络实现靠cnm,所以会导致输出有问题.

# 进入容器内部看Pod IP地址

进到容器的docker命令有docker exec或docker attach，进到容器后再执行ip addr或ifconfig这类常规命令。同一个Pod内所有容器共享网络`namespace`，因此随便找一个有ip或者ifconfig命令的容器就能很容易地查询到Pod IP地址。如果Pod内所有容器都不自带这些命令呢？
在我们这个场景下，进入容器网络namespace的作用等价于进入容器，而且还能使用宿主机的ip或者ifconfig命令。
假设Pod的pause容器ID是6e8147cd2f3d，首先获得该容器在宿主机上映射的PID，如下所示：

```bash
[root@k8s-node1 ~]# docker inspect --format '{{ .State.Pid }}' 446936c59d95
27804
[root@k8s-node1 ~]# nsenter --target 27804 --net
[root@k8s-node1 ~]# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
3: eth0@if16: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
    link/ether de:9c:6b:6c:b9:ca brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 205.205.0.2/16 brd 205.205.255.255 scope global eth0
       valid_lft forever preferred_lft forever

```

这样就可以输出pause容器的ip地址信息.