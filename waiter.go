package startevent

import (
	"context"
	"log"
	"time"
)

type Waiter struct {
	Check     func(time.Time) error
	Durations []time.Duration

	Logger *log.Logger
}

func (w *Waiter) Start(ctx context.Context) <-chan time.Duration {
	finishCH := make(chan time.Duration)
	go func() {
		defer close(finishCH)
		st := time.Now()

		defer func() {
			panicErr := recover()
			if panicErr != nil {
				log.Fatalf("!!%+v", panicErr)
			}
		}()

		for _, duration := range w.Durations {
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
		log.Printf("timeout: %s", time.Now().Sub(st))
	}()
	return finishCH
}

func (w *Waiter) Wait(ctx context.Context) {
	<-w.Start(ctx)
}
