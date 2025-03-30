---
title: k8s核心资源
date: 2023-03-26
draft: false
tags:
  - k8s
---
# Kubernetes集群核心组件全解析：从控制平面到应用生态

Kubernetes作为容器编排领域的核心平台，其架构由多个解耦的组件协同工作。本文将深度剖析集群的三大层级组件，并扩展关键生态工具，帮助读者构建全面的知识体系。

---

## 一、控制平面（Control Plane）组件

### 1. **kube-apiserver**

- **核心枢纽**：提供RESTful API接口，所有集群操作（包括CLI、Dashboard、其他组件）均通过此组件交互
- **安全网关**：集成认证（X.509证书/Bearer Token）、鉴权（RBAC）、准入控制（Admission Webhook）三阶段安全机制
- **状态代理**：作为唯一直接与etcd通信的组件，确保数据一致性

### 2. **etcd**

- **分布式数据库**：采用Raft协议实现高可用，存储集群配置、Secret、Service发现等关键数据
- **运维要点**：推荐奇数节点部署（3/5节点），定期备份快照防止数据丢失

### 3. **kube-scheduler**

- **智能调度器**：基于资源请求、亲和性策略、硬件约束等条件为Pod选择最优节点
- **扩展机制**：支持自定义调度器（Scheduler Extender）实现业务特异性调度

### 4. **kube-controller-manager**

- **集群管家**：运行多个控制循环，包括：
    - **Deployment Controller**：管理副本数滚动更新
    - **Node Controller**：监控节点健康状态
    - **Service Account Controller**：管理命名空间级账户

### 5. **cloud-controller-manager** (云厂商定制)

- **云平台桥梁**：实现负载均衡器配置（AWS ELB/Azure LB）、存储卷动态供给（EBS/Azure Disk）、节点自动扩缩容（如GKE的Node Pool）

---

## 二、工作节点（Worker Node）组件

### 1. **kubelet**

- **节点代理**：
    - 接收PodSpecs并确保容器健康运行
    - 定期向API Server上报节点状态（CPU/Memory/Disk压力）
    - 执行存活探针（Liveness Probe）和就绪探针（Readiness Probe）

### 2. **容器运行时（Container Runtime）**

- **OCI标准实现**：
    - containerd（CNCF毕业项目，K8s默认选择）
    - CRI-O（专为K8s设计的轻量级运行时）
    - Docker（旧版本通过dockershim兼容）

### 3. **kube-proxy**

- **服务网格**：
    - 维护iptables/IPVS规则实现Service的虚拟IP转发
    - 支持三种代理模式：userspace（已弃用）、iptables（默认）、IPVS（高性能）

### 4. **网络与存储插件**

- **CNI网络插件**：
    
    - Flannel（Overlay网络简易方案）
    - Calico（基于BGP的高性能网络策略）
    - Cilium（eBPF驱动的新一代网络方案）
- **CSI存储插件**：
    
    - 对接云存储（AWS EBS, Google Persistent Disk）
    - 本地存储方案（OpenEBS, Rook）

---

## 三、关键应用层组件

### 1. **服务发现与负载均衡**

- **CoreDNS**：集群内DNS解析，将Service名称映射为虚拟IP
- **Ingress Controller**：
    - Nginx Ingress：基于配置生成路由规则
    - Traefik：支持Let's Encrypt自动证书

### 2. **监控与日志**

- **Metrics Server**：收集节点/Pod资源指标，支撑HPA自动扩缩
- **Prometheus + Grafana**：时序数据存储与可视化
- **EFK Stack**：
    - Fluentd：日志收集
    - Elasticsearch：日志存储
    - Kibana：日志展示

### 3. **持续交付与GitOps**

- **Argo CD**：声明式GitOps工具，实现配置即代码（IaC）
- **Tekton**：云原生CI/CD流水线框架

---

## 四、扩展生态系统

### 1. **安全增强**

- **Cert Manager**：自动管理TLS证书（Let's Encrypt集成）
- **OPA Gatekeeper**：策略即代码（Pod安全策略/资源限制）

### 2. **服务网格**

- **Istio**：提供mTLS加密、流量镜像、金丝雀发布等高级功能
- **Linkerd**：轻量级Service Mesh方案

### 3. **集群运维工具**

- **kubeadm**：快速搭建符合最佳实践的集群
- **Lens IDE**：可视化集群管理客户端
- **Velero**：集群备份与迁移工具

---

## 组件协同工作流程

通过一个Pod创建请求的完整生命周期，展示组件协作：

1. 用户提交`kubectl create`请求至**kube-apiserver**
2. **etcd**持久化存储Pod元数据
3. **kube-scheduler**选择合适节点并更新etcd
4. 目标节点的**kubelet**通过**容器运行时**创建容器
5. **kube-proxy**配置Service iptables/IPVS规则
6. **CoreDNS**更新服务域名解析记录