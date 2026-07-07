package model

import "github.com/thenativeweb/esdm/ast"

// Composite-key builders for the maps held by Model. The
// keys encode an entity's full DDD scope plus its bare
// name, so equally-named entities living in different
// scopes coexist - two `register` commands in different
// aggregates are distinct entities, exactly as DDD
// expects. The same applies to actors, events,
// read-models and every other scoped kind.
//
// Callers should never depend on the textual format of
// these keys; they are an internal index detail. Use the
// Lookup* methods on Model to retrieve entities by their
// scope tuple, and view.Name() to read the bare name.

// ScopeText returns scope.<field> as text or "" when the
// field is missing or non-text.
func ScopeText(scope ast.Node, field string) string {
	v, _ := scope.Field(field).Text()
	return v
}

// DomainKey is the key for a domain. Domains carry no
// scope, so the key is just the name.
func DomainKey(name string) string {
	return name
}

// BoundedContextKey is the key for a bounded context.
// BCs are domain-scoped.
func BoundedContextKey(domain, name string) string {
	return domain + "/" + name
}

// SubdomainKey mirrors BoundedContextKey - subdomains are
// domain-scoped.
func SubdomainKey(domain, name string) string {
	return domain + "/" + name
}

// AggregateKey is the key for an aggregate. Aggregates
// are bounded-context-scoped.
func AggregateKey(domain, boundedContext, name string) string {
	return domain + "/" + boundedContext + "/" + name
}

// DynamicConsistencyBoundaryKey mirrors AggregateKey.
func DynamicConsistencyBoundaryKey(domain, boundedContext, name string) string {
	return domain + "/" + boundedContext + "/" + name
}

// CommandKey is the key for a command. Commands belong to
// either an aggregate or a DCB; both occupy the same
// position in the key.
func CommandKey(domain, boundedContext, parent, name string) string {
	return domain + "/" + boundedContext + "/" + parent + "/" + name
}

// EventKey is the key for an event. Aggregate-bound
// events use the aggregate name as parent; free-standing
// BC-scoped events emitted by DCB-bound commands leave
// the parent slot empty.
func EventKey(domain, boundedContext, aggregate, name string) string {
	return domain + "/" + boundedContext + "/" + aggregate + "/" + name
}

// EventHandlerKey is the key for an event handler. Event
// handlers are domain-scoped.
func EventHandlerKey(domain, name string) string {
	return domain + "/" + name
}

// PolicyKey mirrors EventHandlerKey.
func PolicyKey(domain, name string) string {
	return domain + "/" + name
}

// ProcessManagerKey mirrors EventHandlerKey.
func ProcessManagerKey(domain, name string) string {
	return domain + "/" + name
}

// ReadModelKey is the key for a read model. Read models
// are bounded-context-scoped.
func ReadModelKey(domain, boundedContext, name string) string {
	return domain + "/" + boundedContext + "/" + name
}

// QueryKey mirrors ReadModelKey.
func QueryKey(domain, boundedContext, name string) string {
	return domain + "/" + boundedContext + "/" + name
}

// EntityKey is the key for an entity. Entities are
// bounded-context-scoped - they live within one BC, just
// like value-objects.
func EntityKey(domain, boundedContext, name string) string {
	return domain + "/" + boundedContext + "/" + name
}

// ValueObjectKey mirrors ReadModelKey.
func ValueObjectKey(domain, boundedContext, name string) string {
	return domain + "/" + boundedContext + "/" + name
}

// DomainServiceKey mirrors ReadModelKey.
func DomainServiceKey(domain, boundedContext, name string) string {
	return domain + "/" + boundedContext + "/" + name
}

// ActorKey mirrors ReadModelKey.
func ActorKey(domain, boundedContext, name string) string {
	return domain + "/" + boundedContext + "/" + name
}

// ExternalSystemKey is the key for an external system.
// External systems are domain-scoped.
func ExternalSystemKey(domain, name string) string {
	return domain + "/" + name
}

