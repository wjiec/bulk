package signal

import (
	"context"
	"os"
	"os/signal"
)

type Trigger interface {
	Do(context.Context, func(os.Signal))

	private()
}

type observer struct {
	once    bool
	signals []os.Signal
}

func (ob *observer) Do(ctx context.Context, f func(os.Signal)) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, ob.signals...)

	go func() {
		defer func() { signal.Stop(ch); close(ch) }()

		for {
			select {
			case sig := <-ch:
				f(sig)
				if ob.once {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (ob *observer) private() {}

func Once(signals ...os.Signal) Trigger {
	return &observer{
		once:    true,
		signals: signals,
	}
}

func When(signals ...os.Signal) Trigger {
	return &observer{
		once:    false,
		signals: signals,
	}
}

func WithContext(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	go Once(signals...).Do(ctx, func(os.Signal) {
		cancel()
	})
	return ctx, cancel
}
