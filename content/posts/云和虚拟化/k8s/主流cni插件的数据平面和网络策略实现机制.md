---
title: 主流cni插件的数据平面和网络策略实现机制
date: 2025-04-05
draft: true
tags:
  - k8s
  - CNI
---
以下是主流CNI实现（Cilium/eBPF、Calico/BGP、Multus多网卡）的数据平面和网络策略实现机制的详细拆解：

### Cilium/eBPF

- **数据平面**：Cilium基于eBPF技术构建数据平面，eBPF程序直接在内核中运行，无需频繁的内核-用户态切换，大大降低了延迟。它通过在内核中加载eBPF程序，实现了对网络数据包的高效处理，包括流量过滤、负载均衡等功能。eBPF程序可以动态地插入到内核的网络栈中，对网络流量进行实时监控和控制，提高了网络的灵活性和可编程性。

- **网络策略实现机制**：Cilium利用eBPF的高性能键值存储（eBPF Maps）来存储网络策略和端点状态。当网络数据包到达时，eBPF过滤器会根据预定义的策略规则对数据包进行匹配和处理。Cilium的网络策略是基于身份的，它通过分析容器的标签和身份信息来确定网络访问权限。这种基于身份的策略模型使得网络策略的配置更加灵活和直观，能够更好地适应容器化应用的动态变化。

### Calico/BGP

- **数据平面**：Calico提供了两种数据平面模式，一种是传统的基于iptables和BGP路由的模式，另一种是eBPF模式。在传统模式下，Calico通过iptables实现网络策略的执行，并利用BGP协议进行跨节点的路由传播。然而，随着集群规模的增大，iptables规则的复杂性会增加，导致性能下降。在eBPF模式下，Calico绕过了传统的iptables和kube-proxy，直接在内核中通过eBPF程序处理数据包，实现了零拷贝数据平面，避免了内核-用户态上下文切换，从而显著提高了性能。
- 支持：
	-  内置数据加密Built-in data encryption  
	- IPAM功能Advanced IPAM management  
	- overlay和非overlay选项  
	- 数据平面可选：iptables, eBPF, Windows HNS, 或 VPP
    
- **网络策略实现机制**：Calico的网络策略是基于Kubernetes的NetworkPolicy资源定义的。在传统模式下，Calico通过iptables规则来实现网络策略的执行，每个网络策略都会转换为一系列的iptables规则。而在eBPF模式下，Calico利用eBPF程序实现了更高效的网络策略执行。Calico的网络策略可以精确地控制容器之间的网络通信，支持细粒度的访问控制，包括基于IP地址、端口、协议等条件的策略规则。
    

### Multus多网卡

- **数据平面**：Multus是一个多网卡CNI插件，它允许容器同时使用多个网络接口。Multus通过链式调用其他CNI插件来实现多网卡功能。它首先调用主CNI插件（如Flannel、Calico等）为容器分配主网络接口，然后再调用其他CNI插件为容器添加辅助网络接口。Multus的数据平面是基于这些底层CNI插件的数据平面实现的，它本身并不直接处理网络数据包，而是通过组合和协调多个CNI插件来提供多网卡功能。
    
- **网络策略实现机制**：Multus本身不直接实现网络策略，而是依赖于底层CNI插件的网络策略功能。由于Multus支持多种CNI插件的组合使用，因此其网络策略的实现方式会因所使用的底层插件而异。例如，如果Multus与Calico结合使用，那么网络策略将由Calico来实现；如果与Cilium结合使用，则由Cilium来实现网络策略。这种灵活性使得Multus能够适应不同的网络架构和策略需求。