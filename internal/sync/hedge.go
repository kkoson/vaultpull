package sync

import (
	"context"
	"sync"
	"time"
)

// HedgeConfig controls hedged-request behaviour.
type HedgeConfig struct {
	// Delay is how long to wait before issuing the second request.
	Delay time.Duration
	// MaxHedges is the maximum number of additional requests to issue.
	MaxHedges int
}

// DefaultHedgeConfig returns a HedgeConfig with sensible defaults.
func DefaultHedgeConfig() HedgeConfig {
	return HedgeConfig{
		Delay:     200 * time.Millisecond,
		MaxHedges: 1,
	}
}

// HedgeResult carries the outcome of a hedged call.
type HedgeResult struct {
	Value interface{}
	Err   error
}

// Hedge runs fn and, after cfg.Delay, issues up to cfg.MaxHedges additional
// concurrent calls. The first successful result wins; if all fail the last
// error is returned.
func Hedge(ctx context.Context, cfg HedgeConfig, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	if cfg.MaxHedges <= 0 {
		cfg.MaxHedges = 1
	}
	if cfg.Delay <= 0 {
		cfg.Delay = DefaultHedgeConfig().Delay
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultCh := make(chan HedgeResult, cfg.MaxHedges+1)

	launch := func() {
		v, err := fn(ctx)
		resultCh <- HedgeResult{Value: v, Err: err}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		launch()
	}()

	for i := 0; i < cfg.MaxHedges; i++ {
		select {
		case res := <-resultCh:
			if res.Err == nil {
				return res.Value, nil
			}
		case <-time.After(cfg.Delay):
			wg.Add(1)
			go func() {
				defer wg.Done()
				launch()
			}()
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var lastErr error
	for res := range resultCh {
		if res.Err == nil {
			return res.Value, nil
		}
		lastErr = res.Err
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, ctx.Err()
}
