// Package audit provides structured JSON audit logging for vaultpull sync
// operations. Each sync run produces an Entry that records which profile was
// synced, how many secrets were added, updated, removed or left unchanged, and
// whether the run was a dry-run. Entries are written as newline-delimited JSON
// and can be appended to a file or streamed to stdout.
package audit
