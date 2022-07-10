package routine

import (
	"runtime"
	"sync"
)

var allStoreMap sync.Map

// store The underlying datastore of a go-routine.
type store struct {
	g      *gAccessor
	gid    int64
	values map[int]interface{} // it should never be accessed in other routine, unless the go-routine was dead.

	fcnt int8 // for test
}

func (s *store) get(key int) interface{} {
	if Goid() != s.gid {
		panic("Cannot access store in other routine.")
	}
	return s.values[key]
}

func (s *store) set(key int, value interface{}) interface{} {
	if Goid() != s.gid {
		panic("Cannot access store in other routine.")
	}
	oldValue := s.values[key]
	s.values[key] = value
	return oldValue
}

func (s *store) del(key int) interface{} {
	if Goid() != s.gid {
		panic("Cannot access store in other routine.")
	}
	oldValue := s.values[key]
	delete(s.values, key)
	return oldValue
}

// register Register finalizer into goroutine's lifecycle
func (s *store) register() {
	labels := make(map[string]string)
	for k, v := range s.g.Labels() {
		labels[k] = v
	}
	runtime.SetFinalizer(&labels, s.finalize)
	s.g.SetLabels(labels)
}

// finalize Do something when the current g may be dead
func (s *store) finalize(_ interface{}) {
	s.fcnt++
	if s.g == nil {
		return
	}

	// Maybe others (pprof) replaced our labels, register it again.
	if s.g.Status() != GDead {
		go s.register()
		return
	}

	// do final
	s.g = nil
	s.values = nil
	allStoreMap.Delete(s.gid)
}

// loadStore load store of the current goroutine.
func loadStore() (s *store) {
	gid := Goid()

	// reuse the existed store
	val, ok := allStoreMap.Load(gid)
	if ok {
		return val.(*store)
	}

	// create new store, and register it into routine
	s = &store{
		g:      newGAccessor(),
		gid:    gid,
		values: map[int]interface{}{},
	}
	s.register()

	// register the new store
	allStoreMap.Store(gid, s)

	return s
}

// allStore collect all store, only for test.
func allStore() []*store {
	var ss []*store
	allStoreMap.Range(func(_, v interface{}) bool {
		ss = append(ss, v.(*store))
		return true
	})
	return ss
}
