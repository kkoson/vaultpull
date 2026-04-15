package sync

import (
	"testing"
)

func TestDefaultSamplingConfig_Values(t *testing.T) {
	cfg := DefaultSamplingConfig()
	if cfg.Mode != SamplingModeRandom {
		t.Errorf("expected SamplingModeRandom, got %v", cfg.Mode)
	}
	if cfg.Rate != 1.0 {
		t.Errorf("expected rate 1.0, got %v", cfg.Rate)
	}
}

func TestNewSampler_ClampsNegativeRate(t *testing.T) {
	s := NewSampler(SamplingConfig{Mode: SamplingModeRandom, Rate: -0.5})
	if s.cfg.Rate != 0 {
		t.Errorf("expected rate clamped to 0, got %v", s.cfg.Rate)
	}
}

func TestNewSampler_ClampsRateAboveOne(t *testing.T) {
	s := NewSampler(SamplingConfig{Mode: SamplingModeRandom, Rate: 2.5})
	if s.cfg.Rate != 1.0 {
		t.Errorf("expected rate clamped to 1.0, got %v", s.cfg.Rate)
	}
}

func TestSampler_AlwaysMode_ReturnsTrue(t *testing.T) {
	s := NewSampler(SamplingConfig{Mode: SamplingModeAlways})
	for i := 0; i < 20; i++ {
		if !s.Sample("profile-a") {
			t.Fatal("expected Sample to return true in Always mode")
		}
	}
}

func TestSampler_NeverMode_ReturnsFalse(t *testing.T) {
	s := NewSampler(SamplingConfig{Mode: SamplingModeNever})
	for i := 0; i < 20; i++ {
		if s.Sample("profile-a") {
			t.Fatal("expected Sample to return false in Never mode")
		}
	}
}

func TestSampler_RandomMode_RateZero_NeverSamples(t *testing.T) {
	s := NewSampler(SamplingConfig{Mode: SamplingModeRandom, Rate: 0})
	for i := 0; i < 50; i++ {
		if s.Sample("profile-b") {
			t.Fatal("expected no samples with rate 0")
		}
	}
}

func TestSampler_RandomMode_RateOne_AlwaysSamples(t *testing.T) {
	s := NewSampler(SamplingConfig{Mode: SamplingModeRandom, Rate: 1.0})
	for i := 0; i < 20; i++ {
		if !s.Sample("profile-c") {
			t.Fatal("expected all samples with rate 1.0")
		}
	}
}

func TestSampler_Stats_TracksHitsAndMisses(t *testing.T) {
	s := NewSampler(SamplingConfig{Mode: SamplingModeAlways})
	s.Sample("p1")
	s.Sample("p2")
	hits, misses := s.Stats()
	if hits != 2 {
		t2 hits, got %d", hits)
	}
	if misses != 0 {
		t.Errorf("expected 0 misses, got %d", misses)
	}
}

func TestSampler_Stats_TracksNeverMisses(t *testing.T) {
	s := NewSampler(SamplingConfig{Mode: SamplingModeNever})
	s.Sample("p1")
	s.Sample("p2")
	s.Sample("p3")
	hits, misses := s.Stats()
	if hits != 0 {
		t.Errorf("expected 0 hits, got %d", hits)
	}
	if misses != 3 {
		t.Errorf("expected 3 misses, got %d", misses)
	}
}
