package model

import "github.com/thenativeweb/esdm/ast"

// CommandView is the typed view over an ESDM document
// whose kind is "command".
type CommandView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (c CommandView) Scope() ast.Node {
	return c.Field("scope")
}

// Data returns the data field (the command payload schema).
func (c CommandView) Data() ast.Node {
	return c.Field("data")
}

// Publishes returns the publishes field.
func (c CommandView) Publishes() ast.Node {
	return c.Field("publishes")
}

// Actors returns the actors field.
func (c CommandView) Actors() ast.Node {
	return c.Field("actors")
}

// Constraints returns the constraints field.
func (c CommandView) Constraints() ast.Node {
	return c.Field("constraints")
}
