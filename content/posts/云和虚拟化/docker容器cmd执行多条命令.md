---
title: "docker容器cmd执行多条命令"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "云与虚拟化",
]

---

```bash
mysql:
    build:
      context: ./mariadb
      dockerfile: ./Dockerfile
    container_name: mariadb_yoga
    restart: on-failure
    command:
      - /bin/bash
      - -c
      - |
        mysqld --user=root
        mysql -u root -p 123456 < /db/init_db.sql
    volumes:
      - /var/lib/openstack/yoga/mysql:/var/lib/mysql
      - /etc/localtime:/etc/localtime
    ports:
      - "23306:3306"  # host物理直接映射端口为13306
    environment:
      MYSQL_ROOT_PASSWORD: '123456' # root管理员用户密码
    networks:
      network:
        ipv4_address: 177.177.0.13
```