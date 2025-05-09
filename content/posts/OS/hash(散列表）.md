---
title: hash(散列表）
date: 2025-04-16
draft: 
tags:
  - 软件
---
# 散列表

##  散列表的基本概念
  

在线性表和树表的查找中，记录在表中的位置和记录的关键字之间不存在确定关系，因此，在这些表中查找记录时需进行一系列的关键字比较。这类查找方法建立在“比较”的基础上，查找的效率取决于比较的次数

散列函数：一个把查找表中的关键字映射成该关键字对应的地址的函数，记为Hash(key) = Addr(Addr可以是数组下标、索引或内存地址)

散列函数可能会把两个或两个以上的不同关键字映射到同一地址，称这种情况为**冲突**，这些发生碰撞的不同关键字称为**同义词**。

一方面，设计得好的散列函数应尽量减少这样的冲突；
一方面，由于这样的冲突总是不可避免地，所以还要设计好处理冲突的方法。

**散列表**：根据关键字而直接进行访问的数据结构。散列表建立了关键字和存储地址之间的一种直接映射关系。

理想情况下，对散列表进行查找的时间复杂度为O(1)，即与表中元素的个数无关。

##  散列函数的构造方法

在构造散列函数时，需注意：
1. 散列函数的定义域必须包含全部需要存储的关键字，而值域的范围则依赖于散列表的大小或地址范围
2. 散列函数计算出来的地址应该能等概率、均匀地分布在整个地址空间中，从而减少冲突的发生。
3. 散列函数应尽量简单，能够在较短的时间内计算出任一关键字对应的散列地址。
### 1、 直接定址法

散列函数：`H(key)=a*key+b`
a，b是常数
结论：计算最简单，且不会产生冲突。适合关键字的分布基本连续的情况，若关键字分布不连续，空位较多，则会造成存储空间的浪费。
### 2、除留余数法

散列函数：
这是一种最简单、最常用的方法：假定散列表表长m，取一个不大于m但接近或等于m的质数p，
`H(key)=key%p`
关键是选好p，使得每个关键字通过该函数转换后等概率地映射到散列空间上的任一地址，从而尽可能减少冲突的可能性。
### 3、数字分析法

设关键字是r进制数（如十进制数），而r个数码在各位上出现的频率不一定相同，可能在某些位上分布均匀一些，每种数码出现的机会均等；而在某些位上分布不均匀，只有某几种数码经常出现，此时应选取数码分布较为均匀的若干位作为散列地址，此时应选取数码分布较为均匀的若干位作为散列地址。
适用：散列函数适合于已知的关键字集合，若更换了关键字，则需要重新构造新的散列函数。
### 4、平方取中法

取关键字的平方值的中间几位作为散列地址。
适用：散列函数适合于关键字的每位取值都不够均匀或均小于散列地址所需的位数。
### 5、折叠法

将关键字分割成位数相同的几部分（最后一部分位数可以短一些），然后取这几部分的叠加作为散列地址，这种方法称为折叠法。
适用：关键字位数很多，而且关键字中的每位上数字分布大致均匀时，可以采用折叠法得到散列地址。
在实际选择中，采用何种构造散列函数的方法取决于关键字集合的情况，但目标是为了尽量降低产生冲突的可能性。
## 处理冲突的方法

处理的方法：为产生冲突的关键字寻找下一个"空"的Hash地址。
假设已经选定散列函数H(key)，用Hi表示发生冲突后第i次探测的散列地址。

### 1、开放定址法

指：可存放新表项的空闲地址既向它的同义词表项开放，又向它的非同义词表项开放。
数学递推公式为：`Hi = (H(key) + di) % m`
式中，`i=0,1,2,...k(k <= m-1)`；`m`表示散列表表长；`di`为增量序列。
删除时只能进行逻辑删除（做删除标记），不能物理删除，副作用：执行多次删除后，表面上看起来散列表很满，实际上有很多位置未利用，因此需要定期维护散列表，要把删除标记的元素物理删除。

有一下4种取法：

1）线性探测法：冲突发生时，顺序查看表中下一个单元，直到找出一个空闲单元或查遍全表。产生"堆积问题"，降低查找效率。
2）平方探测法（二次探测法）：可以避免出现"堆积问题"，缺点是：不能探测到散列表上的所有单元，但至少能探测一半单元。
3）再散列法（双散列法）：
4）伪随机法
### 2、拉链法（链接法，chaining）

拉链法适用于经常进行插入和删除的情况。
## 散列查找及性能分析


散列表查找效率取决于三个因素：**散列函数**、**处理冲突的方法**和**装填因子**。

装填因子：记为：α。定义为一个表的装满程度。即α=表中记录数n/散列表长度m。

散列表的平均查找长度依赖于散列表的装填因子α，而不直接依赖于n或m。直观的看，α越大，表示装填的记录越"满"，发生冲突的可能性就越大，反之发生冲突的可能性越小。