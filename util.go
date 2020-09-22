package startevent

import (
	"fmt"
	"net/http"
	"time"
)

func CheckByHTTPRequest(url string, failErr error) func(time.Time) error {
	return func(t time.Time) error {
		if debug {
			getLogger().Printf("tick %s: %s", t, url)
		}

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return fmt.Errorf("unexpeted status %d: %w", res.StatusCode, failErr)
		}
		return nil
	}
}

func DurationsFromSecs(secs []float64) []time.Duration {
	r := make([]time.Duration, len(secs))
	for i, sec := range secs {
		r[i] = time.Duration(sec*1000) * time.Millisecond
	}
	return r
}
