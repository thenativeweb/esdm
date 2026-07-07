// Package ast provides the wrapper used by the rest of
// the esdm linter to talk to parsed YAML. It offers
// position-aware navigation over the parsed YAML tree and
// keeps the originating file attached to every node, so
// diagnostics can always report a precise location.
// Navigation over a missing path yields a missing node
// instead of panicking, so deep lookups can be chained
// without guarding each intermediate step.
//
// The wrapper is deliberately thin: it does not reify a
// typed domain model. Typed access to specific ESDM
// entities lives in the model package and is built on top
// of this one.
package ast
