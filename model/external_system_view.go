package model

import "github.com/thenativeweb/esdm/ast"

// ExternalSystemView is the typed view over an ESDM
// document whose kind is "external-system".
type ExternalSystemView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (e ExternalSystemView) Scope() ast.Node {
	return e.Field("scope")
}

// Direction returns the direction field.
func (e ExternalSystemView) Direction() ast.Node {
	return e.Field("direction")
}

// Category returns the category field.
func (e ExternalSystemView) Category() ast.Node {
	return e.Field("category")
}

// Capabilities returns the capabilities field.
func (e ExternalSystemView) Capabilities() ast.Node {
	return e.Field("capabilities")
}
