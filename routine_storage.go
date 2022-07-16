package routine

import (
	"runtime"
	"sync"
)

// We should limit the number of LocalStorage, to avoid performance issue.
// How about 64K at most? nobody would ever need so many storages...
const keyCnt = 1 << 16

// record using key, to support allocating and releasing
var (
	keyMap   = map[int]bool{}
	keyLock  sync.Mutex
	keyIndex = 1000000
)

// allocate an unique key of storage.
func allocateKey() (key int, succ bool) {
	keyLock.Lock()
	defer keyLock.Unlock()
	if len(keyMap) >= keyCnt {
		return 0, false
	}
	for keyMap[keyIndex] {
		keyIndex++
	}
	keyMap[keyIndex] = true
	return keyIndex, true
}

// release the unique key of storage, it could be reused in future.
func releaseKey(id int) {
	keyLock.Lock()
	defer keyLock.Unlock()
	delete(keyMap, id)
}

func keyCount() int {
	keyLock.Lock()
	defer keyLock.Unlock()
	return len(keyMap)
}

type storage struct {
	key int
}

func newStorage() *storage {
	key, succ := allocateKey()
	if !succ {
		panic("Too many storages")
	}

	s := &storage{key: key}
	runtime.SetFinalizer(s, func(s *storage) {
		releaseKey(s.key)
	})
	return s
}

func (t *storage) Get() (v interface{}) {
	return loadStore().get(t.key)
}

func (t *storage) Set(v interface{}) (oldValue interface{}) {
	return loadStore().set(t.key, v)
}

func (t *storage) Del() (v interface{}) {
	return loadStore().del(t.key)
}
