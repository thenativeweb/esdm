package model

import "github.com/thenativeweb/esdm/ast"

// EventHandlerView is the typed view over an ESDM
// document whose kind is "event-handler".
type EventHandlerView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (e EventHandlerView) Scope() ast.Node {
	return e.Field("scope")
}

// DeliveryGuarantee returns the deliveryGuarantee field.
func (e EventHandlerView) DeliveryGuarantee() ast.Node {
	return e.Field("deliveryGuarantee")
}

// Idempotency returns the idempotency field.
func (e EventHandlerView) Idempotency() ast.Node {
	return e.Field("idempotency")
}

// Handles returns the handles field.
func (e EventHandlerView) Handles() ast.Node {
	return e.Field("handles")
}

// Constraints returns the constraints field.
func (e EventHandlerView) Constraints() ast.Node {
	return e.Field("constraints")
}

// SideEffects returns the sideEffects field.
func (e EventHandlerView) SideEffects() ast.Node {
	return e.Field("sideEffects")
}
