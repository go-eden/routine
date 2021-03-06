package routine

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestG(t *testing.T) {
	gp := newGAccessor()
	assert.NotNil(t, gp)
	if gp == nil {
		t.Fatalf("fail to get G.")
	}

	t.Run("G in another goroutine", func(t *testing.T) {
		gp2 := GetG()
		assert.NotNil(t, gp2)
		assert.NotEqual(t, gp, gp2)
	})

	a := newGAccessor()
	assert.True(t, a.Labels() == nil)

	t.Log("goid: ", a.Goid())
	t.Log("labels: ", a.Labels())

	a.SetLabels(map[string]string{"msg": "hello world"})
	assert.True(t, a.Labels() != nil)
	t.Log("labels: ", a.Labels())
}

func TestLabelsGC(t *testing.T) {
	var wg sync.WaitGroup
	var flag int64

	a := newGAccessor()
	wg.Add(1)
	go func() {
		labels := map[string]string{"msg": "hello world"}
		a.SetLabels(labels)
		t.Log("labels init: ", labels)
		runtime.SetFinalizer(&labels, func(v interface{}) {
			t.Log("labels is finalized: ", v)
			atomic.AddInt64(&flag, 1)
		})
		t.Log("sub routine exit")
		wg.Done()
	}()
	wg.Wait()
	runtime.GC()
	nap()
	assert.True(t, atomic.LoadInt64(&flag) > 0, "labels should be finalized")
}

func TestGReusing(t *testing.T) {
	var wg sync.WaitGroup
	var lk sync.Mutex

	// After g was created, it should never be released or finalized.
	const hugeBatch = 1 << 14
	var allgs = make([]G, 0, hugeBatch)
	for i := 0; i < hugeBatch; i++ {
		if i%8000 == 0 {
			nap() // avoid too much concurrency
		}
		wg.Add(1)
		go func() {
			lk.Lock()
			ga := GetG()
			allgs = append(allgs, ga)
			lk.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	runtime.GC()
	nap()
	for i, ga := range allgs {
		assert.True(t, ga.Status() == GDead, "", i)
	}

	// The dead g should be reused, let's see
	gptrMap := map[uintptr]int{}
	wg = sync.WaitGroup{}
	for i := 0; i < hugeBatch; i++ {
		wg.Add(1)
		go func() {
			gp := newGAccessor()
			lk.Lock()
			gptrMap[uintptr(gp.ptr)]++
			lk.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	t.Logf("gotr count: %d", len(gptrMap))

	// protect allgs from gc
	t.Log(len(allgs))
}
