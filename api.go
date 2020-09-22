package startevent

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	URL       string
	Durations []time.Duration

	Waiter *Waiter
	Logger *log.Logger
}

func (c Config) Run(ctx context.Context, sentinel string) {
	url := c.URL
	if c.Waiter == nil {
		c.Waiter = &Waiter{
			Check:     HealthCheck(url, fmt.Errorf("fail")),
			Durations: DurationsFromSecs([]float64{0.1, 0.2, 0.2, 0.4, 0.8, 1.6, 3.2, 6.4, 12.8}),
		}
	}

	ch := c.Waiter.Start(ctx)

	go func() {
		t := <-ch
		c.Release(sentinel, t)
	}()
}

func (c Config) Release(sentinel string, d time.Duration) {
	if err := os.Remove(sentinel); err != nil {
		c.Logger.Println("ng", d, err)
		return
	}
	c.Logger.Println("ok", d, sentinel)
}
