---
title: http长连接和tcp长连接的区别与联系
date: 星期四, 七月 17日 2025, 2:25:01 下午
draft: false
tags:
  - 数通
---
面试回答（背 30 秒版）

一、TCP Keepalive 的三个核心系统参数  

| 参数 | 含义 | 默认值（Linux） |
|---|---|---|
| `net.ipv4.tcp_keepalive_time` | 连接空闲多久后开始发送探测包 | 7200 秒 |
| `net.ipv4.tcp_keepalive_intvl` | 每次探测包之间的间隔 | 75 秒 |
| `net.ipv4.tcp_keepalive_probes` | 连续探测几次无响应后判定连接失效 | 9 次 |

查看/修改示例  
```bash
sysctl net.ipv4.tcp_keepalive_time      # 查看
sysctl -w net.ipv4.tcp_keepalive_time=600
```

代码级设置  
```c
setsockopt(fd, SOL_TCP, TCP_KEEPIDLE,  &600,  sizeof(int)); // 对应 time
setsockopt(fd, SOL_TCP, TCP_KEEPINTVL, &75,   sizeof(int)); // 对应 intvl
setsockopt(fd, SOL_TCP, TCP_KEEPCNT,   &9,    sizeof(int)); // 对应 probes
```

二、TCP Keepalive vs HTTP Keep-Alive  

| 维度 | TCP Keepalive | HTTP Keep-Alive |
|---|---|---|
| 所在层级 | 传输层（内核态） | 应用层（用户态） |
| 目的 | 检测“死”连接并自动回收 | 复用同一 TCP 连接发多个 HTTP 请求 |
| 触发条件 | 连接空闲超过 `tcp_keepalive_time` | 一次 HTTP 事务结束后不立刻 `close` |
| 控制方式 | 内核参数 + `SO_KEEPALIVE` | HTTP 头 `Connection: Keep-Alive` 或 `close` |
| 典型参数 | 如上三参数 | Nginx: `keepalive_timeout`；Gunicorn: `--keep-alive` |

一句话总结  
“TCP Keepalive 是内核的‘心跳保活’，HTTP Keep-Alive 是应用的‘连接复用’，两者名字像，但做的事完全不同。”