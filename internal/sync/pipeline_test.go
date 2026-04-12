package sync

import (
	"context"
	"errors"
	"testing"
)

func TestNewPipeline_Empty(t *testing.T) {
	p := NewPipeline()
	if p.Len() != 0 {
		t.Fatalf("expected 0 stages, got %d", p.Len())
	}
}

func TestPipeline_AddStage(t *testing.T) {
	p := NewPipeline()
	p.AddStage(Stage{Name: "test", Run: func(ctx context.Context, name string) error { return nil }})
	if p.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", p.Len())
	}
}

func TestPipeline_Execute_AllSuccess(t *testing.T) {
	var order []string
	p := NewPipeline(
		Stage{Name: "first", Run: func(ctx context.Context, name string) error {
			order = append(order, "first")
			return nil
		}},
		Stage{Name: "second", Run: func(ctx context.Context, name string) error {
			order = append(order, "second")
			return nil
		}},
	)
	if err := p.Execute(context.Background(), "dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 || order[0] != "first" || order[1] != "second" {
		t.Fatalf("unexpected execution order: %v", order)
	}
}

func TestPipeline_Execute_StopsOnError(t *testing.T) {
	ran := false
	p := NewPipeline(
		Stage{Name: "fail", Run: func(ctx context.Context, name string) error {
			return errors.New("boom")
		}},
		Stage{Name: "should-not-run", Run: func(ctx context.Context, name string) error {
			ran = true
			return nil
		}},
	)
	err := p.Execute(context.Background(), "dev")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if ran {
		t.Fatal("second stage should not have run")
	}
}

func TestPipeline_Execute_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := NewPipeline(
		Stage{Name: "unreachable", Run: func(ctx context.Context, name string) error {
			return nil
		}},
	)
	err := p.Execute(ctx, "dev")
	if err == nil {
		t.Fatal("expected cancellation error, got nil")
	}
}

func TestPipeline_Execute_PassesProfileName(t *testing.T) {
	var got string
	p := NewPipeline(
		Stage{Name: "capture", Run: func(ctx context.Context, name string) error {
			got = name
			return nil
		}},
	)
	_ = p.Execute(context.Background(), "staging")
	if got != "staging" {
		t.Fatalf("expected profile \"staging\", got %q", got)
	}
}
