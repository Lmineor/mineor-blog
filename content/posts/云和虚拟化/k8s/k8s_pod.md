---
title: "《k8s权威指南学习》--Pod"
date: 2023-03-26
draft: false
tags : [                    # 文章所属标签
    "k8s",
]
---


## Pod 生命周期

Pod 在整个生命周期过程中被系统定义为各种状态
Pod 的状态如表2.14 所示。
表2.14 Pod 的状态

|状态值|描述|
|:-|:-|
|Pending| API Server已经创建该Pod,但Pod内还有一个或多个容器的镜像没有创建，包括正在下载镜像的过程|
|Running| Pod 内所有容器均己创建，且至少有一个容器处于运行状态、正在启动状态或正在重启状态|
|Succeeded| Pod 内所有容器均成功执行退出， 且不会再重启|
|Failed| Pod 内所有容器均已退出，但至少有一个容器退出为失败状态|
|Unknown| 由于某种原因无法获取该Pod 的状态， 可能由于网络通信不畅导致|

## Pod重启策略
Pod 的重启策略（ RestartPolicy ）应用于Pod 内的所有容器，井且仅在Pod 所处的Node上由kubelet 进行判断和重启操作。当某个容器异常退出或者健康检查（详见下节）失败时， kubelet将根据RestartPolicy 的设置来进行相应的操作。
Pod 的重启策略包括Always、OnFailure和Never， 默认值为Always：

- Always：当容器失效时，由kubelet自动重启该容器。
- OnFailure：当容器终止运行且退出码不为0时，由kubelet自动重启该容器。
- Never ：不论容器运行状态如何， kubelet 都不会重启该容器。

kubelet 重启失效容器的时间间隔以sync-frequency 乘以2n 来计算；例如1、2 、4 、8 倍等，
最长延时5min ，并且在成功重启后的10min后重置该时间。

## Pod健康检查

可以通过两类探针来检查： LivenessProbe 和ReadinessProbe

- LivenessProbe 探针：用于判断容器是否存活（ running 状态），如果LivenessProbe 探针探测到容器不健康，则kubelet 将杀掉该容器，并根据容器的重启策略做相应的处理。如果一个容器不包含LivenessProbe 探针，那么kubelet 认为该容器的LivenessProbe 探针返回的值永远是“ Success"
- ReadinessProbe 探针：用于判断容器是否启动完成（ ready 状态），可以接收请求。如果ReadinessProbe 探针检测到失败，则Pod 的状态将被修改。Endpoint Con位oiler 将从Service 的Endpoint 中删除包含该容器所在Pod 的Endpoint 。

对于每种探测方式，都需要设置initialDelaySeconds和timeoutSeconds两个参数，它们的含
义分别如下。
- initialDelaySeconds：启动容器后进行首次健康检查的等待时间，单位为s。
- timeoutSeconds：健康检查发送请求后等待响应的超时时间，单位为s 。当超时发生时，kubelet会认为容器己经无法提供服务，将会重启该容器。

## Pod调度

在kubernets中，Pod在大部分场景下都只是容器的载体而已，通常需要通过Deployment、DaemonSet、RC、Job等对象来完成一组Pod的调度与自动控制功能。

1. Deployment/RC ：全自动调度

deployment或RC的主要功能之一就是自动部署一个容器应用的多份副本，以及持续监控副本的数量，在集群内始终维持用户指定的副本数量。


下面是一个Deployment配置的例子， 使用这个配置文件可以创建一个ReplicaSet，这个ReplicaSet会创建3个Nginx应用的Pod 。
```yaml
nginx-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```


除了使用系统自动调度算法完成一组Pod的部署，Kubernetes也提供了多种丰富的调度策略，用户只需在Pod的定义中使用NodeSelector、NodeAffmity、PodAffinity 、Pod 驱逐等更加细粒度的调度策略设置，就能完成对Pod 的精准调度。下面对这些策略进行说明。

### NodeSelector：定向调度

kubernetes master节点上的Scheduler服务(kube-scheduler进程)负责实现Pod的调度

在实际情况中，可将Pod调度到指定的一些Node上，可以通过Node的标签（Label) 和Pod的nodeSelector
属性相匹配，来达到上述目的。

首先通过kubectl label 命令给目标Node打上一些标签：

```bash
kubectl label nodes <node-name> <label -key>=<label-value>
```

需要注意的是，如果我们指定了Pod的nodeSelector条件，且集群中不存在包含相应标签的Node,则即使集群中还有其他可供使用的Node,这个Pod 也无法被成功调度。

## Taints 和Tolerations （污点和容忍）

前面介绍的NodeAffinity节点亲和性，是在Pod上定义的一种属性，使得Pod能够被调度到某些Node上运行（优先选择或强制要求）。Taint 则正好相反一－它让Node拒绝Pod的运行。Taint需要和Toleration配合使用，让Pod避开那些不合适的Node。在Node上设置一个或多个Taint之后，除非Pod明确声明能够容忍这些“污点”，否则无法在这些Node上运行。Toleration是Pod的属性，让Pod能够（ 注意，只是能够，而非必须）运行在标注了Taint的Node上。

可以用`kubectl taint` 命令为Node设置Taint信息：

```bash
$ kubectl taint nodes node1 key=value:NoSchedule
```

这个设置为node1 加上了一个Taint 。该Taint 的键为key ，值为value, Taint 的效果是
NoSchedule 。这意味着除非Pod 明确声明可以容忍这个Taint ，否则就不会被调度到node1上去。

# Pod 探针有几种，分别作用是?

Pod 探针在 Kubernetes 中主要有三种类型：**存活探针（Liveness Probe）**、就绪探针（Readiness Probe）和启动探针（Startup Probe）。

- 存活探针（Liveness Probe）：
作用：存活探针用于检查容器内应用是否还在正常运行，即应用是否“活着”。如果存活探针失败，Kubernetes 会认为容器内的主进程已经死掉或者不再响应，这时 Kubernetes 会杀掉该容器并重新创建一个新的容器实例。
应用场景：例如，当应用由于某种原因卡死或进入了一个无法恢复的错误状态时，存活探针可以帮助系统自动重启容器以恢复服务。
- 就绪探针（Readiness Probe）：
作用：就绪探针用于检查容器内的应用是否已经准备好接受流量。只有当就绪探针成功时，Kubernetes Service 才会将该 Pod 添加到负载均衡池中，开始向其发送请求。如果就绪探针失败，Kubernetes 会从负载均衡池中移除该 Pod，直至探针再次成功。
应用场景：在应用启动初期需要进行一些初始化操作，此时还不适宜接收外部请求，就绪探针可以确保在应用真正准备就绪后再对外提供服务。
- 启动探针（Startup Probe）：
作用：启动探针用于在容器启动初期代替存活探针和就绪探针，直到容器启动成功。它主要用于检查容器内的应用是否已经启动完成。当启动探针成功时，才会开始执行存活探针和就绪探针。
应用场景：对于那些启动时间较长或者启动过程复杂的应用，启动探针可以避免在应用未完全启动时过早触发存活探针而导致的不必要的重启。
