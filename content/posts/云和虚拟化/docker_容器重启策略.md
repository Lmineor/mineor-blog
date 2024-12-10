---
title: "docker容器重启策略"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Docker",
]
---

# 可选策略

在docker run通过 --restart 设置守护机制：

- no: 不自动重新启动容器（默认）
- on-failure: 容器发生error而退出(容器退出状态不为0)重启容器
- unless-stopped:  在容器已经stop掉或Docker stoped/restarted的时候才重启容器
- always: 如果容器停止，总是重新启动容器。如果手动kill容器，则无法自动重启。

# docker update追加命令

运行中的容器,当时没有指定restart可以通过update命令追加

举例： web为正在运行的容器

```bash
docker update --restart=always web
```