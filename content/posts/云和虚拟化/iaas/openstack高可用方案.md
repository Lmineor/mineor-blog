---
title: openstack高可用方案
date: 2025-04-17
draft: false
tags : [                    # 文章所属标签
    "Iaas"
]
---

### 典型架构

```
                       +---------------------+
                       |     HAProxy/Keepalived (VIP)    |
                       +---------------------+
                                  |
          +-----------------------+-----------------------+
          |                       |                       |
+---------------------+ +---------------------+ +---------------------+
|  Controller Node 1  | |  Controller Node 2  | |  Controller Node 3  |
|  - API Services     | |  - API Services     | |  - API Services     |
|  - Galera (MariaDB) | |  - Galera (MariaDB) | |  - Galera (MariaDB) |
|  - RabbitMQ Cluster | |  - RabbitMQ Cluster | |  - RabbitMQ Cluster |
+---------------------+ +---------------------+ +---------------------+
          |                       |                       |
          +-----------------------+-----------------------+
                                  |
                       +---------------------+
                       |     Ceph Cluster    |
                       |  (MON/OSD/MDS)      |
                       +---------------------+
```


### **HAProxy 和 Keepalived 详解**

在 OpenStack 高可用（HA）架构中，**HAProxy** 和 **Keepalived** 是两个关键组件，分别负责 **负载均衡** 和 **VIP（虚拟 IP）管理**，确保服务的高可用性。

---

## **1. HAProxy（高可用代理）**

### **1.1 基本介绍**

HAProxy 是一个高性能的 **TCP/HTTP 负载均衡器**，用于分发客户端请求到多个后端服务节点（如 OpenStack API 服务），避免单点故障。

### **1.2 核心功能**

- **负载均衡**：支持轮询（Round Robin）、最少连接（Least Connections）、源 IP 哈希（Source IP Hash）等算法。
- **健康检查**：定期检测后端服务是否可用，自动剔除故障节点。
- **SSL/TLS 终止**：可卸载 HTTPS 加密，减轻后端服务压力。
- **会话保持（Sticky Sessions）**：确保同一客户端请求始终转发到同一后端节点（适用于有状态服务）。

### **1.3 在 OpenStack 中的应用**

HAProxy 通常用于负载均衡以下 OpenStack 服务：
- **API 服务**：Nova、Neutron、Keystone、Glance、Cinder 等。
- **数据库（Galera）**：分发 MySQL/MariaDB 查询请求。
- **RabbitMQ**：均衡 AMQP 消息队列访问。

### **1.4 HAProxy 高可用**

- **多 HAProxy 节点**：部署多个 HAProxy 实例，避免单点故障。
- **结合 Keepalived**：使用 VIP 确保客户端始终访问可用的 HAProxy。

---

## **2. Keepalived（虚拟 IP 管理）**

### **2.1 基本介绍**

Keepalived 是一个基于 **VRRP（Virtual Router Redundancy Protocol）** 的高可用解决方案，用于管理 **虚拟 IP（VIP）**，确保服务入口的连续性。

### **2.2 核心功能**

- **VIP 漂移**：当主节点故障时，VIP 自动切换到备用节点。
- **健康检查**：可监控 HAProxy 或其他服务，触发故障转移。
- **多节点冗余**：支持主备（Active-Backup）或多主（Active-Active）模式。

### **2.3 在 OpenStack 中的应用**

- **API 入口高可用**：客户端通过 VIP 访问 OpenStack API，无需关心后端节点状态。
- **数据库/消息队列访问**：Galera/RabbitMQ 的负载均衡入口。

### **2.4 Keepalived 高可用**

- **主备模式（Active-Backup）**：
  - 一个节点是 `MASTER`（持有 VIP），另一个是 `BACKUP`。
  - 当 `MASTER` 故障时，`BACKUP` 接管 VIP。
- **多主模式（Active-Active）**：
  - 多个节点可同时提供服务，结合 DNS 轮询或 Anycast 实现负载均衡。

---

## **3. HAProxy + Keepalived 协作流程**

1. **客户端访问 VIP（192.168.1.100）**，请求被转发到当前主节点的 HAProxy。
2. **HAProxy 负载均衡** 请求到后端 OpenStack 服务（如 Nova-API）。
3. **如果主节点故障**：
   - Keepalived 检测到 HAProxy 不可用。
   - VIP 漂移到备用节点，客户端仍可访问服务。
4. **后端服务故障**：
   - HAProxy 健康检查发现故障节点，自动剔除。


---

## **总结**
| 组件       | 作用                          | 关键配置                     |
|------------|-------------------------------|------------------------------|
| **HAProxy** | 负载均衡 API/DB/队列          | `balance`, `server`, `httpchk` |
| **Keepalived** | VIP 管理，故障转移       | `vrrp_instance`, `track_script` |

