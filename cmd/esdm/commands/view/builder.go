package view

import (
	"fmt"
	"sort"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/modelpath"
)

// BuildTree returns the render tree for the given
// resolved model, narrowed to the subtree identified by
// p. An empty Path selects the entire model.
func BuildTree(m *model.Model, p modelpath.Path, withDetails bool) (*Node, error) {
	root := buildAllDomains(m, withDetails)
	if len(p.Segments) == 0 {
		return root, nil
	}

	target, err := narrow(root, p.Segments)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// buildAllDomains constructs the synthetic root node
// containing every domain in the model. When the model
// has exactly one domain the synthetic root is
// transparent - the domain itself becomes the rendered
// root via the empty-path narrowing logic of the
// renderer.
func buildAllDomains(m *model.Model, withDetails bool) *Node {
	root := &Node{Kind: "model", Name: ""}
	for _, d := range sortedByBareName(m.Domains) {
		root.Children = append(root.Children, buildDomain(m, d, withDetails))
	}
	return root
}

func buildDomain(m *model.Model, d model.DomainView, withDetails bool) *Node {
	name, _ := d.Name().Text()
	n := &Node{
		Kind:     "domain",
		Name:     name,
		Key:      name,
		Location: nameLocation(d),
	}

	subdomainCount := 0
	boundedContextCount := 0
	processManagerCount := 0
	eventHandlerCount := 0
	policyCount := 0
	externalSystemCount := 0
	contextMappingCount := 0
	storyCount := 0
	featureCount := 0

	for _, sub := range filterSubdomainsByDomain(m, name) {
		n.Children = append(n.Children, buildSubdomain(sub, withDetails))
		subdomainCount++
	}
	for _, boundedContext := range filterBoundedContextsByDomain(m, name) {
		n.Children = append(n.Children, buildBoundedContext(m, boundedContext, withDetails))
		boundedContextCount++
	}
	for _, processManager := range filterProcessManagersByDomain(m, name) {
		n.Children = append(n.Children, buildProcessManager(processManager, withDetails))
		processManagerCount++
	}
	for _, eventHandler := range filterEventHandlersByDomain(m, name) {
		n.Children = append(n.Children, buildEventHandler(eventHandler, withDetails))
		eventHandlerCount++
	}
	for _, pol := range filterPoliciesByDomain(m, name) {
		n.Children = append(n.Children, buildPolicy(pol, withDetails))
		policyCount++
	}
	for _, externalSystem := range filterExternalSystemsByDomain(m, name) {
		n.Children = append(n.Children, buildExternalSystem(externalSystem, withDetails))
		externalSystemCount++
	}
	for _, contextMapping := range sortedContextMappings(m) {
		if !contextMappingTouchesDomain(contextMapping, name) {
			continue
		}
		n.Children = append(n.Children, buildContextMapping(contextMapping, withDetails))
		contextMappingCount++
	}
	for _, story := range filterStoriesByDomain(m, name) {
		n.Children = append(n.Children, buildStory(story, withDetails))
		storyCount++
	}
	for _, feature := range filterFeaturesByDomain(m, name) {
		n.Children = append(n.Children, buildFeature(feature, withDetails))
		featureCount++
	}

	stats := plural(subdomainCount, "sub")
	if s := plural(boundedContextCount, "bc"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(processManagerCount, "pm"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(eventHandlerCount, "eh"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(policyCount, "pol"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(externalSystemCount, "es"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(contextMappingCount, "cm"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(storyCount, "story"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(featureCount, "feat"); s != "" {
		stats = appendStat(stats, s)
	}
	if stats != "" {
		n.Stats = []string{stats}
	}

	return n
}

func buildSubdomain(s model.SubdomainView, withDetails bool) *Node {
	name, _ := s.Name().Text()
	subType, _ := s.Type().Text()
	domain := scopeText(s.Scope(), "domain")

	n := &Node{
		Kind:     "subdomain",
		Name:     name,
		Tags:     []string{subType},
		Key:      domain + "/" + name,
		Location: nameLocation(s),
	}
	var bcs []string
	for _, item := range s.BoundedContexts().Seq() {
		if v, ok := item.Text(); ok {
			bcs = append(bcs, v)
		}
	}
	if len(bcs) > 0 {
		n.Stats = []string{joinComma(bcs)}
	}
	if withDetails {
		if desc, ok := s.Description().Text(); ok && desc != "" {
			n.Lines = append(n.Lines, "description: "+desc)
		}
	}
	return n
}

func buildBoundedContext(m *model.Model, boundedContext model.BoundedContextView, withDetails bool) *Node {
	name, _ := boundedContext.Name().Text()
	domain := scopeText(boundedContext.Scope(), "domain")

	n := &Node{
		Kind:     "bounded-context",
		Name:     name,
		Key:      domain + "/" + name,
		Location: nameLocation(boundedContext),
	}

	aggregateCount := 0
	dynamicConsistencyBoundaryCount := 0
	readModelCount := 0
	queryCount := 0
	valueObjectCount := 0
	domainServiceCount := 0
	actorCount := 0

	for _, aggregate := range filterAggregatesByBoundedContext(m, domain, name) {
		n.Children = append(n.Children, buildAggregate(m, aggregate, withDetails))
		aggregateCount++
	}
	for _, dcb := range filterDCBsByBoundedContext(m, domain, name) {
		n.Children = append(n.Children, buildDCB(m, dcb, withDetails))
		dynamicConsistencyBoundaryCount++
	}
	for _, readModel := range filterReadModelsByBoundedContext(m, domain, name) {
		n.Children = append(n.Children, buildReadModel(readModel, withDetails))
		readModelCount++
	}
	for _, query := range filterQueriesByBoundedContext(m, domain, name) {
		n.Children = append(n.Children, buildQuery(query, withDetails))
		queryCount++
	}
	for _, valueObject := range filterValueObjectsByBoundedContext(m, domain, name) {
		n.Children = append(n.Children, buildValueObject(valueObject, withDetails))
		valueObjectCount++
	}
	for _, domainService := range filterDomainServicesByBoundedContext(m, domain, name) {
		n.Children = append(n.Children, buildDomainService(domainService, withDetails))
		domainServiceCount++
	}
	for _, a := range filterActorsByBoundedContext(m, domain, name) {
		n.Children = append(n.Children, buildActor(a, withDetails))
		actorCount++
	}

	stats := plural(aggregateCount, "agg")
	if s := plural(dynamicConsistencyBoundaryCount, "dcb"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(readModelCount, "rm"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(queryCount, "qry"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(valueObjectCount, "vo"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(domainServiceCount, "ds"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(actorCount, "act"); s != "" {
		stats = appendStat(stats, s)
	}
	if stats != "" {
		n.Stats = []string{stats}
	}

	if withDetails {
		for _, term := range boundedContext.UbiquitousLanguage().Seq() {
			t, _ := term.Field("term").Text()
			def, _ := term.Field("definition").Text()
			n.Lines = append(n.Lines, fmt.Sprintf("term %q: %s", t, def))
			for _, av := range term.Field("avoid").Seq() {
				avTerm, _ := av.Field("term").Text()
				if reason, ok := av.Field("reason").Text(); ok && reason != "" {
					n.Lines = append(n.Lines, fmt.Sprintf("  avoid %q: %s", avTerm, reason))
					continue
				}
				n.Lines = append(n.Lines, fmt.Sprintf("  avoid %q", avTerm))
			}
		}
	}

	return n
}

func buildAggregate(m *model.Model, aggregate model.AggregateView, withDetails bool) *Node {
	name, _ := aggregate.Name().Text()
	domain := scopeText(aggregate.Scope(), "domain")
	boundedContext := scopeText(aggregate.Scope(), "boundedContext")

	n := &Node{
		Kind:     "aggregate",
		Name:     name,
		Key:      domain + "/" + boundedContext + "/" + name,
		Location: nameLocation(aggregate),
	}

	cmdCount := 0
	for _, cmd := range filterCommandsByAggregate(m, domain, boundedContext, name) {
		n.Children = append(n.Children, buildCommand(cmd, withDetails))
		cmdCount++
	}
	evtCount := 0
	for _, event := range filterEventsByAggregate(m, domain, boundedContext, name) {
		n.Children = append(n.Children, buildEvent(m, event, withDetails))
		evtCount++
	}
	invCount := len(aggregate.Invariants().Seq())

	stats := plural(cmdCount, "cmd")
	if s := plural(evtCount, "evt"); s != "" {
		stats = appendStat(stats, s)
	}
	if s := plural(invCount, "inv"); s != "" {
		stats = appendStat(stats, s)
	}
	if stats != "" {
		n.Stats = []string{stats}
	}

	if withDetails {
		ib := aggregate.IdentifiedBy()
		if src, ok := ib.Field("source").Text(); ok {
			switch src {
			case "state":
				field, _ := ib.Field("field").Text()
				n.Lines = append(n.Lines, fmt.Sprintf("identifiedBy: state.%s", field))
			case "static":
				value, _ := ib.Field("value").Text()
				n.Lines = append(n.Lines, fmt.Sprintf("identifiedBy: static %q", value))
			case "generated":
				gen, _ := ib.Field("generator").Text()
				n.Lines = append(n.Lines, fmt.Sprintf("identifiedBy: generated/%s", gen))
			}
		}
		for _, inv := range aggregate.Invariants().Seq() {
			invName, _ := inv.Field("name").Text()
			rule, _ := inv.Field("rule").Text()
			n.Lines = append(n.Lines, fmt.Sprintf("invariant %q: %s", invName, rule))
		}
	}
	return n
}

func buildDCB(m *model.Model, dcb model.DynamicConsistencyBoundaryView, withDetails bool) *Node {
	name, _ := dcb.Name().Text()
	domain := scopeText(dcb.Scope(), "domain")
	boundedContext := scopeText(dcb.Scope(), "boundedContext")

	n := &Node{
		Kind:     "dynamic-consistency-boundary",
		Name:     name,
		Key:      domain + "/" + boundedContext + "/" + name,
		Location: nameLocation(dcb),
	}

	cmdCount := 0
	for _, cmd := range filterCommandsByDCB(m, domain, boundedContext, name) {
		n.Children = append(n.Children, buildCommand(cmd, withDetails))
		cmdCount++
	}
	consults := len(dcb.Consults().Seq())
	stats := plural(cmdCount, "cmd")
	if s := plural(consults, "consult"); s != "" {
		stats = appendStat(stats, s)
	}
	if stats != "" {
		n.Stats = []string{stats}
	}
	if withDetails {
		for _, inv := range dcb.Invariants().Seq() {
			invName, _ := inv.Field("name").Text()
			rule, _ := inv.Field("rule").Text()
			n.Lines = append(n.Lines, fmt.Sprintf("invariant %q: %s", invName, rule))
		}
	}
	return n
}

func buildCommand(cmd model.CommandView, withDetails bool) *Node {
	name, _ := cmd.Name().Text()
	domain := scopeText(cmd.Scope(), "domain")
	boundedContext := scopeText(cmd.Scope(), "boundedContext")
	parent := scopeText(cmd.Scope(), "aggregate")
	if parent == "" {
		parent = scopeText(cmd.Scope(), "dynamicConsistencyBoundary")
	}

	n := &Node{
		Kind:     "command",
		Name:     name,
		Key:      domain + "/" + boundedContext + "/" + parent + "/" + name,
		Location: nameLocation(cmd),
	}

	var publishes []string
	for _, item := range cmd.Publishes().Seq() {
		if v, ok := item.Text(); ok {
			publishes = append(publishes, v)
		}
	}
	var actors []string
	for _, item := range cmd.Actors().Seq() {
		if v, ok := item.Text(); ok {
			actors = append(actors, v)
		}
	}
	if len(publishes) > 0 {
		n.Stats = append(n.Stats, "→ "+joinComma(publishes))
	}
	if len(actors) > 0 {
		n.Stats = append(n.Stats, "by "+joinComma(actors))
	}

	if withDetails {
		schema := schemaSummary(cmd.Data())
		if schema != "" {
			n.Lines = append(n.Lines, "data: "+schema)
		}
		for _, c := range cmd.Constraints().Seq() {
			cName, _ := c.Field("name").Text()
			rule, _ := c.Field("rule").Text()
			n.Lines = append(n.Lines, fmt.Sprintf("constraint %q: %s", cName, rule))
		}
	}
	return n
}

func buildEvent(m *model.Model, event model.EventView, withDetails bool) *Node {
	name, _ := event.Name().Text()
	domain := scopeText(event.Scope(), "domain")
	boundedContext := scopeText(event.Scope(), "boundedContext")
	aggregate := scopeText(event.Scope(), "aggregate")

	n := &Node{
		Kind:     "event",
		Name:     name,
		Key:      domain + "/" + boundedContext + "/" + aggregate + "/" + name,
		Location: nameLocation(event),
	}
	if publishers := filterCommandsPublishingEvent(m, event); len(publishers) > 0 {
		n.Stats = append(n.Stats, "← "+joinComma(publishers))
	}
	if withDetails {
		schema := schemaSummary(event.Data())
		if schema != "" {
			n.Lines = append(n.Lines, "data: "+schema)
		}
	}
	return n
}

func buildReadModel(readModel model.ReadModelView, withDetails bool) *Node {
	name, _ := readModel.Name().Text()
	domain := scopeText(readModel.Scope(), "domain")
	boundedContext := scopeText(readModel.Scope(), "boundedContext")

	n := &Node{
		Kind:     "read-model",
		Name:     name,
		Key:      domain + "/" + boundedContext + "/" + name,
		Location: nameLocation(readModel),
	}
	projections := readModel.Projections().Seq()
	if len(projections) > 0 {
		n.Stats = []string{fmt.Sprintf("← %s", plural(len(projections), "evt"))}
	}
	if withDetails {
		for _, proj := range projections {
			pBC, _ := proj.Field("boundedContext").Text()
			pAgg, _ := proj.Field("aggregate").Text()
			pEvent, _ := proj.Field("event").Text()
			rule, _ := proj.Field("rule").Text()
			ref := pBC + "/" + pAgg + "/" + pEvent
			if pAgg == "" {
				ref = pBC + "//" + pEvent
			}
			n.Lines = append(n.Lines, fmt.Sprintf("projects %s - %s", ref, rule))
		}
		if para, ok := readModel.Paradigm().Text(); ok && para != "" {
			n.Lines = append(n.Lines, "paradigm: "+para)
		}
	}
	return n
}

func buildQuery(query model.QueryView, withDetails bool) *Node {
	name, _ := query.Name().Text()
	domain := scopeText(query.Scope(), "domain")
	boundedContext := scopeText(query.Scope(), "boundedContext")

	n := &Node{
		Kind:     "query",
		Name:     name,
		Key:      domain + "/" + boundedContext + "/" + name,
		Location: nameLocation(query),
	}
	if rmName, ok := query.ReadModel().Text(); ok {
		n.Stats = append(n.Stats, "→ "+rmName)
	}
	var actors []string
	for _, item := range query.Actors().Seq() {
		if v, ok := item.Text(); ok {
			actors = append(actors, v)
		}
	}
	if len(actors) > 0 {
		n.Stats = append(n.Stats, "by "+joinComma(actors))
	}
	if withDetails {
		for _, c := range query.Constraints().Seq() {
			cName, _ := c.Field("name").Text()
			rule, _ := c.Field("rule").Text()
			n.Lines = append(n.Lines, fmt.Sprintf("constraint %q: %s", cName, rule))
		}
	}
	return n
}

func buildValueObject(valueObject model.ValueObjectView, withDetails bool) *Node {
	name, _ := valueObject.Name().Text()
	domain := scopeText(valueObject.Scope(), "domain")
	boundedContext := scopeText(valueObject.Scope(), "boundedContext")
	n := &Node{
		Kind:     "value-object",
		Name:     name,
		Key:      domain + "/" + boundedContext + "/" + name,
		Location: nameLocation(valueObject),
	}
	invariantCount := len(valueObject.Invariants().Seq())
	if s := plural(invariantCount, "inv"); s != "" {
		n.Stats = []string{s}
	}
	if withDetails {
		schema := schemaSummary(valueObject.Schema())
		if schema != "" {
			n.Lines = append(n.Lines, "schema: "+schema)
		}
		for _, inv := range valueObject.Invariants().Seq() {
			invName, _ := inv.Field("name").Text()
			rule, _ := inv.Field("rule").Text()
			n.Lines = append(n.Lines, fmt.Sprintf("invariant %q: %s", invName, rule))
		}
	}
	return n
}

func buildDomainService(domainService model.DomainServiceView, withDetails bool) *Node {
	name, _ := domainService.Name().Text()
	domain := scopeText(domainService.Scope(), "domain")
	boundedContext := scopeText(domainService.Scope(), "boundedContext")
	n := &Node{
		Kind:     "domain-service",
		Name:     name,
		Key:      domain + "/" + boundedContext + "/" + name,
		Location: nameLocation(domainService),
	}
	functionCount := len(domainService.Functions().Seq())
	if s := plural(functionCount, "fn"); s != "" {
		n.Stats = []string{s}
	}
	if withDetails {
		for _, fn := range domainService.Functions().Seq() {
			fnName, _ := fn.Field("name").Text()
			n.Lines = append(n.Lines, "function "+fnName)
		}
	}
	return n
}

func buildActor(a model.ActorView, withDetails bool) *Node {
	name, _ := a.Name().Text()
	domain := scopeText(a.Scope(), "domain")
	boundedContext := scopeText(a.Scope(), "boundedContext")
	atype, _ := a.Type().Text()
	n := &Node{
		Kind:     "actor",
		Name:     name,
		Tags:     []string{atype},
		Key:      domain + "/" + boundedContext + "/" + name,
		Location: nameLocation(a),
	}
	if withDetails {
		for _, r := range a.Responsibilities().Seq() {
			if v, ok := r.Text(); ok {
				n.Lines = append(n.Lines, "- "+v)
			}
		}
	}
	return n
}

func buildProcessManager(processManager model.ProcessManagerView, withDetails bool) *Node {
	name, _ := processManager.Name().Text()
	domain := scopeText(processManager.Scope(), "domain")
	n := &Node{
		Kind:     "process-manager",
		Name:     name,
		Key:      domain + "/" + name,
		Location: nameLocation(processManager),
	}
	if deliveryGuarantee, ok := processManager.DeliveryGuarantee().Text(); ok {
		n.Stats = append(n.Stats, deliveryGuarantee)
	}
	if withDetails {
		for _, r := range processManager.Reactions().Seq() {
			rule, _ := r.Field("rule").Text()
			when := r.Field("when")
			if t, ok := when.Field("timer").Text(); ok {
				n.Lines = append(n.Lines, fmt.Sprintf("on timer %q: %s", t, rule))
				continue
			}
			event, _ := when.Field("event").Text()
			n.Lines = append(n.Lines, fmt.Sprintf("on event %s: %s", event, rule))
		}
	}
	return n
}

func buildEventHandler(eventHandler model.EventHandlerView, withDetails bool) *Node {
	name, _ := eventHandler.Name().Text()
	domain := scopeText(eventHandler.Scope(), "domain")
	n := &Node{
		Kind:     "event-handler",
		Name:     name,
		Key:      domain + "/" + name,
		Location: nameLocation(eventHandler),
	}
	if deliveryGuarantee, ok := eventHandler.DeliveryGuarantee().Text(); ok {
		n.Stats = append(n.Stats, deliveryGuarantee)
	}
	if withDetails {
		for _, h := range eventHandler.Handles().Seq() {
			event, _ := h.Field("event").Text()
			n.Lines = append(n.Lines, "handles "+event)
		}
	}
	return n
}

func buildPolicy(p model.PolicyView, withDetails bool) *Node {
	name, _ := p.Name().Text()
	domain := scopeText(p.Scope(), "domain")
	n := &Node{
		Kind:     "policy",
		Name:     name,
		Key:      domain + "/" + name,
		Location: nameLocation(p),
	}
	if deliveryGuarantee, ok := p.DeliveryGuarantee().Text(); ok {
		n.Stats = append(n.Stats, deliveryGuarantee)
	}
	return n
}

func buildExternalSystem(externalSystem model.ExternalSystemView, withDetails bool) *Node {
	name, _ := externalSystem.Name().Text()
	domain := scopeText(externalSystem.Scope(), "domain")
	n := &Node{
		Kind:     "external-system",
		Name:     name,
		Key:      domain + "/" + name,
		Location: nameLocation(externalSystem),
	}
	if d, ok := externalSystem.Direction().Text(); ok {
		n.Stats = append(n.Stats, d)
	}
	if cat, ok := externalSystem.Category().Text(); ok && cat != "" {
		n.Stats = append(n.Stats, cat)
	}
	return n
}

func buildContextMapping(contextMapping model.ContextMappingView, withDetails bool) *Node {
	name, _ := contextMapping.Name().Text()
	cmType, _ := contextMapping.Type().Text()
	n := &Node{
		Kind:     "context-mapping",
		Name:     name,
		Tags:     []string{cmType},
		Key:      name,
		Location: nameLocation(contextMapping),
	}
	return n
}

func buildFeature(feature model.FeatureView, withDetails bool) *Node {
	name, _ := feature.Name().Text()
	scope := feature.Scope()
	domain := scopeText(scope, "domain")

	n := &Node{
		Kind:     "feature",
		Name:     name,
		Key:      domain + "/" + name,
		Location: nameLocation(feature),
	}

	variant, target := featureVariantAndTarget(scope)
	if variant != "" {
		n.Tags = append(n.Tags, variant)
	}
	if target != "" {
		n.Tags = append(n.Tags, target)
	}

	scenarios := feature.Scenarios().Seq()
	if len(scenarios) > 0 {
		word := "scenario"
		if len(scenarios) != 1 {
			word = "scenarios"
		}
		n.Stats = append(n.Stats, fmt.Sprintf("%d %s", len(scenarios), word))
	}

	if withDetails {
		if desc, ok := feature.Description().Text(); ok && desc != "" {
			n.Lines = append(n.Lines, desc)
		}
		for _, scenario := range scenarios {
			sName, _ := scenario.Field("name").Text()
			n.Lines = append(n.Lines, "- "+sName)
		}
	}
	return n
}

// featureVariantAndTarget returns the discriminator name
// (aggregate / dynamic-consistency-boundary / process-
// manager / read-model) of the feature's scope plus the
// bare name of the targeted unit. Both default to "" when
// no recognized discriminator is present.
func featureVariantAndTarget(scope ast.Node) (string, string) {
	switch {
	case scope.Field("aggregate").Exists():
		t, _ := scope.Field("aggregate").Text()
		return "aggregate", t
	case scope.Field("dynamicConsistencyBoundary").Exists():
		t, _ := scope.Field("dynamicConsistencyBoundary").Text()
		return "dynamic-consistency-boundary", t
	case scope.Field("processManager").Exists():
		t, _ := scope.Field("processManager").Text()
		return "process-manager", t
	case scope.Field("readModel").Exists():
		t, _ := scope.Field("readModel").Text()
		return "read-model", t
	}
	return "", ""
}

func buildStory(story model.DomainStoryView, withDetails bool) *Node {
	name, _ := story.Name().Text()
	domain := scopeText(story.Scope(), "domain")

	n := &Node{
		Kind:     "domain-story",
		Name:     name,
		Key:      domain + "/" + name,
		Location: nameLocation(story),
	}
	if pointInTime, ok := story.PointInTime().Text(); ok {
		n.Tags = append(n.Tags, pointInTime)
	}
	if g, ok := story.Granularity().Text(); ok {
		n.Tags = append(n.Tags, g)
	}
	if domainPurity, ok := story.DomainPurity().Text(); ok {
		n.Tags = append(n.Tags, domainPurity)
	}
	sentenceCount := len(story.Sentences().Seq())
	if sentenceCount > 0 {
		word := "sentence"
		if sentenceCount != 1 {
			word = "sentences"
		}
		n.Stats = append(n.Stats, fmt.Sprintf("%d %s", sentenceCount, word))
	}
	if withDetails {
		if desc, ok := story.Description().Text(); ok && desc != "" {
			n.Lines = append(n.Lines, desc)
		}
	}
	return n
}

// Filter helpers - pull the entries of a particular
// kind that live inside a given scope. They iterate the
// model values (not the composite map keys) and read
// the scope from each view.

func filterSubdomainsByDomain(m *model.Model, domain string) []model.SubdomainView {
	var out []model.SubdomainView
	for _, v := range m.Subdomains {
		if scopeText(v.Scope(), "domain") == domain {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterBoundedContextsByDomain(m *model.Model, domain string) []model.BoundedContextView {
	var out []model.BoundedContextView
	for _, v := range m.BoundedContexts {
		if scopeText(v.Scope(), "domain") == domain {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterAggregatesByBoundedContext(m *model.Model, domain, boundedContext string) []model.AggregateView {
	var out []model.AggregateView
	for _, v := range m.Aggregates {
		if scopeText(v.Scope(), "domain") == domain && scopeText(v.Scope(), "boundedContext") == boundedContext {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterDCBsByBoundedContext(m *model.Model, domain, boundedContext string) []model.DynamicConsistencyBoundaryView {
	var out []model.DynamicConsistencyBoundaryView
	for _, v := range m.DynamicConsistencyBoundaries {
		if scopeText(v.Scope(), "domain") == domain && scopeText(v.Scope(), "boundedContext") == boundedContext {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterReadModelsByBoundedContext(m *model.Model, domain, boundedContext string) []model.ReadModelView {
	var out []model.ReadModelView
	for _, v := range m.ReadModels {
		if scopeText(v.Scope(), "domain") == domain && scopeText(v.Scope(), "boundedContext") == boundedContext {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterQueriesByBoundedContext(m *model.Model, domain, boundedContext string) []model.QueryView {
	var out []model.QueryView
	for _, v := range m.Queries {
		if scopeText(v.Scope(), "domain") == domain && scopeText(v.Scope(), "boundedContext") == boundedContext {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterValueObjectsByBoundedContext(m *model.Model, domain, boundedContext string) []model.ValueObjectView {
	var out []model.ValueObjectView
	for _, v := range m.ValueObjects {
		if scopeText(v.Scope(), "domain") == domain && scopeText(v.Scope(), "boundedContext") == boundedContext {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterDomainServicesByBoundedContext(m *model.Model, domain, boundedContext string) []model.DomainServiceView {
	var out []model.DomainServiceView
	for _, v := range m.DomainServices {
		if scopeText(v.Scope(), "domain") == domain && scopeText(v.Scope(), "boundedContext") == boundedContext {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterActorsByBoundedContext(m *model.Model, domain, boundedContext string) []model.ActorView {
	var out []model.ActorView
	for _, v := range m.Actors {
		if scopeText(v.Scope(), "domain") == domain && scopeText(v.Scope(), "boundedContext") == boundedContext {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterCommandsByAggregate(m *model.Model, domain, boundedContext, aggregate string) []model.CommandView {
	var out []model.CommandView
	for _, v := range m.Commands {
		if scopeText(v.Scope(), "domain") != domain {
			continue
		}
		if scopeText(v.Scope(), "boundedContext") != boundedContext {
			continue
		}
		if scopeText(v.Scope(), "aggregate") != aggregate {
			continue
		}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterCommandsByDCB(m *model.Model, domain, boundedContext, dcb string) []model.CommandView {
	var out []model.CommandView
	for _, v := range m.Commands {
		if scopeText(v.Scope(), "domain") != domain {
			continue
		}
		if scopeText(v.Scope(), "boundedContext") != boundedContext {
			continue
		}
		if scopeText(v.Scope(), "dynamicConsistencyBoundary") != dcb {
			continue
		}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

// filterCommandsPublishingEvent returns the names of
// every command that publishes ev. Because command.publishes
// carries bare event names, the publisher must sit in the
// same scope as the event - the same aggregate for
// aggregate-owned events, the same DCB for free-standing
// events. The result is sorted alphabetically so the
// renderer's output is deterministic.
func filterCommandsPublishingEvent(m *model.Model, event model.EventView) []string {
	evName, _ := event.Name().Text()
	evDomain := scopeText(event.Scope(), "domain")
	evBC := scopeText(event.Scope(), "boundedContext")
	evAgg := scopeText(event.Scope(), "aggregate")

	var out []string
	for _, cmd := range m.Commands {
		if scopeText(cmd.Scope(), "domain") != evDomain {
			continue
		}
		if scopeText(cmd.Scope(), "boundedContext") != evBC {
			continue
		}
		if evAgg != "" {
			if scopeText(cmd.Scope(), "aggregate") != evAgg {
				continue
			}
		} else {
			if scopeText(cmd.Scope(), "dynamicConsistencyBoundary") == "" {
				continue
			}
		}
		for _, item := range cmd.Publishes().Seq() {
			v, ok := item.Text()
			if !ok || v != evName {
				continue
			}
			if cName, ok := cmd.Name().Text(); ok {
				out = append(out, cName)
			}
			break
		}
	}
	sort.Strings(out)
	return out
}

func filterEventsByAggregate(m *model.Model, domain, boundedContext, aggregate string) []model.EventView {
	var out []model.EventView
	for _, v := range m.Events {
		if scopeText(v.Scope(), "domain") != domain {
			continue
		}
		if scopeText(v.Scope(), "boundedContext") != boundedContext {
			continue
		}
		if scopeText(v.Scope(), "aggregate") != aggregate {
			continue
		}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterProcessManagersByDomain(m *model.Model, domain string) []model.ProcessManagerView {
	var out []model.ProcessManagerView
	for _, v := range m.ProcessManagers {
		if scopeText(v.Scope(), "domain") == domain {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterEventHandlersByDomain(m *model.Model, domain string) []model.EventHandlerView {
	var out []model.EventHandlerView
	for _, v := range m.EventHandlers {
		if scopeText(v.Scope(), "domain") == domain {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterPoliciesByDomain(m *model.Model, domain string) []model.PolicyView {
	var out []model.PolicyView
	for _, v := range m.Policies {
		if scopeText(v.Scope(), "domain") == domain {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterExternalSystemsByDomain(m *model.Model, domain string) []model.ExternalSystemView {
	var out []model.ExternalSystemView
	for _, v := range m.ExternalSystems {
		if scopeText(v.Scope(), "domain") == domain {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func sortedContextMappings(m *model.Model) []model.ContextMappingView {
	var out []model.ContextMappingView
	for _, v := range m.ContextMappings {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func contextMappingTouchesDomain(contextMapping model.ContextMappingView, domain string) bool {
	endpoints := []model.DocumentViewBase{
		{Node: contextMapping.Customer()}, {Node: contextMapping.Supplier()},
		{Node: contextMapping.Conformist()}, {Node: contextMapping.Upstream()}, {Node: contextMapping.Downstream()},
		{Node: contextMapping.Host()}, {Node: contextMapping.Consumer()}, {Node: contextMapping.Publisher()},
	}
	for _, ep := range endpoints {
		if d, ok := ep.Field("domain").Text(); ok && d == domain {
			return true
		}
	}
	for _, p := range contextMapping.Participants().Seq() {
		if d, ok := p.Field("domain").Text(); ok && d == domain {
			return true
		}
	}
	return false
}

func filterStoriesByDomain(m *model.Model, domain string) []model.DomainStoryView {
	var out []model.DomainStoryView
	for _, v := range m.Extensions.DomainStorytelling.Stories {
		if scopeText(v.Scope(), "domain") == domain {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

func filterFeaturesByDomain(m *model.Model, domain string) []model.FeatureView {
	var out []model.FeatureView
	for _, v := range m.Extensions.GivenWhenThen.Features {
		if scopeText(v.Scope(), "domain") == domain {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}
