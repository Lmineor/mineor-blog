---
title: "《k8s权威指南学习》--ConfigMap"
date: 2023-03-26
draft: false
tags : [                    # 文章所属标签
    "k8s",
]
categories : [              # 文章所属标签
    "技术",
]
---


# ConfigMap


使用ConfigMap 的限制条件如下。

- ConfigMap 必须在Pod 之前创建。
- ConfigMap 受Namespace 限制，只有处于相同Namespaces 中的Pod 可以引用它。
- ConfigMap 中的配额管理还未能实现。
- kubelet 只支持可以被API Server 管理的Pod 使用ConfigMap 。kubelet 在本Node 上通过 `--manifest-url`或`--config` 自动创建的静态P od 将无法引用Conf1gMap

- 在Pod 对ConfigMap 进行挂载（ volumeMount ）操作时，容器内部只能挂载为“目录”，无法挂载为“文件”。在挂载到容器内部后，目录中将包含Co nfigMap 定义的每个item,如果该目录下原来还有其他文件，则容器内的该目录将会被挂载的ConfigMap 覆盖。如果应用程序需要保留原来的其他文件，则需要进行额外的处理。可以将ConfigMap挂载到容器内部的临时目录，再通过启动脚本将配置文件复制或者链接到（ cp 或link命令）应用所用的实际配置目录下。