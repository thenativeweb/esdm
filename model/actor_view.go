package model

import "github.com/thenativeweb/esdm/ast"

// ActorView is the typed view over an ESDM document
// whose kind is "actor".
type ActorView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (a ActorView) Scope() ast.Node {
	return a.Field("scope")
}

// Type returns the type field (human or system).
func (a ActorView) Type() ast.Node {
	return a.Field("type")
}

// Responsibilities returns the responsibilities field.
func (a ActorView) Responsibilities() ast.Node {
	return a.Field("responsibilities")
}

// BackedBy returns the backedBy field.
func (a ActorView) BackedBy() ast.Node {
	return a.Field("backedBy")
}
