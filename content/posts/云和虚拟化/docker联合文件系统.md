---
title: "docker联合文件系统"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "云与虚拟化",
]
---


> 参考:https://www.jianshu.com/p/5c1f152ac4a6

# 1. 前言

(Union filesystem)联合文件系统允许我们把多个文件系统逻辑上合并成一个文件系统，组成Union filesystem的文件系统不必相同(它们可以是ext2/3/4,vfat,ntfs,jffs...)。overlay是联合文件系统的一种(a ufs...)，overlay文件系统构建于其他文件系统之上，overlay其实更像是个挂载系统(mount system)，功能是把不同的文件系统挂载到统一的路径。

# 2. overlay文件系统

overlay是个分层的文件系统，底层文件系统通常叫lower，顶层文件系统系统通常叫upper，两者通常合并挂载到merged目录，最终用户看到的就是merged中的文件。

lower文件系统是readonly，对merged中所有的修改都只对upper操作，记住这点很重要。下面我们在linux上创建一个overlay文件系统，用以说明overlay文件系统挂载，文件读写，文件新增和删除。

# 3. 创建overlay文件系统
创建lower upper merged work目录，把lower和upper挂载到merged，work是空目录，必须和merged的文件系统类型一样。

```bash
(base) [root@lex learn_cgroup]# tree upper/
upper/
├── books
│   ├── readher.txt
│   └── readyou.txt
├── readme.txt
└── read.txt

(base) [root@lex learn_cgroup]# tree lower/
lower/
├── books
│   ├── readme.txt
│   └── readyou.txt
└── readme.txt

```

挂载: `mount -t overlay overlay -o lowerdir=./lower,upperdir=./upper,workdir=./work merged`，挂载后merged目录结构如下:

```bash
(base) [root@lex learn_cgroup]# tree merged/
merged/
├── books
│   ├── readher.txt
│   ├── readme.txt
│   └── readyou.txt
├── readme.txt
└── read.txt
```

可以看到lower和upper中的文件合并到了merged中，当lower和upper有相同路径的文件时，merged中只显示upper中的。也就是说upper会遮住lower中同名的文件(同路径下)。

# 4. overlay文件系统读写

## 新建文件:merged层

merged中新建文件,只在upper和merged中可见,lower中并不可见

> 原文:在upper中新建文件，lower只读,经过测试后觉得不正确

## 修改文件:merged层

修改merged中的文件,

- 文件在upper和lower中都存在,则会修改merged和upper中的文件,lower不会修改:

```bash
(base) [root@lex learn_cgroup]# cat merged/readme.txt
readme in upper edit in merged
(base) [root@lex learn_cgroup]# cat upper/readme.txt
readme in upper edit in merged
(base) [root@lex learn_cgroup]# cat lower/readme.txt
readme in lower         
```

- 文件只在upper中存在,则会修改merged中和upper中的文件
- 文件只在lower中存在,则从lower中复制文件到upper，再修改upper中的复制品
```bash
(base) [root@lex learn_cgroup]# cat merged/books/readme.txt
books readme in lower edit in merged
(base) [root@lex learn_cgroup]# cat upper/books/readme.txt
books readme in lower edit in merged
(base) [root@lex learn_cgroup]# cat lower/books/readme.txt
books readme in lower
```


## 读取文件:merged层

upper中有该文件则读取upper的，否则读取lower中的
whiteout是个字符设备文件，主次设备号为0，用来屏蔽对lower的访问。

## 删除文件:merged层

删除merged层的readne.txt[只在merged和lower层中存在]

如果文件在upper中存在，则删除;
如果文件在upper中不存在，在lower中存在，则在upper新建一个同名的whiteout文件

```bash
(base) [root@lex learn_cgroup]# tree upper/
upper/
├── books
│   ├── readher.txt
│   ├── readme.txt
│   └── readyou.txt
├── newfile.txt
└── readme.txt

(base) [root@lex learn_cgroup]# cd upper/books/
(base) [root@lex books]# ll
total 8
-rw-r--r-- 1 root root   23 Dec 14 16:12 readher.tx
c--------- 1 root root 0, 0 Dec 14 16:31 readme.txt # whiteout文件
-rw-r--r-- 1 root root   23 Dec 14 16:12 readyou.txt 
```

