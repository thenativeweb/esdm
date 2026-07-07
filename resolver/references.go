package resolver

import (
	"fmt"
	"sort"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/hint"
	"github.com/thenativeweb/esdm/model"
)

// resolveReferences walks every entity in the populated
// Model and verifies its cross-references.
//
// Two reference families are checked here:
//
//  1. Scope triples (domain -> boundedContext ->
//     aggregate / dynamicConsistencyBoundary) on every
//     scoped entity. Each present layer must exist in the
//     index and must belong to the parent it claims; on
//     the first failing layer the check stops to avoid
//     cascading noise.
//
//  2. Bare-name references that piggyback on the
//     surrounding entity's scope to fix their resolution
//     context: subdomain.boundedContexts,
//     actor.backedBy, query.readModel, query.actors,
//     command.publishes, command.actors. Each item is
//     looked up by name within the appropriate scope.
//
// All failures emit esdm/structure/unresolved-reference
// diagnostics; close-but-not-equal names get a
// Levenshtein-based "did you mean?" hint as a Related
// entry.
func resolveReferences(m *model.Model) []diag.Diagnostic {
	var out []diag.Diagnostic

	out = append(out, checkScopes(m, m.Subdomains)...)
	out = append(out, checkScopes(m, m.BoundedContexts)...)
	out = append(out, checkScopes(m, m.Aggregates)...)
	out = append(out, checkScopes(m, m.DynamicConsistencyBoundaries)...)
	out = append(out, checkScopes(m, m.Commands)...)
	out = append(out, checkScopes(m, m.Events)...)
	out = append(out, checkScopes(m, m.EventHandlers)...)
	out = append(out, checkScopes(m, m.Policies)...)
	out = append(out, checkScopes(m, m.ProcessManagers)...)
	out = append(out, checkScopes(m, m.ReadModels)...)
	out = append(out, checkScopes(m, m.Queries)...)
	out = append(out, checkScopes(m, m.Entities)...)
	out = append(out, checkScopes(m, m.ValueObjects)...)
	out = append(out, checkScopes(m, m.DomainServices)...)
	out = append(out, checkScopes(m, m.Actors)...)
	out = append(out, checkScopes(m, m.ExternalSystems)...)

	for _, e := range m.Subdomains {
		out = append(out, checkSubdomainBoundedContexts(m, e)...)
	}
	for _, e := range m.Actors {
		out = append(out, checkActorBackedBy(m, e)...)
	}
	for _, e := range m.Queries {
		out = append(out, checkQueryReadModel(m, e)...)
		out = append(out, checkQueryActors(m, e)...)
	}
	for _, e := range m.Commands {
		out = append(out, checkCommandPublishes(m, e)...)
		out = append(out, checkCommandActors(m, e)...)
	}

	for _, e := range m.EventHandlers {
		domain := scopeField(e.Scope(), "domain")
		for _, ref := range e.Handles().Seq() {
			out = append(out, checkEventReference(m, ref, domain)...)
		}
	}
	for _, e := range m.Policies {
		domain := scopeField(e.Scope(), "domain")
		for _, ref := range e.Handles().Seq() {
			out = append(out, checkEventReference(m, ref, domain)...)
		}
		for _, ref := range e.Emits().Seq() {
			out = append(out, checkCommandReference(m, ref, domain)...)
		}
	}
	for _, e := range m.ProcessManagers {
		domain := scopeField(e.Scope(), "domain")
		for _, ref := range e.StartsWhen().Seq() {
			out = append(out, checkEventReference(m, ref, domain)...)
		}
		for _, reaction := range e.Reactions().Seq() {
			when := reaction.Field("when")
			if when.HasField("event") {
				out = append(out, checkEventReference(m, when, domain)...)
			}
			for _, ref := range reaction.Field("emits").Seq() {
				out = append(out, checkCommandReference(m, ref, domain)...)
			}
		}
	}
	for _, e := range m.DynamicConsistencyBoundaries {
		domain := scopeField(e.Scope(), "domain")
		for _, ref := range e.Consults().Seq() {
			out = append(out, checkEventReference(m, ref, domain)...)
		}
	}
	for _, e := range m.ReadModels {
		domain := scopeField(e.Scope(), "domain")
		for _, ref := range e.Projections().Seq() {
			out = append(out, checkEventReference(m, ref, domain)...)
		}
	}

	for _, e := range m.ContextMappings {
		out = append(out, checkMappingEndpoint(m, e.Customer())...)
		out = append(out, checkMappingEndpoint(m, e.Supplier())...)
		out = append(out, checkMappingEndpoint(m, e.Conformist())...)
		out = append(out, checkMappingEndpoint(m, e.Upstream())...)
		out = append(out, checkMappingEndpoint(m, e.Downstream())...)
		out = append(out, checkMappingEndpoint(m, e.Host())...)
		out = append(out, checkMappingEndpoint(m, e.Consumer())...)
		out = append(out, checkMappingEndpoint(m, e.Publisher())...)
		for _, p := range e.Participants().Seq() {
			out = append(out, checkMappingEndpoint(m, p)...)
		}
	}

	for _, e := range m.Aggregates {
		out = append(out, checkAggregateIdentifiedBy(e)...)
	}
	for _, e := range m.ProcessManagers {
		out = append(out, checkProcessManagerTimers(e)...)
		out = append(out, checkProcessManagerCorrelation(m, e)...)
	}
	for _, e := range m.DynamicConsistencyBoundaries {
		out = append(out, checkDCBIdentifiedBy(m, e)...)
	}

	out = append(out, checkScopes(m, m.Extensions.DomainStorytelling.Stories)...)

	return out
}

