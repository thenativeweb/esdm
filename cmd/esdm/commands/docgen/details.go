package docgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/glossary"
	"github.com/thenativeweb/esdm/model"
)

// renderDetails writes the kind-specific detail sections of a page,
// aligned with what `esdm view --with-details` shows for each kind.
// Sections that name other elements link to them relative to the
// current page when the target is part of this output; otherwise they
// fall back to the element's reference so the page never carries a
// broken link.
func renderDetails(b *strings.Builder, n *node, index map[string]string) {
	fromFile := n.filePath()
	domain := n.segments[0]

	switch v := n.view.(type) {
	case model.BoundedContextView:
		renderUbiquitousLanguage(b, v)

	case model.SubdomainView:
		renderScalar(b, "Type", v.Type())
		renderRefList(b, "Bounded Contexts", fromFile, domainScopedPaths(domain, scalarNames(v.BoundedContexts().Seq())), index)

	case model.AggregateView:
		renderIdentity(b, v.IdentifiedBy())
		renderFields(b, "State", v.State())
		renderNameRules(b, "Invariants", v.Invariants().Seq())

	case model.DynamicConsistencyBoundaryView:
		renderConsults(b, fromFile, domain, v.Consults().Seq(), index)
		renderNameRules(b, "Invariants", v.Invariants().Seq())

	case model.CommandView:
		renderFields(b, "Payload", v.Data())
		renderRefList(b, "Publishes", fromFile, publishedEvents(n, v), index)
		renderRefList(b, "Actors", fromFile, boundedContextScopedPaths(domain, n.segments[1], scalarNames(v.Actors().Seq())), index)
		renderNameRules(b, "Constraints", v.Constraints().Seq())

	case model.EventView:
		renderFields(b, "Payload", v.Data())

	case model.ReadModelView:
		renderScalar(b, "Paradigm", v.Paradigm())
		renderProjections(b, fromFile, domain, v.Projections().Seq(), index)

	case model.QueryView:
		renderQueryReadModel(b, fromFile, domain, n.segments[1], v, index)
		renderRefList(b, "Actors", fromFile, boundedContextScopedPaths(domain, n.segments[1], scalarNames(v.Actors().Seq())), index)
		renderNameRules(b, "Constraints", v.Constraints().Seq())

	case model.EntityView:
		renderFields(b, "Schema", v.Schema())
		renderIdentity(b, v.IdentifiedBy())
		renderNameRules(b, "Invariants", v.Invariants().Seq())

	case model.ValueObjectView:
		renderFields(b, "Schema", v.Schema())
		renderNameRules(b, "Invariants", v.Invariants().Seq())

	case model.DomainServiceView:
		renderFunctions(b, v.Functions().Seq())

	case model.ActorView:
		renderScalar(b, "Type", v.Type())
		renderBullets(b, "Responsibilities", scalarNames(v.Responsibilities().Seq()))

	case model.PolicyView:
		renderScalar(b, "Delivery Guarantee", v.DeliveryGuarantee())
		renderRefList(b, "Handles", fromFile, eventReferences(domain, v.Handles().Seq()), index)
		renderRefList(b, "Emits", fromFile, commandReferences(domain, v.Emits().Seq()), index)

	case model.EventHandlerView:
		renderScalar(b, "Delivery Guarantee", v.DeliveryGuarantee())
		renderRefList(b, "Handles", fromFile, eventReferences(domain, v.Handles().Seq()), index)

	case model.ProcessManagerView:
		renderScalar(b, "Delivery Guarantee", v.DeliveryGuarantee())
		renderReactions(b, fromFile, domain, v.Reactions().Seq(), index)
		renderNameRules(b, "Invariants", v.Invariants().Seq())

	case model.ExternalSystemView:
		renderScalar(b, "Direction", v.Direction())
		renderScalar(b, "Category", v.Category())

	case model.FeatureView:
		renderScenarios(b, v.Scenarios().Seq())

	case model.DomainStoryView:
		renderStory(b, v)

	case model.ContextMappingView:
		renderScalar(b, "Type", v.Type())
		renderEndpoints(b, fromFile, v, index)
	}
}

