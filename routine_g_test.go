package routine

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
	"time"
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
	var flag bool

	a := newGAccessor()
	wg.Add(1)
	go func() {
		labels := map[string]string{"msg": "hello world"}
		a.SetLabels(labels)
		t.Log("labels init: ", labels)
		runtime.SetFinalizer(&labels, func(v interface{}) {
			t.Log("labels is finalized: ", v)
			flag = true
		})
		t.Log("sub routine exit")
		wg.Done()
	}()
	wg.Wait()
	runtime.GC()
	time.Sleep(time.Millisecond * 10)
	assert.True(t, flag, "labels should be finalized")
}

func TestGReusing(t *testing.T) {
	var wg sync.WaitGroup
	var lk sync.Mutex

	// After g was created, it should never be released or finalized.
	const hugeBatch = 1 << 20
	var allgs = make([]G, 0, hugeBatch)
	for i := 0; i < hugeBatch; i++ {
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
	time.Sleep(time.Millisecond * 100)
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
