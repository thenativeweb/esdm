package model

import "github.com/thenativeweb/esdm/ast"

// BoundedContextView is the typed view over an ESDM
// document whose kind is "bounded-context".
type BoundedContextView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (b BoundedContextView) Scope() ast.Node {
	return b.Field("scope")
}

// UbiquitousLanguage returns the ubiquitousLanguage field.
func (b BoundedContextView) UbiquitousLanguage() ast.Node {
	return b.Field("ubiquitousLanguage")
}
