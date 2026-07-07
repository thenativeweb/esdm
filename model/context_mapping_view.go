package model

import "github.com/thenativeweb/esdm/ast"

// ContextMappingView is the typed view over an ESDM
// document whose kind is "context-mapping". The schema
// uses `oneOf` to discriminate by the `type` field;
// every possible per-pattern field is exposed here, and
// callers query each via Exists() to find out which
// pattern is in use.
type ContextMappingView struct {
	DocumentViewBase
}

// Type returns the type field (the mapping pattern).
func (c ContextMappingView) Type() ast.Node {
	return c.Field("type")
}

// Customer returns the customer endpoint
// (customer-supplier mappings only).
func (c ContextMappingView) Customer() ast.Node {
	return c.Field("customer")
}

// Supplier returns the supplier endpoint
// (customer-supplier mappings only).
func (c ContextMappingView) Supplier() ast.Node {
	return c.Field("supplier")
}

// Conformist returns the conformist endpoint
// (conformist mappings only).
func (c ContextMappingView) Conformist() ast.Node {
	return c.Field("conformist")
}

// Upstream returns the upstream endpoint (conformist and
// anti-corruption-layer mappings).
func (c ContextMappingView) Upstream() ast.Node {
	return c.Field("upstream")
}

// Downstream returns the downstream endpoint
// (anti-corruption-layer mappings only).
func (c ContextMappingView) Downstream() ast.Node {
	return c.Field("downstream")
}

// Host returns the host endpoint (open-host-service
// mappings only).
func (c ContextMappingView) Host() ast.Node {
	return c.Field("host")
}

// Consumer returns the consumer endpoint
// (open-host-service and published-language mappings).
func (c ContextMappingView) Consumer() ast.Node {
	return c.Field("consumer")
}

// Publisher returns the publisher endpoint
// (published-language mappings only).
func (c ContextMappingView) Publisher() ast.Node {
	return c.Field("publisher")
}

// Participants returns the participants list (symmetric
// mappings: shared-kernel, partnership, separate-ways).
func (c ContextMappingView) Participants() ast.Node {
	return c.Field("participants")
}
