// Package sync provides utilities for syncing secrets from Vault into local
// environment files.
//
// # Pipeline
//
// Pipeline allows composing multiple sync stages into an ordered execution
// chain. Each Stage has a name and a Run function that receives a context and
// the target profile name.
//
// Stages are executed sequentially. If any stage returns an error the pipeline
// halts immediately and returns a wrapped error that includes the stage name
// and profile name for easy debugging.
//
// Example usage:
//
//	p := sync.NewPipeline(
//		sync.Stage{Name: "validate",  Run: validateFn},
//		sync.Stage{Name: "fetch",     Run: fetchFn},
//		sync.Stage{Name: "write",     Run: writeFn},
//	)
//	if err := p.Execute(ctx, "production"); err != nil {
//		log.Fatal(err)
//	}
package sync
