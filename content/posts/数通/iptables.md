---
title: Table
date: 星期四, 七月 17日 2025, 2:57:32 下午
draft: true
tags:
  - 数通
---
# Table

## Mangle表

主要用于修改数据包的TOS（Type Of Service，服务类型）、TTL（Time To Live，生存周期）指以及为数据包设置Mark标记，以实现Qos(Quality Of Service，服务质量)调整以及策略路由等应用，由于需要相应的路由设备支持，因此应用并不广泛。包含五个规则链——PREROUTING，POSTROUTING，INPUT，OUTPUT，FORWARD。

强烈建议你不要在这个表里做任何过滤，不管是DANT，SNAT或者Masquerade。

## Filter表

主要用于对数据包进行过滤，根据具体的规则决定是否放行该数据包（如DROP、ACCEPT、REJECT、LOG）。filter 表对应的内核模块为iptable_filter，包含三个规则链：

INPUT链：INPUT针对那些目的地是本地的包

FORWARD链：FORWARD过滤所有不是本地产生的并且目的地不是本地(即本机只是负责转发)的包

OUTPUT链：OUTPUT是用来过滤所有本地生成的包

## Nat表

主要用于修改数据包的IP地址、端口号等信息（网络地址转换，如SNAT、DNAT、MASQUERADE、REDIRECT）。属于一个流的包(因为包的大小限制导致数据可能会被分成多个数据包)只会经过这个表一次。如果第一个包被允许做NAT或Masqueraded，那么余下的包都会自动地被做相同的操作，也就是说，余下的包不会再通过这个表。表对应的内核模块为 iptable_nat，包含三个链：

PREROUTING链：作用是在包刚刚到达防火墙时改变它的目的地址

OUTPUT链：改变本地产生的包的目的地址

POSTROUTING链：在包就要离开防火墙之前改变其源地址

## Raw表

是自1.2.9以后版本的iptables新增的表，主要用于决定数据包是否被状态跟踪机制处理。在匹配数据包时，raw表的规则要优先于其他表。包含两条规则链——OUTPUT、PREROUTING

iptables中数据包和4种被跟踪连接的4种不同状态：

NEW：该包想要开始一个连接（重新连接或将连接重定向）

RELATED：该包是属于某个已经建立的连接所建立的新连接。例如：FTP的数据传输连接就是控制连接所 RELATED出来的连接。

`--icmp-type 0`( ping 应答) 就是`--icmp-type 8`(ping 请求)所RELATED出来的。

ESTABLISHED ：只要发送并接到应答，一个数据连接从NEW变为ESTABLISHED,而且该状态会继续匹配这个连接的后续数据包。

INVALID：数据包不能被识别属于哪个连接或没有任何状态比如内存溢出，收到不知属于哪个连接的ICMP错误信息，一般应该DROP这个状态的任何数据。

## 规则表之间的优先顺序

Raw>>>>mangle>>>>nat>>>>filter

![iptables](../../../images/network/iptables.png)