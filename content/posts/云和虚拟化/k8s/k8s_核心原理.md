---
title: "《k8s权威指南学习》--k8s核心原理"
date: 2023-03-26
draft: false
tags : [                    # 文章所属标签
    "k8s",
]
categories : [              # 文章所属标签
    "技术",
]
---


## k8s核心原理

### API Server

总体来看， Kubemetes API Server 的核心功能是提供了Kubemetes 各类资源对象（如Pod 、RC 、Service 等〉的增、删、改、查及Watch 等HTTP Rest 接口，成为集群内各个功能模块之间数据交互和通信的中心枢纽，是整个系统的数据总线和数据中心。
除此之外，它还有以下一些功能特性。

1. 是集群管理的API入口。
2. 是资源配额控制的入口。
3. 提供了完备的集群安全机制。