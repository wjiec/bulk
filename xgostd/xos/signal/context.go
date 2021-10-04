package signal

import (
	"context"
	"os"
)

// WithContext returns a copy of parent based on the specified signal and
// cancels the context when the signal occurs
func WithContext(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	Once(signals...).Do(parent, func(os.Signal) {
		cancel()
	})
	return ctx, cancel
}
