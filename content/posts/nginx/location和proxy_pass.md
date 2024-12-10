---
title: "proxy_pass的末尾带不带/的作用"
date: 2024-06-30
draft: false
tags : [                    # 文章所属标签
    "nginx",
]
categories : [              # 文章所属标签
    "技术",
]
---




在 Nginx 的 proxy_pass 指令中，末尾的斜线（/）对代理行为有重要的影响。这主要涉及到请求 URI 的处理方式。

# 末尾不带斜线

如果你配置 `proxy_pass` 不带末尾的斜线，Nginx 会将请求的 URI 完整地传递给后端服务器，包括任何路径信息。

例如：

```nginx
location /somepath/ {  
    proxy_pass http://localhost:8080;  
}
```

当客户端请求 `/somepath/foo` 时，Nginx 会将请求代理到 `http://localhost:8080/somepath/foo`。

末尾带斜线
如果你配置 `proxy_pass` 带有末尾的斜线，Nginx 会修改请求的 URI，移除与 location 块匹配的部分，然后将剩余部分（如果有的话）传递给后端服务器。

例如：

```nginx
location /somepath/ {  
    proxy_pass http://localhost:8080/;  
}
```


当客户端请求 `/somepath/foo` 时，Nginx 会将请求代理到 `http://localhost:8080/foo`（注意 `/somepath/` 部分被移除了）。

注意事项

当你在 location 块中使用正则表达式时，末尾的斜线通常不需要（或不建议）使用，因为正则表达式已经能够捕获到请求的 URI 部分。
如果你的后端服务器期望接收到完整的 URI（包括 location 块中定义的路径），那么你应该在 proxy_pass 中省略末尾的斜线。
如果你希望 Nginx 剥离 location 块中定义的路径部分，只将剩余的 URI 部分传递给后端服务器，那么你应该在 proxy_pass 中包含末尾的斜线。
示例
假设你有一个后端服务，它只响应根路径（/）下的请求。你可以使用以下配置来将所有到 /api/ 的请求代理到这个服务的根路径：

```nginx
location /api/ {  
    proxy_pass http://localhost:8080/;  
}
```

这样，无论客户端请求 `/api/foo` 还是 `/api/bar`，请求都会被代理到 `http://localhost:8080/`，并且后端服务可以在接收到请求后自行解析 `/foo` 或 `/bar` 部分。
