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
	wg, cnt := sync.WaitGroup{}, 0
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	When(sigUsr1, sigUsr2).Do(ctx, func(signal os.Signal) {
		if cnt%2 == 0 {
			assert.Equal(t, sigUsr1, signal)
		} else {
			assert.Equal(t, sigUsr2, signal)
		}

		cnt++
		wg.Done()
	})

	wg.Add(1)
	if err := sendSignal(sigUsr2); err != nil {
		assert.FailNow(t, "send signal", err)
	}

	wg.Add(1)
	if err := sendSignal(sigUsr1); err != nil {
		assert.FailNow(t, "send signal", err)
	}

	wg.Wait()
	assert.Equal(t, 2, cnt)
}

func TestWithContext(t *testing.T) {
	ctx, cancel := WithContext(context.Background(), sigUsr1)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
			assert.Fail(t, "context uncanceled")
		}
	}()

	if err := sendSignal(sigUsr1); err != nil {
		assert.NoError(t, err)
	}
}
