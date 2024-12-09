---
title: "修改rpm包"
date: 2023-06-18
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

基本思路：
rpm包没有修改工具，只能是把rpm包解压、修改（包括增删）其中的文件，然后重新制作rpm包。

注：制作rpm包，需要原rpm包的spec文件。

所需工具：
rpmrebuild 
它主要是用来提取原rpm包中的spec文件。

rpmrebuild有如下两种安装方式（建议第1中）：

 1）下载安装rpmrebuild rpm包

下载地址：

http://rpmfind.net/linux/rpm2html/search.php?query=rpmrebuild

rpm -ivh rpmrebuild-2.11-3.el7.noarch.rpm  安装后，可直接使用 rpmrebuild 命令。

2）下载解压tar包：

下载地址：https://jaist.dl.sourceforge.net/project/rpmrebuild/rpmrebuild/2.15/rpmrebuild-2.15.tar.gz

解压后，使用./rpmrebuild.sh 脚本，和 rpmrebuild命令一样。

rpmbuild
它主要是用来制作rpm包的。

步骤：
0）创建一个临时目录

mkdir -p /root/test_rpm_dir

cp mlnx-ofa_kernel-5.2-OFED.5.2.1.0.4.1.rhel7u3.x86_64.rpm /root/test_rpm_dir

cd  /root/test_rpm_dir

1）解压原rpm包

rpm2cpio mlnx-ofa_kernel-5.2-OFED.5.2.1.0.4.1.rhel7u3.x86_64.rpm | cpio -div

2）修改内容 

按自己需求修改内容，或增删文件

3）提取原rpm包spec文件 

#注：下面命令中的 --spec-only=test.spec 中的test.spec是要保存spec的文件路径（即把提取的spec文件保存为当前路径下的test.spec）。

有两种提取spec文件的方式：

a）从指定的rpm包文件中提取 

#注：-p, 即 --package 就是指使用rpm包文件，而不是系统中已安装的rpm。

#注：-n, 即 --notest-install，不要执行一个测试性的安装（do not perform a test install）。

#注：-s，即 --spec-only=<specfile> ，指只提取创建spec文件（generate specfile only）。

rpmrebuild  -p -n -s test.spec mlnx-ofa_kernel-5.2-OFED.5.2.1.0.4.1.rhel7u3.x86_64.rpm 

或

rpmrebuild  --package --notest-install --spec-only=test.spec mlnx-ofa_kernel-5.2-OFED.5.2.1.0.4.1.rhel7u3.x86_64.rpm 

a）从系统中安装的rpm中提取 

如果rpm包已经安装到系统中（且你已经从本地删除了该rpm包源文件），可以执行如下命令提取：

rpmrebuild -s test.spec -n mlnx-ofa_kernel-5.6-OFED.5.6.2.0.9.1.rhel7u4.x86_64

4）修改spec文件

如果有增删的文件，则需要在spec文件中体现

5）重新制作rpm

使用rpmbuild  通过指定--buildroot  和 提取的spec 重新制作包。

rpmbuild -ba  --buildroot /root/test_rpm_dir  test.spec

注：如果是已经删除了rpm包文件，只能从系统中已安装的rpm路径下打包对应文件来制作新的rpm包，此时应去掉 --buildroot /root/test_rpm_dir
————————————————
版权声明：本文为CSDN博主「billbonaparte1」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/billbonaparte1/article/details/125998193

