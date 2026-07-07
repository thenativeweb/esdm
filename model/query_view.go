package model

import "github.com/thenativeweb/esdm/ast"

// QueryView is the typed view over an ESDM document
// whose kind is "query".
type QueryView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (q QueryView) Scope() ast.Node {
	return q.Field("scope")
}

// ReadModel returns the readModel field.
func (q QueryView) ReadModel() ast.Node {
	return q.Field("readModel")
}

// Paradigm returns the paradigm field.
func (q QueryView) Paradigm() ast.Node {
	return q.Field("paradigm")
}

// Result returns the result field.
func (q QueryView) Result() ast.Node {
	return q.Field("result")
}

// Parameters returns the parameters field.
func (q QueryView) Parameters() ast.Node {
	return q.Field("parameters")
}

// Actors returns the actors field.
func (q QueryView) Actors() ast.Node {
	return q.Field("actors")
}

// Constraints returns the constraints field.
func (q QueryView) Constraints() ast.Node {
	return q.Field("constraints")
}
