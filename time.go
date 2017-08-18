package gofaster

import (
	"sync/atomic"
	"time"
)

/***Variables***/
var rTime atomic.Value

/***Functions***/
func AproxTimeNow() time.Time {
	return *rTime.Load().(*time.Time)
}

/***Init***/
func init() {
	// get time
	t := time.Now().Truncate(time.Second)
	rTime.Store(&t)
	// get time (every second after)
	go func() {
		for {
			time.Sleep(time.Second)
			t := time.Now().Truncate(time.Second)
			rTime.Store(&t)
		}
	}()
}
