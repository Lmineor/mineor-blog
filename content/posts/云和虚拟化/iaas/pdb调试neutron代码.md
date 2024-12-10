---
title: "pdb调试neutron代码"
date: 2024-10-20
draft: false
tags : [                    # 文章所属标签
    "Iaas"
]
---


### 1、停服务

```bash
systemctl stop neutron-server.service
```

### 2、插入调试代码段

```bash
import pdb;pdb.set_trace()
```

代码路径 /usr/lib/python2.7/site-packages/neutron/plugins/

### 3、查看服务状态，手动停服务，打断点，手动启服务

- `systemctl status/stop neutron-server.service`
- `pgrep neutron-server`
- `pkill neutron-server`

查看服务启动项/加载项

`cat /usr/lib/systemd/system/neutron-server.service`

```bash
[root@controller ~]# cat /usr/lib/systemd/system/neutron-server.service
[Unit]
Description=OpenStack Neutron Server
After=syslog.target network.target

[Service]
Type=notify
User=neutron
ExecStart=/usr/bin/neutron-server --config-file /usr/share/neutron/neutron-dist.conf --config-dir /usr/share/neutron/server --config-file /etc/neutron/neutron.conf --config-file /etc/neutron/plugin.ini --config-dir /etc/neutron/conf.d/common --config-dir /etc/neutron/conf.d/neutron-server --log-file /var/log/neutron/server.log
PrivateTmp=true
NotifyAccess=all
KillMode=process
Restart=on-failure
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

**手动启服务[neutron]**

`su -s /bin/sh -c 'ExecStart' User`

`''`替换加载项(ExecStart)，**带单引号**；`User`替换对应节点的服务的用户

### 4、`neutron cli/web`进行跟进

### 5、举例

```bash
su -s /bin/sh -c '/usr/bin/neutron-server --config-file /usr/share/neutron/neutron-dist.conf --config-dir /usr/share/neutron/server --config-file /etc/neutron/neutron.conf --config-file /etc/neutron/plugin.ini --config-dir /etc/neutron/conf.d/common --config-dir /etc/neutron/conf.d/neutron-server --log-file /var/log/neutron/server.log' neutron
```