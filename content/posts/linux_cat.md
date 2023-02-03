---
title: "linux cat"
date: 2023-02-03
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

# 利用cat 给文件写内容,追加的方式

```bash
cat >> proxy.sh <<EOF
export http_proxy=http://99.0.85.1:808
export https_proxy=http://99.0.85.1:808
EOF
```

# cat给文件写内容,覆盖的方式

```bash
cat > proxy.sh <<EOF
export http_proxy=http://99.0.85.1:808
export https_proxy=http://99.0.85.1:808
EOF
```