package sync

import (
	"hash/fnv"
	"sync"

	"github.com/yourusername/vaultpull/internal/config"
)

// DefaultAffinityConfig returns an AffinityConfig with sensible defaults.
func DefaultAffinityConfig() AffinityConfig {
	return AffinityConfig{
		TagKey:    "affinity",
		SlotCount: 16,
	}
}

// AffinityConfig controls the behaviour of AffinityRouter.
type AffinityConfig struct {
	// TagKey is the profile tag used to look up an affinity group.
	TagKey string
	// SlotCount is the total number of logical affinity slots.
	SlotCount uint32
}

// AffinityRouter assigns each profile a stable numeric slot.
type AffinityRouter struct {
	cfg   AffinityConfig
	mu    sync.Mutex
	cache map[string]uint32
}

// NewAffinityRouter creates a router using the supplied config.
// Zero-value fields fall back to DefaultAffinityConfig values.
func NewAffinityRouter(cfg AffinityConfig) *AffinityRouter {
	def := DefaultAffinityConfig()
	if cfg.TagKey == "" {
		cfg.TagKey = def.TagKey
	}
	if cfg.SlotCount == 0 {
		cfg.SlotCount = def.SlotCount
	}
	return &AffinityRouter{
		cfg:   cfg,
		cache: make(map[string]uint32),
	}
}

// Assign returns a stable slot index in [0, SlotCount) for the profile.
// Profiles sharing the same affinity tag receive the same slot.
func (r *AffinityRouter) Assign(p config.Profile) uint32 {
	key := r.affinityKey(p)

	r.mu.Lock()
	defer r.mu.Unlock()

	if slot, ok := r.cache[key]; ok {
		return slot
	}

	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	slot := h.Sum32() % r.cfg.SlotCount
	r.cache[key] = slot
	return slot
}

// affinityKey returns the value of the affinity tag, falling back to the
// profile name when the tag is absent.
func (r *AffinityRouter) affinityKey(p config.Profile) string {
	for _, tag := range p.Tags {
		if len(tag) > len(r.cfg.TagKey)+1 && tag[:len(r.cfg.TagKey)+1] == r.cfg.TagKey+":" {
			return tag[len(r.cfg.TagKey)+1:]
		}
	}
	return p.Name
}
