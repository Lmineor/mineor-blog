---
title: "shell中=和～的用法"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

shell 中 =~ 的用法

我们先看一个脚本，该脚本的功能是搜索当前目录下文件中的指定字符串

```bash
#!/bin/bash
apath=$1;acontent=$2;aexp=$3;

if [[ $aexp =~ all ]] ;then                                                                                                                                                        
    atype=''
else
    atype=".$aexp"
fi

find $apath  -name  "*"$atype -type f -print0 | xargs -0 grep --color -rn "$acontent"

```

`if [[ $aexp =~ all ]]`
其中 ~是对后面的正则表达式匹配的意思，如果匹配就输出1，不匹配就输出0

