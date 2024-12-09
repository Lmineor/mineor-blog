---
title: "k8s_给你的服务签发证书"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "k8s",
]
categories : [              # 文章所属标签
    "技术",
]
---


使用场景
当我们需要进行服务端认证，甚至双向认证时，我们需要生成密钥对和服务信息，并使用ca对公钥和服务信息进行批准签发，生成一个证书。

我们简单描述下单向认证和双向认证的场景流程

在单向认证场景中：

服务端会将自己的证书和公钥告知客户端
客户端向CA查询该证书的合法性，确认合法后会记录服务端公钥
客户端会与服务端明文通信确认加密方式，
客户端确认加密方式后，会生成随机码作为对称加密密钥，以服务端的公钥对对称加密密钥进行加密，告知服务端，
服务端以自己的私钥解密得到对称加密密钥，
之后， 客户端与服务端之间使用对称加密密钥进行加密通信。
在双向认证的场景中:

服务端会将自己的证书和公钥告知客户端
客户端向CA查询该证书的合法性，确认合法后会记录服务端公钥
客户端会将自己的证书和公钥发给服务端，
服务端发现客户端的证书也可以通过CA认证，则服务端会记录客户端的公钥
然后客户端会与服务端明文通信确认加密方式，
但服务端会用客户端的公钥将加密方式进行加密
客户端使用自己的私钥解密得到加密方式，会生成随机码作为对称加密密钥，以服务端的公钥对对称加密密钥进行加密，告知服务端，
服务端以自己的私钥解密得到对称加密密钥，
之后， 客户端与服务端之间使用对称加密密钥进行加密通信。
那么，我们如果要开放自己的https服务，或者给kubelet创建可用的客户端证书，就需要：

生成密钥对
生成使用方的信息
使用ca对使用方的公钥和其他信息进行审核，签发，生成一个使用方证书。
这里使用方可以是客户端（kubelet）或服务端（比如一个我们自己开发的webhook server）

手动签发证书
k8s集群部署时会自动生成一个CA（证书认证机构），当然这个CA是我们自动生成的，并不具有任何合法性。k8s还提供了一套api，用于对用户自主创建的证书进行认证签发。

准备
安装k8s集群
安装cfssl工具，从这里下载cfssl和cfssljson
创建你的证书
执行下面的命令，生成server.csr和server-key.pem。

```
cat <<EOF | cfssl genkey - | cfssljson -bare server
{
  "hosts": [
    "my-svc.my-namespace.svc.cluster.local",
    "my-pod.my-namespace.pod.cluster.local",
    "192.0.2.24",
    "10.0.34.2"
  ],
  "CN": "kubernetes",
  "key": {
    "algo": "ecdsa",
    "size": 256
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
EOF
```

这里你可以修改文件里的内容，主要是：

hosts。 服务地址，你可以填入service的域名，service的clusterIP，podIP等等
CN 。对于 SSL 证书，一般为网站域名；而对于代码签名证书则为申请单位名称；而对于客户端证书则为证书申请者的姓名.k8s会将CN的内容视为正式使用者的User Name
key。加密算法和长度。一般有ecdsa算法和rsa算法，rsa算法的size一般是2048或1024
names。证书申请者的信息，比如位置、组织等。k8s的RBAC会将证书中的O视为证书使用者所在的Group
这一步生成的server-key.pem是服务端的私钥，而server.csr则含有公钥、组织信息、个人信息(域名)。

创建一个CSR资源
执行如下脚本：

```
cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: my-svc.my-namespace
spec:
  request: $(cat server.csr | base64 | tr -d '\n')
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF
```

在k8s集群中创建一个csr资源。注意要将第一步中创建的server.csr内容进行base64编码，去掉换行后填入spec.request中。spec.usages中填入我们对证书的要求，包括数字签名、密钥加密、服务器验证。一般填这三个就够了。

之后我们通过kubectl describe csr my-svc.my-namespace 可以看到：
```
Name:                   my-svc.my-namespace
Labels:                 <none>
Annotations:            <none>
CreationTimestamp:      Tue, 21 Mar 2017 07:03:51 -0700
Requesting User:        yourname@example.com
Status:                 Pending
Subject:
        Common Name:    my-svc.my-namespace.svc.cluster.local
        Serial Number:
Subject Alternative Names:
        DNS Names:      my-svc.my-namespace.svc.cluster.local
        IP Addresses:   192.0.2.24
                        10.0.34.2
Events: <none>
```
认证csr
注意到，csr的status是pending，说明还没有被CA认证。在k8s集群中，如果是node上kubelet创建的CSR，kube-controller-manager会自动进行认证，而我们手动创建的证书，需要进行手动认证：
```
kubectl certificate approve
```
也可以拒绝：`kubectl certificate deny`

之后我们再检查csr,发现已经是approved了：

```
kubectl get csr
NAME                  AGE       REQUESTOR               CONDITION
my-svc.my-namespace   10m       yourname@example.com    Approved,Issued
```

我们可以通过

```
kubectl get csr my-svc.my-namespace -o jsonpath='{.status.certificate}' | base64 --decode > server.crt
```

命令，得到server的证书。之后你就可以使用server.crt和server-key.pem作为你的服务的https认证

流程总结
- 编写一个json文件，描述server的信息，包括域名（或IP），CN，加密方式
- 执行cfssl命令生成server的密钥，和认证请求文件server.csr
- 将server.csr内容编码，在k8s中创建一个server的CSR资源
- 手动对该CSR资源进行认证签发
- 将k8s生成的server.crt 即服务端证书拷贝下来。
- server.crt 和server-key.pem 即server的https服务配置