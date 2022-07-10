// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

// Package g exposes goroutine struct g to user space.
package g

import (
	"unsafe"
)

// getgp returns the pointer to the current runtime.g.
func getgp() unsafe.Pointer

// getg0 returns the value of runtime.g0.
func getg0() interface{}

// G returns current g (the goroutine struct) to user space.
//go:nosplit
func G() unsafe.Pointer {
	return getgp()
}

// G0 returns the g0, the main goroutine.
//go:nosplit
func G0() interface{} {
	return getg0()
}
