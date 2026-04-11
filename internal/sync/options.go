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

// WithDryRun returns a copy of the Options with DryRun set to the given value.
func (o Options) WithDryRun(dryRun bool) Options {
	o.DryRun = dryRun
	return o
}

// WithVerbose returns a copy of the Options with Verbose set to the given value.
func (o Options) WithVerbose(verbose bool) Options {
	o.Verbose = verbose
	return o
}

// WithOverwriteExisting returns a copy of the Options with OverwriteExisting set to the given value.
func (o Options) WithOverwriteExisting(overwrite bool) Options {
	o.OverwriteExisting = overwrite
	return o
}
