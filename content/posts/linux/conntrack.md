---
title: "conntrack相关字段"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---


参考连接[conntrack](https://docs.openvswitch.org/en/latest/tutorials/ovs-conntrack/#definitions)

OVS supports following match fields related to conntrack:

1. `ct_state`: The state of a connection matching the packet. 可能的值有:
```bash
new
est
rel
rpl
inv
trk
snat
dnat
```

上述字段
如果前面有`+`号, 则说明是必须设置的标记;
如果前面有`-`号, 则说明必须取消的标记.

同样也支持多个字段同时设置, 例如 ct_state=+trk+new


参考[ovs-fields](http://openvswitch.org/support/dist-docs/ovs-fields.7.txt)了解更多

2. `ct_zone`:16bit字段用作另一个flow entry的匹配field
3. `ct_mark`:给当前连接的包打上32bit的metadata数据
4. `ct_label`给当前连接的包打上128bit的label
5. ...