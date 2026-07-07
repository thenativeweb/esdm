package model

// Model is the resolved index of an ESDM source tree. It
// groups entities by kind and keys them by name, so rules
// can query the model by type and reach each entity in
// constant time.
//
// The index is populated by the resolver package. Every
// kind enumerated by the core schema has its own map at
// the top level. Entities defined by extension schemas
// live under the Extensions namespace so that kind names
// can never collide across extensions, or between an
// extension and the core.
type Model struct {
	Domains                      map[string]DomainView
	Subdomains                   map[string]SubdomainView
	BoundedContexts              map[string]BoundedContextView
	ContextMappings              map[string]ContextMappingView
	Aggregates                   map[string]AggregateView
	DynamicConsistencyBoundaries map[string]DynamicConsistencyBoundaryView
	Commands                     map[string]CommandView
	Events                       map[string]EventView
	EventHandlers                map[string]EventHandlerView
	Policies                     map[string]PolicyView
	ProcessManagers              map[string]ProcessManagerView
	ReadModels                   map[string]ReadModelView
	Queries                      map[string]QueryView
	Entities                     map[string]EntityView
	ValueObjects                 map[string]ValueObjectView
	DomainServices               map[string]DomainServiceView
	Actors                       map[string]ActorView
	ExternalSystems              map[string]ExternalSystemView

	Extensions Extensions
}

// Extensions holds the entity indices produced by
// extension schemas. Each extension has its own named
// field; adding a new extension means adding a new field
// here plus a dispatch case in the resolver.
type Extensions struct {
	DomainStorytelling DomainStorytellingIndex
	GivenWhenThen      GivenWhenThenIndex
}

// DomainStorytellingIndex groups every entity kind the
// domain-storytelling extension schema defines.
// Currently only Stories (kind: domain-story).
type DomainStorytellingIndex struct {
	Stories map[string]DomainStoryView
}

// GivenWhenThenIndex groups every entity kind the
// given-when-then extension schema defines. Currently
// only Features (kind: feature).
type GivenWhenThenIndex struct {
	Features map[string]FeatureView
}

// NewModel returns an empty Model with all its maps
// initialized, ready to be populated by the resolver.
func NewModel() *Model {
	return &Model{
		Domains:                      make(map[string]DomainView),
		Subdomains:                   make(map[string]SubdomainView),
		BoundedContexts:              make(map[string]BoundedContextView),
		ContextMappings:              make(map[string]ContextMappingView),
		Aggregates:                   make(map[string]AggregateView),
		DynamicConsistencyBoundaries: make(map[string]DynamicConsistencyBoundaryView),
		Commands:                     make(map[string]CommandView),
		Events:                       make(map[string]EventView),
		EventHandlers:                make(map[string]EventHandlerView),
		Policies:                     make(map[string]PolicyView),
		ProcessManagers:              make(map[string]ProcessManagerView),
		ReadModels:                   make(map[string]ReadModelView),
		Queries:                      make(map[string]QueryView),
		Entities:                     make(map[string]EntityView),
		ValueObjects:                 make(map[string]ValueObjectView),
		DomainServices:               make(map[string]DomainServiceView),
		Actors:                       make(map[string]ActorView),
		ExternalSystems:              make(map[string]ExternalSystemView),
		Extensions: Extensions{
			DomainStorytelling: DomainStorytellingIndex{
				Stories: make(map[string]DomainStoryView),
			},
			GivenWhenThen: GivenWhenThenIndex{
				Features: make(map[string]FeatureView),
			},
		},
	}
}
