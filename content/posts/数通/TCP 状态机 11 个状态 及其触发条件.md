---
title: TCP 状态机 11 个状态 及其触发条件
date: 星期四, 七月 17日 2025, 2:16:55 下午
draft: true
tags:
  - 数通
---
TCP 状态机 **11 个状态** 及其触发条件（面试可直接背诵）：

| 状态          | 含义                   | 典型触发          |
| ----------- | -------------------- | ------------- |
| CLOSED      | 初始/最终                | 无任何连接         |
| LISTEN      | 服务器等待连接              | 调用 listen()   |
| SYN_SENT    | 已发 SYN               | 客户端 connect() |
| SYN_RCVD    | 已收 SYN，待 ACK         | 服务器收到 SYN     |
| ESTABLISHED | 连接建立                 | 三次握手完成        |
| FIN_WAIT_1  | 主动发 FIN              | 主动 close()    |
| FIN_WAIT_2  | 已收 ACK，等对端 FIN       | 收到对端 ACK      |
| CLOSE_WAIT  | 被动收 FIN，待应用 close()  | 服务器收到 FIN     |
| LAST_ACK    | 被动发 FIN，等 ACK        | 服务器调用 close() |
| CLOSING     | 同时关闭，收到 FIN 但 ACK 未达 | 罕见            |
| TIME_WAIT   | 等 2MSL 后彻底关闭         | 主动方收到最后 ACK   |

背诵口诀  
> **CLOSED → LISTEN/SYN_SENT → SYN_RCVD → ESTABLISHED → FIN_WAIT_1/2 → TIME_WAIT/CLOSE_WAIT → LAST_ACK → CLOSED**