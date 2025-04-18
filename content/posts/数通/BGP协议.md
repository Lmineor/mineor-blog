---
title: BGP协议
date: 2025-04-06
draft: false
tags:
  - 数通
---
边界网关协议（Border Gateway Protocol，BGP）是一种用来在路由选择域之间交换网络层可达性信息（Network Layer Reachability Information，NLRI）的路由选择协议。
由于不同的管理机构分别控制着他们各自的路由选择域，因此，路由选择域经常被称为自治系统AS（Autonomous System）。

现在的Internet是一个由多个自治系统相互连接构成的大网络，BGP作为事实上的Internet外部路由协议标准，被广泛应用于ISP（Internet Service Provider）之间。 

早期发布的三个版本分别是BGP-1、BGP-2和BGP-3，主要用于交换AS之间的可达路由信息，构建AS域间的传播路径，防止路由环路的产生，并在AS级别应用一些路由策略。当前使用的版本是BGP-4。