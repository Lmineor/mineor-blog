---
title: "k8s downward API"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Docker",
    "Go", 
    "Python",
    "Linux",
    "k8s",
    "Cloud",
    "SDN"
]
categories : [              # 文章所属标签
    "技术",
    "美食",
    "生活",
    "阅读",
]
---

前面我们从pod的原理到生命周期介绍了pod的一些使用，作为kubernetes中最核心的对象，最基本的调度单元，我们可以发现pod中的属性还是非常繁多的，前面我们使用过一个volumes的属性，表示声明一个数据卷，我们可以通过命令`kubectl explain pod.spec.volumes`去查看该对象下面的属性非常多，前面我们只是简单的使用了hostpath和empryDir{}这两种模式，其中还有一种叫做downwardAPI这个模式和其他模式不一样的地方在于它不是为了存放容器的数据也不是用来进行容器和宿主机的数据交换的，而是让pod里的容器能够直接获取到这个pod对象本身的一些信息。

https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/downward-api/


downwardAPI提供了两种方式用于将pod的信息注入到容器内部：

- 环境变量： 用于单个变量，可以将pod信息和容器信息直接注入容器内部

- volume挂载：将pod信息生成为文件，直接挂载到容器内部中去

# 环境变量

```yaml
[root@master1 ~]# cat  env-pod.yaml 
apiVersion: v1
kind: Pod
metadata:
  name: env-pod
  namespace: kube-system
spec:
  containers:
    - name: env-pod
      image: busybox
      command: ["/bin/sh", "-c","env"]
      env:
      - name: POD_NAME
        valueFrom:
          fieldRef:
            fieldPath: metadata.name
      - name: POD_NAMESPACE
        valueFrom:
          fieldRef:
            fieldPath: metadata.namespace
      - name: POD_IP
        valueFrom:
          fieldRef:
            fieldPath: status.podIP
```

打印

```bash
[root@k8s-master ~]# kubectl logs env-pod -n kube-system | grep POD
POD_IP=205.205.0.2
POD_NAME=env-pod
POD_NAMESPACE=kube-system
```

可以看到pod的ip,name, namespace都是通过环境变量打印出来的.

# volume挂载

downward API除了提供环境变量外,还提供通过volume挂载的方式去获取pod的基本信息,

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: volume-pod
  namespace: kube-system
  labels:
    k8s-app: test-volume
    node-env: test
  annotations:
    own: wangmuniangniang
    bulid: test
spec:
  volumes:
  - name: podinfo
    downwardAPI:
        items:
        - path: labels
          fieldRef:
            fieldPath: metadata.annotations
        - path: anntations
          fieldRef:
             fieldPath: metadata.annotations
 
  containers:
  - name: volume-pod
    image: busybox
    args:
    - sleep
    - "3600"
    volumeMounts:
    - name: podinfo
      mountPath: /etc/podinfo
```


查看

```bash
[root@master1 ~]# kubectl   exec  -it  volume-pod   /bin/sh  -n  kube-system
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
/ # ls  /etc/podinfo/
anntations  labels
/ # ls
bin   dev   etc   home  proc  root  sys   tmp   usr   var
/ # cat  /etc/podinfo/labels 
bulid="test"
kubernetes.io/config.seen="2022-04-27T03:23:02.840574876-04:00"
kubernetes.io/config.source="api"
own="wangmuniangniang"/ # ^C
/ # 

```

# 注意

downwardAPI能够获取到的信息,一定是pod里的容器进程启动之前就能够确定下来的信息,而如果想要获取pod容器运行之后才会出现的信息,比如容器进程PID,则不能使用downwardAPI,而应该考虑在pod里定义一个sidecar容器来获取了.