通过 **HAProxy + Keepalived**，可以构建可靠的 OpenStack 高可用入口，确保服务持续可用。

### **OpenStack Neutron API 请求流程（HAProxy + Keepalived 高可用架构）**
当用户或服务（如 Nova）通过 **HAProxy + Keepalived** 访问 Neutron 的 `network` 接口时，完整的请求流程如下：

---

# 访问neutron的network接口请求示例

## **1. 客户端发起请求**
假设：
- **VIP（虚拟 IP）**: `192.168.1.100`
- **Neutron API 端口**: `9696`
- **后端 Neutron Server 节点**: `controller1:9696`, `controller2:9696`, `controller3:9696`

客户端（如 `curl` 或 Nova-Compute）发起请求：
```bash
curl -X GET http://192.168.1.100:9696/v2.0/networks -H "X-Auth-Token: <TOKEN>"
```

---

## **2. Keepalived 处理 VIP 路由**
1. **VIP 绑定**：
   - Keepalived 的 `MASTER` 节点（如 `controller1`）当前持有 VIP `192.168.1.100`。
   - 客户端请求发送到该 VIP，由 `controller1` 的网卡接收。

2. **故障转移场景**：
   - 如果 `controller1` 宕机，Keepalived 检测到故障（通过 VRRP 协议）。
   - VIP 漂移到 `BACKUP` 节点（如 `controller2`），客户端请求自动切换到新主节点。

---

## **3. HAProxy 负载均衡**
请求到达 HAProxy 后：
1. **监听端口**：
   - HAProxy 监听 `192.168.1.100:9696`（VIP 端口）。
   - 配置示例（`/etc/haproxy/haproxy.cfg`）：
     ```ini
     frontend neutron-api
         bind 192.168.1.100:9696
         default_backend neutron-server

     backend neutron-server
         balance roundrobin
         option httpchk GET /v2.0 HTTP/1.0  # 健康检查
         server controller1 192.168.1.101:9696 check inter 10s
         server controller2 192.168.1.102:9696 check inter 10s
         server controller3 192.168.1.103:9696 check inter 10s
     ```

2. **负载均衡逻辑**：
   - HAProxy 根据策略（如 `roundrobin`）将请求转发到后端 Neutron Server（如 `controller2:9696`）。
   - 如果 `controller2` 不可用（健康检查失败），HAProxy 自动剔除该节点，将请求发给 `controller1` 或 `controller3`。

---

## **4. Neutron Server 处理请求**
1. **API 处理**：
   - 选中的 Neutron Server（如 `controller2`）接收请求，解析 HTTP 方法（如 `GET /v2.0/networks`）。
   - 验证 `X-Auth-Token`（通过 Keystone 校验权限）。

2. **数据库/消息队列交互**：
   - Neutron 查询 **MySQL/Galera 集群**（如 `192.168.1.101:3306`）获取网络数据。
   - 可能通过 **RabbitMQ 集群**（如 `controller1:5672`）通知 L2 Agent（如 Open vSwitch）。

3. **返回响应**：
   - Neutron 生成响应（如网络列表 JSON），通过 HAProxy 原路返回客户端。

---

## **5. 客户端接收响应**
```json
{
  "networks": [
    {"id": "net1", "name": "private"},
    {"id": "net2", "name": "public"}
  ]
}
```

---

## **关键组件协作流程图**

```text
+----------+       +-------------+       +-----------+       +-----------------+
|          |       |             |       |           |       |                 |
|  Client  | ----> |  VIP:       | ----> |  HAProxy  | ----> |  Neutron Server |
|          |       |  192.168.1.100       |           |       |  (controller1/2/3) |
+----------+       +-------------+       +-----------+       +-----------------+
                                      ↑
                                      | VRRP (Keepalived)
                                +-----+-----+
                                |           |
                                | Galera DB |
                                | RabbitMQ  |
                                +-----------+
```

---

## **故障场景处理**
| 故障点               | 处理机制                                                                 |
|----------------------|--------------------------------------------------------------------------|
| **HAProxy 节点宕机** | Keepalived 将 VIP 漂移到备用节点，客户端无感知。                         |
| **Neutron Server 宕机** | HAProxy 健康检查剔除故障节点，请求转发到其他存活节点。                   |
| **数据库/队列故障** | Galera/RabbitMQ 集群自动切换，Neutron 重试或报错（依赖应用层容错）。     |

---

## **总结**
1. **客户端 → VIP**：通过 Keepalived 确保入口高可用。
2. **VIP → HAProxy**：负载均衡到多个 Neutron Server。
3. **Neutron Server → DB/Queue**：依赖集群化存储和消息中间件。
4. **全程无单点故障**：VIP、负载均衡、服务、存储均冗余部署。

这种架构是 OpenStack 生产环境的典型高可用方案，适用于所有核心 API（Nova、Neutron、Keystone 等）。