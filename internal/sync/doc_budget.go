// Package sync — ErrorBudget
//
// ErrorBudget implements a sliding-window error-budget tracker. It records
// boolean outcomes (success / failure) for sync operations and reports whether
// the configured maximum error rate has been exceeded across the most recent
// window of samples.
//
// Usage:
//
//	budget := sync.NewErrorBudget(sync.DefaultBudgetConfig())
//
//	// after each sync attempt:
//	budget.Record(err == nil)
//
//	if budget.Exhausted() {
//		log.Println("error budget exhausted:", budget.Stats())
//	}
//
// Configuration:
//
//	sync.BudgetConfig{
//		MaxErrorRate:  0.5,  // halt when >50 % of recent ops fail
//		WindowSize:    10,   // track last 10 outcomes
//		MinSampleSize: 3,    // need at least 3 samples before enforcing
//	}
package sync
