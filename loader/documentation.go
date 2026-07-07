// Package loader discovers the set of ESDM source files
// that belong to a model. Given a directory, it walks the
// directory tree recursively and collects every file whose
// name ends in ".esdm.yaml".
//
// The loader intentionally does not produce diagnostics -
// it returns plain Go errors for failures such as
// unreadable directories or empty model sets. Higher-level
// callers translate these into whatever user-facing form
// is appropriate.
package loader
