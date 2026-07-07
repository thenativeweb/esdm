package model

import "github.com/thenativeweb/esdm/ast"

// ReadModelView is the typed view over an ESDM document
// whose kind is "read-model".
type ReadModelView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (r ReadModelView) Scope() ast.Node {
	return r.Field("scope")
}

// Projections returns the projections field.
func (r ReadModelView) Projections() ast.Node {
	return r.Field("projections")
}

// Paradigm returns the paradigm field.
func (r ReadModelView) Paradigm() ast.Node {
	return r.Field("paradigm")
}

// Schema returns the schema field.
func (r ReadModelView) Schema() ast.Node {
	return r.Field("schema")
}
