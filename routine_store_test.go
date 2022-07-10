package routine

import (
	"context"
	"github.com/stretchr/testify/assert"
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

const key = 1000

func TestStoreBase(t *testing.T) {
	s := loadStore()

	for i := 0; i < 100; i++ {
		str := "hello"
		s.set(key, str)
		p := s.get(key)
		assert.True(t, p.(string) == str)
	}

	v := s.del(key)
	assert.True(t, v != nil)

	v = s.get(key)
	assert.True(t, v == nil)

	// other routine cannot access s
	go func() {
		assert.Panics(t, func() {
			s.set(key, "test")
		})
		assert.Panics(t, func() {
			s.get(key)
		})
		assert.Panics(t, func() {
			s.del(key)
		})
	}()
	time.Sleep(time.Millisecond * 10)
}

// after routine dead and gc, store should be clean up
func TestStoreFinalize(t *testing.T) {
	var wg sync.WaitGroup
	var ss []*store

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			s := loadStore()
			s.set(key, "test")
			ss = append(ss, s)
			wg.Done()
		}()
	}
	wg.Wait()
	runtime.GC()
	time.Sleep(time.Millisecond)

	for _, s := range ss {
		assert.True(t, s.values == nil, "s.values is nil")
		assert.True(t, s.g == nil, "s.g is nil")
	}

	// after routine exit, all store should be clean up.
	const round = 10
	const concurrency = 1000
	for i := 0; i < round; i++ {
		var wg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				loadStore().set(key, "test")
				wg.Done()
			}()
		}
		wg.Wait()
		runtime.GC()
		time.Sleep(time.Millisecond)

		ss = allStore()
		assert.True(t, len(ss) == 0)
	}

	// after labels was occupied by others (pprof), we should register finalizer again
	s := loadStore()
	s.set(key, "test")

	// mock pprof, check gc status
	time.Sleep(time.Millisecond)
	pprof.SetGoroutineLabels(context.Background())
	runtime.GC()
	assert.True(t, s.values != nil, "s.values isn't nil")
	assert.True(t, s.fcnt == 1, "", s.fcnt)

	// check allStoreMap
	ss = allStore()
	assert.True(t, len(ss) == 1)
}

func TestStoreOccupyLabels(t *testing.T) {
	s := loadStore()

	labels := map[string]string{
		"name": "sulin",
	}
	s.g.SetLabels(labels)
	runtime.GC() // store should occupy pointer back
	time.Sleep(time.Millisecond * 10)
	assert.True(t, s.fcnt == 1)
	assert.True(t, s.g.Labels()["name"] == "sulin")
}

// BenchmarkLoadStore-12    	52496490	        20.48 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLoadStore(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = loadStore()
	}
}

// BenchmarkStoreGet-12    	424059255	         2.943 ns/op	       0 B/op	       0 allocs/op
func BenchmarkStoreGet(b *testing.B) {
	s := loadStore()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.get(key)
	}
}

// BenchmarkStoreSet-12    	82243570	        13.17 ns/op	       0 B/op	       0 allocs/op
func BenchmarkStoreSet(b *testing.B) {
	s := loadStore()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.set(key, key)
	}
}
