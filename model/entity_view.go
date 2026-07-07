package model

import "github.com/thenativeweb/esdm/ast"

// EntityView is the typed view over an ESDM document
// whose kind is "entity" - the DDD modeling element that
// carries identity without a consistency container of
// its own.
type EntityView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (e EntityView) Scope() ast.Node {
	return e.Field("scope")
}

// Schema returns the schema field.
func (e EntityView) Schema() ast.Node {
	return e.Field("schema")
}

// IdentifiedBy returns the identifiedBy field.
func (e EntityView) IdentifiedBy() ast.Node {
	return e.Field("identifiedBy")
}

// Invariants returns the invariants field.
func (e EntityView) Invariants() ast.Node {
	return e.Field("invariants")
}
