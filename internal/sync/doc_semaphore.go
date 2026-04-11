// Package sync provides orchestration utilities for syncing Vault secrets
// into local .env files across multiple profiles.
//
// # Semaphore
//
// The Semaphore type controls how many profile sync operations may run
// concurrently. This is useful when syncing a large number of profiles to
// avoid overwhelming the Vault server or exhausting local resources.
//
// Example usage:
//
//	sem, err := sync.NewSemaphore(5)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, p := range profiles {
//		if err := sem.Acquire(ctx); err != nil {
//			break // context cancelled
//		}
//		go func(profile string) {
//			defer sem.Release()
//			// ... sync profile ...
//		}(p)
//	}
package sync
