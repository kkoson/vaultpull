package sync

import "sync"

// Deduplicator suppresses duplicate sync runs for the same profile
// that are triggered concurrently. If a sync for a given key is already
// in-flight, subsequent calls block and share the same result.
type Deduplicator struct {
	mu    sync.Mutex
	group map[string]*dedupCall
}

type dedupCall struct {
	wg  sync.WaitGroup
	err error
}

// NewDeduplicator returns an initialised Deduplicator.
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{
		group: make(map[string]*dedupCall),
	}
}

// Do executes fn for key exactly once while it is in-flight.
// Concurrent callers with the same key wait for the first call to finish
// and receive its error. Once the call completes the key is removed so
// that future calls can trigger a new execution.
func (d *Deduplicator) Do(key string, fn func() error) error {
	d.mu.Lock()
	if call, ok := d.group[key]; ok {
		d.mu.Unlock()
		call.wg.Wait()
		return call.err
	}

	call := &dedupCall{}
	call.wg.Add(1)
	d.group[key] = call
	d.mu.Unlock()

	call.err = fn()
	call.wg.Done()

	d.mu.Lock()
	delete(d.group, key)
	d.mu.Unlock()

	return call.err
}

// InFlight returns the number of keys currently being executed.
func (d *Deduplicator) InFlight() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.group)
}
