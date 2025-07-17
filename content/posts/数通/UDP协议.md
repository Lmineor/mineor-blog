---
title: UDP首部
date: 星期四, 七月 17日 2025, 2:17:09 下午
draft: true
tags:
  - 数通
---

UDP是一个简单的面向数据报的运输层协议
UDP封装如图所示
![UDP](static/images/network/udp1.png)

UDP不提供可靠性：他把应用程序传给IP层的数据发送出去，但是并不保证他们能到达目的地。

# UDP首部

![udp首部](static/images/network/udp2.png)
- 端口号表示发送进程和接收进程
- UDP长度表示：UDP首部和UIDP数据的字节的长度（8+数据字节）