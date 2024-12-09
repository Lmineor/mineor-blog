---
title: "SNAT和DNAT"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
    "k8s",
    "Cloud",
    "SDN",
]
categories : [              # 文章所属标签
    "技术",
]
---

[参考](https://docs.openstack.org/neutron/latest/admin/intro-nat.html)

# SNAT

In Source Network Address Translation (SNAT), the NAT router modifies the IP address of the sender in IP packets. SNAT is commonly used to enable hosts with private addresses to communicate with servers on the public Internet.

RFC 1918 reserves the following three subnets as private addresses:

10.0.0.0/8

172.16.0.0/12

192.168.0.0/16

# DNAT

In Destination Network Address Translation (DNAT), the NAT router modifies the IP address of the destination in IP packet headers.

OpenStack uses DNAT to route packets from instances to the OpenStack metadata service. Applications running inside of instances access the OpenStack metadata service by making HTTP GET requests to a web server with IP address 169.254.169.254. In an OpenStack deployment, there is no host with this IP address. Instead, OpenStack uses DNAT to change the destination IP of these packets so they reach the network interface that a metadata service is listening on.