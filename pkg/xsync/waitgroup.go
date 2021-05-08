package xsync

import "sync"

type WaitGroup struct {
	sync.WaitGroup
}

func (wg *WaitGroup) Do(fn func(args ...interface{}), args ...interface{}) {
	wg.Add(1)
	go func() {
		fn(args...)
		wg.Done()
	}()
}
