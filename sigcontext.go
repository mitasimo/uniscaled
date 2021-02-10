package main

import (
	"context"
	"os"
	"os/signal"
)

// WithSignal creates a context that canceled when a signals is received
func WithSignal(parentCtx context.Context, s ...os.Signal) (context.Context, context.CancelFunc) {
	cancelContext, cancelFunc := context.WithCancel(parentCtx)
	signalContext := &signalContext{Context: cancelContext}
	signalContext.processSignal(cancelFunc, s...)
	return context.Context(signalContext), func() { cancelFunc() }
}

//
type signalContext struct {
	context.Context // parent context
}

// processSignal starts a goroutine that terminates
// the context when the OS receives a signal
func (c *signalContext) processSignal(cancel context.CancelFunc, s ...os.Signal) {
	chanSig := make(chan os.Signal, 1)
	signal.Notify(chanSig, s...)

	go func() {
		select {
		case <-chanSig:
			cancel()
		case <-c.Done():
			//cancel()
		}
		signal.Stop(chanSig)
	}()
}
