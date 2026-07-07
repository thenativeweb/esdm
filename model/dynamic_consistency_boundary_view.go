package model

import "github.com/thenativeweb/esdm/ast"

// DynamicConsistencyBoundaryView is the typed view over
// an ESDM document whose kind is
// "dynamic-consistency-boundary".
type DynamicConsistencyBoundaryView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (d DynamicConsistencyBoundaryView) Scope() ast.Node {
	return d.Field("scope")
}

// IdentifiedBy returns the identifiedBy field.
func (d DynamicConsistencyBoundaryView) IdentifiedBy() ast.Node {
	return d.Field("identifiedBy")
}

// Consults returns the consults field.
func (d DynamicConsistencyBoundaryView) Consults() ast.Node {
	return d.Field("consults")
}

// Invariants returns the invariants field.
func (d DynamicConsistencyBoundaryView) Invariants() ast.Node {
	return d.Field("invariants")
}
