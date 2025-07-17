---
title: QUIC
date: 2025-07-17
draft: true
tags:
  - 数通
---
QUIC（Quick UDP Internet Connections）是 Google 2013 年提出、IETF 2022 年正式标准化的 **新一代传输层协议**，它把传统“TCP+TLS”两层栈压缩成 **“UDP+QUIC+TLS1.3”** 一层，专为 **更快、更稳、更安全** 的互联网而生。

一句话定位  
> **QUIC 就是“用 UDP 跑的 TCP+TLS 超融合升级版”，HTTP/3 默认跑在它上面。**

---

### 1 核心特征（面试速背 5 点）
| 特征 | 一句话解释 |
|---|---|
| **0/1-RTT 建链** | 首次 1-RTT，复用会话 0-RTT 即可发数据。 |
| **连接迁移** | 4G/WiFi 切换时，**IP 和端口变了** 也能保持连接，因为用 **64-bit Connection ID** 标识。 |
| **多路复用无队头阻塞** | 同一 UDP 通道里多 **Stream** 并发，丢包只影响单条 Stream。 |
| **内置 TLS 1.3** | 握手与加密一次完成，无法被中间设备“降级”。 |
| **应用层实现** | 拥塞控制算法可在用户态更新，**无需改内核**。 |

---

### 2 与 HTTP 的关系
- HTTP/1.1、HTTP/2 跑在 **TCP** 上。  
- HTTP/3 直接跑在 **QUIC** 上，因此又叫 **HTTP-over-QUIC**。

---

### 3 如何肉眼识别
用 Chrome DevTools → Network → **Protocol 列**  
- `h3` 或 `h3-29` → 就是 QUIC/HTTP3。

---

### 4 典型使用场景
- **网页加速**：Chrome 对支持 h3 站点自动用 QUIC。  
- **实时媒体**：RTMP、MoQ 跑在 QUIC 上，弱网不卡。  
- **CDN/边缘**：阿里云、腾讯云 CDN 已提供 **QUIC 开关**。