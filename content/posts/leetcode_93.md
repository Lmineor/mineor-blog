---
title: "复原 IP 地址"
date: 2022-12-05
draft: true
tags : [                    # 文章所属标签
    "Go", "Leetcode", "回溯"
]
categories : [              # 文章所属标签
    "技术",
]
---

有效 IP 地址 正好由四个整数（每个整数位于 0 到 255 之间组成，且不能含有前导 0），整数之间用 '.' 分隔。

例如："0.1.2.201" 和 "192.168.1.1" 是 有效 IP 地址，但是 "0.011.255.245"、"192.168.1.312" 和 "192.168@1.1" 是 无效 IP 地址。
给定一个只包含数字的字符串 s ，用以表示一个 IP 地址，返回所有可能的有效 IP 地址，这些地址可以通过在 s 中插入 '.' 来形成。你 不能 重新排序或删除 s 中的任何数字。你可以按 任何 顺序返回答案。

原题链接[复原 IP 地址](https://leetcode.cn/problems/restore-ip-addresses/description/)

```go
var res []string

func restoreIpAddresses(s string) []string {
    res = make([]string, 0)
    path := make([]string, 0)
    tracingBack(s, 0, path)
    return res
}


func tracingBack(s string, startIndex int, path []string){
    if startIndex >= len(s) && len(path) == 4{
        // 遍历结束且候选长度为4，则为一种解法
        res = append(res, strings.Join(path, ".")) 
        return
    }
    for index:= startIndex+1; index < len(s)+1; index++{
        if isValidIpAddress(s[startIndex:index]){
            // 该子序列为有效的序列，则加入候选进行遍历
            path = append(path, s[startIndex:index])
            tracingBack(s, index, path)
        }
        if len(path) > 0{
            // case 1: 本次遍历的子序列为无效子序列，需将上一次分割结果抛弃
            // case 2: 上一次分割结束，已经产生一种可能结果，需将上一次分割结果抛弃，开始回溯下一种解法

            // 不管上次结果如何，需进行回溯，抛弃上一次的结果。
            path = path[0:len(path)-1]
        }
    }
}

func isValidIpAddress(s string)bool{
    a, _ := strconv.Atoi(s)
    if a <= 255 && strconv.Itoa(a) == s{
        return true
    }
    return false
}
```