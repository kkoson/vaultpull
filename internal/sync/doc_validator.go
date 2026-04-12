// Package sync provides synchronisation primitives for pulling secrets
// from HashiCorp Vault into local .env files.
//
// # Validator
//
// Validator performs pre-flight configuration checks against one or all
// profiles defined in a Config before a sync run is started.
//
// Basic usage:
//
//	v := sync.NewValidator(cfg)
//
//	// Validate every profile at once.
//	results, err := v.ValidateAll()
//	for _, r := range results {
//		if !r.Valid() {
//			fmt.Printf("profile %s has errors: %v\n", r.Profile, r.Errors)
//		}
//	}
//
//	// Validate a single profile by name.
//	result, err := v.ValidateProfile("production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !result.Valid() {
//		fmt.Printf("profile has errors: %v\n", result.Errors)
//	}
package sync
