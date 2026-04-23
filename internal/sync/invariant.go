package sync

import (
	"fmt"
	"sync"
)

// InvariantRule defines a named check applied to a secret map.
type InvariantRule struct {
	Name  string
	Check func(secrets map[string]string) error
}

// InvariantChecker validates a set of rules against synced secrets.
type InvariantChecker struct {
	mu    sync.RWMutex
	rules []InvariantRule
}

// NewInvariantChecker returns an empty InvariantChecker.
func NewInvariantChecker() *InvariantChecker {
	return &InvariantChecker{}
}

// Register adds a named invariant rule. Nil checks are ignored.
func (c *InvariantChecker) Register(rule InvariantRule) {
	if rule.Check == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rules = append(c.rules, rule)
}

// Len returns the number of registered rules.
func (c *InvariantChecker) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.rules)
}

// Check runs all registered rules against secrets.
// It returns the first violation encountered, or nil if all pass.
func (c *InvariantChecker) Check(secrets map[string]string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, rule := range c.rules {
		if err := rule.Check(secrets); err != nil {
			return fmt.Errorf("invariant %q violated: %w", rule.Name, err)
		}
	}
	return nil
}

// RequireKeys returns a rule that fails if any of the given keys are absent.
func RequireKeys(keys ...string) InvariantRule {
	return InvariantRule{
		Name: "require-keys",
		Check: func(secrets map[string]string) error {
			for _, k := range keys {
				if _, ok := secrets[k]; !ok {
					return fmt.Errorf("required key %q is missing", k)
				}
			}
			return nil
		},
	}
}

// ForbidKeys returns a rule that fails if any of the given keys are present.
func ForbidKeys(keys ...string) InvariantRule {
	return InvariantRule{
		Name: "forbid-keys",
		Check: func(secrets map[string]string) error {
			for _, k := range keys {
				if _, ok := secrets[k]; ok {
					return fmt.Errorf("forbidden key %q is present", k)
				}
			}
			return nil
		},
	}
}
