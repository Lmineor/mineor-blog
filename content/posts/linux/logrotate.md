---
title: "linux日志轮转"
date: 2023-12-01
draft: false
tags : [
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---


`/etc/logrotate.conf` 文件是系统级别的 logrotate 配置文件，它通常用于配置全局的日志轮转规则和选项。如果您希望针对您的项目或应用程序设置特定的日志轮转规则，
可以创建一个单独的 logrotate 配置文件并将其放置在 /etc/logrotate.d/ 目录中。

下面是一些步骤来创建并使用一个针对您的项目的 logrotate 配置文件：

# 创建配置文件：

在 `/etc/logrotate.d/` 目录中创建一个新的文件，可以使用您项目的名称或相关描述来命名该文件，例如 `my_project_logrotate`。

# 编辑配置文件：

使用文本编辑器（如vi或nano）编辑该文件，添加您项目的日志文件、轮转规则和选项。以下是一个简单的示例：

```bash
/path/to/your/project/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 <user> <group>
    su <user> <group>
}
```

- daily：每天轮转日志。
- rotate 7：保留最多7个旧日志文件。
- compress：压缩轮转后的日志文件。
- delaycompress：在下一次轮转时再压缩旧的日志文件。
- missingok：如果日志文件丢失，不发出警告消息。
- notifempty：如果日志文件为空，则不轮转。
- create 0644 <user> <group>：创建新的日志文件并设置权限。
- su <user> <group>：指定轮转时使用的用户和组。

将上述配置中的 /path/to/your/project/logs/*.log 替换为实际项目日志文件的路径和模式。

# 测试和激活配置：

可以使用 `logrotate` 命令进行手动测试您的配置文件：

```bash
logrotate -d /etc/logrotate.d/my_project_logrotate
```

上述命令会模拟运行 `logrotate`，并显示它将执行的轮转操作。如果一切正常，您可以使用以下命令来手动轮转日志：

```bash
logrotate -f /etc/logrotate.d/my_project_logrotate
```

这将强制执行日志轮转。

# 定期运行logrotate：

为了定期自动运行 logrotate，通常会有一个 cron 作业来执行该操作。logrotate本身通常在 cron 中预先配置好，以便在系统中定期运行。您可以检查 `/etc/cron.daily/logrotate` 或类似的位置，以确保logrotate被定期执行。

通过这些步骤，您可以为您的项目创建一个独立的 logrotate 配置文件，并定期对日志文件进行轮转和管理。
