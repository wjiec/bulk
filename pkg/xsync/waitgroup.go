package xsync

import "sync"

type WaitGroup struct {
	sync.WaitGroup
}

func (wg *WaitGroup) Do(fn func()) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}
