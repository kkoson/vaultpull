package sync

import (
	"sync"
	"time"
)

// ShedderMode controls how load shedding decisions are made.
type ShedderMode int

const (
	// ShedRandom drops requests randomly based on load.
	ShedRandom ShedderMode = iota
	// ShedTail drops the newest requests when overloaded.
	ShedTail
)

// DefaultShedderConfig returns a ShedderConfig with sensible defaults.
func DefaultShedderConfig() ShedderConfig {
	return ShedderConfig{
		Mode:        ShedRandom,
		MaxLoad:     100,
		WindowSize:  5 * time.Second,
		DropPercent: 0.5,
	}
}

// ShedderConfig configures the load shedder.
type ShedderConfig struct {
	Mode        ShedderMode
	MaxLoad     int
	WindowSize  time.Duration
	DropPercent float64
}

// LoadShedder tracks in-flight requests and sheds load when overloaded.
type LoadShedder struct {
	mu      sync.Mutex
	cfg     ShedderConfig
	current int
	shed    int
	total   int
	window  time.Time
}

// NewLoadShedder creates a LoadShedder with the given config.
// Zero values fall back to defaults.
func NewLoadShedder(cfg ShedderConfig) *LoadShedder {
	def := DefaultShedderConfig()
	if cfg.MaxLoad <= 0 {
		cfg.MaxLoad = def.MaxLoad
	}
	if cfg.WindowSize <= 0 {
		cfg.WindowSize = def.WindowSize
	}
	if cfg.DropPercent <= 0 {
		cfg.DropPercent = def.DropPercent
	}
	return &LoadShedder{cfg: cfg, window: time.Now()}
}

// Admit returns true if the request should be allowed through.
// It increments the in-flight counter and resets the window when expired.
func (s *LoadShedder) Admit() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if time.Since(s.window) > s.cfg.WindowSize {
		s.current = 0
		s.shed = 0
		s.total = 0
		s.window = time.Now()
	}
	s.total++
	if s.current >= s.cfg.MaxLoad {
		s.shed++
		return false
	}
	s.current++
	return true
}

// Release decrements the in-flight counter.
func (s *LoadShedder) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.current > 0 {
		s.current--
	}
}

// Stats returns a snapshot of shedder counters.
func (s *LoadShedder) Stats() (current, shed, total int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current, s.shed, s.total
}