// checkScopes runs checkScope over every entity in a
// model map. It captures the loop that would otherwise be
// repeated once per scoped kind, so resolveReferences can
// read as a flat list of what gets checked rather than
// how.
func checkScopes[V interface{ Scope() ast.Node }](m *model.Model, entities map[string]V) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, entity := range entities {
		out = append(out, checkScope(m, entity.Scope())...)
	}
	return out
}

// checkScope inspects the four possible layers of a
// scope object - domain, boundedContext, and either
// aggregate or dynamicConsistencyBoundary - and verifies
// every present layer resolves correctly.
func checkScope(m *model.Model, scope ast.Node) []diag.Diagnostic {
	if !scope.Exists() {
		return nil
	}

	domainNode := scope.Field("domain")
	domainName, ok := domainNode.Text()
	if !ok {
		return nil
	}
	if _, exists := m.LookupDomain(domainName); !exists {
		return []diag.Diagnostic{
			unresolvedReferenceDiag(domainNode, "domain", domainName, suggestDomain(m, domainName)),
		}
	}

	bcNode := scope.Field("boundedContext")
	if !bcNode.Exists() {
		return nil
	}
	bcName, ok := bcNode.Text()
	if !ok {
		return nil
	}
	if _, exists := m.LookupBoundedContext(domainName, bcName); !exists {
		// Fall back: the BC may exist in another domain.
		// Surface a mismatched-parent diagnostic when so;
		// otherwise an unresolved-reference.
		if other, otherDomain, found := findBoundedContextAnyDomain(m, bcName); found {
			return []diag.Diagnostic{mismatchedParentDiag(bcNode, "bounded-context", bcName, "domain", domainName, otherDomain, other.Name().Location())}
		}
		return []diag.Diagnostic{unresolvedReferenceDiag(bcNode, "bounded-context", bcName, suggestBoundedContext(m, bcName))}
	}

	aggNode := scope.Field("aggregate")
	if aggNode.Exists() {
		return checkAggregateLayer(m, aggNode, domainName, bcName)
	}

	dcbNode := scope.Field("dynamicConsistencyBoundary")
	if dcbNode.Exists() {
		return checkDCBLayer(m, dcbNode, domainName, bcName)
	}

	return nil
}

func checkAggregateLayer(m *model.Model, aggNode ast.Node, domain, expectedBC string) []diag.Diagnostic {
	aggName, ok := aggNode.Text()
	if !ok {
		return nil
	}
	if _, exists := m.LookupAggregate(domain, expectedBC, aggName); exists {
		return nil
	}
	if other, otherBC, found := findAggregateInDomain(m, domain, aggName); found {
		return []diag.Diagnostic{mismatchedParentDiag(aggNode, "aggregate", aggName, "bounded-context", expectedBC, otherBC, other.Name().Location())}
	}
	return []diag.Diagnostic{unresolvedReferenceDiag(aggNode, "aggregate", aggName, suggestAggregate(m, aggName))}
}

func checkDCBLayer(m *model.Model, dcbNode ast.Node, domain, expectedBC string) []diag.Diagnostic {
	dcbName, ok := dcbNode.Text()
	if !ok {
		return nil
	}
	if _, exists := m.LookupDynamicConsistencyBoundary(domain, expectedBC, dcbName); exists {
		return nil
	}
	if other, otherBC, found := findDCBInDomain(m, domain, dcbName); found {
		return []diag.Diagnostic{mismatchedParentDiag(dcbNode, "dynamic-consistency-boundary", dcbName, "bounded-context", expectedBC, otherBC, other.Name().Location())}
	}
	return []diag.Diagnostic{unresolvedReferenceDiag(dcbNode, "dynamic-consistency-boundary", dcbName, suggestDCB(m, dcbName))}
}

// scopeField is a tiny helper that pulls a string field
// from a scope object and tolerates missing or non-string
// values by returning "".
func scopeField(scope ast.Node, field string) string {
	v, _ := scope.Field(field).Text()
	return v
}

// checkSubdomainBoundedContexts verifies every entry of
// subdomain.boundedContexts resolves to a BoundedContext
// in the subdomain's domain.
func checkSubdomainBoundedContexts(m *model.Model, sub model.SubdomainView) []diag.Diagnostic {
	expectedDomain := scopeField(sub.Scope(), "domain")
	if expectedDomain == "" {
		return nil
	}

	var out []diag.Diagnostic
	for _, item := range sub.BoundedContexts().Seq() {
		name, ok := item.Text()
		if !ok {
			continue
		}

		if _, exists := m.LookupBoundedContext(expectedDomain, name); exists {
			continue
		}
		if other, otherDomain, found := findBoundedContextAnyDomain(m, name); found {
			out = append(out, mismatchedParentDiag(item, "bounded-context", name, "domain", expectedDomain, otherDomain, other.Name().Location()))
			continue
		}
		out = append(out, unresolvedReferenceDiag(item, "bounded-context", name, suggestBoundedContext(m, name)))
	}
	return out
}

