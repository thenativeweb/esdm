package model

import "github.com/thenativeweb/esdm/ast"

// AggregateView is the typed view over an ESDM document
// whose kind is "aggregate".
type AggregateView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (a AggregateView) Scope() ast.Node {
	return a.Field("scope")
}

// IdentifiedBy returns the identifiedBy field.
func (a AggregateView) IdentifiedBy() ast.Node {
	return a.Field("identifiedBy")
}

// State returns the state field.
func (a AggregateView) State() ast.Node {
	return a.Field("state")
}

// Invariants returns the invariants field.
func (a AggregateView) Invariants() ast.Node {
	return a.Field("invariants")
}