// renderScalar writes a single-value section when the value is present.
func renderScalar(b *strings.Builder, heading string, value ast.Node) {
	content, ok := value.Text()
	if !ok || content == "" {
		return
	}
	fmt.Fprintf(b, "\n## %s\n\n%s\n", heading, content)
}

// renderFields writes the sorted property names of a JSON schema as a
// list, matching the field summary `esdm view` shows.
func renderFields(b *strings.Builder, heading string, schema ast.Node) {
	names := schemaFieldNames(schema)
	if len(names) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## %s\n\n", heading)
	for _, name := range names {
		fmt.Fprintf(b, "- `%s`\n", name)
	}
}

// renderBullets writes a plain bullet list section.
func renderBullets(b *strings.Builder, heading string, items []string) {
	if len(items) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## %s\n\n", heading)
	for _, item := range items {
		fmt.Fprintf(b, "- %s\n", item)
	}
}

// renderNameRules writes a section for name/rule pairs, used by both
// invariants and constraints.
func renderNameRules(b *strings.Builder, heading string, items []ast.Node) {
	if len(items) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## %s\n\n", heading)
	for _, item := range items {
		fmt.Fprintf(b, "- **%s**: %s\n", text(item, "name"), text(item, "rule"))
	}
}

// renderRefList writes a section that links to each target element.
func renderRefList(b *strings.Builder, heading, fromFile string, targets [][]string, index map[string]string) {
	if len(targets) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## %s\n\n", heading)
	for _, segments := range targets {
		fmt.Fprintf(b, "- %s\n", referenceMarkdown(fromFile, segments, index))
	}
}

// renderIdentity describes how instances of an aggregate or entity are
// identified.
func renderIdentity(b *strings.Builder, identifiedBy ast.Node) {
	source := text(identifiedBy, "source")
	var line string
	switch source {
	case "state":
		line = fmt.Sprintf("By its `%s` field, from the state.", text(identifiedBy, "field"))
	case "schema":
		line = fmt.Sprintf("By its `%s` field, from the schema.", text(identifiedBy, "field"))
	case "static":
		line = fmt.Sprintf("Statically, as `%s`.", text(identifiedBy, "value"))
	case "generated":
		line = fmt.Sprintf("Generated by `%s`.", text(identifiedBy, "generator"))
	default:
		return
	}
	fmt.Fprintf(b, "\n## Identity\n\n%s\n", line)
}

// renderProjections lists the events a read model projects, with the
// projection rule.
func renderProjections(b *strings.Builder, fromFile, domain string, projections []ast.Node, index map[string]string) {
	if len(projections) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Projections\n\n")
	for _, projection := range projections {
		fmt.Fprintf(b, "- %s: %s\n", referenceMarkdown(fromFile, eventReferenceSegments(domain, projection), index), text(projection, "rule"))
	}
}

// renderConsults lists the events a dynamic consistency boundary
// consults, with the criteria.
func renderConsults(b *strings.Builder, fromFile, domain string, consults []ast.Node, index map[string]string) {
	if len(consults) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Consults\n\n")
	for _, consult := range consults {
		fmt.Fprintf(b, "- %s: %s\n", referenceMarkdown(fromFile, eventReferenceSegments(domain, consult), index), text(consult, "criteria"))
	}
}

// renderReactions lists a process manager's reactions to events and
// timers.
func renderReactions(b *strings.Builder, fromFile, domain string, reactions []ast.Node, index map[string]string) {
	if len(reactions) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Reactions\n\n")
	for _, reaction := range reactions {
		rule := text(reaction, "rule")
		when := reaction.Field("when")
		if timer, ok := when.Field("timer").Text(); ok && timer != "" {
			fmt.Fprintf(b, "- On timer `%s`: %s\n", timer, rule)
			continue
		}
		fmt.Fprintf(b, "- On %s: %s\n", referenceMarkdown(fromFile, eventReferenceSegments(domain, when), index), rule)
	}
}

