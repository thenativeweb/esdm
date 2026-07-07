// Package lint provides the `esdm lint` subcommand.
// It runs the linter pipeline against a directory and
// formats the collected diagnostics for human or machine
// consumption.
//
// Reporting lint findings and reporting a failed command
// are deliberately kept apart. A run that loads, parses,
// resolves, and reports without trouble is a successful
// command, even when the inspected model has findings;
// only a genuine failure, such as an unreadable directory
// or an invalid flag, is surfaced as a command error.
// Whether the model was clean travels instead on a
// separate exit code: zero when no error-severity finding
// was produced, non-zero otherwise, with warnings
// optionally escalated to that same non-zero outcome. That
// code lives as package state, reset at the start of every
// run and read back by the process entry point after the
// command returns, so a successful invocation and a
// has-findings signal never get conflated.
package lint
