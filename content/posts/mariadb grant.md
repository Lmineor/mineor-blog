---
title: "Mariadb Grant命令"
date: 2023-11-26
draft: false
tags : [                    # 文章所属标签
    "Maridb"
]
categories : [              # 文章所属标签
    "技术",
]
---

```bash
$ mysql -u root -p
Create the keystone database:

MariaDB [(none)]> CREATE DATABASE keystone;
Grant proper access to the keystone database:

MariaDB [(none)]> GRANT ALL PRIVILEGES ON keystone.* TO 'keystone'@'localhost' IDENTIFIED BY 'KEYSTONE_DBPASS';
MariaDB [(none)]> GRANT ALL PRIVILEGES ON keystone.* TO 'keystone'@'%' IDENTIFIED BY 'KEYSTONE_DBPASS';
```