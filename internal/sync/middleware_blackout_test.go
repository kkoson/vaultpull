package sync

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildBlackoutProfile() Profile {
	return Profile{Name: "test", VaultPath: "secret/test", OutputFile: "/tmp/test.env"}
}

func TestWithBlackout_PanicsOnNil(t *testing.T) {
	assert.Panics(t, func() {
		WithBlackout(nil)(context.Background(), buildBlackoutProfile(), func(_ context.Context, _ Profile) error {
			return nil
		})
	})
}

func TestWithBlackout_AllowsWhenNotBlackedOut(t *testing.T) {
	cfg := DefaultBlackoutConfig()
	bm := NewBlackoutManager(cfg)

	called := false
	next := func(_ context.Context, _ Profile) error {
		called = true
		return nil
	}

	err := WithBlackout(bm)(context.Background(), buildBlackoutProfile(), next)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestWithBlackout_SkipsWhenBlackedOut(t *testing.T) {
	cfg := DefaultBlackoutConfig()
	now := time.Now()
	start := now.Add(-1 * time.Hour).Format("15:04")
	end := now.Add(1 * time.Hour).Format("15:04")
	cfg.Windows = []BlackoutWindow{{Start: start, End: end}}
	bm := NewBlackoutManager(cfg)

	called := false
	next := func(_ context.Context, _ Profile) error {
		called = true
		return nil
	}

	err := WithBlackout(bm)(context.Background(), buildBlackoutProfile(), next)
	require.NoError(t, err)
	assert.False(t, called)
}

func TestWithBlackout_PropagatesInnerError(t *testing.T) {
	cfg := DefaultBlackoutConfig()
	bm := NewBlackoutManager(cfg)

	expected := errors.New("inner error")
	next := func(_ context.Context, _ Profile) error {
		return expected
	}

	err := WithBlackout(bm)(context.Background(), buildBlackoutProfile(), next)
	assert.ErrorIs(t, err, expected)
}
