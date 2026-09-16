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
// A bounded context's vocabulary is written in one language,
// and its terms may carry translations. By default every
// bounded context is rendered in its own language. With a
// language selected, a bounded context written in that
// language is rendered as is and every other one through its
// translations, each with the translation's own rejected
// alternatives. A term without a translation into the
// selected language keeps its primary form and is marked, so
// the glossary shows where the language is incomplete
// instead of dropping the term.
//
// Output is plain Markdown on stdout so it can be redirected
// into a file; the command therefore has no color option.
package glossary
