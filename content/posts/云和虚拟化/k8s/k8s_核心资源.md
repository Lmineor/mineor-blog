---
title: k8s核心资源
date: 2023-03-26
draft: false
tags:
  - k8s
---

# API Server

总体来看， Kubemetes API Server 的核心功能是提供了Kubemetes 各类资源对象（如Pod 、RC 、Service 等〉的增、删、改、查及Watch 等HTTP Rest 接口，成为集群内各个功能模块之间数据交互和通信的中心枢纽，是整个系统的数据总线和数据中心。
除此之外，它还有以下一些功能特性。

1. 是集群管理的API入口。
2. 是资源配额控制的入口。
3. 提供了完备的集群安全机制。
# Controller Manager

智能系统和自动系统通常会通过一个操作系统来不断修正系统的工作状态。在k8s中，每个Controller都是这样一个操作系统，它们通过API Server提供的（`List-Watch`）接口实时监控集群中特定资源的状态变化，当发生各种故障导致某资源对象的状态发生变化时，Controller会尝试将其状态调整为期望状态。

Kubernetes 中的 `Controller Manager` 是集群管理中的核心组件之一，它是一个聚合了多个控制器功能的进程，运行在 Master 节点上，负责维护集群的**期望状态**（`desired state`）与**实际状态**（`current state`）的一致性。

`Controller Manager`内部包含`Replication Controller`、`Node Controller`、`ResourceQuota Controller`、`Namespace Controller`、`ServiceAccount Controller`、`Token Controller`、`Service Controller`及`Endpoint Controoler`八种。而`Controller Manager`正式这些Controller的核心管理者。

## Replication Controller

`Replication Controller` 核心作用确保任何时候集群中的RC关联的Pod副本数量都保持预设值，多退少补原则。需要注意的是，只有当重启策略是Always时，才起作用。

最好不要越过RC直接创建Pod，因为RC会管理Pod的副本，这样可以提高系统的容灾能力。

- 确保在当前集群中有且仅有N个Pod实例，N是在RC中定义的Pod副本数量。
- 调整RC的`spec.replicas`属性值来实现系统扩缩容。
- 改变RC中的Pod模版（镜像版本）来实现系统的滚动升级。

`Replication Controller`使用场景：

1. 重新调度
2. 弹性伸缩
3. 滚动更新

## Node Controller

Kubelet 进程在启动时通过API Server注册自身的节点信息，并定时向API Server汇报状态，API Server收到这些信息后会通过到`ETCD`中。其信息包括节点健康状况、节点资源、节点名称、节点地址信息、操作系统版本、docker版本、kubelet版本等。

`Node Controller`通过API Server实时获取Node的相关信息，实现管理和监控集群中的各个Node的相关控制功能。

