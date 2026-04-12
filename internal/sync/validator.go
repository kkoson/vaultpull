package sync

import (
	"errors"
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
)

// ValidationResult holds the outcome of a profile validation check.
type ValidationResult struct {
	Profile string
	Errors  []string
}

// Valid returns true when no validation errors were recorded.
func (r ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

// Validator checks a set of profiles for configuration correctness
// before a sync run is attempted.
type Validator struct {
	cfg *config.Config
}

// NewValidator creates a Validator backed by the provided Config.
func NewValidator(cfg *config.Config) *Validator {
	return &Validator{cfg: cfg}
}

// ValidateAll validates every profile defined in the config and returns
// one ValidationResult per profile. An error is returned only when the
// config itself is nil.
func (v *Validator) ValidateAll() ([]ValidationResult, error) {
	if v.cfg == nil {
		return nil, errors.New("validator: config must not be nil")
	}

	results := make([]ValidationResult, 0, len(v.cfg.Profiles))
	for _, p := range v.cfg.Profiles {
		results = append(results, v.validateProfile(p))
	}
	return results, nil
}

// ValidateProfile validates a single named profile and returns its result.
// An error is returned when the profile does not exist in the config.
func (v *Validator) ValidateProfile(name string) (ValidationResult, error) {
	p, err := v.cfg.GetProfile(name)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("validator: %w", err)
	}
	return v.validateProfile(*p), nil
}

func (v *Validator) validateProfile(p config.Profile) ValidationResult {
	result := ValidationResult{Profile: p.Name}
	if err := p.Validate(); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}
	if p.OutputFile == "" {
		result.Errors = append(result.Errors, "output_file must not be empty")
	}
	return result
}
