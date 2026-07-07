package model

import "github.com/thenativeweb/esdm/ast"

// ValueObjectView is the typed view over an ESDM
// document whose kind is "value-object".
type ValueObjectView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (v ValueObjectView) Scope() ast.Node {
	return v.Field("scope")
}

// Schema returns the schema field.
func (v ValueObjectView) Schema() ast.Node {
	return v.Field("schema")
}

// Invariants returns the invariants field.
func (v ValueObjectView) Invariants() ast.Node {
	return v.Field("invariants")
}
