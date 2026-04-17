package sync

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRunner struct {
	fn func(ctx context.Context, name string) error
}

func (m *mockRunner) RunProfile(ctx context.Context, name string) error {
	return m.fn(ctx, name)
}

func TestNewCascadeRunner_PanicsOnNil(t *testing.T) {
	assert.Panics(t, func() { NewCascadeRunner(nil) })
}

func TestCascadeRunner_NoDeps_RunsPrimary(t *testing.T) {
	var called []string
	r := &mockRunner{fn: func(_ context.Context, name string) error {
		called = append(called, name)
		return nil
	}}
	c := NewCascadeRunner(r)
	err := c.Run(context.Background(), "primary")
	require.NoError(t, err)
	assert.Equal(t, []string{"primary"}, called)
}

func TestCascadeRunner_PrimaryFails_DepsNotRun(t *testing.T) {
	var depCalled atomic.Bool
	r := &mockRunner{fn: func(_ context.Context, name string) error {
		if name == "primary" {
			return errors.New("primary error")
		}
		depCalled.Store(true)
		return nil
	}}
	c := NewCascadeRunner(r)
	c.AddDependency("primary", "dep1")
	err := c.Run(context.Background(), "primary")
	assert.Error(t, err)
	assert.False(t, depCalled.Load())
}

func TestCascadeRunner_DepsRunAfterPrimary(t *testing.T) {
	var calls atomic.Int32
	r := &mockRunner{fn: func(_ context.Context, _ string) error {
		calls.Add(1)
		return nil
	}}
	c := NewCascadeRunner(r)
	c.AddDependency("primary", "dep1")
	c.AddDependency("primary", "dep2")
	err := c.Run(context.Background(), "primary")
	require.NoError(t, err)
	assert.Equal(t, int32(3), calls.Load())
}

func TestCascadeRunner_DepFails_ReturnsError(t *testing.T) {
	r := &mockRunner{fn: func(_ context.Context, name string) error {
		if name == "dep1" {
			return errors.New("dep error")
		}
		return nil
	}}
	c := NewCascadeRunner(r)
	c.AddDependency("primary", "dep1")
	err := c.Run(context.Background(), "primary")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dep1")
}
