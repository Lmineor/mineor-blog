---
title: "sed替换命令收集"
date: 2023-01-30
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

普通操作可以使用冒号(:)井号(#)正斜杠(/)来作为分隔符

```bash
sed -i 's#abc#def#g'  a.file  #将文件a.file中的abc替换成def
sed -i 's/^abc.*/abc=def/' a.file # 将a.file中以abc开头的一行替换成abc=def
sed -i '/ABC/,$d' a.file # 将a.file中从ABC开始（包括ABC）以后的所有行删除
sed -i '$a aabbccdd' a.file # 给a.file追加aabbccdd
cat geng.file | sed  's/abc/def/g'   ## 打印文件geng，并将其中的abc替换成def
```

参考：https://blog.csdn.net/genghongsheng/article/details/120432010
