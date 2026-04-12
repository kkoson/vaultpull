// Package sync provides the core synchronisation engine for vaultpull.
//
// # Event Bus
//
// EventBus implements a lightweight publish/subscribe mechanism for sync
// lifecycle events. Consumers register an EventHandler via Subscribe and
// receive an Event value each time the bus publishes.
//
// Supported event types:
//
//	- EventSyncStarted    — emitted once before any profile is processed
//	- EventProfileStarted — emitted when a profile sync begins
//	- EventProfileSuccess — emitted when a profile sync completes without error
//	- EventProfileFailed  — emitted when a profile sync returns an error
//	- EventSyncCompleted  — emitted after all profiles have been processed
//
// Example:
//
//	bus := sync.NewEventBus()
//	bus.Subscribe(func(e sync.Event) {
//	    fmt.Printf("[%s] %s\n", e.Type, e.Profile)
//	})
package sync
