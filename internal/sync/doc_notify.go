// Package sync provides synchronisation primitives for vaultpull.
//
// # Notifier
//
// Notifier emits human-readable messages whenever a profile sync completes.
// It supports three verbosity levels:
//
//   - NotifyNone    – silent; no output is produced.
//   - NotifyFailure – only failed syncs are reported.
//   - NotifyAll     – both successful and failed syncs are reported.
//
// Example:
//
//	n := sync.NewNotifier(os.Stderr, sync.NotifyAll)
//	n.Notify(sync.NotifyEvent{
//		Profile:  "production",
//		Success:  true,
//		Changes:  3,
//		Duration: 120 * time.Millisecond,
//	})
package sync