// checkActorBackedBy verifies every entry of
// actor.backedBy resolves to an ExternalSystem in the
// actor's domain.
func checkActorBackedBy(m *model.Model, actor model.ActorView) []diag.Diagnostic {
	expectedDomain := scopeField(actor.Scope(), "domain")
	if expectedDomain == "" {
		return nil
	}

	var out []diag.Diagnostic
	for _, item := range actor.BackedBy().Seq() {
		name, ok := item.Text()
		if !ok {
			continue
		}

		if _, exists := m.LookupExternalSystem(expectedDomain, name); exists {
			continue
		}
		if other, otherDomain, found := findExternalSystemAnyDomain(m, name); found {
			out = append(out, mismatchedParentDiag(item, "external-system", name, "domain", expectedDomain, otherDomain, other.Name().Location()))
			continue
		}
		out = append(out, unresolvedReferenceDiag(item, "external-system", name, suggestExternalSystem(m, name)))
	}
	return out
}

// checkQueryReadModel verifies query.readModel resolves
// to a ReadModel in the query's bounded context.
func checkQueryReadModel(m *model.Model, q model.QueryView) []diag.Diagnostic {
	expectedDomain := scopeField(q.Scope(), "domain")
	expectedBC := scopeField(q.Scope(), "boundedContext")
	if expectedBC == "" {
		return nil
	}

	readModelNode := q.ReadModel()
	if !readModelNode.Exists() {
		return nil
	}
	name, ok := readModelNode.Text()
	if !ok {
		return nil
	}

	if _, exists := m.LookupReadModel(expectedDomain, expectedBC, name); exists {
		return nil
	}
	if other, otherBC, found := findReadModelInDomain(m, expectedDomain, name); found {
		return []diag.Diagnostic{mismatchedParentDiag(readModelNode, "read-model", name, "bounded-context", expectedBC, otherBC, other.Name().Location())}
	}
	return []diag.Diagnostic{unresolvedReferenceDiag(readModelNode, "read-model", name, suggestReadModel(m, name))}
}

// checkQueryActors verifies query.actors[i] each resolves
// to an Actor in the query's bounded context.
func checkQueryActors(m *model.Model, q model.QueryView) []diag.Diagnostic {
	expectedDomain := scopeField(q.Scope(), "domain")
	expectedBC := scopeField(q.Scope(), "boundedContext")
	if expectedBC == "" {
		return nil
	}
	return checkActorList(m, q.Actors(), expectedDomain, expectedBC)
}

// checkCommandActors verifies command.actors[i] each
// resolves to an Actor in the command's bounded context.
func checkCommandActors(m *model.Model, cmd model.CommandView) []diag.Diagnostic {
	expectedDomain := scopeField(cmd.Scope(), "domain")
	expectedBC := scopeField(cmd.Scope(), "boundedContext")
	if expectedBC == "" {
		return nil
	}
	return checkActorList(m, cmd.Actors(), expectedDomain, expectedBC)
}

// checkActorList factors out the actor-by-name lookup
// shared by command.actors and query.actors.
func checkActorList(m *model.Model, list ast.Node, expectedDomain, expectedBC string) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, item := range list.Seq() {
		name, ok := item.Text()
		if !ok {
			continue
		}

		if _, exists := m.LookupActor(expectedDomain, expectedBC, name); exists {
			continue
		}
		if other, otherBC, found := findActorInDomain(m, expectedDomain, name); found {
			out = append(out, mismatchedParentDiag(item, "actor", name, "bounded-context", expectedBC, otherBC, other.Name().Location()))
			continue
		}
		out = append(out, unresolvedReferenceDiag(item, "actor", name, suggestActor(m, name)))
	}
	return out
}

// checkCommandPublishes verifies every entry of
// command.publishes resolves to an Event the command may
// emit. Aggregate-bound commands check strictly: the
// event must live in the same aggregate. DCB-bound
// commands check leniently: same bounded context is
// enough, since DCB writers deliberately span aggregates.
func checkCommandPublishes(m *model.Model, cmd model.CommandView) []diag.Diagnostic {
	scope := cmd.Scope()
	expectedDomain := scopeField(scope, "domain")
	expectedBC := scopeField(scope, "boundedContext")
	if expectedBC == "" {
		return nil
	}
	expectedAggregate := scopeField(scope, "aggregate")
	isAggregateBound := expectedAggregate != ""

	var out []diag.Diagnostic
	for _, item := range cmd.Publishes().Seq() {
		name, ok := item.Text()
		if !ok {
			continue
		}

		if isAggregateBound {
			if _, exists := m.LookupEvent(expectedDomain, expectedBC, expectedAggregate, name); exists {
				continue
			}
			// Maybe the event lives in another aggregate of the same BC.
			if other, otherAgg, found := findEventInBC(m, expectedDomain, expectedBC, name); found {
				display := otherAgg
				if display == "" {
					display = "(BC-scoped, no aggregate)"
				}
				out = append(out, mismatchedParentDiag(item, "event", name, "aggregate", expectedAggregate, display, other.Name().Location()))
				continue
			}
			// Not in this BC at all; report unresolved or BC mismatch.
			if other, otherBC, found := findEventAnyBC(m, expectedDomain, name); found {
				out = append(out, mismatchedParentDiag(item, "event", name, "bounded-context", expectedBC, otherBC, other.Name().Location()))
				continue
			}
			out = append(out, unresolvedReferenceDiag(item, "event", name, suggestEvent(m, name)))
			continue
		}

		// DCB-bound: free-standing BC-scoped event accepted.
		if _, exists := m.LookupEvent(expectedDomain, expectedBC, "", name); exists {
			continue
		}
		// Or any aggregate-bound event in the same BC - DCBs deliberately span aggregates.
		if _, _, found := findEventInBC(m, expectedDomain, expectedBC, name); found {
			continue
		}
		if other, otherBC, found := findEventAnyBC(m, expectedDomain, name); found {
			out = append(out, mismatchedParentDiag(item, "event", name, "bounded-context", expectedBC, otherBC, other.Name().Location()))
			continue
		}
		out = append(out, unresolvedReferenceDiag(item, "event", name, suggestEvent(m, name)))
	}
	return out
}

