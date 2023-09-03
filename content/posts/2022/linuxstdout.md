---
title: "输入输出定向"
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

# 标准输入重定向

```bash
1>： 以覆盖的方法将「正确的数据」输出到指定的文件或装置上；
1>>：以追加的方法将「正确的数据」输出到指定的文件或装置上；
2>： 以覆盖的方法将「错误的数据」输出到指定的文件或装置上；
2>>：以追加的方法将「错误的数据」输出到指定的文件或装置上；
```


要注意：「1>>」以及「2>>」中间无空格。只有「>」代表1。
## 将正确与错误的输出分别定向到不同的文件中

使用范例：
```bash
[root@centos ~]# find /home -name .bashrc > list_rigth 2> list_err
```

此时屏幕不会出现任何信息，因为刚刚执行的结果中，标准正确输出被重定向到了`list_right`文件中，而标准错误输出被重定向到了`list_err`文件中。

## 将正确的输出在屏幕上显示，错误输出不显示
利用黑洞 `/dev/null`，这个东西可以吃掉任何导向这个设备的数据。
```bash
[root@centos ~]# find /home -name .bashrc 2> /dev/null
/home/gc/.bashrc
/home/lex/.bashrc
/home/www/.bashrc
```
正确与错误的均写入同一个文件
```bash
[root@centos ~]# find /home -name .bashrc >list 2>&1
```

## &&与||
|指令下达|说明|
|:-|:-|
|`cmd1&&cmd2` |1.若cmd1执行完毕且正确执行（`$?=0`），则开始执行cmd2. 2. 若cmd1执行完毕且为错误(`$?!=0`)，则cmd2不执行。|
|`cmd1｜｜cmd2` |1.若cmd1执行完毕且正确执行，则cmd2不执行。2.若cmd1执行完毕且错误，则开始执行cmd2.|

