HTTP 状态码速记表（面试用）

| 类别           | 常用码 | 短语                    | 一句话场景                  |
| ------------ | --- | --------------------- | ---------------------- |
| **1xx 信息**   | 100 | Continue              | 客户端先问“能发吗？”服务器回“继续”    |
| **2xx 成功**   | 200 | OK                    | 请求成功，正常返回              |
|              | 201 | Created               | 新建资源成功（POST/PUT）       |
|              | 204 | No Content            | 成功但无实体返回（DELETE 成功）    |
| **3xx 重定向**  | 301 | Moved Permanently     | 永久搬家，Location 给新地址     |
|              | 302 | Found                 | 临时搬家，浏览器自动 GET         |
|              | 304 | Not Modified          | 缓存有效，直接用本地副本           |
| **4xx 客户端错** | 400 | Bad Request           | 报文语法/参数错误              |
|              | 401 | Unauthorized          | 未认证，需 WWW-Authenticate |
|              | 403 | Forbidden             | 已认证但无权限                |
|              | 404 | Not Found             | 资源不存在                  |
|              | 409 | Conflict              | 资源冲突（如版本冲突）            |
|              | 429 | Too Many Requests     | 触发限流                   |
| **5xx 服务端错** | 500 | Internal Server Error | 服务器内部异常                |
|              | 502 | Bad Gateway           | 网关/代理收到无效响应            |
|              | 503 | Service Unavailable   | 服务暂时不可用（维护/过载）         |
|              | 504 | Gateway Timeout       | 网关/代理等待后端超时            |

一句话口诀  
“1 继续 2 成功 3 重定向 4 客户端错 5 服务端哭”。