// checkEventReference verifies that an event-reference
// object - used by event-handler.handles,
// policy.handles, process-manager.startsWhen and
// reactions[].when, dcb.consults, and
// read-model.projections - points at an event whose
// scope matches the reference's
// {boundedContext [, aggregate], event} shape.
//
// The reference shape itself is schema-validated; this
// function focuses on the semantic link: does the event
// named by the reference actually live where the
// reference claims it lives? The function accepts
// references enriched with extra prose fields such as
// `criteria` or `rule` - those are ignored here.
//
// `domain` is the consumer's domain - piggybacking on
// the consumer's scope, mirroring how the reference
// itself omits a `domain` field.
func checkEventReference(m *model.Model, ref ast.Node, domain string) []diag.Diagnostic {
	if !ref.Exists() {
		return nil
	}

	eventNode := ref.Field("event")
	eventName, ok := eventNode.Text()
	if !ok {
		return nil
	}

	expectedBC, _ := ref.Field("boundedContext").Text()
	if expectedBC == "" {
		return nil
	}

	aggNode := ref.Field("aggregate")
	if aggNode.Exists() {
		expectedAgg, _ := aggNode.Text()
		if _, exists := m.LookupEvent(domain, expectedBC, expectedAgg, eventName); exists {
			return nil
		}
		if other, otherAgg, found := findEventInBC(m, domain, expectedBC, eventName); found {
			display := otherAgg
			if display == "" {
				display = "(BC-scoped, no aggregate)"
			}
			return []diag.Diagnostic{mismatchedParentDiag(eventNode, "event", eventName, "aggregate", expectedAgg, display, other.Name().Location())}
		}
		if other, otherBC, found := findEventAnyBC(m, domain, eventName); found {
			return []diag.Diagnostic{mismatchedParentDiag(eventNode, "event", eventName, "bounded-context", expectedBC, otherBC, other.Name().Location())}
		}
		return []diag.Diagnostic{unresolvedReferenceDiag(eventNode, "event", eventName, suggestEvent(m, eventName))}
	}

	// BC-only shape: the event must be free-standing
	// (no aggregate) within expectedBC.
	if _, exists := m.LookupEvent(domain, expectedBC, "", eventName); exists {
		return nil
	}
	// Maybe the event has an aggregate - shape mismatch.
	if other, otherAgg, found := findEventInBC(m, domain, expectedBC, eventName); found {
		return []diag.Diagnostic{mismatchedParentDiag(eventNode, "event", eventName, "aggregate", "(free-standing, none)", otherAgg, other.Name().Location())}
	}
	if other, otherBC, found := findEventAnyBC(m, domain, eventName); found {
		return []diag.Diagnostic{mismatchedParentDiag(eventNode, "event", eventName, "bounded-context", expectedBC, otherBC, other.Name().Location())}
	}
	return []diag.Diagnostic{unresolvedReferenceDiag(eventNode, "event", eventName, suggestEvent(m, eventName))}
}

