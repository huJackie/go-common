package http_test

import (
	"sync"
	"testing"
	"time"

	"github.com/cloudadrd/go-common/rate/limit"
)

func TestLimit(t *testing.T) {
	var (
		burst int32 = 10
		l           = limit.NewLimiter(time.Second, burst)
		wg    sync.WaitGroup
	)

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		ii := i
		go func() {
			defer wg.Done()
			current := l.Get()
			if current > burst {
				t.Fatalf("burst !!! real num [%d] want [%d]", current, burst)
			}
			if l.Increase() {
				t.Logf("get!!! 序号[%d]当前次数:%d\n", ii, l.Get())
			}
		}()
		if i%2 == 1 {
			time.Sleep(time.Microsecond)
		}
	}
	wg.Wait()

	time.Sleep(time.Second)
	b := l.Get()
	if b != 0 {
		t.Fatal("un release...", b)
	}
	t.Log("success")
}
