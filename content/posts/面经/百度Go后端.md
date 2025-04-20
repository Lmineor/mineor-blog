---
title: 百度智能运维岗
date: 2025-04-05
draft: true
tags:
  - 面经
---

// go的goroutine和python线程的区别
// 1. goroutine是轻量级的线程，python的线程是重的线程
// 2. goroutine是协作式的，python的线程是抢占式的
// 3. goroutine是非阻塞的，python的线程是阻塞的
// 4. goroutine是非抢占式的，python的线程是抢占式的

// k8s编辑yaml，更新镜像后，apply该yaml， pod的变化状态都有什么，用户正在访问的话会发生什么，该如何避免这种情况的发生
// flannel的流量怎么走的(pod间互访）
// 客户k8s集群的规模
// k8s核心原理
// k8s从create一个pod的yaml，到pod的创建都发生了什么
// www.baidu.com从回车到显示都发生了什么
// go 通道的底层原理
// k8s的调度器

# 数据量很大，达到内存放不下，怎么办？

1. 首先需要明确，数据量为什么很大，那么可以从这几个方面考虑：
	1. 缓存的数据是否合理
	2. 数据是否都设置了有效时间
	3. 数据有效期是否合理
2. 如果上述都完成后确实因为业务量过大，那么则需要考虑扩容

# map和slice哪个是线程安全的