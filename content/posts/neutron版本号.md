---
title: "获取neutron版本号"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "SDN",
    "OpenStack"
]
categories : [              # 文章所属标签
    "技术",
    "美食",
    "生活",
    "阅读",
]
---

```bash
python -c  "import neutron.version;print(neutron.version.version_info)" 
```
