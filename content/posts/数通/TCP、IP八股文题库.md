【网络/后端“八股文”题库 · 共 60 题】  
分六大板块，每题给出答题要点（≈ 20 秒可答出）。面试官可随机抽 10-15 题，或让候选人自选 3-4 题深入。

---

### 1. Linux 基础（10 题）
1. [**进程 vs 线程** 的本质区别？内核分别用什么结构体表示？](content/posts/OS/进程vs线程的本质区别)
2. **fork、vfork、clone** 三者在实现和用途上的差异？  
3. **进程地址空间布局**（低地址→高地址）？哪些区域会随 fork 写时复制？  
4. **页表级数** 在 x86_64 上默认几级？每级页表占用多少字节？  
5. **epoll 水平触发与边缘触发** 区别？ET 下为何必须一次读完？  
6. **select/poll/epoll** 时间复杂度各是多少？epoll 的回调唤醒机制？  
7. [Linux OOM killer打分机制？如何保护关键进程？](content/posts/linux/OOM_killer打分机制)
8. **系统调用流程**（用户态→内核态）触发哪两条 CPU 指令？  
9. **CPU 亲和性** 设置命令及 C 接口？  
10. **/proc 与 /sys** 区别？列举 3 个调优常用节点。

---

### 2. 语言/工具链（10 题）
11. **C++ RAII** 原理，举一个必须自己实现析构的类。  
12. **shared_ptr 循环引用** 如何检测/解决？  
13. **Python GIL** 是什么？为何多线程 CPU 密集任务跑不满？  
14. **Python 协程** await 关键字背后的状态机？  
15. **C++11 move 语义** 解决了什么问题？完美转发模板怎么写？  
16. **gdb 调试 core** 三步法？如何用 gdb 打印 STL 容器？  
17. **strace -f -e trace=network** 会输出哪些系统调用？  
18. **perf top 火焰图** 红色、蓝色、绿色分别含义？  
19. **tcpdump 抓取 SYN Flood** 的典型过滤表达式？  
20. **objdump/readelf** 常用选项？如何查看符号表？

---

### 3. TCP/IP 协议（10 题）
21. [**三次握手** 各发送哪些标志位？为何初始序列号 ISN 要随机？](content/posts/数通/三次握手、四次挥手、同时打开、关闭、半关闭)
22. [**四次挥手** TIME_WAIT 2MSL 的作用？如何避免？](content/posts/数通/三次握手、四次挥手、同时打开、关闭、半关闭)
23. [**TCP 拥塞控制** 四个阶段？CUBIC 与 BBR 核心思想？](content/posts/数通/TCP协议)
24. [**TCP Keepalive** 三个系统参数？与 HTTP Keep-Alive 区别？](content/posts/数通/http长连接和tcp长连接的区别与联系)
25. **滑动窗口** 与 **拥塞窗口** 关系？  
26. **TCP_NODELAY** 与 **TCP_CORK** 使用场景？  
27. **Nagle 算法** 触发条件？如何关闭？  
28. **MSS、MTU、Window Scale** 计算示例？  
29. **RST 包** 出现的 4 种场景？  
30. **半连接队列（SYN Queue）** 与 **全连接队列（Accept Queue）** 溢出会发生什么？

---

### 4. 应用层/高并发（10 题）
31. **HTTP/1.1 与 HTTP/2** 最大区别？HPACK 压缩原理？  
32. **HTTP 状态码 301、302、307、308** 语义差异？  
33. **Cookie 与 Session** 实现机制？如何防止会话固定攻击？  
34. **WebSocket 握手** 头字段？如何兼容 80/443？  
35. **Nginx 11 个处理阶段** 顺序？access 阶段为何放 rewrite 后？  
36. **Nginx location 匹配优先级**？正则 location 如何终止搜索？  
37. **LVS 三种模式**（NAT/DR/TUN）原理、性能、场景？  
38. **HAProxy stick-table** 解决什么问题？  
39. **惊群效应** 在 accept/epoll 的表现？Nginx 如何解决？  
40. **零拷贝 sendfile** 流程？DMA gather-copy 条件？

---

### 5. 网络虚拟化/云原生（10 题）
41. **VXLAN 报文格式**？VNI 多少位？UDP 端口？  
42. **GRE 与 IPIP 隧道** 区别？  
43. **Linux Bridge、OVS、TC Flower** 转发路径差异？  
44. **Namespace + Veth** 创建一对虚拟网卡的命令？  
45. **iptables 四表五链** 顺序？DNAT 发生在哪条链？  
46. **conntrack 五元组** 超时时间可调吗？  
47. **IPVS 与 iptables** 性能差距？  
48. **Calico BGP** 模式与 VXLAN 模式优劣？  
49. **OpenStack Neutron** 三大核心插件？  
50. **eBPF 程序类型** 中，可用于 socket filter 的是哪一类？

---

### 6. 综合/场景题（10 题）
51. **线上 CPU 飙高 100%**，如何 3 步定位？  
52. **内存泄漏** 排查三板斧？  
53. **“大量 CLOSE_WAIT”** 如何脚本快速统计？  
54. **“SYN Flood 攻击”** 临时缓解 2 条命令？  
55. **“TIME_WAIT 过多”** 调哪两个内核参数？  
56. **“Docker 容器无法解析域名”** 排查顺序？  
57. **“curl: (56) Recv failure”** 可能原因 3 种？  
58. **“Nginx 502 Bad Gateway”** 排查思路？  
59. **“LVS DR 模式后端无法收到包”** 最常见错误？  
60. **“K8s Pod 跨节点不通”** 先查 CNI 还是路由？

---

### 使用建议
- **快速面试**：随机点 5-8 题，每题 30 秒答要点。  
- **深度面试**：选 2-3 题让候选人展开，追问实现细节或源码。  
- **反向提问**：让候选人挑 2 题他觉得不熟的，现场推导。