---
title: "收藏??就叫收藏吧"
date: 2022-11-08
draft: false
tags : [                    # 文章所属标签
    "collection",
]
categories : [              # 文章所属标签
    "收藏",
]
---

# Nginx 配置生成器

NginxWebUI是一款方便实用的nginx 网页配置工具，可以使用 WebUI 配置 Nginx 的各项功能，包括端口转发，反向代理，ssl 证书配置，负载均衡等，最终生成「nginx.conf」配置文件并覆盖目标配置文件，完成 nginx 的功能配置。

项目地址：https://gitee.com/cym1102/nginxWebUI
官方网站：https://nginxwebui.gitee.io

NginxWebUI功能说明

-该项目是基于springBoot的web系统，数据库使用sqlite，因此服务器上不需要安装任何数据库。
- 本项目可管理多个nginx服务器集群, 随时一键切换到对应服务器上进行nginx配置, 也可以一键将某台服务器配置同步到其他服务器, 方便集群管理。
- nginx本身功能复杂, 本项目并不能涵盖nginx所有功能, 只能配置常用功能, 更高级的功能配置仍然需要在最终生成的nginx.conf中进行手动编写。
- 部署此项目后, 配置nginx再也不用上网各种搜索, 再也不用手动申请和配置ssl证书, 只需要在本项目中进行增删改查就可方便的配置nginx。
