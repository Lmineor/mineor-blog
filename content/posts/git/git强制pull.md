---
title: "git强制pull"
date: 2023-11-12
draft: false
tags : [                    # 文章所属标签
    "Git",
]
---


> 版权声明：本文为CSDN博主「我想要身体健康」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上> 原文出处链接及本声明。
> 原文链接：https://blog.csdn.net/m0_57236802/article/details/131249491


如果你想强制 git pull 来覆盖本地的更改，你需要注意这个过程会删除所有你在本地做的更改，并将你的本地分支同步到远程分支。如果你想要这样做，可以使用以下的命令：

```bash
git fetch --all
git reset --hard origin/<branch_name>
```

在上面的 `<branch_name>` 中填写你想要同步的远程分支的名字。比如，如果你想要同步的分支是 master，你可以运行：

```bash
git fetch --all
git reset --hard origin/master
```

第一条命令 `git fetch --all` 会从远程仓库获取所有分支的最新更改，但是并不会修改你的本地仓库。

第二条命令 `git reset --hard origin/<branch_name>` 会将你的本地分支重置到远程分支的状态，这会删除所有的本地更改。

再次提醒，这个过程会丢失所有未提交的本地更改，所以在使用这个命令之前一定要确认你是否真的需要这样做。