// checkCommandReference verifies that a command-reference
// object - used by policy.emits and
// process-manager.reactions[].emits - points at a
// command whose scope matches the reference's
// {boundedContext, aggregate | dynamicConsistencyBoundary,
// command} shape.
func checkCommandReference(m *model.Model, ref ast.Node, domain string) []diag.Diagnostic {
	if !ref.Exists() {
		return nil
	}

	cmdNode := ref.Field("command")
	cmdName, ok := cmdNode.Text()
	if !ok {
		return nil
	}

	expectedBC, _ := ref.Field("boundedContext").Text()
	if expectedBC == "" {
		return nil
	}

	aggNode := ref.Field("aggregate")
	dcbNode := ref.Field("dynamicConsistencyBoundary")

	switch {
	case aggNode.Exists():
		expectedAgg, _ := aggNode.Text()
		if _, exists := m.LookupCommand(domain, expectedBC, expectedAgg, cmdName); exists {
			return nil
		}
		if other, otherParent, found := findCommandInBC(m, domain, expectedBC, cmdName); found {
			display := otherParent
			if display == "" {
				display = "(DCB-bound, no aggregate)"
			}
			return []diag.Diagnostic{mismatchedParentDiag(cmdNode, "command", cmdName, "aggregate", expectedAgg, display, other.Name().Location())}
		}
		if other, otherBC, found := findCommandAnyBC(m, domain, cmdName); found {
			return []diag.Diagnostic{mismatchedParentDiag(cmdNode, "command", cmdName, "bounded-context", expectedBC, otherBC, other.Name().Location())}
		}
		return []diag.Diagnostic{unresolvedReferenceDiag(cmdNode, "command", cmdName, suggestCommand(m, cmdName))}
	case dcbNode.Exists():
		expectedDCB, _ := dcbNode.Text()
		if _, exists := m.LookupCommand(domain, expectedBC, expectedDCB, cmdName); exists {
			return nil
		}
		if other, otherParent, found := findCommandInBC(m, domain, expectedBC, cmdName); found {
			display := otherParent
			if display == "" {
				display = "(aggregate-bound, no DCB)"
			}
			return []diag.Diagnostic{mismatchedParentDiag(cmdNode, "command", cmdName, "dynamic-consistency-boundary", expectedDCB, display, other.Name().Location())}
		}
		if other, otherBC, found := findCommandAnyBC(m, domain, cmdName); found {
			return []diag.Diagnostic{mismatchedParentDiag(cmdNode, "command", cmdName, "bounded-context", expectedBC, otherBC, other.Name().Location())}
		}
		return []diag.Diagnostic{unresolvedReferenceDiag(cmdNode, "command", cmdName, suggestCommand(m, cmdName))}
	}

	return nil
}

// checkMappingEndpoint verifies that a mappingEndpoint
// object - the oneOf between scopeBoundedContext and
// externalSystemReference used by every endpoint slot in
// context-mapping - resolves to the claimed entity.
//
// The discriminator is the presence of the
// `boundedContext` field (scopeBoundedContext variant)
// or the `externalSystem` field (externalSystemReference
// variant). When neither is present the node is
// malformed per the schema; the parser's structural pass
// already catches that, so we return silently.
func checkMappingEndpoint(m *model.Model, endpoint ast.Node) []diag.Diagnostic {
	if !endpoint.Exists() {
		return nil
	}

	switch {
	case endpoint.HasField("boundedContext"):
		return checkScope(m, endpoint)
	case endpoint.HasField("externalSystem"):
		return checkExternalSystemReference(m, endpoint)
	default:
		return nil
	}
}

// checkExternalSystemReference verifies that an
// externalSystemReference - {domain, externalSystem} -
// resolves: domain must be a known Domain, the external
// system must exist by name, and its own scope.domain
// must equal the reference's domain.
func checkExternalSystemReference(m *model.Model, ref ast.Node) []diag.Diagnostic {
	domainNode := ref.Field("domain")
	domainName, ok := domainNode.Text()
	if !ok {
		return nil
	}
	if _, exists := m.LookupDomain(domainName); !exists {
		return []diag.Diagnostic{unresolvedReferenceDiag(domainNode, "domain", domainName, suggestDomain(m, domainName))}
	}

	externalSystemNode := ref.Field("externalSystem")
	externalSystemName, ok := externalSystemNode.Text()
	if !ok {
		return nil
	}
	if _, exists := m.LookupExternalSystem(domainName, externalSystemName); exists {
		return nil
	}
	if other, otherDomain, found := findExternalSystemAnyDomain(m, externalSystemName); found {
		return []diag.Diagnostic{mismatchedParentDiag(externalSystemNode, "external-system", externalSystemName, "domain", domainName, otherDomain, other.Name().Location())}
	}
	return []diag.Diagnostic{unresolvedReferenceDiag(externalSystemNode, "external-system", externalSystemName, suggestExternalSystem(m, externalSystemName))}
}

// checkAggregateIdentifiedBy verifies that an aggregate
// whose identifier is drawn from its own state declares
// the referenced field in state.properties. Other
// identifiedBy shapes (static, generated) do not carry a
// field reference and are left alone.
func checkAggregateIdentifiedBy(agg model.AggregateView) []diag.Diagnostic {
	ib := agg.IdentifiedBy()
	if source, _ := ib.Field("source").Text(); source != "state" {
		return nil
	}

	fieldNode := ib.Field("field")
	fieldName, ok := fieldNode.Text()
	if !ok {
		return nil
	}

	state := agg.State()
	if fieldExistsInSchema(state, fieldName) {
		return nil
	}

	return []diag.Diagnostic{unresolvedReferenceDiag(fieldNode, "state field", fieldName, suggestSchemaField(state, fieldName))}
}

// checkProcessManagerTimers verifies that each timer with
// an `at` reference points at a field declared in the
// process manager's own state schema. Timers with an
// `after` relative delay do not carry a field reference
// and are left alone.
func checkProcessManagerTimers(pm model.ProcessManagerView) []diag.Diagnostic {
	state := pm.State()

	var out []diag.Diagnostic
	for _, timer := range pm.Timers().Seq() {
		atNode := timer.Field("at")
		if !atNode.Exists() {
			continue
		}
		fieldName, ok := atNode.Text()
		if !ok {
			continue
		}
		if fieldExistsInSchema(state, fieldName) {
			continue
		}

		out = append(out, unresolvedReferenceDiag(atNode, "state field", fieldName, suggestSchemaField(state, fieldName)))
	}
	return out
}