// ContextMappingKey is the key for a context mapping.
// Context mappings have no scope; their name alone is
// the modeling-convention identifier.
func ContextMappingKey(name string) string {
	return name
}

// KeyForCoreDocument returns the composite key for any
// core-schema document, dispatching on kind. Returns
// ("", false) when the document has no readable name.
// Used by the resolver to populate Model maps and the
// duplicate-name index in one place, consistent with the
// scope semantics each kind carries.
func KeyForCoreDocument(view DocumentViewBase, kind string) (string, bool) {
	name, ok := view.Name().Text()
	if !ok {
		return "", false
	}
	scope := view.Field("scope")
	domain := ScopeText(scope, "domain")
	boundedContext := ScopeText(scope, "boundedContext")

	switch kind {
	case "domain":
		return DomainKey(name), true
	case "subdomain":
		return SubdomainKey(domain, name), true
	case "bounded-context":
		return BoundedContextKey(domain, name), true
	case "context-mapping":
		return ContextMappingKey(name), true
	case "aggregate":
		return AggregateKey(domain, boundedContext, name), true
	case "dynamic-consistency-boundary":
		return DynamicConsistencyBoundaryKey(domain, boundedContext, name), true
	case "command":
		parent := ScopeText(scope, "aggregate")
		if parent == "" {
			parent = ScopeText(scope, "dynamicConsistencyBoundary")
		}
		return CommandKey(domain, boundedContext, parent, name), true
	case "event":
		aggregate := ScopeText(scope, "aggregate")
		return EventKey(domain, boundedContext, aggregate, name), true
	case "event-handler":
		return EventHandlerKey(domain, name), true
	case "policy":
		return PolicyKey(domain, name), true
	case "process-manager":
		return ProcessManagerKey(domain, name), true
	case "read-model":
		return ReadModelKey(domain, boundedContext, name), true
	case "query":
		return QueryKey(domain, boundedContext, name), true
	case "entity":
		return EntityKey(domain, boundedContext, name), true
	case "value-object":
		return ValueObjectKey(domain, boundedContext, name), true
	case "domain-service":
		return DomainServiceKey(domain, boundedContext, name), true
	case "actor":
		return ActorKey(domain, boundedContext, name), true
	case "external-system":
		return ExternalSystemKey(domain, name), true
	}
	return "", false
}

// LookupDomain returns the domain by name and whether it
// exists.
func (m *Model) LookupDomain(name string) (DomainView, bool) {
	v, ok := m.Domains[DomainKey(name)]
	return v, ok
}

// LookupSubdomain returns the subdomain by domain and
// name.
func (m *Model) LookupSubdomain(domain, name string) (SubdomainView, bool) {
	v, ok := m.Subdomains[SubdomainKey(domain, name)]
	return v, ok
}

// LookupBoundedContext returns the bounded context by
// domain and name.
func (m *Model) LookupBoundedContext(domain, name string) (BoundedContextView, bool) {
	v, ok := m.BoundedContexts[BoundedContextKey(domain, name)]
	return v, ok
}

// LookupContextMapping returns the context mapping by
// name. Mappings have no scope.
func (m *Model) LookupContextMapping(name string) (ContextMappingView, bool) {
	v, ok := m.ContextMappings[ContextMappingKey(name)]
	return v, ok
}

// LookupAggregate returns the aggregate by domain, BC,
// and name.
func (m *Model) LookupAggregate(domain, boundedContext, name string) (AggregateView, bool) {
	v, ok := m.Aggregates[AggregateKey(domain, boundedContext, name)]
	return v, ok
}

// LookupDynamicConsistencyBoundary returns the DCB by
// domain, BC, and name.
func (m *Model) LookupDynamicConsistencyBoundary(domain, boundedContext, name string) (DynamicConsistencyBoundaryView, bool) {
	v, ok := m.DynamicConsistencyBoundaries[DynamicConsistencyBoundaryKey(domain, boundedContext, name)]
	return v, ok
}

