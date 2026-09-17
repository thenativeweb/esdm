// Package refgen derives the reference snippets that the documentation
// site embeds via pymdownx.snippets: per-kind excerpts of the embedded
// ESDM schemas, and the rule entries of the Linter Rules pages.
//
// Both kinds of snippet exist for the same reason: the documentation
// must describe what the binary actually does, and the only way to
// keep that true over time is to derive the description from the
// code. For the schemas, refgen strips the bookkeeping fields, inlines
// every internal reference, and pins each snippet to the one kind it
// describes, so a reader sees the resolved shape of a kind directly.
// For the rules, refgen renders the metadata of every rule in the
// catalog, so the ID, severity, and description on the page are the
// ones the linter ships with; the rule's extension decides which page
// it lands on, its category decides the section.
//
// The package exposes the full snippet set as a single map, which
// serves two consumers: the refgen command under cmd writes the
// entries to disk, and the sync test in the documentation package
// compares them against the committed files, which is what keeps the
// docs from drifting away from the schemas and the rules.
package refgen
