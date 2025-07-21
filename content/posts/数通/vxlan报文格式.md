---
title: vxlan报文格式
date: 2025-07-21
draft: true
tags:
  - 数通
---
VXLAN 报文格式要点（面试速答）

1. VXLAN Header  
   • Flags：8 bit（I 位必须为 1，表示 VNI 有效）  
   • **VNI：24 bit** —— 可支持 2²⁴ ≈ 1677 万个虚拟网络   
   • Reserved：24 bit + 8 bit（置 0）

2. UDP 封装  
   • **目的端口（Dst Port）：4789**（IANA 标准）  
   • 源端口（Src Port）：动态随机或由原始帧哈希生成，范围 49152-65535（RFC 7348 建议）

3. 整体封装顺序  
   Outer Ethernet → Outer IP → UDP → VXLAN Header → Inner Ethernet → Original L2 Frame