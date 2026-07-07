package model

import "github.com/thenativeweb/esdm/ast"

// EventView is the typed view over an ESDM document whose
// kind is "event". It embeds DocumentViewBase to inherit
// the common fields and adds the event-specific ones.
type EventView struct {
	DocumentViewBase
}

// Scope returns the scope field (either a scopeAggregate
// or a scopeBoundedContext triple).
func (e EventView) Scope() ast.Node {
	return e.Field("scope")
}

// Data returns the data field (a JSON Schema describing
// the event's immutable payload).
func (e EventView) Data() ast.Node {
	return e.Field("data")
}
