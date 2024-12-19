---
title: "docker format常用选项"
date: 2022-11-30
draft: false
tags : [                    # 文章所属标签
    "云与虚拟化", 
]

---


格式化选项(--format)
.ID 容器ID
.Image 镜像ID
.Command Quoted command
.CreatedAt 创建容器的时间点.
.RunningFor 从容器创建到现在过去的时间.
.Ports 暴露的端口.
.Status 容器状态.
.Size 容器占用硬盘大小.
.Names 容器名称.
.Labels 容器所有的标签.
.Label 指定label的值 例如'{{.Label “com.docker.swarm.cpu”}}’
.Mounts 挂载到这个容器的数据卷名称

用例：

删除已经死掉的容器，但 `docker ps -a`还有且状态为Exit...的容器

```bash
docker ps -a --format "{{.ID}}\t{{.Status}}" | grep Ex | awk '{print$1}' | xargs docker rm
```

删除docker images的tag是none的镜像

```bash
docker images --format "{{.Repository}}\t{{.ID}}" | grep none | awk '{print $2}' | xargs docker rmi
```


https://blog.csdn.net/aiwangtingyun/article/details/123380626