package sync

import "context"

// WithRedact returns a StageFunc middleware that redacts sensitive keys
// from the secrets map before passing them to the next stage.
// The redacted map is written back into the context under redactedSecretsKey
// so downstream stages receive safe-to-log values, while the original
// secrets map (used for writing) is left untouched in the caller.
//
// Usage:
//
//	protected := WithRedact(redactor)(myStage)
func WithRedact(r *Redactor) func(StageFunc) StageFunc {
	return func(next StageFunc) StageFunc {
		return func(ctx context.Context, profile string) error {
			ctx = contextWithRedactor(ctx, r)
			return next(ctx, profile)
		}
	}
}

type contextKey int

const redactorKey contextKey = iota

// contextWithRedactor stores r in ctx.
func contextWithRedactor(ctx context.Context, r *Redactor) context.Context {
	return context.WithValue(ctx, redactorKey, r)
}

// RedactorFromContext retrieves the Redactor stored by WithRedact.
// Returns nil if none was stored.
func RedactorFromContext(ctx context.Context) *Redactor {
	v, _ := ctx.Value(redactorKey).(*Redactor)
	return v
}