// checkProcessManagerCorrelation verifies that the
// correlation field declared by the process manager
// exists in the data payload of every event the PM
// references - both in startsWhen and in
// reactions[].when. Events not present in the index are
// skipped; those are already flagged by the
// event-reference pass.
func checkProcessManagerCorrelation(m *model.Model, pm model.ProcessManagerView) []diag.Diagnostic {
	fieldNode := pm.CorrelatedBy().Field("field")
	fieldName, ok := fieldNode.Text()
	if !ok {
		return nil
	}

	pmDomain := scopeField(pm.Scope(), "domain")

	var out []diag.Diagnostic
	for _, ref := range referencedEventRefs(pm) {
		eventView, ok := resolveEventByRef(m, ref, pmDomain)
		if !ok {
			continue
		}
		eventName, _ := eventView.Name().Text()
		data := eventView.Data()
		if fieldExistsInSchema(data, fieldName) {
			continue
		}

		d := diag.Diagnostic{
			RuleID:   "esdm/structure/unresolved-reference",
			Severity: diag.SeverityError,
			Message:  fmt.Sprintf("correlation field %q not found in data of event %q", fieldName, eventName),
			Location: fieldNode.Location(),
		}
		if s := suggestSchemaField(data, fieldName); s != nil {
			d.Related = []diag.Related{
				{
					Message:  fmt.Sprintf("did you mean %q?", s.name),
					Location: s.location,
				},
			}
		}
		out = append(out, d)
	}
	return out
}

// checkDCBIdentifiedBy verifies, for every identifiedBy
// item whose source is command-payload, that the
// referenced field exists in the data payload of every
// command that targets this DCB via its scope. Static
// and generated variants carry no field reference.
func checkDCBIdentifiedBy(m *model.Model, dcb model.DynamicConsistencyBoundaryView) []diag.Diagnostic {
	dcbName, ok := dcb.Name().Text()
	if !ok {
		return nil
	}
	dcbDomain := scopeField(dcb.Scope(), "domain")
	dcbBC := scopeField(dcb.Scope(), "boundedContext")

	targetingCommands := commandsTargetingDCB(m, dcbDomain, dcbBC, dcbName)

	var out []diag.Diagnostic
	for _, item := range dcb.IdentifiedBy().Seq() {
		if source, _ := item.Field("source").Text(); source != "command-payload" {
			continue
		}
		fieldNode := item.Field("field")
		fieldName, ok := fieldNode.Text()
		if !ok {
			continue
		}

		for _, entry := range targetingCommands {
			cmdName := entry.name
			data := entry.view.Data()
			if fieldExistsInSchema(data, fieldName) {
				continue
			}

			d := diag.Diagnostic{
				RuleID:   "esdm/structure/unresolved-reference",
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("identifier field %q not found in data of command %q", fieldName, cmdName),
				Location: fieldNode.Location(),
			}
			if s := suggestSchemaField(data, fieldName); s != nil {
				d.Related = []diag.Related{
					{
						Message:  fmt.Sprintf("did you mean %q?", s.name),
						Location: s.location,
					},
				}
			}
			out = append(out, d)
		}
	}
	return out
}

type namedCommand struct {
	name string
	view model.CommandView
}

