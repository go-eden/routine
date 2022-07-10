package routine

import (
	"fmt"
	"github.com/go-eden/routine/internal/g"
	"reflect"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var (
	goidOffset   uintptr
	labelsOffset uintptr
	statusOffset uintptr
)

func init() {
	offset := func(t reflect.Type, f string) uintptr {
		if field, found := t.FieldByName(f); found {
			return field.Offset
		}
		panic(fmt.Sprintf("init routine failed, cannot find g.%s, version=%s", f, runtime.Version()))
	}
	gt := reflect.TypeOf(g.G0())
	goidOffset = offset(gt, "goid")
	labelsOffset = offset(gt, "labels")
	statusOffset = offset(gt, "atomicstatus")
}

type gAccessor struct {
	ptr unsafe.Pointer
	gid int64
}

func newGAccessor() *gAccessor {
	return &gAccessor{
		ptr: g.G(),
		gid: Goid(),
	}
}

func (g *gAccessor) Goid() int64 {
	goidPtr := (*int64)(unsafe.Pointer(uintptr(g.ptr) + goidOffset))
	return atomic.LoadInt64(goidPtr)
}

func (g *gAccessor) Labels() map[string]string {
	labelsPPtr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(g.ptr) + labelsOffset))
	labelsPtr := atomic.LoadPointer(labelsPPtr)
	if labelsPtr == nil {
		return nil
	}
	// see SetGoroutineLabels, labelsPtr is `*labelMap`
	return *(*map[string]string)(labelsPtr)
}

func (g *gAccessor) Status() GStatus {
	if g.gid != g.Goid() {
		return GDead
	}
	statusPtr := (*uint32)(unsafe.Pointer(uintptr(g.ptr) + statusOffset))
	return GStatus(atomic.LoadUint32(statusPtr))
}

func (g *gAccessor) SetLabels(labels map[string]string) {
	labelsPtr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(g.ptr) + labelsOffset))
	atomic.StorePointer(labelsPtr, unsafe.Pointer(&labels))
}
