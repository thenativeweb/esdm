package model

import "github.com/thenativeweb/esdm/ast"

// DomainServiceView is the typed view over an ESDM
// document whose kind is "domain-service".
type DomainServiceView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (d DomainServiceView) Scope() ast.Node {
	return d.Field("scope")
}

// Functions returns the functions field.
func (d DomainServiceView) Functions() ast.Node {
	return d.Field("functions")
}
