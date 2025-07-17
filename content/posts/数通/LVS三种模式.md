---
title: LVS三种工作模式详解（NAT/DR/TUN）
date: 2025-07-17
draft: true
tags:
---
# LVS三种工作模式详解（NAT/DR/TUN）

## 一、NAT模式（Network Address Translation）

### 原理
- **请求流程**：
  1. 客户端访问VIP（Virtual IP）
  2. Load Balancer修改目标IP为RIP（Real Server IP）并转发
  3. Real Server返回数据给Load Balancer
  4. Load Balancer将源IP改回VIP后发给客户端

- **关键特征**：
  - 双向流量都经过LB
  - 需要修改IP包头（源/目标地址）
  - Real Server需配置网关指向LB

### 性能
- **吞吐量**：受限于LB的NAT处理能力（通常5-8万并发）
- **延迟**：较高（两次NAT转换）
- **瓶颈**：LB成为带宽和性能瓶颈

### 适用场景
- 所有Real Server位于同一局域网
- 需要端口映射（如外网80转内网8080）
- 对性能要求不高的内部系统
- Real Server操作系统无限制（任何OS都支持）

## 二、DR模式（Direct Routing）

### 原理
- **请求流程**：
  1. 客户端访问VIP
  2. LB修改目标MAC为Real Server MAC（不修改IP）
  3. Real Server直接响应客户端（源IP仍为VIP）
  4. 响应流量不经过LB

- **关键特征**：
  - 仅请求经过LB，响应直连客户端
  - Real Server需配置VIP在lo接口（arp_ignore/arp_announce）
  - 要求LB与RS在同一个二层网络

### 性能
- **吞吐量**：极高（可支持百万级并发）
- **延迟**：最低（响应不经过LB）
- **瓶颈**：LB的调度能力和网络带宽

### 适用场景
- 高并发Web服务（电商、门户网站）
- 对响应速度要求苛刻的场景
- Real Server与LB同机房且二层互通
- Real Server需支持ARP抑制（Linux/Unix）

## 三、TUN模式（IP Tunneling）

### 原理
- **请求流程**：
  1. 客户端访问VIP
  2. LB封装原始IP包为新IP包（目标IP为RIP）
  3. Real Server解封装后处理请求
  4. Real Server直接响应客户端

- **关键特征**：
  - 通过IP隧道跨网络传输
  - Real Server需支持隧道协议
  - 可跨越不同机房部署

### 性能
- **吞吐量**：中等（受隧道封装开销影响）
- **延迟**：较高（封装/解封装开销）
- **瓶颈**：隧道协议处理能力

### 适用场景
- Real Server分布在多个数据中心
- 需要跨公网部署负载均衡
- DR模式无法满足时的替代方案
- Real Server支持IP隧道（需内核支持）

## 三模式对比表

| 特性        | NAT模式               | DR模式                 | TUN模式                |
|------------|-----------------------|------------------------|------------------------|
| **LB压力** | 高（处理进出流量）      | 低（仅处理入站）        | 中（需封装隧道）        |
| **性能**    | 低（~8万并发）         | 高（百万级并发）        | 中（~20万并发）         |
| **网络要求**| 同局域网               | 同二层网络              | 可跨三层网络            |
| **Real Server配置** | 需改网关       | 需配置ARP参数           | 需支持IP隧道            |
| **IP头修改**| 修改源/目标IP          | 仅改MAC地址             | 添加新IP头封装          |
| **典型应用**| 企业内部系统           | 高并发Web服务           | 异地多活架构            |

## 选型建议

1. **首选DR模式**：当满足同机房、二层可达条件时（90%生产环境选择）
2. **次选NAT模式**：需要端口转换或Real Server无法配置ARP参数时
3. **特殊选TUN模式**：必须跨机房部署且无法使用DR时

> **补充**：现代云环境常结合DSR（Direct Server Return）技术，类似DR模式但支持更灵活的部署方式，如AWS的NLB就采用类似原理。