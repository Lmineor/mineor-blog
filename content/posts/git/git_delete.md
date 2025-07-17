---
title: git_delete
date: 星期四, 七月 17日 2025, 2:15:13 下午
draft: false
tags:
  - Git
---

如果不配置.gitignore的文件，带push代码的时候就会把一写不必要的文件push到远程仓库，如.idea文件。如果不小心出现此文件在远程仓库可以通过一下步骤delete此文件：

1.配置.gitignore文件(新建/编辑）

```bash
vim .gitignore-->i--->添加.idea--->esc--->:wq

# 将.gitignore文件上传到远程仓库

git pull

git add .gitignore

git commit -m 'edit .gitignore'

git push origin master
```

3.删除git的.idea文件

```bash
git rm --cached -r .idea
```

4.同步到远程仓库

```bash
git commit -m 'delete .idea'
git push origin master
```

完成之后就可以发现git仓库中的.idea文件已经被删除，而且之后push代码也不会再把.idea文件上传