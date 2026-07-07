// Package refgen extracts per-kind schema excerpts from the embedded
// ESDM schemas. The Reference markdown pages embed those excerpts via
// pymdownx.snippets, so the schema visible in the docs is always the
// schema the binary actually validates against.
//
// The excerpts are derived, never hand-written: refgen strips the
// schema's bookkeeping fields, inlines every internal reference, and
// pins each snippet to the one kind it describes, so a reader sees the
// resolved shape of a kind directly, without having to chase $refs or
// wade through hosting metadata. The package exposes the full snippet
// set as a single map, which serves two consumers: the refgen command
// under cmd writes the entries to disk, and the sync test in the
// documentation package compares them against the committed files,
// which is what keeps the docs from drifting away from the schemas.
package refgen
