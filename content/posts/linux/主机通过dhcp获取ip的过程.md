---
title: "主机通过dhcp获取ip的过程"
date: 2023-12-06
draft: false
tags : [
    "Linux",
]
categories : [
    "技术",
]
---

虚拟机通过DHCP获取IP的过程通常包括以下步骤：

DHCP Discover：虚拟机启动时，会发送一个DHCP Discover广播包到网络中，该包中包含了对DHCP服务器的请求。

```
14:31:53.831654 fa:16:3e:d6:62:84 > ff:ff:ff:ff:ff:ff, ethertype IPv4 (0x0800), length 342: (tos 0x10, ttl 128, id 0, offset 0, flags [none], proto UDP (17), length 328)
    0.0.0.0.68 > 255.255.255.255.67: [udp sum ok] BOOTP/DHCP, Request from fa:16:3e:d6:62:84, length 300, xid 0x6265ab6d, Flags [none] (0x0000)
          Client-Ethernet-Address fa:16:3e:d6:62:84
          Vendor-rfc1048 Extensions
            Magic Cookie 0x63825363
            DHCP-Message Option 53, length 1: Discover
            Requested-IP Option 50, length 4: 1.1.1.13
            Parameter-Request Option 55, length 13:
              Subnet-Mask, BR, Time-Zone, Classless-Static-Route
              Domain-Name, Domain-Name-Server, Hostname, YD
              YS, NTP, MTU, Option 119
              Default-Gateway
```

DHCP Offer：网络中的DHCP服务器接收到DHCP Discover广播包后，会向虚拟机发送一个DHCP Offer广播包，其中包含了可用的IP地址、子网掩码、网关、DNS服务器等信息。

```
14:31:53.832156 fa:16:3e:4c:de:f2 > fa:16:3e:d6:62:84, ethertype IPv4 (0x0800), length 385: (tos 0xc0, ttl 64, id 11227, offset 0, flags [none], proto UDP (17), length 371)
    1.1.1.2.67 > 1.1.1.13.68: [udp sum ok] BOOTP/DHCP, Reply, length 343, xid 0x6265ab6d, Flags [none] (0x0000)
          Your-IP 1.1.1.13
          Server-IP 1.1.1.2
          Client-Ethernet-Address fa:16:3e:d6:62:84
          Vendor-rfc1048 Extensions
            Magic Cookie 0x63825363
            DHCP-Message Option 53, length 1: Offer
            Server-ID Option 54, length 4: 1.1.1.2
            Lease-Time Option 51, length 4: 86400
            RN Option 58, length 4: 43200
            RB Option 59, length 4: 75600
            Subnet-Mask Option 1, length 4: 255.255.255.0
            BR Option 28, length 4: 1.1.1.255
            Domain-Name-Server Option 6, length 4: 1.1.1.2
            Domain-Name Option 15, length 14: "openstacklocal"
            Hostname Option 12, length 13: "host-1-1-1-13"
            MTU Option 26, length 2: 1450
            Default-Gateway Option 3, length 4: 1.1.1.1
            Classless-Static-Route Option 121, length 14: (169.254.169.254/32:1.1.1.2),(default:1.1.1.1)

```

DHCP Request：虚拟机收到DHCP Offer后，会选择其中一个提供的IP地址，并发送一个DHCP Request广播包到网络中，请求该IP地址。

```
14:31:53.859117 fa:16:3e:d6:62:84 > ff:ff:ff:ff:ff:ff, ethertype IPv4 (0x0800), length 342: (tos 0x10, ttl 128, id 0, offset 0, flags [none], proto UDP (17), length 328)
    0.0.0.0.68 > 255.255.255.255.67: [udp sum ok] BOOTP/DHCP, Request from fa:16:3e:d6:62:84, length 300, xid 0x6265ab6d, Flags [none] (0x0000)
          Client-Ethernet-Address fa:16:3e:d6:62:84
          Vendor-rfc1048 Extensions
            Magic Cookie 0x63825363
            DHCP-Message Option 53, length 1: Request
            Server-ID Option 54, length 4: 1.1.1.2
            Requested-IP Option 50, length 4: 1.1.1.13
            Parameter-Request Option 55, length 13:
              Subnet-Mask, BR, Time-Zone, Classless-Static-Route
              Domain-Name, Domain-Name-Server, Hostname, YD
              YS, NTP, MTU, Option 119
              Default-Gateway
```

DHCP Acknowledgement：DHCP服务器收到虚拟机的DHCP Request后，会向虚拟机发送一个DHCP Acknowledgement广播包，确认虚拟机获取到了该IP地址，并在包中携带了一些配置信息。

```
14:31:53.872377 fa:16:3e:4c:de:f2 > fa:16:3e:d6:62:84, ethertype IPv4 (0x0800), length 385: (tos 0xc0, ttl 64, id 11234, offset 0, flags [none], proto UDP (17), length 371)
    1.1.1.2.67 > 1.1.1.13.68: [udp sum ok] BOOTP/DHCP, Reply, length 343, xid 0x6265ab6d, Flags [none] (0x0000)
          Your-IP 1.1.1.13
          Server-IP 1.1.1.2
          Client-Ethernet-Address fa:16:3e:d6:62:84
          Vendor-rfc1048 Extensions
            Magic Cookie 0x63825363
            DHCP-Message Option 53, length 1: ACK
            Server-ID Option 54, length 4: 1.1.1.2
            Lease-Time Option 51, length 4: 86400
            RN Option 58, length 4: 43200
            RB Option 59, length 4: 75600
            Subnet-Mask Option 1, length 4: 255.255.255.0
            BR Option 28, length 4: 1.1.1.255
            Domain-Name-Server Option 6, length 4: 1.1.1.2
            Domain-Name Option 15, length 14: "openstacklocal"
            Hostname Option 12, length 13: "host-1-1-1-13"
            MTU Option 26, length 2: 1450
            Default-Gateway Option 3, length 4: 1.1.1.1
            Classless-Static-Route Option 121, length 14: (169.254.169.254/32:1.1.1.2),(default:1.1.1.1)
```

IP地址分配：虚拟机接收到DHCP Acknowledgement后，会配置自身的网络接口，使用DHCP服务器提供的IP地址、子网掩码、网关和DNS服务器等信息，完成IP地址的获取和配置。