// renderQueryReadModel links a query to the read model it reads from.
func renderQueryReadModel(b *strings.Builder, fromFile, domain, boundedContext string, query model.QueryView, index map[string]string) {
	name, ok := query.ReadModel().Text()
	if !ok || name == "" {
		return
	}
	fmt.Fprintf(b, "\n## Read Model\n\n%s\n", referenceMarkdown(fromFile, []string{domain, boundedContext, name}, index))
}

// renderFunctions lists a domain service's functions.
func renderFunctions(b *strings.Builder, functions []ast.Node) {
	if len(functions) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Functions\n\n")
	for _, function := range functions {
		fmt.Fprintf(b, "- `%s`\n", text(function, "name"))
	}
}

// renderScenarios lists a feature's scenarios by name.
func renderScenarios(b *strings.Builder, scenarios []ast.Node) {
	if len(scenarios) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Scenarios\n\n")
	for _, scenario := range scenarios {
		fmt.Fprintf(b, "- %s\n", text(scenario, "name"))
	}
}

// renderStory lists a domain story's classifying properties.
func renderStory(b *strings.Builder, story model.DomainStoryView) {
	var lines []string
	if pointInTime, ok := story.PointInTime().Text(); ok && pointInTime != "" {
		lines = append(lines, "- Point in time: "+pointInTime)
	}
	if granularity, ok := story.Granularity().Text(); ok && granularity != "" {
		lines = append(lines, "- Granularity: "+granularity)
	}
	if domainPurity, ok := story.DomainPurity().Text(); ok && domainPurity != "" {
		lines = append(lines, "- Domain purity: "+domainPurity)
	}
	if len(lines) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Story\n\n%s\n", strings.Join(lines, "\n"))
}

// renderUbiquitousLanguage renders a bounded context's terms, reusing
// the glossary command's term extraction.
func renderUbiquitousLanguage(b *strings.Builder, boundedContext model.BoundedContextView) {
	terms := glossary.CollectTerms(boundedContext.UbiquitousLanguage())
	if len(terms) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Ubiquitous Language\n\n")
	for _, term := range terms {
		fmt.Fprintf(b, "- **%s**: %s\n", term.Term, term.Definition)
		for _, avoid := range term.Avoid {
			if avoid.Reason != "" {
				fmt.Fprintf(b, "  - avoid *%s*: %s\n", avoid.Term, avoid.Reason)
				continue
			}
			fmt.Fprintf(b, "  - avoid *%s*\n", avoid.Term)
		}
	}
}

// renderEndpoints links a context mapping to its two endpoints.
func renderEndpoints(b *strings.Builder, fromFile string, contextMapping model.ContextMappingView, index map[string]string) {
	roles := []struct {
		label string
		node  ast.Node
	}{
		{"Customer", contextMapping.Customer()},
		{"Supplier", contextMapping.Supplier()},
		{"Conformist", contextMapping.Conformist()},
		{"Upstream", contextMapping.Upstream()},
		{"Downstream", contextMapping.Downstream()},
		{"Host", contextMapping.Host()},
		{"Consumer", contextMapping.Consumer()},
		{"Publisher", contextMapping.Publisher()},
	}

	var lines []string
	for _, role := range roles {
		if link, ok := endpointLink(fromFile, role.node, index); ok {
			lines = append(lines, fmt.Sprintf("- %s: %s", role.label, link))
		}
	}
	for _, participant := range contextMapping.Participants().Seq() {
		if link, ok := endpointLink(fromFile, participant, index); ok {
			lines = append(lines, "- Participant: "+link)
		}
	}

	if len(lines) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Endpoints\n\n%s\n", strings.Join(lines, "\n"))
}

// endpointLink builds a link to a context-mapping endpoint, which is
// either a bounded context or an external system in some domain.
func endpointLink(fromFile string, endpoint ast.Node, index map[string]string) (string, bool) {
	if !endpoint.Exists() {
		return "", false
	}
	domain := text(endpoint, "domain")
	if boundedContext := text(endpoint, "boundedContext"); boundedContext != "" {
		return referenceMarkdown(fromFile, []string{domain, boundedContext}, index), true
	}
	if externalSystem := text(endpoint, "externalSystem"); externalSystem != "" {
		return referenceMarkdown(fromFile, []string{domain, externalSystem}, index), true
	}
	return "", false
}

