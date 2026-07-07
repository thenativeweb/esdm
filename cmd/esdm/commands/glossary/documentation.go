// Package glossary implements the `esdm glossary` command.
// It reads the ubiquitous language declared on the bounded
// contexts of an ESDM model and writes it to stdout as a
// human-readable Markdown glossary.
//
// Without arguments the command emits the glossary for the
// whole model in the current directory. An optional path
// argument narrows the output: a single segment selects a
// domain and emits the glossary for every bounded context
// inside it, two segments select a single bounded context.
// The path follows the model hierarchy the same way the
// `esdm view` path does, and an unknown segment is rejected
// as invalid input.
//
// Output is plain Markdown on stdout so it can be redirected
// into a file; the command therefore has no color option.
package glossary
