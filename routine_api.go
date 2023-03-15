package routine

import (
	"github.com/go-eden/routine/internal/g"
	"unsafe"
)

// GStatus represents the real status of runtime.g
type GStatus uint32

const (
	GIdle      GStatus = 0 // see runtime._Gidle
	GRunnable  GStatus = 1 // see runtime._Grunnable
	GRunning   GStatus = 2 // see runtime._Grunning
	GSyscall   GStatus = 3 // see runtime._Gsyscall
	GWaiting   GStatus = 4 // see runtime._Gwaiting
	GMoribund  GStatus = 5 // see runtime._Gmoribund_unused
	GDead      GStatus = 6 // see runtime._Gdead
	GEnqueue   GStatus = 7 // see runtime._Genqueue_unused
	GCopystack GStatus = 8 // see runtime._Gcopystack
	GPreempted GStatus = 9 // see runtime._Gpreempted
)

// G defines some usefull api to access the underlying go-routine.
type G interface {

	// Goid returns g.goid
	Goid() int64

	// Labels returns g.labels, which's real type is map[string]string
	Labels() map[string]string

	// Status returns g.atomicstatus
	Status() GStatus
}

// LocalStorage provides goroutine-local variables.
type LocalStorage[T interface{}] interface {

	// Get returns the value in the current goroutine's local storage, if it was set before.
	Get() (value T)

	// Set copy the value into the current goroutine's local storage, and return the old value.
	Set(value T) (oldValue T)

	// Del delete the value from the current goroutine's local storage, and return it.
	Del() (oldValue T)
}

// ImmutableContext represents all local allStoreMap of one goroutine.
type ImmutableContext struct {
	gid    int64
	values map[int64]interface{}
}

// Go start a new goroutine, and copy all local allStoreMap from current goroutine.
func Go(f func()) {
	ic := BackupContext()
	go func() {
		InheritContext(ic)
		f()
	}()
}

// BackupContext copy all local allStoreMap into an ImmutableContext instance.
func BackupContext() *ImmutableContext {
	s := loadStore()
	data := map[int64]interface{}{}
	for k, v := range s.values {
		data[k] = v
	}
	return &ImmutableContext{gid: s.gid, values: data}
}

// InheritContext load the specified ImmutableContext instance into the local storage of current goroutine.
func InheritContext(ic *ImmutableContext) {
	if ic == nil || ic.values == nil {
		return
	}
	s := loadStore()
	for k, v := range ic.values {
		s.values[k] = v
	}
}

// NewLocalStorage create and return a new LocalStorage instance.
func NewLocalStorage[T interface{}]() LocalStorage[T] {
	return newStorage[T]()
}

// Goid get the unique goid of the current routine.
func Goid() int64 {
	return *(*int64)(unsafe.Pointer(uintptr(g.G()) + goidOffset))
}

// GetG returns the accessor of the current routine.
func GetG() G {
	return newGAccessor()
}