// referenceMarkdown links to the element at segments when its page is part of
// this output; otherwise it renders the element's reference so the
// link is never broken.
func referenceMarkdown(fromFile string, segments []string, index map[string]string) string {
	logical := strings.Join(segments, "/")
	if file, ok := index[logical]; ok {
		return fmt.Sprintf("[%s](%s)", segments[len(segments)-1], relLink(fromFile, file))
	}
	return fmt.Sprintf("`esdm:%s`", logical)
}

// eventReferenceSegments turns an event-reference triple into a containment
// path, collapsing to three segments for a free-standing event. The
// domain is the referencing element's, since triples omit it.
func eventReferenceSegments(domain string, reference ast.Node) []string {
	boundedContext := text(reference, "boundedContext")
	aggregate := text(reference, "aggregate")
	event := text(reference, "event")
	if aggregate == "" {
		return []string{domain, boundedContext, event}
	}
	return []string{domain, boundedContext, aggregate, event}
}

// commandReferenceSegments turns a command-reference triple into a containment
// path, whose parent is either an aggregate or a dynamic consistency
// boundary.
func commandReferenceSegments(domain string, reference ast.Node) []string {
	boundedContext := text(reference, "boundedContext")
	parent := text(reference, "aggregate")
	if parent == "" {
		parent = text(reference, "dynamicConsistencyBoundary")
	}
	return []string{domain, boundedContext, parent, text(reference, "command")}
}

// eventReferences maps event-reference triples to containment paths.
func eventReferences(domain string, references []ast.Node) [][]string {
	var out [][]string
	for _, reference := range references {
		out = append(out, eventReferenceSegments(domain, reference))
	}
	return out
}

// commandReferences maps command-reference triples to containment paths.
func commandReferences(domain string, references []ast.Node) [][]string {
	var out [][]string
	for _, reference := range references {
		out = append(out, commandReferenceSegments(domain, reference))
	}
	return out
}

// publishedEvents resolves a command's published event names to
// containment paths. An aggregate's command publishes that aggregate's
// events; a command on a dynamic consistency boundary publishes
// free-standing events.
func publishedEvents(n *node, command model.CommandView) [][]string {
	domain := n.segments[0]
	boundedContext := n.segments[1]
	parentIsAggregate := text(command.Scope(), "aggregate") != ""

	var out [][]string
	for _, event := range scalarNames(command.Publishes().Seq()) {
		if parentIsAggregate {
			out = append(out, []string{domain, boundedContext, n.segments[2], event})
			continue
		}
		out = append(out, []string{domain, boundedContext, event})
	}
	return out
}

// domainScopedPaths builds domain-scoped containment paths from bare names.
func domainScopedPaths(domain string, names []string) [][]string {
	var out [][]string
	for _, name := range names {
		out = append(out, []string{domain, name})
	}
	return out
}

// boundedContextScopedPaths builds bounded-context-scoped containment paths from bare
// names.
func boundedContextScopedPaths(domain, boundedContext string, names []string) [][]string {
	var out [][]string
	for _, name := range names {
		out = append(out, []string{domain, boundedContext, name})
	}
	return out
}

// scalarNames reads a sequence of scalar names.
func scalarNames(seq []ast.Node) []string {
	var out []string
	for _, item := range seq {
		if v, ok := item.Text(); ok {
			out = append(out, v)
		}
	}
	return out
}

// schemaFieldNames returns the sorted top-level property names of a
// JSON schema node.
func schemaFieldNames(schema ast.Node) []string {
	properties := schema.Field("properties")
	if !properties.Exists() {
		return nil
	}
	var names []string
	for _, entry := range properties.Entries() {
		if key, ok := entry.Key.Text(); ok {
			names = append(names, key)
		}
	}
	sort.Strings(names)
	return names
}

// text reads a scalar field, returning "" when it is absent.
func text(n ast.Node, field string) string {
	value, _ := n.Field(field).Text()
	return value
}
