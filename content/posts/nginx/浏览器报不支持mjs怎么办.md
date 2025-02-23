---
title: "浏览器报不支持mjs怎么办"
date: 2025-02-23
draft: false
tags : [                    # 文章所属标签
    "nginx",
]

---

# 浏览器报不支持mjs怎么办

需要对nginx默认的支持文件进行修改

```
# `nginx` 路径根据需要调整
vim /etc/nginx/mime.types
# 宝塔
vim /www/server/nginx/conf/mime.types
```

修改文件

```
# 在下面这一行的 `js` 后面加上 `mjs`
application/javascript                 js;
# 如
application/javascript                 js mjs;
# 保存并退出
:wq
# 重启nginx
nginx -s reload
```