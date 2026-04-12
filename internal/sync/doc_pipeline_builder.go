// Package sync provides synchronisation primitives for vaultpull.
//
// # PipelineBuilder
//
// PipelineBuilder constructs a ready-to-run [Pipeline] for a given
// [config.Profile]. It wires together the two mandatory stages:
//
//  1. validate – confirms the profile is structurally valid before any
//     network call is made.
//  2. sync – fetches secrets from Vault and writes them to the output file.
//
// Usage:
//
//	pb := sync.NewPipelineBuilder(runner)
//	pl, err := pb.Build(ctx, profile)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if err := pl.Execute(ctx); err != nil {
//	    log.Fatal(err)
//	}
package sync