![10ee71341309e250f7de66bfb528c452.jpeg](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/39d99ac8bc8641e28444535a6770f48a~tplv-k3u1fbpfcp-jj-mark:3024:0:0:0:q75.awebp#?w=716&h=500&s=167987&e=png&b=fdfdfd)

**初始化与启动**：

- `Node Controller` 是 Kubernetes `Controller Manager` 中的一部分，在 `Controller Manager` 启动时会被初始化和启动。
- `Node Controlle`r 初始化时，会根据给定的配置构建自身的大结构体，完成必要的参数配置。

**监听节点状态**：

- `Node Controller` 通过 API Server 获取集群中所有 Node（节点）的实时状态信息。
- 为 `nodeInformer` 配置回调函数（AddFunc, UpdateFunc, DeleteFunc），当有节点新增、更新或删除时，`Node Controller` 能够得到通知。

**节点心跳监测**：

- `Node Controller` 通过检查节点发送的心跳（通常是通过节点定期更新其 Lease 对象实现）来确定节点是否正常工作。
- 如果节点长时间未发送心跳，`Node Controller` 将判断节点可能已宕机或失去联系，并采取相应措施。

**节点健康状态处理**：

- 当节点未响应或报告为非健康状态时，`Node Controller` 会先尝试标记节点为不可调度（Taint节点，防止新 Pod 被调度到该节点上）。
- 若节点长时间未恢复，`Node Controller` 会进一步操作，如驱逐（Evict）节点上运行的所有 Pods，让它们在其他健康节点上重新调度。

**节点清理**：

- 对于长时间未响应且确认无法恢复的节点，`Node Controller` 可能会从集群中彻底删除该节点记录，以便释放相关资源和名称空间。

**节点配置更新同步**：

- 当节点的配置发生变化时，如节点容量（`capacity`）更新或标签（`labels`）变化，`Node Controller` 会确保集群中的相关信息得到及时更新和同步。

**与其它控制器协作**：

- `Node Controller` 还与其他控制器交互，例如 `DaemonSet Controller`，确保守护进程集在每个节点上的正确部署和管理。

## ResourceQuota Controller

资源配额管理确保了指定的资源对象在任何时候都不会超量占用系统物理资源，避免了由于某些业务进程的设计和实现的缺陷导致整个系统瘫痪。

目前k8s支持如下三个层次的资源配额

1. 容器级别
2. Pod级别
3. Namespace级别
    - Pod数量
    - Replication Controller数量
    - Service数量
    - ResourceQuota数量
    - Secret数量
    - PV数量

其实现方式通过 `Admission Control`（准入控制）来控制的，当前提供了两种方式的配额约束

- `LimitRanger`：作用于Pod和Container。
- `ResourceQuota`：作用于Namespace。

如果Pod定义中声明了`LimitRanger`，则用户通过API Server请求创建或者修改资源时，`Admission Control`会计算当前配合的使用情况，如果不符合条件则创建失败。

对应定义了`ResourceQuota`的`Namespace`，`ResourceQuota Controller`组件会定义统计和生成该Namespace下的所有资源使用量，将统计结果写入ETCD。随后统计数据被`Admission Control`使用，以确保当前`Namespace`下的资源配额总量不会超过`ResourceQuota`中的限定值。

![591aee0a0d9dda0033785236749fd260.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/c453c601af4342638b2fe2b94039b8df~tplv-k3u1fbpfcp-jj-mark:3024:0:0:0:q75.awebp#?w=901&h=708&s=40500&e=png&b=fefdfd)

## Namespace Controller

API Server可以创建新的Namespace并将其保存到`ETCD`中，`Namespace Controller`定时通过API Server读取这些信息。如果Namespace被API标记为删除，则将该Namespace状态设置成Terminating并保存到ETCD中。同时，Namespace Controller删除其下的所有资源，最后对Namespace执行`finalize`操作，删除`spec.finalizers` 域中的信息。

## Service Controller和Endpoints Controller

Service、Endpoints与Pod关系，如下图：

![9da1d529d73a95a83322b8aea10f5dbc.jpeg](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/b5d9dd859adb49b18899fcc2c31611e5~tplv-k3u1fbpfcp-jj-mark:3024:0:0:0:q75.awebp#?w=500&h=537&s=20895&e=jpg&b=fcfcfc)

`Endpoints`表示一个Service对应的所有Pod副本的访问地址，`Endpoints Controller`负责生成和维护所有Endpoints对象的控制器。它负责监听Service和对应的Pod副本变化，如果监听到Service被删除，则删除和该Service同名的Endpoints对象。如果检测到被创建或则修改，则根据该Service信息获得相关的Pod列表，然后创建或者更新Endpoints对象。更新或者身处Pod，同理。

Endpoints对象在哪里被使用呢？答案是Node上的kube-proxy进程，它获取每个Endpoints，实现Service的负载均衡。

  

作者：翠竹静斋  
链接：https://juejin.cn/post/7359138355180847131  
来源：稀土掘金  
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。