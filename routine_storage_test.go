package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

func TestStorage(t *testing.T) {
	s := newStorage[interface{}]()

	for i := 0; i < 1000; i++ {
		str := "hello"
		s.Set(str)
		p := s.Get()
		assert.True(t, p.(string) == str)
		tmp := s.Del()
		assert.True(t, tmp == str)
		tmp = s.Del()
		assert.True(t, tmp == nil)
	}
	assert.True(t, s.Get() == nil)
}

func TestStorageConcurrency(t *testing.T) {
	const concurrency = 100
	const loopTimes = 100000
	var wg sync.WaitGroup

	s := newStorage[interface{}]()

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			v := rand.Uint64()
			for i := 0; i < loopTimes; i++ {
				s.Set(v)
				tmp := s.Get()
				assert.True(t, tmp.(uint64) == v)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// BenchmarkStorage-12    	 8102013	       133.6 ns/op	      16 B/op	       1 allocs/op
// BenchmarkStorage-10    	12186906	        87.35 ns/op	      16 B/op	       1 allocs/op
func BenchmarkStorage(b *testing.B) {
	s := newStorage[interface{}]()
	var variable = "hello world"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Get()
		s.Set(variable)
		s.Del()
	}
}
