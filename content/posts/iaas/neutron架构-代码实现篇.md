---
title: "neutron架构-代码实现篇"
date: 2024-08-15
draft: true
tags : [                    # 文章所属标签
    "OpenStack",
]
categories : [              # 文章所属标签
    "技术",
]
---



# Neutron核心相关

# Neutron核心概念

## Network

- 一个L2 二层网络单元
- 租户可通过Neutron API创建自己的私有/外部网络
[代码模型定义](https://github.com/openstack/neutron/blob/master/neutron/db/models_v2.py#L313)

## Subnet

- 一个L3的IPv4/IPv6地址段
- 为VM提供私网或公网IP地址

[代码模型定义](https://github.com/openstack/neutron/blob/master/neutron/db/models_v2.py#L195)

## Router

- 三层路由器
- 为租户的VM提供路由功能，连接Internet

[代码模型定义](https://github.com/openstack/neutron/blob/master/neutron/db/models/l3.py#L47)


## Port

- 虚拟交换器上的端口
- 管理VM的网卡

# Neutron核心插件

- Open vSwitch

# Neutron核心服务

通常情况下，Neutron组件在Openstack架构中常以单独的Node形式提供网络服务，作为网络节点。提供了多种服务：

- neutron-server
    
    提供REST API服务，后端使用关系数据库。neutron-server是一个守护进程，用来提供外部调用的API和与其他组件交互的接口。
    
- Message Queue
    
    neutron-server使用Message Queue与其他Neutron agents进行交换消息
    
- L2 Agent
    
    负责连接端口（ports）设备，使他们处于共享的广播域（broadcast domain）。通常运行在Hypervisor上
    
- DHCP Agent
    
    用于配置虚拟主机的网络。DHCP代理，给租户网络提供动态主机配置服务，主要用途是为租户网络内的虚拟机动态地分配IP地址
    
- L3 Agent
    
    负责连接tenant（租户）网络到数据中心，或连接到internet。L3代理，提供三层网络功能和网络地址转换（NAT）功能，来让租户的虚拟机可以与外部网络通信。
    
- plugin-in agent
    
    插件代理，需要部署在每一个运行hypervisor的主机上，它提供本地的vSwitch配置，一般用OpenvSwitch。
    
- metering agent
    
    计量代理，为租户网络提供三层网络流量计量服务。
    

## RPC

1. *# 消息队列的生产者类(xxxxNotifyAPI)和对应的消费者类（xxxxRpcCallback）定义有相同的接口函数，*
2. *# 生产者类中的函数主要作用是rpc调用消费者类中的同名函数，消费者类中的函数执行实际的动作。*
3. *# 如：xxxNotifyAPI类中定义有network_delete()函数，则xxxRpcCallback类中也会定义有network_delete()函数。*
4. *# xxxNotifyAPI::network_delete()通过rpc调用xxxRpcCallback::network_delete()函数，*
5. *# xxxRpcCallback::network_delete()执行实际的network delete删除动作*

## Ml2相关介绍

[openstack-neutron-ML2_chenxiangui88的博客-CSDN博客](https://blog.csdn.net/chenxiangui88/article/details/78093318)

## Neutron网络架构分析

**查/var/lib/neutron/dhcp**  

得 dhcp缓存信息（实际上为绑定的network id）：

**6a5d2c52-20f5-428e-85bc-775552172c84** 

**78eb67f0-44f1-4d21-8e00-452a6c17f2dc** 

→ 即相应的network中的dhcp信息