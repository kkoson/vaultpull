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
// The verbosity level can be parsed from a string using ParseNotifyLevel,
// which is useful when reading configuration from environment variables or
// command-line flags.
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
//
// Parsing a level from a string:
//
//	level, err := sync.ParseNotifyLevel("failure")
//	if err != nil {
//		log.Fatal(err)
//	}
//	n := sync.NewNotifier(os.Stderr, level)
package sync
