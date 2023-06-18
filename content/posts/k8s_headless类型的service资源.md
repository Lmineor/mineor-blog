---
title: "k8s_headless类型的service资源"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "k8s",
]
categories : [              # 文章所属标签
    "技术",
]
---

Service对象隐藏了各Pod资源，并负责将客户端的请求流量调度至该组pod对象之上。不过，偶尔也会存在这样一类需求：客户端需要直接访问Service资源后端的所有pod资源，这时就应该向客户端暴露每个pod资源的IP地址，而不再是中间层Service对象的ClusterIP，这种类型的Service资源便称为Headless Service。

　**　Headless Service对象没有ClusterIP**，于是kube-proxy便无须处理此类请求，也就更没有了负载均衡或代理它的需要。在前端应用拥有自有的其他服务发现机制时，Headless Service即可省去定义ClusterIP的需求。至于如何为此类Service资源配置IP地址，则取决于它的标签选择器的定义。

　　具有标签选择器：端点控制器（Endpoints Controller）会在API中为其创建Endpoints记录，并将ClusterDNS服务中的A记录直接解析到此Service后端的各Pod对象的ip地址上。

　　没有标签选择器：端点控制器不会在API中为其创建Endpoints记录，ClusterDNS的配置分为两种情形：对ExternalName类型的服务创建CNAME记录，对其他三种类型来说，为那些与当前Service共享名称的所有Endpoints对象创建一条记录。

# 1. 创建Headless Service资源

配置Service资源配置清单时,只需要将`ClusterIP`字段的值设置为`None`,即可将其定义为Headless类型.
如下示例:

```yaml
[root@k8s-master1 service]# cat headless-svc.yaml
kind: Service
apiVersion: v1
metadata:
  name: my-nginx-headless-svc
spec:
  clusterIP: None
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
    name: httpport
  selector:
    run: my-nginx
 
[root@k8s-master1 service]# kubectl apply -f headless-svc.yaml
service/my-nginx-headless-svc created
[root@k8s-master1 service]# kubectl get svc -o wide
NAME                    TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE   SELECTOR
kubernetes              ClusterIP   10.96.0.1    <none>        443/TCP   44d   <none>
my-nginx-headless-svc   ClusterIP   None         <none>        80/TCP    9s    run=my-nginx
```

使用资源创建命令完成资源创建后，查看获取的Service资源相关信息便可看到，它没有ClusterIP，不过如果标签选择器能够匹配到相关的pod资源，它便拥有EndPoints记录，这些EndPoints对象会作为DNS资源记录名称my-nginx-headless-svc查询时的A记录解析结果。

```bash
[root@k8s-master1 service]# kubectl describe svc my-nginx-headless-svc
Name:              my-nginx-headless-svc
Namespace:         default
Labels:            <none>
Annotations:       <none>
Selector:          run=my-nginx
Type:              ClusterIP
IP Families:       <none>
IP:                None
IPs:               None
Port:              httpport  80/TCP
TargetPort:        80/TCP
Endpoints:         10.244.36.116:80,10.244.36.118:80
Session Affinity:  None
Events:            <none>
[root@k8s-master1 service]# kubectl get pods -o wide --show-labels
NAME                        READY   STATUS    RESTARTS   AGE     IP              NODE        NOMINATED NODE   READINESS GATES   LABELS
my-nginx-6f6bcdf657-98h45   1/1     Running   0          4m22s   10.244.36.116   k8s-node1   <none>           <none>            pod-template-hash=6f6bcdf657,run=my-nginx,version=v1
my-nginx-6f6bcdf657-r866w   1/1     Running   0          4m22s   10.244.36.118   k8s-node1   <none>           <none>            pod-template-hash=6f6bcdf657,run=my-nginx,version=v1
```

# 2.Pod资源发现

根据Headless Service的工作特性可知，它记录于ClusterDNS的A记录的相关解析结果是后端pod资源的ip地址，这就意味着客户端通过此Service资源的名称发现的是各pod资源。下面依然创建一个专用的测试pod对象进行测试：

```bash
[root@k8s-master1 service]# cat dig-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: dig
  namespace: default
spec:
  containers:
  - name: dig
    image:  test/dig:latest
    imagePullPolicy: IfNotPresent
    command:
      - sleep
      - "3600"
  restartPolicy: Always
You have new mail in /var/spool/mail/root
[root@k8s-master1 service]# kubectl apply -f dig-pod.yaml
pod/dig created
[root@k8s-master1 service]# kubectl exec -it dig -- nslookup my-nginx-headless-svc
Server:         10.96.0.10
Address:        10.96.0.10#53
 
Name:   my-nginx-headless-svc.default.svc.cluster.local
Address: 10.244.36.116
Name:   my-nginx-headless-svc.default.svc.cluster.local
Address: 10.244.36.118
```

查看与Headless Service标签选择器关联的pod资源信息：

```bash
[root@k8s-master1 service]# kubectl get pods -l run=my-nginx -o wide
NAME                        READY   STATUS    RESTARTS   AGE   IP              NODE        NOMINATED NODE   READINESS GATES
my-nginx-6f6bcdf657-98h45   1/1     Running   0          20m   10.244.36.116   k8s-node1   <none>           <none>
my-nginx-6f6bcdf657-r866w   1/1     Running   0          20m   10.244.36.118   k8s-node1   <none>           <none>
```

其解析结果正是Headless Service通过标签选择器关联到到的所有pod资源的ip地址。于是，客户端向此Service对象发起请求将直接接入到pod资源中的应用之上，而不再由Service资源进行代理转发，它每次接入的pod资源则是由DNS服务器接收到查询请求时以轮询的方式返回的IP地址。