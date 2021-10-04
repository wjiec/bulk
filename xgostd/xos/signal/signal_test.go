package signal

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnce(t *testing.T) {
	wg, cnt := sync.WaitGroup{}, 0

	Once(sigUsr1, sigUsr2).Do(context.Background(), func(signal os.Signal) {
		cnt++
		assert.Equal(t, sigUsr1, signal)
		wg.Done()
	})

	wg.Add(1)
	if err := sendSignal(sigUsr1); err != nil {
		assert.FailNow(t, "send signal", err)
	}

	wg.Wait()
	if err := sendSignal(sigUsr1); err != nil {
		assert.FailNow(t, "send signal", err)
	}
	assert.Equal(t, 1, cnt)
}

func TestWhen(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	When(sigUsr1, sigUsr2).Do(ctx, func(signal os.Signal) {
		assert.Contains(t, []os.Signal{sigUsr1, sigUsr2}, signal)
		wg.Done()
	})

	wg.Add(1)
	if err := sendSignal(sigUsr1); err != nil {
		assert.FailNow(t, "send signal", err)
	}

	wg.Add(1)
	if err := sendSignal(sigUsr2); err != nil {
		assert.FailNow(t, "send signal", err)
	}

	wg.Wait()
}
