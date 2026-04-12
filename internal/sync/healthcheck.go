package sync

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthStatus represents the current health state of a component.
type HealthStatus string

const (
	HealthOK      HealthStatus = "ok"
	HealthDegraded HealthStatus = "degraded"
	HealthDown    HealthStatus = "down"
)

// HealthResult holds the result of a single health check probe.
type HealthResult struct {
	Name      string
	Status    HealthStatus
	Message   string
	CheckedAt time.Time
}

// ProbeFunc is a function that performs a health check and returns an error if unhealthy.
type ProbeFunc func(ctx context.Context) error

// HealthChecker runs named probes and aggregates their results.
type HealthChecker struct {
	mu     sync.RWMutex
	probes map[string]ProbeFunc
}

// NewHealthChecker creates a new HealthChecker with no registered probes.
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		probes: make(map[string]ProbeFunc),
	}
}

// Register adds a named probe to the checker. Overwrites any existing probe with the same name.
func (h *HealthChecker) Register(name string, probe ProbeFunc) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.probes[name] = probe
}

// RunAll executes all registered probes concurrently and returns aggregated results.
func (h *HealthChecker) RunAll(ctx context.Context) []HealthResult {
	h.mu.RLock()
	names := make([]string, 0, len(h.probes))
	for name := range h.probes {
		names = append(names, name)
	}
	h.mu.RUnlock()

	resultCh := make(chan HealthResult, len(names))
	var wg sync.WaitGroup

	for _, name := range names {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			h.mu.RLock()
			probe := h.probes[n]
			h.mu.RUnlock()

			result := HealthResult{Name: n, CheckedAt: time.Now()}
			if err := probe(ctx); err != nil {
				result.Status = HealthDown
				result.Message = err.Error()
			} else {
				result.Status = HealthOK
				result.Message = "healthy"
			}
			resultCh <- result
		}(name)
	}

	wg.Wait()
	close(resultCh)

	results := make([]HealthResult, 0, len(names))
	for r := range resultCh {
		results = append(results, r)
	}
	return results
}

// Overall returns a single aggregated status across all results.
func Overall(results []HealthResult) HealthStatus {
	if len(results) == 0 {
		return HealthOK
	}
	for _, r := range results {
		if r.Status == HealthDown {
			return HealthDown
		}
	}
	return HealthOK
}

// Summary returns a human-readable summary string.
func Summary(results []HealthResult) string {
	return fmt.Sprintf("health: %s (%d probes)", Overall(results), len(results))
}
