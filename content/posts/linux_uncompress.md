---
title: "解压缩"
date: 2022-10-23
draft: false
tags : [                    # 文章所属标签
    "linux",
]
categories : [              # 文章所属标签
    "技术",
]
---


最简单的使用tar的命令：

```bash
# 压缩
tar -jcv -f filename.tar.bz2 要被压缩的文件或目录名称

# 查询
tar -jtv -f filename.tar.bz2

# 解压缩
tar -jxv -f filename.tar.bz2 -C 欲解压缩的目录
```
