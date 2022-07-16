# routine

[![Build Status](https://travis-ci.com/go-eden/routine.svg?branch=main)](https://travis-ci.com/github/go-eden/routine)
[![codecov](https://codecov.io/gh/go-eden/routine/branch/main/graph/badge.svg?token=R4GC2IuGoh)](https://codecov.io/gh/go-eden/routine)
[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/go-eden/routine)

> [English Version](README_zh.md)

`routine`封装并提供了一些易用、高性能的`goroutine`上下文访问接口，它可以帮助你更优雅地访问协程上下文信息，但你也可能就此打开了潘多拉魔盒。

# 介绍

`Golang`语言从设计之初，就一直在不遗余力地向开发者屏蔽协程上下文的概念，包括协程`goid`的获取、进程内部协程状态、协程上下文存储等。

如果你使用过其他语言如`C++/Java`等，那么你一定很熟悉`ThreadLocal`，而在开始使用`Golang`之后，你一定会为缺少类似`ThreadLocal`的便捷功能而深感困惑与苦恼。 当然你可以选择使用`Context`
，让它携带着全部上下文信息，在所有函数的第一个输入参数中出现，然后在你的系统中到处穿梭。

而`routine`的核心目标就是开辟另一条路：将`goroutine local storage`引入`Golang`世界，同时也将协程信息暴露出来，以满足某些人可能有的需求。

# 使用演示

此章节简要介绍如何安装与使用`routine`库。

## 安装

```bash
go get github.com/go-eden/routine
```

## 使用`goid`

以下代码简单演示了`routine.Goid()`的使用：

```go
package main

import (
	"fmt"
	"github.com/go-eden/routine"
	"time"
)

func main() {
	go func() {
		time.Sleep(time.Second)
	}()
	goid := routine.Goid()
	fmt.Printf("curr goid: %d\n", goid)
}
```

此例中`main`函数启动了一个新的协程，因此`Goid()`返回了主协程`1`:

```text
curr goid: 1
```

## 使用`LocalStorage`

以下代码简单演示了`LocalStorage`的创建、设置、获取、跨协程传播等：

```go
package main

import (
	"fmt"
	"github.com/go-eden/routine"
	"time"
)

var nameVar = routine.NewLocalStorage()

func main() {
	nameVar.Set("hello world")
	fmt.Println("name: ", nameVar.Get())

	// 其他协程不能读取前面Set的"hello world"
	go func() {
		fmt.Println("name1: ", nameVar.Get())
	}()

	// 但是可以通过Go函数启动新协程，并将当前main协程的全部协程上下文变量赋值过去
	routine.Go(func() {
		fmt.Println("name2: ", nameVar.Get())
	})

	// 或者，你也可以手动copy当前协程上下文至新协程，Go()函数的内部实现也是如此
	ic := routine.BackupContext()
	go func() {
		routine.InheritContext(ic)
		fmt.Println("name3: ", nameVar.Get())
	}()

	time.Sleep(time.Second)
}
```

执行结果为：

```text
name:  hello world
name1:  <nil>
name3:  hello world
name2:  hello world
```

# API文档

此章节详细介绍了`routine`库封装的全部接口，以及它们的核心功能、实现方式等。

## `Goid() (id int64)`

获取当前`goroutine`的`goid`，`Goid()`通过`go_tls`的方式直接获取，此操作性能极高，耗时通常只相当于`rand.Int()`的五分之一。

## `NewLocalStorage()`:

创建一个新的`LocalStorage`实例，它的设计思路与用法和其他语言中的`ThreadLocal`非常相似。

## `BackupContext() *ImmutableContext`

备份当前协程上下文的`local storage`数据，它只是一个便于上下文数据传递的不可变结构体。

## `InheritContext(ic *ImmutableContext)`

主动继承备份到的上下文`local storage`数据，它会将其他协程`BackupContext()`的数据复制入当前协程上下文中，从而支持**跨协程的上下文数据传播**。

## `Go(f func())`

启动一个新的协程，同时自动将当前协程的全部上下文`local storage`数据复制至新协程，它的内部实现由`BackupContext()`和`InheritContext()`组成。

## `LocalStorage`

表示协程上下文变量，支持的函数包括：

+ `Get() (value interface{})`：获取当前协程已设置的变量值，若未设置则为`nil`
+ `Set(v interface{}) interface{}`：设置当前协程的上下文变量值，返回之前已设置的旧值
+ `Del() (v interface{})`：删除当前协程的上下文变量值，返回已删除的旧值

**提示：`Get/Set/Del`的内部实现采用无锁设计，在大部分情况下，它的性能表现都应该非常稳定且高效。**

# 垃圾回收

在`v1`版本中，`routine`会通过一个后台定时任务，通过轮询的方式扫描已退出的协程，并主动清理掉相关的`LocalStorage`数据。

在`v2`版本中，`routine`会主动监听`runtime.g`的生命周期，在协程退出后，系统执行垃圾回收时，通过`runtime`的`finalizer`机制，主动将无用的`LocalStorage`数据清理掉，
从而避免内存的泄露。

# License

MIT