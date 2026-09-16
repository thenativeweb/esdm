// Package view implements the `esdm view` command. It
// renders the containment tree of an ESDM model, built by
// the tree package, as a hierarchical, opinionated terminal
// summary with box-drawing connectors.
//
// The command runs the resolver and rule pipeline
// implicitly and inline-marks any diagnostic-affected
// node with a severity glyph (warning / error). The
// complete tree is annotated before it is narrowed, so a
// diagnostic from outside the selected subtree never lands
// on a node inside it. Without arguments the command
// summarizes the whole model in the current directory; with
// a path (e.g. sample/context-one/widget) it narrows to the
// matching subtree, rendering every node that matches the
// path when same-named siblings of different kinds exist.
// The --with-details flag swaps the compact, skeleton-only
// format for an extended one that includes data schemas,
// invariants and rule prose.
package view
