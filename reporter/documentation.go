// Package reporter turns collected diagnostics into concrete
// output. It provides a thread-safe collector that gathers the
// diagnostics emitted by concurrently running rules, plus
// formatters that render them as human-readable terminal
// output or as structured JSON.
//
// Diagnostics are collected during the run and formatted
// afterwards. Output is always sorted deterministically by
// file, line, column, and rule ID so that repeated runs on
// unchanged input produce identical output.
package reporter
