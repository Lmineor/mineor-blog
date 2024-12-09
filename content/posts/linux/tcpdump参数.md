---
title: "tcpdump参数"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

Tcpdump命令参数详解 
## tcpdump的选项介绍:

tcpdump[ -adeflnNOpqStvx ] [-c 数量] [-F 文件名] [-i 网络接口] [-r 文件名] [-s snaplen] [-T 类型] [-w 文件名] [表达式]

各参数说明如下：

    -a    将网络地址和广播地址转变成名字；

    -b    在数据-链路层上选择协议，包括ip、arp、rarp、ipx都是这一层的。tcpdump -b arp 将只显示网络中的arp即地址转换协议信息；

    -c    在收到指定数目的包后，tcpdump就会停止；

    -d    将匹配信息包的代码以人们能够理解的汇编格式给出；

    -dd   将匹配信息包的代码以c语言程序段的格式给出；

    -ddd  将匹配信息包的代码以十进制的形式给出；

    -e    在输出行打印出数据链路层的头部信息；

    -f    将外部的Internet地址以数字的形式打印出来；

    -F    从指定的文件中读取表达式,忽略其它的表达式；

    -i    指定监听的网络接口；

    -l    使标准输出变为缓冲行形式,如tcpdump -l >tcpcap.txt将得到的数据存入tcpcap.txt文件中；

    -n    不进行IP地址到主机名的转换；

    -N    不打印出默认的域名

    -nn   n不进行端口名称的转换；

    -O    不进行匹配代码的优化，当怀疑某些bug是由优化代码引起的, 此选项将很有用；

    -r    从指定的文件中读取包(这些包一般通过-w选项产生)；

    -s    抓取数据包时默认抓取长度为68字节。加上 -s 0 后可以抓到完整的数据包

    -t    在输出的每一行不打印UNIX时间戳，也就是不显示时间；

    -T    将监听到的包直接解释为指定的类型的报文，常见的类型有rpc(远程过程调用)和snmp；

    -tt   打印原始的、未格式化过的时间；

    -v    输出一个稍微详细的信息，例如在ip包中可以包括ttl和服务类型的信息；

    -vv   输出详细的报文信息；

    -w    直接将包写入文件中，并不分析和打印出来；



tcpdump [-i 网卡] -nnAX '表达式'

    -i：   interface 监听的网卡。

    -nn：  表示以ip和port的方式显示来源主机和目的主机，而不是用主机名和服务。

    -A：   以ascii的方式显示数据包，抓取web数据时很有用。

    -X：   数据包将会以16进制和ascii的方式显示。

    表达式：表达式有很多种，常见的有：host 主机；port 端口；src host 发包主机；dst host 收包主机。多个条件可以用and、or组合，取反可以使用!，更多的使用可以查看man 7 pcap-filter。



## 以下是tcpdump的其他一些示例

1、抓取包含10.10.10.122的数据包 

tcpdump -i eth0 -vnn host 10.10.10.122



2、抓取包含10.10.10.0/24网段的数据包

tcpdump -i eth0 -vnn net 10.10.10.0/24



3、抓取包含端口22的数据包

tcpdump -i eth0 -vnn port 22



4、抓取udp协议的数据包

tcpdump -i eth0 -vnn  udp



5、抓取icmp协议的数据包

tcpdump -i eth0 -vnn icmp



6、抓取arp协议的数据包

tcpdump -i eth0 -vnn arp



7、抓取ip协议的数据包

tcpdump -i eth0 -vnn ip



8、抓取源ip是10.10.10.122数据包。

tcpdump -i eth0 -vnn src host 10.10.10.122



9、抓取目的ip是10.10.10.122数据包

tcpdump -i eth0 -vnn dst host 10.10.10.122



10、抓取源端口是22的数据包

tcpdump -i eth0 -vnn src port 22



11、抓取源ip是10.10.10.253且目的ip是22的数据包

tcpdump -i eth0 -vnn src host 10.10.10.253 and dst port 22

            

12、抓取源ip是10.10.10.122或者包含端口是22的数据包

tcpdump -i eth0 -vnn src host 10.10.10.122 or port 22



13、抓取源ip是10.10.10.122且端口不是22的数据包

tcpdump -i eth0 -vnn src host 10.10.10.122 and not port 22



14、抓取源ip是10.10.10.2且目的端口是22，或源ip是10.10.10.65且目的端口是80的数据包。

tcpdump -i eth0 -vnn \( src host 10.10.10.2 and dst port 22 \) or \( src host 10.10.10.65 and dst port 80 \)



15、抓取源ip是10.10.10.59且目的端口是22，或源ip是10.10.10.68且目的端口是80的数据包。

tcpdump -i  eth0 -vnn 'src host 10.10.10.59 and dst port 22' or  ' src host 10.10.10.68 and dst port 80 '



16、把抓取的数据包记录存到/tmp/fill文件中，当抓取100个数据包后就退出程序。

tcpdump –i eth0 -vnn -w  /tmp/fil1 -c 100



17、从/tmp/fill记录中读取tcp协议的数据包

tcpdump –i eth0 -vnn -r  /tmp/fil1 tcp



18、从/tmp/fill记录中读取包含10.10.10.58的数据包

tcpdump –i eth0 -vnn -r /tmp/fil1 host 10.10.10.58



19、假如要抓vlan 1的包，命令格式如下：

tcpdump -i eth0 port 80 and vlan 1 -w /tmp/vlan.cap



20、在后台抓eth0在80端口的包，命令格式如下：

nohup tcpdump -i eth0 port 80 -w /tmp/temp.cap &



21、ARP包的tcpdump输出信息

tcpdump arp -nvv



22、使用tcpdump抓取与主机192.168.43.23或着与主机192.168.43.24通信报文，并且显示在控制台上

tcpdump -X -s 1024 -i eth0 host \(192.168.43.23 or 192.168.43.24\) and host 172.16.70.35



23、常用命令收藏

tcpdump -i eth0 -nn 'dst host 172.100.6.231'

tcpdump -i eth0 -nn 'src host 172.100.6.12'

tcpdump -i eth0 -nnA 'port 80'

tcpdump -i eth0 -XnnA 'port 22'

tcpdump -i eth0 -nnA 'port 80 and src host 192.168.1.231'

tcpdump -i eth0 -nnA '!port 22' and 'src host 172.100.6.230'

tcpdump -i eth0 -nnA '!port 22'
