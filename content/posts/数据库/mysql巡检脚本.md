---
title: mysql巡检脚本
date: 2023-06-18
draft: false
tags:
  - 数据库
categories:
  - 技术
---

```bash
#!/bin/bash
# MySQL巡检脚本
# 设置MySQL用户名和密码（请将它们设置为适当的值）
MYSQL_USER="root"
MYSQL_PASSWORD="123456"
# 获取MySQL版本信息
MYSQL_VERSION=$(mysql -u ${MYSQL_USER} -p${MYSQL_PASSWORD} -e "SELECT VERSION();" | awk 'NR==2{print $1}')
# 获取MySQL运行状态信息
STATUS=$(systemctl status mysql.service)
# 获取MySQL进程列表
PROCESS_LIST=$(mysql -u ${MYSQL_USER} -p${MYSQL_PASSWORD} -e "SHOW PROCESSLIST;" | awk '{print $1,$2,$3,$4,$5,$6}')
# 检查MySQL是否在运行
if [[ "$STATUS" =~ "active (running)" ]]; then
  MYSQL_RUNNING="YES"
else
  MYSQL_RUNNING="NO"
fi

# 检查MySQL进程是否存在
if [[ -z "$PROCESS_LIST" ]]; then
  MYSQL_PROCESS="NO"
else
  MYSQL_PROCESS="YES"
fi

# 生成报告

echo "MySQL巡检报告"
echo "----------------"
echo "MySQL版本: $MYSQL_VERSION"
echo "MySQL运行状态: $MYSQL_RUNNING"
echo "MySQL进程存在: $MYSQL_PROCESS"
echo ""
echo "MySQL进程列表"
echo "----------------"
echo "$PROCESS_LIST"

```