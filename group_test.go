package par

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestParGroup(t *testing.T) {
	var g Group

	var counter int32 = 0

	incrementCounter := func(interface{}) {
		atomic.AddInt32(&counter, 1)
	}

	incrementAfterDelay := func(interface{}) {
		time.Sleep(100 * time.Millisecond)
		g.Add(incrementCounter, nil)
	}

	g.Add(incrementCounter, nil)
	g.Add(incrementCounter, nil)
	g.Add(incrementAfterDelay, nil)

	g.Wait()

	if counter != 3 {
		t.Fail()
	}
}
