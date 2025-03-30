---
title: k8s面试相关
date: 2023-06-18
draft: true
tags:
  - k8s
---
为了应对涉及 Kubernetes（K8s）编排框架及二次开发 的面试，你需要系统性地准备以下技术要点和潜在问题：

---

 一、核心技术要点
 1. Kubernetes 基础与架构
- 核心组件：  
  - 掌握 API Server、Controller Manager、Scheduler、kubelet、etcd 的作用及交互流程。  
  - 理解 声明式 API 和 控制器模式（如 Deployment 如何通过 ReplicaSet 管理 Pod）。  
- 资源对象：  
  - 熟悉 Pod、Deployment、StatefulSet、Service、Ingress、ConfigMap/Secret、PV/PVC、CRD 的使用场景。  
  - 问题示例：  
    > _“如何设计一个高可用的 StatefulSet 服务？”_  
    > _“Service 的 ClusterIP 和 NodePort 有什么区别？如何通过 Ingress 实现灰度发布？”_

 2. 编排与调度机制
- 调度器（Scheduler）：  
  - 了解默认调度流程（预选/优选策略）、如何通过 `nodeSelector`/`affinity`/`taint` 控制调度。  
  - 二次开发场景：  
    - 自定义调度插件（Scheduler Framework）或扩展调度器（如实现基于 GPU 资源的调度）。  
- 问题示例：  
  > _“如何实现跨可用区的 Pod 均匀分布？”_  
  > _“如果让你优化调度器，你会考虑哪些指标？”_

 3. 网络与存储
- 网络模型：  
  - CNI 原理（如 Calico、Flannel）、Service 的 iptables/IPVS 实现、Ingress Controller（如 Nginx）。  
- 存储：  
  - CSI 驱动开发流程、动态卷供给（StorageClass）。  
- 问题示例：  
  > _“如何排查 Pod 间网络不通的问题？”_  
  > _“如何设计一个支持快照功能的 CSI 插件？”_

 4. 监控与运维
- 可观测性：  
  - Prometheus + Grafana 监控体系、kube-state-metrics 和 cAdvisor 的作用。  
- 故障排查：  
  - 熟悉 `kubectl debug`、`kubectl logs/describe/exec`、`kubelet` 日志分析。  
- 问题示例：  
  > _“如何快速定位一个 CrashLoopBackOff 的 Pod 问题？”_  
  > _“如何设计一个自动化的集群健康检查工具？”_

 5. 安全与多租户
- RBAC：角色绑定、ServiceAccount 权限控制。  
- 安全策略：PodSecurityPolicy（或替代方案如 OPA Gatekeeper）、网络策略（NetworkPolicy）。  
- 问题示例：  
  > _“如何限制某个 Namespace 的 Pod 只能访问特定外部 IP？”_

 6. 二次开发能力
- Operator/CRD 开发：  
  - 使用 Kubebuilder 或 Operator SDK 开发自定义控制器（Controller）。  
  - 理解 Informer 机制、WorkQueue 的使用。  
- API 扩展：  
  - 熟悉 Aggregated API Server 或修改 K8s 源码（如添加新的资源类型）。  
- 问题示例：  
  > _“如何设计一个 Operator 来自动管理数据库集群？”_  
  > _“如果让你扩展 K8s API，你会如何设计？”_

 7. 生态工具链
- CI/CD 集成：ArgoCD、Tekton、Jenkins X。  
- GitOps：熟悉 FluxCD 或 ArgoCD 的部署流程。  

---

 二、高频面试问题分类
 1. 基础理论
- “简述 K8s 的架构和组件交互。”  
- “Deployment 和 StatefulSet 的区别是什么？”  

 2. 场景设计
- “如何设计一个支持滚动升级且零宕机的服务？”  
- “如何实现跨集群的配置同步？”  

 3. 故障排查
- “Pod 一直处于 Pending 状态，可能的原因有哪些？”  
- “如何分析 API Server 的性能瓶颈？”  

 4. 二次开发
- “你如何扩展 K8s 的功能？需要修改哪些模块？”  
- “写过自定义 Controller 吗？描述 Informer 的工作机制。”  

 5. 底层原理
- “etcd 在 K8s 中的作用是什么？如何备份？”  
- “kube-proxy 的 IPVS 模式是如何工作的？”  

---

 三、准备建议
1. 动手实践：  
   - 部署一个多节点集群（可用 kubeadm 或 kind），尝试修改源码并编译组件（如自定义 Scheduler）。  
   - 开发一个简单的 Operator（例如管理 Nginx 配置的 CRD）。  
2. 阅读源码：  
   - 重点阅读 。  
3. 模拟面试：  
   - 针对上述问题自问自答，或使用工具（如 ChatGPT）模拟技术对话。  

---

通过系统性地覆盖这些领域，你不仅能应对技术面试，还能展现对 K8s 生态的深度理解。如果需要更具体的某个方向（如网络或调度），可以进一步深入！