## 补充
linux4.0以后，overlay文件系统支持多层lower挂载，挂载方式如下：
sudo mount -t overlay overlay -o lowerdir=./dir1:./dir2:./dir3,upperdir=./upper,workdir=./work merged/
dir1在lower的顶层，dir3在lower的底层，debian上docker就是采用这种方式。
overlay提供了对只读文件系统的读写功能，适合用在需要维持一个只读镜像，又需要提供读写功能的系统中，比如openwrt和docker

# 5. docker
docker的基础镜像其实就是个readonly的根文件系统，从基础镜像构建的镜像其实都只是把修改部分和基础镜像合并重新打包，我们从ubuntu镜像构建一个具有golang环境的镜像，用来说明overlay在docker中的应用。
Dockerfile如下：
```Dockerfile
## golang dockerfile
## version 1.0
## 以ubuntu为基础镜像构建新镜像
FROM ubuntu:latest

## 维护者
MAINTAINER "joker"

## 新增用户go
RUN ["useradd", "-m", "-s", "/bin/bash", "go"]

## 指定工作目录
WORKDIR "/home/go"

## 把golang源码包打包到镜像，并解压到/opt路径
ADD go1.16.6.linux-amd64.tar.gz /opt

## 设置GOROOT/GOPATH/GOPROXY环境变量，并把对应的bin目录加到系统PATH环境变量
ENV GOROOT="/opt/go" GOPATH="/home/go/golang" GOPROXY="https://goproxy.cn,direct" GO111MODULE="on"
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

## 以下命令以用户go执行
USER go

## docker run执行的命令
ENTRYPOINT ["/bin/bash"]

## 挂载/home/go到一个匿名路径，后续能看到具体是哪个路径，可以通过-v覆盖匿名路径
VOLUME ["/home/go"]

## 添加一些元数据
LABEL author="joker" email="cyp_fly@126.com"

```

构建镜像

```bash
$$ docker build -t ubuntu:golang .
Sending build context to Docker daemon  129.1MB         
Step 1/11 : FROM ubuntu:latest                        
 ---> ba6acccedd29                                             
Step 2/11 : MAINTAINER "joker"                                                    
 ---> Running in 1906ecf1e073                
Removing intermediate container 1906ecf1e073
 ---> b10d3bffdec9                                                                          
Step 3/11 : RUN ["useradd", "-m", "-s", "/bin/bash", "go"]
 ---> Running in f5fc6f4d31c2                               
Removing intermediate container f5fc6f4d31c2                  
 ---> 8985b91b0827                                                     
Step 4/11 : WORKDIR "/home/go"                 
 ---> Running in f4f4b66d61fe                                                                   
Removing intermediate container f4f4b66d61fe      
 ---> a6e334f8073f                                          
Step 5/11 : ADD go1.16.6.linux-amd64.tar.gz /opt          
 ---> 142a1460d16c                          
Step 6/11 : ENV GOROOT="/opt/go" GOPATH="/home/go/golang" GOPROXY="https://goproxy.cn,direct" GO111MODULE="on"
 ---> Running in a0b88f65dc77                         
Removing intermediate container a0b88f65dc77
 ---> 331d4c34ae5f                                         
Step 7/11 : ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH               
 ---> Running in 3909af77862a                                                                                  
Removing intermediate container 3909af77862a                                                                   
 ---> 55035165f7c3                                                                          
Step 8/11 : USER go                                           
 ---> Running in 0840158531d4                           
Removing intermediate container 0840158531d4
 ---> 68459adca191                                     
Step 9/11 : ENTRYPOINT ["/bin/bash"]           
 ---> Running in fe9897b8d155                            
Removing intermediate container fe9897b8d155   
 ---> d17e20160cfd                                       
Step 10/11 : VOLUME ["/home/go"]                    
 ---> Running in dce2482115c0                     
Removing intermediate container dce2482115c0                                                
 ---> 7252f0757d57                                     
Step 11/11 : LABEL author="joker" email="cyp_fly@126.com"                            
 ---> Running in c0a3975854ba                                                                                       
Removing intermediate container c0a3975854ba                                                    
 ---> d4e4da229727                                           
Successfully built d4e4da229727                                               
Successfully tagged ubuntu:golang 

```

