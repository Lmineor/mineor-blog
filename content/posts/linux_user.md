---
title: "建立一个名为agetest的账号，该账号第一次登陆后使用默认密码，但必须更改密码后使用新密码才能够登陆系统使用bash环境"
date: 2022-11-03
draft: false
tags : [                    # 文章所属标签
    "linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

> 建立一个名为agetest的账号，该账号第一次登陆后使用默认密码，但必须更改密码后使用新密码才能够登陆系统使用bash环境。

增加用户agetest，且设置初始密码和用户名保持一致

```bash
[root@centos dashboard]# useradd agetest
[root@centos dashboard]# echo "agetest" | passwd --stdin agetest
Changing password for user agetest.
passwd: all authentication tokens updated successfully.
[root@centos dashboard]# chage -d 0 agetest
[root@centos dashboard]# chage -l agetest | head -n 3
Last password change					: password must be changed
Password expires					: password must be changed
Password inactive					: password must be changed
[root@centos dashboard]# logout

```

使用agetest登陆linux，便会强制让用户更改密码

```bash
[root@study ~]#ssh agetest@ip
agetest@ip's password:
You are required to change your password immediately (administrator enforced)

WARNING: Your password has expired.
You must change your password now and login again!
Changing password for user agetest.
Current password:
New password:
Retype new password:
passwd: all authentication tokens updated successfully.
Connection to ip closed.
```

over