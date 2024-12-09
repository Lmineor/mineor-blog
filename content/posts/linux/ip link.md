---
title: "ip link"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

ip link命令格式

命令实例
`ip link add link eth0 name eth0.10 type vlan id 10`

|命令实例|解释|
|:-:|:-:|
|ip link add link eth0 name eth0.10 type vlan id 10|在设备eth0上创建新的vlan设备eth0.10|
|||
|ip link set eth0 up 或：ifconfig eth0 up|开启eth0网卡|
|ip link set eth0 down或：ifconfig eth0 down|关闭eth0网卡|
|ip link set eth0 promisc on|开启网卡的混合模式|
|ip link set eth0 promisc off|关闭网卡的混合模式|
|ip link set eth0 txqueuelen 1200|设置网卡队列长度|
|ip link set eth0 mtu 1400|设置网卡最大传输单元|
|||
|ip link show|显示网络接口信息|
|ip link show eht0|显示eth0网卡的网络接口信息|
|ip link show type vlan|显示vlan类型设备|
|||
|ip link delete dev eth0.10|删除设备|

```bash
[root@k8s-master ~]# ip link help link
Usage: ip link add [link DEV] [ name ] NAME
                   [ txqueuelen PACKETS ]
                   [ address LLADDR ]
                   [ broadcast LLADDR ]
                   [ mtu MTU ] [index IDX ]
                   [ numtxqueues QUEUE_COUNT ]
                   [ numrxqueues QUEUE_COUNT ]
                   type TYPE [ ARGS ]

       ip link delete { DEVICE | dev DEVICE | group DEVGROUP } type TYPE [ ARGS ]

       ip link set { DEVICE | dev DEVICE | group DEVGROUP }
                          [ { up | down } ]
                          [ type TYPE ARGS ]
                          [ arp { on | off } ]
                          [ dynamic { on | off } ]
                          [ multicast { on | off } ]
                          [ allmulticast { on | off } ]
                          [ promisc { on | off } ]
                          [ trailers { on | off } ]
                          [ carrier { on | off } ]
                          [ txqueuelen PACKETS ]
                          [ name NEWNAME ]
                          [ address LLADDR ]
                          [ broadcast LLADDR ]
                          [ mtu MTU ]
                          [ netns { PID | NAME } ]
                          [ link-netnsid ID ]
                          [ alias NAME ]
                          [ vf NUM [ mac LLADDR ]
                                   [ vlan VLANID [ qos VLAN-QOS ] [ proto VLAN-PROTO ] ]
                                   [ rate TXRATE ]
                                   [ max_tx_rate TXRATE ]
                                   [ min_tx_rate TXRATE ]
                                   [ spoofchk { on | off} ]
                                   [ query_rss { on | off} ]
                                   [ state { auto | enable | disable} ] ]
                                   [ trust { on | off} ] ]
                                   [ node_guid { eui64 } ]
                                   [ port_guid { eui64 } ]
                          [ xdp { off |
                                  object FILE [ section NAME ] [ verbose ] |
                                  pinned FILE } ]
                          [ master DEVICE ][ vrf NAME ]
                          [ nomaster ]
                          [ addrgenmode { eui64 | none | stable_secret | random } ]
                          [ protodown { on | off } ]

       ip link show [ DEVICE | group GROUP ] [up] [master DEV] [vrf NAME] [type TYPE]

       ip link xstats type TYPE [ ARGS ]

       ip link afstats [ dev DEVICE ]

       ip link help [ TYPE ]

TYPE := { vlan | veth | vcan | dummy | ifb | macvlan | macvtap |
          bridge | bond | team | ipoib | ip6tnl | ipip | sit | vxlan |
          gre | gretap | ip6gre | ip6gretap | vti | nlmon | team_slave |
          bond_slave | ipvlan | geneve | bridge_slave | vrf | macsec }
```
