package model

import "github.com/thenativeweb/esdm/ast"

// ProcessManagerView is the typed view over an ESDM
// document whose kind is "process-manager".
type ProcessManagerView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (p ProcessManagerView) Scope() ast.Node {
	return p.Field("scope")
}

// DeliveryGuarantee returns the deliveryGuarantee field.
func (p ProcessManagerView) DeliveryGuarantee() ast.Node {
	return p.Field("deliveryGuarantee")
}

// Idempotency returns the idempotency field.
func (p ProcessManagerView) Idempotency() ast.Node {
	return p.Field("idempotency")
}

// CorrelatedBy returns the correlatedBy field.
func (p ProcessManagerView) CorrelatedBy() ast.Node {
	return p.Field("correlatedBy")
}

// State returns the state field.
func (p ProcessManagerView) State() ast.Node {
	return p.Field("state")
}

// Invariants returns the invariants field.
func (p ProcessManagerView) Invariants() ast.Node {
	return p.Field("invariants")
}

// Constraints returns the constraints field.
func (p ProcessManagerView) Constraints() ast.Node {
	return p.Field("constraints")
}

// StartsWhen returns the startsWhen field.
func (p ProcessManagerView) StartsWhen() ast.Node {
	return p.Field("startsWhen")
}

// EndsWhen returns the endsWhen field.
func (p ProcessManagerView) EndsWhen() ast.Node {
	return p.Field("endsWhen")
}

// Timers returns the timers field.
func (p ProcessManagerView) Timers() ast.Node {
	return p.Field("timers")
}

// Reactions returns the reactions field.
func (p ProcessManagerView) Reactions() ast.Node {
	return p.Field("reactions")
}
