package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoid(t *testing.T) {
	t.Log(Goid())
}

func TestGoStorage(t *testing.T) {
	var variable = "hello world"
	stg := NewLocalStorage()
	stg.Set(variable)
	Go(func() {
		v := stg.Get()
		assert.True(t, v != nil && v.(string) == variable)
	})
	time.Sleep(time.Millisecond)
}

// BenchmarkGoid-12    	1000000000	         1.036 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGoid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Goid()
	}
}

// BenchmarkCurr-12    	45119158	        26.00 ns/op	      16 B/op	       1 allocs/op
func BenchmarkCurr(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetG()
	}
}