// LookupCommand returns the command by domain, BC, parent
// (aggregate or DCB name), and bare command name.
func (m *Model) LookupCommand(domain, boundedContext, parent, name string) (CommandView, bool) {
	v, ok := m.Commands[CommandKey(domain, boundedContext, parent, name)]
	return v, ok
}

// LookupEvent returns the event by domain, BC, aggregate
// (empty for free-standing BC-scoped events), and bare
// event name.
func (m *Model) LookupEvent(domain, boundedContext, aggregate, name string) (EventView, bool) {
	v, ok := m.Events[EventKey(domain, boundedContext, aggregate, name)]
	return v, ok
}

// LookupEventHandler returns the event handler by domain
// and name.
func (m *Model) LookupEventHandler(domain, name string) (EventHandlerView, bool) {
	v, ok := m.EventHandlers[EventHandlerKey(domain, name)]
	return v, ok
}

// LookupPolicy returns the policy by domain and name.
func (m *Model) LookupPolicy(domain, name string) (PolicyView, bool) {
	v, ok := m.Policies[PolicyKey(domain, name)]
	return v, ok
}

// LookupProcessManager returns the process manager by
// domain and name.
func (m *Model) LookupProcessManager(domain, name string) (ProcessManagerView, bool) {
	v, ok := m.ProcessManagers[ProcessManagerKey(domain, name)]
	return v, ok
}

// LookupReadModel returns the read model by domain, BC,
// and name.
func (m *Model) LookupReadModel(domain, boundedContext, name string) (ReadModelView, bool) {
	v, ok := m.ReadModels[ReadModelKey(domain, boundedContext, name)]
	return v, ok
}

// LookupQuery returns the query by domain, BC, and name.
func (m *Model) LookupQuery(domain, boundedContext, name string) (QueryView, bool) {
	v, ok := m.Queries[QueryKey(domain, boundedContext, name)]
	return v, ok
}

// LookupEntity returns the entity by domain, BC, and
// name.
func (m *Model) LookupEntity(domain, boundedContext, name string) (EntityView, bool) {
	v, ok := m.Entities[EntityKey(domain, boundedContext, name)]
	return v, ok
}

// LookupValueObject returns the value object by domain,
// BC, and name.
func (m *Model) LookupValueObject(domain, boundedContext, name string) (ValueObjectView, bool) {
	v, ok := m.ValueObjects[ValueObjectKey(domain, boundedContext, name)]
	return v, ok
}

// LookupDomainService returns the domain service by
// domain, BC, and name.
func (m *Model) LookupDomainService(domain, boundedContext, name string) (DomainServiceView, bool) {
	v, ok := m.DomainServices[DomainServiceKey(domain, boundedContext, name)]
	return v, ok
}

// LookupActor returns the actor by domain, BC, and name.
func (m *Model) LookupActor(domain, boundedContext, name string) (ActorView, bool) {
	v, ok := m.Actors[ActorKey(domain, boundedContext, name)]
	return v, ok
}

// LookupExternalSystem returns the external system by
// domain and name.
func (m *Model) LookupExternalSystem(domain, name string) (ExternalSystemView, bool) {
	v, ok := m.ExternalSystems[ExternalSystemKey(domain, name)]
	return v, ok
}

// FindEventsByName returns every event in the model
// whose bare name equals the given name, regardless of
// scope. Used when a reference's scope information is
// incomplete (e.g. when looking up a publishing target
// to detect cross-aggregate emission) or when generating
// "did you mean?" suggestions across the full set.
func (m *Model) FindEventsByName(name string) []EventView {
	var out []EventView
	for _, v := range m.Events {
		if n, ok := v.Name().Text(); ok && n == name {
			out = append(out, v)
		}
	}
	return out
}

// FindCommandsByName mirrors FindEventsByName for
// commands.
func (m *Model) FindCommandsByName(name string) []CommandView {
	var out []CommandView
	for _, v := range m.Commands {
		if n, ok := v.Name().Text(); ok && n == name {
			out = append(out, v)
		}
	}
	return out
}
