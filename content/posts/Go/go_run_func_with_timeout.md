---
title: "超时控制下执行函数"
date: 2023-01-30
draft: false
tags : [                    # 文章所属标签
    "Go",
]
---

go中实现超时控制下执行函数功能

```go
func RunWithTimeout(fun func() error, timeout time.Duration) error {
	finished := make(chan struct{})
	var err error
	go func() {
		err = fun()
		finished <- struct{}{}
	}()

	select {
	case <-finished:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	}
}
```