func commandsTargetingDCB(m *model.Model, domain, bc, dcbName string) []namedCommand {
	var out []namedCommand
	for _, cmd := range m.Commands {
		scope := cmd.Scope()
		if scopeField(scope, "domain") != domain {
			continue
		}
		if scopeField(scope, "boundedContext") != bc {
			continue
		}
		if scopeField(scope, "dynamicConsistencyBoundary") != dcbName {
			continue
		}
		name, _ := cmd.Name().Text()
		out = append(out, namedCommand{name: name, view: cmd})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

// referencedEventRefs returns the eventReference nodes
// of every event a process manager consumes - both in
// startsWhen and in reactions[].when (when the trigger is
// an event, not a timer).
func referencedEventRefs(pm model.ProcessManagerView) []ast.Node {
	out := append([]ast.Node(nil), pm.StartsWhen().Seq()...)
	for _, reaction := range pm.Reactions().Seq() {
		when := reaction.Field("when")
		if when.HasField("event") {
			out = append(out, when)
		}
	}
	return out
}

// resolveEventByRef looks up the event a reference
// points at, applying the same scope semantics
// checkEventReference uses. Returns the view and true on
// success; (zero, false) when the reference does not
// resolve.
func resolveEventByRef(m *model.Model, ref ast.Node, domain string) (model.EventView, bool) {
	eventName, ok := ref.Field("event").Text()
	if !ok {
		return model.EventView{}, false
	}
	bc, ok := ref.Field("boundedContext").Text()
	if !ok {
		return model.EventView{}, false
	}
	if aggNode := ref.Field("aggregate"); aggNode.Exists() {
		agg, _ := aggNode.Text()
		v, exists := m.LookupEvent(domain, bc, agg, eventName)
		return v, exists
	}
	v, exists := m.LookupEvent(domain, bc, "", eventName)
	return v, exists
}

// findBoundedContextAnyDomain searches every BC by bare
// name and returns the first match plus the domain it
// belongs to.
func findBoundedContextAnyDomain(m *model.Model, name string) (model.BoundedContextView, string, bool) {
	for _, v := range m.BoundedContexts {
		if n, _ := v.Name().Text(); n == name {
			domain, _ := v.Scope().Field("domain").Text()
			return v, domain, true
		}
	}
	return model.BoundedContextView{}, "", false
}

// findAggregateInDomain returns the first aggregate with
// the given bare name in the given domain (any BC).
func findAggregateInDomain(m *model.Model, domain, name string) (model.AggregateView, string, bool) {
	for _, v := range m.Aggregates {
		if n, _ := v.Name().Text(); n != name {
			continue
		}
		if d, _ := v.Scope().Field("domain").Text(); d != domain {
			continue
		}
		bc, _ := v.Scope().Field("boundedContext").Text()
		return v, bc, true
	}
	return model.AggregateView{}, "", false
}

// findDCBInDomain mirrors findAggregateInDomain for
// dynamic-consistency-boundaries.
func findDCBInDomain(m *model.Model, domain, name string) (model.DynamicConsistencyBoundaryView, string, bool) {
	for _, v := range m.DynamicConsistencyBoundaries {
		if n, _ := v.Name().Text(); n != name {
			continue
		}
		if d, _ := v.Scope().Field("domain").Text(); d != domain {
			continue
		}
		bc, _ := v.Scope().Field("boundedContext").Text()
		return v, bc, true
	}
	return model.DynamicConsistencyBoundaryView{}, "", false
}

// findExternalSystemAnyDomain returns the first external
// system with the given bare name and the domain it
// lives in.
func findExternalSystemAnyDomain(m *model.Model, name string) (model.ExternalSystemView, string, bool) {
	for _, v := range m.ExternalSystems {
		if n, _ := v.Name().Text(); n == name {
			domain, _ := v.Scope().Field("domain").Text()
			return v, domain, true
		}
	}
	return model.ExternalSystemView{}, "", false
}

// findReadModelInDomain returns the first read model
// with the given bare name in the given domain (any BC).
func findReadModelInDomain(m *model.Model, domain, name string) (model.ReadModelView, string, bool) {
	for _, v := range m.ReadModels {
		if n, _ := v.Name().Text(); n != name {
			continue
		}
		if d, _ := v.Scope().Field("domain").Text(); d != domain {
			continue
		}
		bc, _ := v.Scope().Field("boundedContext").Text()
		return v, bc, true
	}
	return model.ReadModelView{}, "", false
}

// findActorInDomain mirrors findReadModelInDomain for
// actors.
func findActorInDomain(m *model.Model, domain, name string) (model.ActorView, string, bool) {
	for _, v := range m.Actors {
		if n, _ := v.Name().Text(); n != name {
			continue
		}
		if d, _ := v.Scope().Field("domain").Text(); d != domain {
			continue
		}
		bc, _ := v.Scope().Field("boundedContext").Text()
		return v, bc, true
	}
	return model.ActorView{}, "", false
}

// findEventInBC returns the first event in the given
// (domain, BC) with the given bare name and its
// aggregate ("" for free-standing).
func findEventInBC(m *model.Model, domain, bc, name string) (model.EventView, string, bool) {
	for _, v := range m.Events {
		if n, _ := v.Name().Text(); n != name {
			continue
		}
		scope := v.Scope()
		if d, _ := scope.Field("domain").Text(); d != domain {
			continue
		}
		if b, _ := scope.Field("boundedContext").Text(); b != bc {
			continue
		}
		agg, _ := scope.Field("aggregate").Text()
		return v, agg, true
	}
	return model.EventView{}, "", false
}

// findEventAnyBC returns the first event in the given
// domain with the given bare name and the BC it lives
// in.
func findEventAnyBC(m *model.Model, domain, name string) (model.EventView, string, bool) {
	for _, v := range m.Events {
		if n, _ := v.Name().Text(); n != name {
			continue
		}
		scope := v.Scope()
		if d, _ := scope.Field("domain").Text(); d != domain {
			continue
		}
		bc, _ := scope.Field("boundedContext").Text()
		return v, bc, true
	}
	return model.EventView{}, "", false
}

// findCommandInBC returns the first command in the given
// (domain, BC) with the given bare name and its parent
// (aggregate or DCB name; "" if neither is set).
func findCommandInBC(m *model.Model, domain, bc, name string) (model.CommandView, string, bool) {
	for _, v := range m.Commands {
		if n, _ := v.Name().Text(); n != name {
			continue
		}
		scope := v.Scope()
		if d, _ := scope.Field("domain").Text(); d != domain {
			continue
		}
		if b, _ := scope.Field("boundedContext").Text(); b != bc {
			continue
		}
		parent, _ := scope.Field("aggregate").Text()
		if parent == "" {
			parent, _ = scope.Field("dynamicConsistencyBoundary").Text()
		}
		return v, parent, true
	}
	return model.CommandView{}, "", false
}

// findCommandAnyBC mirrors findEventAnyBC for commands.
func findCommandAnyBC(m *model.Model, domain, name string) (model.CommandView, string, bool) {
	for _, v := range m.Commands {
		if n, _ := v.Name().Text(); n != name {
			continue
		}
		scope := v.Scope()
		if d, _ := scope.Field("domain").Text(); d != domain {
			continue
		}
		bc, _ := scope.Field("boundedContext").Text()
		return v, bc, true
	}
	return model.CommandView{}, "", false
}

// fieldExistsInSchema reports whether a JSON-Schema
// fragment (a mapping node with a `properties` child
// that is itself a mapping) declares a property of the
// given name.
func fieldExistsInSchema(schemaNode ast.Node, field string) bool {
	return schemaNode.Field("properties").HasField(field)
}

// suggestSchemaField returns a Levenshtein-close
// property name from the given JSON-Schema fragment's
// properties, or nil when no close-enough candidate
// exists.
func suggestSchemaField(schemaNode ast.Node, name string) *suggestion {
	props := schemaNode.Field("properties")
	if !props.Exists() {
		return nil
	}

	candidates := make([]string, 0)
	for _, entry := range props.Entries() {
		if k, ok := entry.Key.Text(); ok {
			candidates = append(candidates, k)
		}
	}
	sort.Strings(candidates)

	best, ok := hint.Best(name, candidates)
	if !ok {
		return nil
	}
	return &suggestion{
		name:     best,
		location: props.Field(best).Location(),
	}
}

// suggestion holds a single "did you mean?" hint: the
// name of the closest match plus the location of its
// definition, so the diagnostic can point the reader
// straight at the corrected entity.
type suggestion struct {
	name     string
	location diag.Location
}

// suggest returns the best Levenshtein match for name
// among candidates, provided it falls within the
// suggestion threshold. It returns nil otherwise so
// callers can pass the result directly into the
// diagnostic builders.
func suggest(name string, candidates []candidate) *suggestion {
	names := make([]string, 0, len(candidates))
	seen := make(map[string]bool)
	locByName := make(map[string]diag.Location)
	for _, c := range candidates {
		if seen[c.name] {
			continue
		}
		seen[c.name] = true
		names = append(names, c.name)
		locByName[c.name] = c.location
	}
	sort.Strings(names)

	best, ok := hint.Best(name, names)
	if !ok {
		return nil
	}
	return &suggestion{name: best, location: locByName[best]}
}

// candidate is a (bare name, location) tuple used for
// suggestion building. Several Model views may share a
// bare name (different scopes); suggest dedupes them so
// the user sees one canonical "did you mean?" entry.
type candidate struct {
	name     string
	location diag.Location
}

// candidatesFromMap turns a Model map into a list of
// (bare name, definition location) tuples. Used to feed
// suggest with the bare-name set the user is most likely
// thinking of.
func candidatesFromMap[V interface {
	Name() ast.Node
}](m map[string]V) []candidate {
	out := make([]candidate, 0, len(m))
	for _, v := range m {
		nameNode := v.Name()
		if name, ok := nameNode.Text(); ok {
			out = append(out, candidate{name: name, location: nameNode.Location()})
		}
	}
	return out
}

// Per-kind suggest wrappers that pull the bare-name
// candidate set from the matching Model map.

func suggestBoundedContext(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.BoundedContexts))
}

