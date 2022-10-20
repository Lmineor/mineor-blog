---
title: "Go语言的%d,%p,%v等占位符的使用"
date: 2022-10-19
draft: false
tags : [                    # 文章所属标签
    "占位符",
]
categories : [              # 文章所属标签
    "go",
]

---

# Go语言的%d,%p,%v等占位符的使用


## 1.占位符分别代表了什么？

golang 的fmt 包实现了格式化I/O函数，类似于C的 printf 和 scanf。
定义示例类型和变量
```go
type Human struct {
	Name string
} 

var people = Human{Name:"zhangsan"}
```

## 2 普通占位符
|占位符|说明|举例|输出|
|:-:|:-:|:-:|:-:|
|%v|相应值的默认格式。|Printf("%v", people)|{zhangsan}|
|%+v|打印结构体时，会添加字段名|Printf("%+v", people)|{Name:zhangsan}|
|%#v|相应值的Go语法表示|Printf("#v", people)|main.Human{Name:"zhangsan"}|
|%T|相应值的类型的Go语法表示|Printf("%T", people)|main.Human|	
|%%|字面上的百分号，并非值的占位符|Printf("%%")|%|

## 3 布尔占位符

|占位符|说明|举例|输出|
|:-:|:-:|:-:|:-:|
|%t|true或false|Printf("%t", true)|true|

## 4 整数占位符
|占位符|说明|举例|输出|
|:-:|:-:|:-:|:-:|
|%b|二进制表示|Printf("%b", 5)|101|
|%c|相应Unicode码点锁表示的字符|Printf("%c", 0x4e2d)|中|
|%d|十进制表示|Printf("%d", 0x12)|18|

