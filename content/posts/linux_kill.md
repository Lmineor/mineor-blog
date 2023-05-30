---
title: "kill 命令的用途"
date: 2023-05-30
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

# Linux kill

用途：kill – terminate or signal a process

kill 是向进程发送信号的命令。当然我们可以向进程发送一个终止运行的信号，此时的 kill 命令才是名至实归。事实上如果我们不给 kill 命令传递信号参数，它默认传递终止进程运行的信号给进程！这是 kill 命令最主要的用法，也是本文要介绍的内容。

一般情况下，终止一个前台进程使用 Ctrl + C 就可以了。对于一个后台进程就须用 kill 命令来终止。我们会先使用 ps、top 等命令获得进程的 PID，然后使用 kill 命令来杀掉该进程。


## kill命令格式

```bash
kill [options] <pid> [...]

<pid> […] : 把信号发送给列出的所有进程。

options :

    -<signal> : 指定发送给进程的信号，指定信号的名称或号码都可以。

    -l : 列出所有信号的名称和号码。

```

## 有哪些信号可以发送给进程

```bash
[root@centos ~]# kill -l
 1) SIGHUP	 2) SIGINT	 3) SIGQUIT	 4) SIGILL	 5) SIGTRAP
 6) SIGABRT	 7) SIGBUS	 8) SIGFPE	 9) SIGKILL	10) SIGUSR1
11) SIGSEGV	12) SIGUSR2	13) SIGPIPE	14) SIGALRM	15) SIGTERM
16) SIGSTKFLT	17) SIGCHLD	18) SIGCONT	19) SIGSTOP	20) SIGTSTP
21) SIGTTIN	22) SIGTTOU	23) SIGURG	24) SIGXCPU	25) SIGXFSZ
26) SIGVTALRM	27) SIGPROF	28) SIGWINCH	29) SIGIO	30) SIGPWR
31) SIGSYS	34) SIGRTMIN	35) SIGRTMIN+1	36) SIGRTMIN+2	37) SIGRTMIN+3
38) SIGRTMIN+4	39) SIGRTMIN+5	40) SIGRTMIN+6	41) SIGRTMIN+7	42) SIGRTMIN+8
43) SIGRTMIN+9	44) SIGRTMIN+10	45) SIGRTMIN+11	46) SIGRTMIN+12	47) SIGRTMIN+13
48) SIGRTMIN+14	49) SIGRTMIN+15	50) SIGRTMAX-14	51) SIGRTMAX-13	52) SIGRTMAX-12
53) SIGRTMAX-11	54) SIGRTMAX-10	55) SIGRTMAX-9	56) SIGRTMAX-8	57) SIGRTMAX-7
58) SIGRTMAX-6	59) SIGRTMAX-5	60) SIGRTMAX-4	61) SIGRTMAX-3	62) SIGRTMAX-2
63) SIGRTMAX-1	64) SIGRTMAX
```

这些信号中只有第9中信号（SIGKILL）才可以无条件的终止进程，其他信号进程都有权忽略。

几个常用信号的含义

|代号|名称|内容|
|:-|:-|:-|
|1|SIGHUP|启动被终止的程序，可让该进程重新读取自己的配置文件，类似重新启动。|
|2|SIGINT|相当于用键盘输入 [ctrl]-c 来中断一个程序的进行。|
|9|SIGKILL|代表强制中断一个程序的进行，如果该程序进行到一半，那么尚未完成的部分可能会有“半产品”产生，类似 vim会有 .filename.swp 保留下来。|
|15|SIGTERM|以正常的方式来终止该程序。由于是正常的终止，所以后续的动作会将他完成。不过，如果该程序已经发生问题，就是无法使用正常的方法终止时，输入这个 signal 也是没有用的。|
|19|SIGSTOP|相当于用键盘输入 [ctrl]-z 来暂停一个程序的进行。|

