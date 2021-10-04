package signal

import (
	"context"
	"os"
	"os/signal"
)

// Observer is a system signal listener that can be used to
// listen to any number of signals and perform callback functions,
// either single or multiple times
type Observer struct {
	once    bool
	signals []os.Signal
}

// Do will create the listener and call the callback function when a signal occurs,
// and will decide whether to close the listener depending on the context
func (obs *Observer) Do(ctx context.Context, handle func(os.Signal)) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, obs.signals...)

	go func() {
		defer close(ch)
		defer signal.Stop(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case sig := <-ch:
				go handle(sig)
				if obs.once {
					return
				}
			}
		}
	}()
}

// Once returns an observer to wait for the specified signal to occur,
// will only trigger once
func Once(signals ...os.Signal) *Observer {
	return &Observer{
		once:    true,
		signals: signals,
	}
}

// When returns an observer to wait for the specified signal to occur,
// which will trigger multiple times
func When(signals ...os.Signal) *Observer {
	return &Observer{
		once:    false,
		signals: signals,
	}
}
