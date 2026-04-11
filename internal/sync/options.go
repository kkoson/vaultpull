package sync

// Options configures the behaviour of a Syncer run.
type Options struct {
	// DryRun prints the diff without writing any changes to disk.
	DryRun bool

	// Verbose enables additional output during the sync.
	Verbose bool

	// OverwriteExisting controls whether existing keys in the .env file are
	// overwritten with values fetched from Vault. When false the local value
	// is kept and only new keys are added.
	OverwriteExisting bool
}

// DefaultOptions returns an Options struct populated with sensible defaults.
func DefaultOptions() Options {
	return Options{
		DryRun:            false,
		Verbose:           false,
		OverwriteExisting: true,
	}
}
