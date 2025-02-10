---
title: "Mariadb Grant命令"
date: 2023-11-26
draft: false
tags : [                    # 文章所属标签
    "数据库"
]
categories : [              # 文章所属标签
    "技术",
]
---

```bash
$ mysql -u root -p
# Create the keystone database:

MariaDB [(none)]> CREATE DATABASE keystone;
# Grant proper access to the keystone database:

GRANT ALL PRIVILEGES ON yong.* TO 'yong'@'localhost' IDENTIFIED BY '123456';
GRANT ALL PRIVILEGES ON yong.* TO 'yong'@'%' IDENTIFIED BY '123456';
```