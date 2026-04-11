// Package rotate provides AppRole secret-id rotation for vaultpull profiles.
//
// When using AppRole authentication, Vault secret-ids have a finite TTL.
// The rotate package automates generating a fresh secret-id via the Vault API
// and persisting it to the local token cache so subsequent syncs continue
// to authenticate without manual intervention.
package rotate
