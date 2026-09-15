// Package view implements the `esdm view` command. It
// renders a hierarchical, opinionated summary of an ESDM
// model: domain -> subdomains -> bounded contexts ->
// consistency units (aggregates, DCBs), free-standing
// events, read models, queries, entities, value objects,
// domain services and actors -> commands / events, plus
// the integration layer (process managers, event
// handlers, policies, external systems, context mappings)
// and the extension documents: domain stories under their
// domain, features under the consistency unit they are
// about.
//
// Placement follows one rule: an element sits at the
// position its own `scope` names, and relationships
// appear as annotations on the node, never as placement.
// An event is rendered under the aggregate its scope
// names, or directly under the bounded context when its
// scope names no aggregate; the command that publishes it
// is an annotation, not its parent. A feature is rendered
// under the aggregate, DCB, process manager, or read model
// its scope names, not under the domain its scope starts
// with. The one exception is
// `context-mapping`, which has no scope by design because
// its endpoints may straddle domains; it is rendered under
// every domain it touches. The rule keeps every document
// reachable through exactly the path its scope spells out,
// which is what path narrowing and the reference notation
// rely on.
//
// The command runs the resolver and rule pipeline
// implicitly and inline-marks any diagnostic-affected
// node with a severity glyph (warning / error). Without
// arguments it summarizes the whole model in the current
// directory; with a path (e.g. sample/context-one/widget)
// it narrows to the matching subtree, rendering every
// node that matches the path when same-named siblings of
// different kinds exist. The --with-details flag swaps
// the compact, skeleton-only format for an extended one
// that includes data schemas, invariants and rule prose.
package view
