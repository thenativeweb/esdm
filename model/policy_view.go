package model

import "github.com/thenativeweb/esdm/ast"

// PolicyView is the typed view over an ESDM document
// whose kind is "policy".
type PolicyView struct {
	DocumentViewBase
}

// Scope returns the scope field.
func (p PolicyView) Scope() ast.Node {
	return p.Field("scope")
}

// DeliveryGuarantee returns the deliveryGuarantee field.
func (p PolicyView) DeliveryGuarantee() ast.Node {
	return p.Field("deliveryGuarantee")
}

// Idempotency returns the idempotency field.
func (p PolicyView) Idempotency() ast.Node {
	return p.Field("idempotency")
}

// Handles returns the handles field.
func (p PolicyView) Handles() ast.Node {
	return p.Field("handles")
}

// Emits returns the emits field.
func (p PolicyView) Emits() ast.Node {
	return p.Field("emits")
}

// Constraints returns the constraints field.
func (p PolicyView) Constraints() ast.Node {
	return p.Field("constraints")
}
