---
title: "《k8s权威指南学习》--Service"
date: 2023-03-26
draft: false
tags : [                    # 文章所属标签
    "k8s",
]
---

# Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app.kubernetes.io/name: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
```