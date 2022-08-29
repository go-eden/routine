# routine

[![Build Status](https://travis-ci.com/go-eden/routine.svg?branch=main)](https://travis-ci.com/github/go-eden/routine)
[![codecov](https://codecov.io/gh/go-eden/routine/branch/main/graph/badge.svg?token=R4GC2IuGoh)](https://codecov.io/gh/go-eden/routine)
[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/go-eden/routine)

> [中文版](README_zh.md)

`routine` encapsulates and provides some easy-to-use, high-performance `goroutine` context access interfaces, which can
help you access coroutine context information more elegantly, but you may also open Pandora's Box.

# Introduce

The `Golang` language has been sparing no effort to shield developers from the concept of coroutine context from the
beginning of its design, including the acquisition of coroutine `goid`, the state of the coroutine within the process,
and the storage of coroutine context.

If you have used other languages such as `C++/Java/...`, then you must be familiar with `ThreadLocal`, and after
starting to use `Golang`, you will definitely feel confused and distressed by the lack of convenient functions similar
to `ThreadLocal` . Of course, you can choose to use `Context`, let it carry all the context information, appear in the
first input parameter of all functions, and then shuttle around in your system.

The core goal of `routine` is to open up another path: to introduce `goroutine local storage` into the world of `Golang`
, and at the same time expose the coroutine information to meet the needs of some people.

# Usage & Demo

This chapter briefly introduces how to install and use the `routine` library.

## Install

```bash
go get github.com/go-eden/routine
```

## Use `goid`

The following code simply demonstrates the use of `routine.Goid()`:

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

In this example, the `main` function starts a new coroutine, so `Goid()` returns the main coroutine `1`:

```text
curr goid: 1
```

## Use `LocalStorage`

The following code simply demonstrates `NewLocalStorage()`, `Set()`, `Get()`, and cross-coroutine propagation
of `LocalStorage`:

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

	// other goroutine cannot read nameVar
	go func() {
		fmt.Println("name1: ", nameVar.Get())
	}()

	// but, the new goroutine could inherit/copy all local data from the current goroutine like this:
	routine.Go(func() {
		fmt.Println("name2: ", nameVar.Get())
	})

	// or, you could copy all local data manually
	ic := routine.BackupContext()
	go func() {
		routine.InheritContext(ic)
		fmt.Println("name3: ", nameVar.Get())
	}()

	time.Sleep(time.Second)
}
```

The results of the upper example are:

```text
name:  hello world
name1:  <nil>
name3:  hello world
name2:  hello world
```

# API

This chapter introduces in detail all the interfaces encapsulated by the `routine` library, as well as their core
functions and implementation methods.

## `Goid() (id int64)`

Get the `goid` of the current `goroutine`.

## `NewLocalStorage()`:

Create a new instance of `LocalStorage`, its design idea is very similar to the usage of `ThreadLocal` in other
languages.

## `BackupContext() *ImmutableContext`

Back up the `local storage` data of the current coroutine context. It is just an immutable structure that facilitates
the transfer of context data.

## `InheritContext(ic *ImmutableContext)`

Actively inherit the backed-up context `local storage` data, it will copy the data of other coroutines `BackupContext()`
into the current coroutine context, thus supporting the contextual data propagation across coroutines.

## `Go(f func())`

Start a new coroutine and automatically copy all the context `local storage` data of the current coroutine to the new
coroutine. Its internal implementation consists of `BackupContext()` and `InheritContext()`.

## `LocalStorage`

Represents the context variable of the coroutine, and the supported functions include:

+ `Get() (value interface{})`: Get the variable value that has been set by the current coroutine.
+ `Set(v interface{}) interface{}`: Set the value of the context variable of the current coroutine, and return the old
  value that has been set before.
+ `Del() (v interface{})`: Delete the context variable value of the current coroutine and return the deleted old value.

**Tip: The internal implementation of `Get/Set/Del` adopts a lock-free design. In most cases, its performance should be
very stable and efficient.**

# Garbage Collection

Before the `v1.0.0` version, `routine` will setup a backgrount timer to scan all go-routines intervally, and find the exited routine to clean the related `LocalStorage` data.

After the `v1.0.0` version, `routine` will register a `finalizer` to listen the lifecycle of `runtime.g`. 

After the coroutine exits, when runtime's GC running, the `finalizer` mechanism of `runtime` will actively remove the useless `LocalStorage` `Data clean up, So as to avoid memory leaks.

# Thanks

The internal model `internal/g` is from other repos, mainly two functions:

+ `G()`, from https://github.com/huandu/go-tls
+ `G0()`, from https://github.com/timandy/routine

# License

MIT
