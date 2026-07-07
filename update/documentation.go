// Package update implements the once-a-day version check
// that runs after every esdm command in interactive
// terminals. It fetches a small JSON document from the
// documentation site, compares the version it advertises
// with the running binary's version, and renders a hint
// to stderr when a newer release is available.
//
// The package is deliberately defensive: every failure
// path - network error, malformed JSON, corrupt cache -
// leaves the user's command unaffected. The check never
// returns an error to its caller.
package update
