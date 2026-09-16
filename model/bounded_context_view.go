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

// Language returns the BCP 47 tag of the language the
// ubiquitous language is written in.
func (b BoundedContextView) Language() ast.Node {
	return b.Field("language")
}

// UbiquitousLanguage returns the ubiquitousLanguage field.
func (b BoundedContextView) UbiquitousLanguage() ast.Node {
	return b.Field("ubiquitousLanguage")
}
