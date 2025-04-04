---
title: "cu"
date: 2024-08-08
draft: false
categories: ["技术"]
tags: ["Linux"]
---
cu
===

用于连接另一个系统主机

## 补充说明

**cu命令** 用于连接另一个系统主机。cu(call up)指令可连接另一台主机，并采用类似拨号终端机的接口工作，也可执行简易的文件传输作业。

###  语法

```shell
cu [dehnotv][-a<通信端口>][-c<电话号码>][-E<脱离字符>][-I<设置文件>][-l<外围设备代号>]
[-s<连线速率>][-x<排错模式>][-z<系统主机>][--help][-nostop][--parity=none][<系统主机>/<电话号码>]
```

###  选项

```shell
-a<通信端口>或-p<通信端口>或--port<通信端口> 使用指定的通信端口进行连线。
-c<电话号码>或--phone<电话号码> 拨打该电话号码。
-d 进入排错模式。
-e或--parity=even 使用双同位检查。
-E<脱离字符>或--escape<脱离字符> 设置脱离字符。
-h或--halfduple 使用半双工模式。
-I<配置文件>或--config<配置文件> 指定要使用的配置文件。
-l<外围设备代号>或--line<外围设备代号> 指定某项外围设备，作为连接的设备。
-n或--prompt 拨号时等待用户输入电话号码。
-o或--parity=odd 使用单同位检查。
-s<连线速率>或--speed<连线速率>或--baud<连线速率>或-<连线速率> 设置连线的速率，单位以鲍率计算。
-t或--maper 把CR字符置换成LF+CR字符。
-v或--version 显示版本信息。
-x<排错模式>或--debug<排错模式> 使用排错模式。
-z<系统主机>或--system<系统主机> 连接该系统主机。
--help 在线帮助。
--nostop 关闭Xon/Xoff软件流量控制。
--parity=none 不使用同位检查。
```

### 实例

与远程主机连接

```shell
cu -c 0102377765
cu -s 38400 9=12015551234
```




