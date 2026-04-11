package sync

import (
	"testing"
)

func TestDefaultOptions_DryRun(t *testing.T) {
	opts := DefaultOptions()
	if opts.DryRun {
		t.Error("expected DryRun to be false by default")
	}
}

func TestDefaultOptions_Verbose(t *testing.T) {
	opts := DefaultOptions()
	if opts.Verbose {
		t.Error("expected Verbose to be false by default")
	}
}

func TestDefaultOptions_OverwriteExisting(t *testing.T) {
	opts := DefaultOptions()
	if !opts.OverwriteExisting {
		t.Error("expected OverwriteExisting to be true by default")
	}
}

func TestOptions_DryRunToggle(t *testing.T) {
	opts := DefaultOptions()
	opts.DryRun = true
	if !opts.DryRun {
		t.Error("expected DryRun to be true after toggle")
	}
}

func TestOptions_OverwriteExistingFalse(t *testing.T) {
	opts := DefaultOptions()
	opts.OverwriteExisting = false
	if opts.OverwriteExisting {
		t.Error("expected OverwriteExisting to be false after setting")
	}
}

func TestOptions_VerboseToggle(t *testing.T) {
	opts := DefaultOptions()
	opts.Verbose = true
	if !opts.Verbose {
		t.Error("expected Verbose to be true after toggle")
	}
}
