---
title: "docker启动命令逆向输出"
date: 2022-10-19
draft: false
tags : [                    # 文章所属标签
    "docker","容器与虚拟化",
]
categories : [              # 文章所属标签
    "技术",
]
---

主要依靠的一个开源软件`runlike`

See [runlike-github](https://www.notion.so/docker-12c712016d55490ca8716dc603401d11#3c6a4c361b34403782a86aee662fab7e)

```bash

yum clean all
yum makecache

yum install python-pip python-setuptools -y

pip install runlike

# 若此时正在运行的一个容器为`redis`
# 则
runlike redis > redis_start.sh
# 此时redis_start.sh即为docker启动redis容器的命令

```