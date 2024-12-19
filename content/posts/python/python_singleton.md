---
title: "python单例模式"
date: 2022-10-23
draft: false
tags : [                    # 文章所属标签
    "Python",
]
categories : [              # 文章所属标签
    "技术",
]
---

# 单例模式

## 简介

单例模式是指在内存中只会创建一次对象的设计模式，在程序中多次使用同一个对象且作用相同的时候，为了防止频繁的创建对象，单例模式可以让程序在内存中创建一个对象，让所有调用者都共享这一单例对象。


## 几种单例模式的demo

### 线程安全的单例模式

```py
import threading


class Singleton(object):
    _instance_lock = threading.Lock()
    
    def __init__(self):
        pass
    
    def __new__(cls, *args, **kwargs):
        if not hasattr(Singleton, "_instance"):
            with Singleton._instance_lock:
                if not hasattr(Singleton, "_instance"):
                    Singleton._instance = object.__new__(cls)
        return Singleton._instance

```

使用

```py
obj1 = Singleton()
obj2 = Singleton()

print(obj1,obj2)

# <__main__.Singleton object at 0x7fc6c6a57898> <__main__.Singleton object at 0x7fc6c6a57898>

```


多线程使用

```py
def task(arg):
    obj = Singleton()
    print(obj)
for i in range(10):
    t = threading.Thread(target=task,args=[i,])
    t.start()

# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>
# <__main__.Singleton object at 0x7fc6c6a57898>

```