启动容器

```bash
$$ docker run -itd --name golang -u go -v /var/golang:/home/go ubuntu:golang
```

查看容器

```
$$ docker inspect golang
```

查看GraphDriver
```
        "GraphDriver": {
            "Data": {
                "LowerDir": "/var/lib/docker/overlay2/144b8379bf7b144ab9ec4d8de712b8ec56d87fd0f5d74b2c938c2b215e1b9ebd-init/diff:/var/li
b/docker/overlay2/1d65e86e54373a5b01afe28d0878a953fa2e9eb7cd14552a17d4a8628b476978/diff:/var/lib/docker/overlay2/b56e1dd54abdf8300b2f972
6e3e92d3b32b78c8bf5b5cd807e27be78c671af40/diff:/var/lib/docker/overlay2/44c045232b8f3510d28965fbebcaad458263855785cc38e1a784a6731df7433d
/diff",
                "MergedDir": "/var/lib/docker/overlay2/144b8379bf7b144ab9ec4d8de712b8ec56d87fd0f5d74b2c938c2b215e1b9ebd/merged",
                "UpperDir": "/var/lib/docker/overlay2/144b8379bf7b144ab9ec4d8de712b8ec56d87fd0f5d74b2c938c2b215e1b9ebd/diff",
                "WorkDir": "/var/lib/docker/overlay2/144b8379bf7b144ab9ec4d8de712b8ec56d87fd0f5d74b2c938c2b215e1b9ebd/work"
            },
            "Name": "overlay2"
        },

```
从上面可以看到docker采用的是overlay2文件系统，LowerDir有多层。


查看每一层都有什么:

```bash
ls /var/lib/docker/overlay2/144b8379bf7b144ab9ec4d8de712b8ec56d87fd0f5d74b2c938c2b215e1b9ebd-init/diff
dev  etc
```

```bash
ls /var/lib/docker/overlay2/1d65e86e54373a5b01afe28d0878a953fa2e9eb7cd14552a17d4a8628b476978/diff
opt
```

```bash
ls /var/lib/docker/overlay2/b56e1dd54abdf8300b2f9726e3e92d3b32b78c8bf5b5cd807e27be78c671af40/diff
etc  home  var
```

```bash
ls /var/lib/docker/overlay2/44c045232b8f3510d28965fbebcaad458263855785cc38e1a784a6731df7433d/diff
bin  boot  dev  etc  home  lib  lib32  lib64  libx32  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
```

为什么会有这么多层，其实很好理解，查看Dockerfile中写了哪些规则:

1. 基于ubuntu基础镜像，所以最底层必然是ubuntu的镜像文件，待会可以查证
2. 创建一个用户，必然会修改/etc，新用户在/home下创建了工作目录，导致LowerDir增加一层
3. 上传了golang源码包并解压到/opt，并然导致/opt修改，导致LowerDir增加一层
4. docker把/etc自动添加一个init层，docker官方认为/etc下的修改一般会影响kernel，从而影响所有使用该镜像的用户，其他用户可能不希望有这些修改。
5. UpperDir初始化是empty的，用户有修改就会在UpperDir中创建


现在我们核对下最下层的镜像是否是ubuntu的镜像文件:

```bash
$$ docker inspect ubuntu
        "GraphDriver": {
            "Data": {
                "MergedDir": "/var/lib/docker/overlay2/44c045232b8f3510d28965fbebcaad458263855785cc38e1a784a6731df7433d/merged",
                "UpperDir": "/var/lib/docker/overlay2/44c045232b8f3510d28965fbebcaad458263855785cc38e1a784a6731df7433d/diff",
                "WorkDir": "/var/lib/docker/overlay2/44c045232b8f3510d28965fbebcaad458263855785cc38e1a784a6731df7433d/work"
            },
            "Name": "overlay2"
        },

```

ubuntu是基础镜像，没有LowerDir, UpperDir就是我新建ubuntu:golang镜像LowerDir的最底层文件系统
