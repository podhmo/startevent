package startevent

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Config struct {
	URL       string
	Durations []time.Duration

	Waiter *Waiter
	Logger logger
}

func (c Config) Run(ctx context.Context, sentinel string) {
	url := c.URL
	if c.Waiter == nil {
		c.Waiter = &Waiter{
			Check:     HealthCheck(url, fmt.Errorf("fail")),
			Durations: DurationsFromSecs([]float64{0.1, 0.2, 0.2, 0.4, 0.8, 1.6, 3.2, 6.4, 12.8}),
		}
	}
	if c.Logger == nil {
		c.Logger = getLogger()
	}

	ch := c.Waiter.Start(ctx)

	go func() {
		t := <-ch
		c.Release(sentinel, t)
	}()
}

func (c Config) Release(sentinel string, d time.Duration) {
	if err := os.Remove(sentinel); err != nil {
		c.Logger.Printf("ng	duration=%s	err=%s", d, err)
		return
	}
	c.Logger.Printf("ok	duration=%s	file=%s", d, sentinel)
}
