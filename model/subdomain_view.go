package model

import "github.com/thenativeweb/esdm/ast"

// SubdomainView is the typed view over an ESDM document
// whose kind is "subdomain".
type SubdomainView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (s SubdomainView) Scope() ast.Node {
	return s.Field("scope")
}

// Type returns the type field (core, supporting, generic).
func (s SubdomainView) Type() ast.Node {
	return s.Field("type")
}

// BoundedContexts returns the boundedContexts field.
func (s SubdomainView) BoundedContexts() ast.Node {
	return s.Field("boundedContexts")
}
