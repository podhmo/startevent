package startevent

import (
	"context"
	"time"
)

type Waiter struct {
	Check     func(time.Time) error
	Durations []time.Duration

	Logger logger
}

func (w *Waiter) Start(ctx context.Context) <-chan time.Duration {
	finishCH := make(chan time.Duration)
	go func() {
		defer close(finishCH)
		st := time.Now()

		defer func() {
			panicErr := recover()
			if panicErr != nil {
				w.Logger.Fatalf("!!%+v", panicErr)
			}
		}()

		for _, duration := range w.Durations {
			if debug {
				w.Logger.Printf("wait: %s", duration)
			}
			select {
			case <-ctx.Done():
				return
			case t := <-time.After(duration):
				if err := w.Check(t); err != nil {
					continue
				}
				finishCH <- time.Now().Sub(st)
			}
		}

		// timeout
		w.Logger.Printf("timeout: %s", time.Now().Sub(st))
	}()
	return finishCH
}

func (w *Waiter) Wait(ctx context.Context) {
	<-w.Start(ctx)
}
