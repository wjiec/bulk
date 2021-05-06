package xcontext

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var state struct {
	sync.Mutex
	handler map[os.Signal][]func()
}

func init() {
	state.handler = make(map[os.Signal][]func())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	go func() {
		for sig := range signals {
			state.Lock()
			if cbs, ok := state.handler[sig]; ok {
				for _, cb := range cbs {
					go cb()
				}
			}
			state.Unlock()
		}
	}()
}

func register(fn func(), signals ...os.Signal) {
	state.Lock()
	for _, sig := range signals {
		state.handler[sig] = append(state.handler[sig], fn)
	}
	state.Unlock()
}

// WithSignal returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when any of the signals are occurs
// or when the parent context's Done channel is closed, whichever happens first
func WithSignal(parent context.Context, signals ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(parent)
	register(cancel, signals...)

	return ctx
}

// WithExitSignal wrapped WithSignal and pass common exit signals
func WithExitSignal(parent context.Context) context.Context {
	return WithSignal(parent, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)
}
