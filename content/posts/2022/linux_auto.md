---
title: "linux自动化交互工具：expect示例"
date: 2022-11-28
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

# 安装

```bash
yum -y install expect
```

# 使用

```bash
#!/usr/bin/expect
# 本脚本需要安装 
# yum -y install expect
# yum install rsync -y

# 该脚本的作用为将指定名称的文件同步到其他节点
# params
# fname: 要同步的文件名
# hostip：目的节点
# user：目的节点用户名
# passwd：目的节点密码

# 设置变量
set fname [lindex $argv 0]
set hostip  [lindex $argv 1]
set user [lindex $argv 2]
set passwd [lindex $argv 3]
puts "fname=$fname"
puts "hostip=$hostip"
puts "user=$user"
puts "passwd=$passwd"

spawn rsync -rvl $fname $user@$hostip:/root/
expect {
 "*yes/no" {send "yes\r";exp_continue}
 "*password*:" {send "$passwd\r"}
}

expect eof
```

