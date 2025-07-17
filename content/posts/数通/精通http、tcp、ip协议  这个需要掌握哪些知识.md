---
title: 精通http、tcp、ip协议  这个需要掌握哪些知识
date: 星期四, 七月 17日 2025, 2:15:53 下午
draft: true
tags:
  - 数通
---
“精通 HTTP、TCP/IP 协议” 不是背八股，而是能在面试/排障/架构/编码四个场景里 **立刻拿出可落地的知识与工具**。下面给出一张 **“精通检查表”**，分 5 个层级、60+ 个知识点，全部打勾即可宣称精通。

---

### ① 体系与分层（必须能徒手画图）
- OSI七层：物理->数据链路层->网络->传输->会话->表示->应用
- 五层模型：物理→链路→网络→传输→应用  
- TCP/IP 四层对应：链路→网际→传输→应用  
- 数据封装/解封装流程：HTTP 报文 → TCP 段 → IP 包 → 以太网帧  
- MTU、MSS（Maximum Segment Size）、分片、重组现场计算（给出 1500 字节 MTU 场景）

---

### ② IP & 路由（排障必备）
- IPv4 首部长度、TTL 含义、DF/MF 位  
- 子网划分：/28、/30、CIDR 速算主机数  
- ARP、NDP（IPv6）报文格式 & 攻击场景  
- 路由表匹配算法（最长前缀匹配）  
- traceroute、mtr 原理（TTL 递增、三次探测）  
- ip rule / ip route / iptables mangle 做策略路由实例  
- GRE、IPIP、VXLAN 头格式各 1 张图  
- VPN 场景：IPsec ESP/AH、OpenVPN tun/tap、WireGuard 握手一次往返

---

### ③ TCP 深度（背不下来就别说精通）
- 20 字节首部 + 选项：MSS、WS、TS、SACK 字段位置  
- 三次握手、四次挥手、同时打开/关闭、半关闭  
- 11 种状态机（CLOSED → SYN_SENT → … → TIME_WAIT）  
- 重传机制：超时 RTO、快速重传、SACK、FACK、ER  
- 拥塞控制：Reno/CUBIC/BBR 算法公式 & 调整参数  
- 滑动窗口 vs 拥塞窗口、Nagle、Delayed ACK、CORK、TCP_NODELAY  
- keepalive 三系统参数、半连接队列 vs 全连接队列溢出场景  
- SYN Flood、RST 攻击、窗口缩放攻击、Sockstress 防御  
- ss -i、tcpdump 过滤表达式现场写 3 条  
- 内核调优：net.ipv4.tcp_* 必背 10 个参数及典型值

---

### ④ HTTP 全栈（从报文到性能）
- 报文格式：起始行 + 首部 + 实体，CRLF 分隔  
- 方法语义：GET/POST/PUT/PATCH/DELETE/HEAD/OPTIONS/CONNECT  
- 状态码 1xx-5xx，重点背 301 vs 302 vs 307 vs 308  
- 首部：Cache-Control、ETag、If-None-Match、Range、Expect:100-continue  
- Cookie、Set-Cookie、SameSite、Secure、HttpOnly  
- HTTPS 握手：TLS1.2/1.3 四次/三次握手，ALPN 选 h2/h3  
- HTTP/1.1 vs 2 vs 3：二进制帧、HPACK/QPACK、多路复用、0-RTT  
- 性能：Keep-Alive、管线化、分块传输、压缩（gzip/br）  
- Nginx/OpenResty 11 个阶段、location 优先级、变量 $uri vs $request_uri  
- CDN 原理：DNS 调度、回源、缓存键、Range 回源、SNI 分片

---

### ⑤ 工具与实战（必须现场能敲）
- 抓包：tcpdump、Wireshark、BPF 过滤 5 条常用  
- 测性能：wrk/ab/hey、iperf3、netperf  
- 看内核：ss、netstat -s、nstat、/proc/net/*  
- 调延迟：tc qdisc（netem）、ip route add … rtt 50ms  
- 调试 HTTP：curl -v、httpie、openssl s_client -alpn h2  
- 高并发：epoll 边缘/水平触发、惊群、SO_REUSEPORT  
- 容器网络：veth、bridge、macvlan、iptables DNAT/SNAT  
- 故障剧本：  
  – “大量 CLOSE_WAIT” 3 步定位脚本  
  – “SYN Flood” 临时防护 2 条命令  
  – “HTTP 502” 排查流程图

---

### 面试 60 秒自我介绍模板
> “我能徒手画 TCP 三次握手、RTO 计算、HTTP/2 二进制帧；线上用 ss + tcpdump 10 分钟定位过大量 TIME_WAIT；Nginx/OpenResty 写过 Lua 插件做灰度；GRE/VXLAN 打通跨机房网络；内核调优把 RTT 从 100 ms 降到 30 ms。”