func suggestExternalSystem(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.ExternalSystems))
}

func suggestReadModel(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.ReadModels))
}

func suggestActor(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.Actors))
}

func suggestEvent(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.Events))
}

func suggestCommand(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.Commands))
}

func suggestDomain(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.Domains))
}

func suggestAggregate(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.Aggregates))
}

func suggestDCB(m *model.Model, name string) *suggestion {
	return suggest(name, candidatesFromMap(m.DynamicConsistencyBoundaries))
}

// unresolvedReferenceDiag builds the diagnostic for "the
// referenced entity does not exist in the index". The
// optional suggestion becomes a Related entry pointing
// at the suggested entity's definition.
func unresolvedReferenceDiag(refNode ast.Node, kind, name string, s *suggestion) diag.Diagnostic {
	d := diag.Diagnostic{
		RuleID:   "esdm/structure/unresolved-reference",
		Severity: diag.SeverityError,
		Message:  fmt.Sprintf("unresolved %s %q", kind, name),
		Location: refNode.Location(),
	}
	if s != nil {
		d.Related = []diag.Related{
			{
				Message:  fmt.Sprintf("did you mean %q?", s.name),
				Location: s.location,
			},
		}
	}
	return d
}

// mismatchedParentDiag builds the diagnostic for "the
// referenced entity exists, but it lives under a
// different parent than the scope claims". The Related
// entry points at the entity's actual definition so the
// reader can compare directly.
func mismatchedParentDiag(refNode ast.Node, kind, name, parentKind, expectedParent, actualParent string, definitionLocation diag.Location) diag.Diagnostic {
	return diag.Diagnostic{
		RuleID:   "esdm/structure/unresolved-reference",
		Severity: diag.SeverityError,
		Message:  fmt.Sprintf("%s %q exists but its %s is %q, not %q", kind, name, parentKind, actualParent, expectedParent),
		Location: refNode.Location(),
		Related: []diag.Related{
			{
				Message:  fmt.Sprintf("%s %q defined here", kind, name),
				Location: definitionLocation,
			},
		},
	}
}
