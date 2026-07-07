// Package view implements the `esdm view` command. It
// renders a hierarchical, opinionated summary of an ESDM
// model: domain -> subdomains -> bounded contexts ->
// consistency units (aggregates, DCBs) -> commands /
// events, plus the integration layer (process managers,
// context mappings) and any domain-storytelling stories.
//
// The command runs the resolver and rule pipeline
// implicitly and inline-marks any diagnostic-affected
// node with a severity glyph (warning / error). Without
// arguments it summarizes the whole model in the current
// directory; with a path (e.g. sample/context-one/widget)
// it narrows to the matching subtree. The --with-details
// flag swaps the compact, skeleton-only format for an
// extended one that includes data schemas, invariants
// and rule prose.
package view
