package docgen

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/modelpath"
)

// Page is one rendered Markdown file: a relative slash path within the
// output tree and its content.
type Page struct {
	Path    string
	Content string
}

// documented is the shared surface every model view exposes through
// its embedded base: a name and a description.
type documented interface {
	Name() ast.Node
	Description() ast.Node
}

// scoped is a documented view that additionally carries a scope, which
// is every kind except domain and context-mapping.
type scoped interface {
	documented
	Scope() ast.Node
}

// node is one element in the containment tree. A node with children is
// written as a directory with a README.md index; a node without
// children is written as <name>.md.
type node struct {
	kind        string
	name        string
	description string
	segments    []string
	view        any
	children    []*node
}

// Build turns the resolved model into the set of pages to write,
// narrowed to the region identified by p. The empty path selects the
// whole model and adds a root index; a non-empty path selects a single
// element and its subtree, keeping the full containment paths so the
// pages still line up with their references. An unknown or too-deep
// path segment is rejected as invalid input, mirroring esdm view.
func Build(m *model.Model, p modelpath.Path) ([]Page, error) {
	top := buildTopLevel(m)

	var roots []*node
	isWholeModel := len(p.Segments) == 0
	if isWholeModel {
		roots = top
	} else {
		target, err := narrow(top, p.Segments)
		if err != nil {
			return nil, err
		}
		roots = []*node{target}
	}

	var emitted []*node
	flatten(roots, &emitted)

	// The path index maps an element's containment path to the file
	// its page lives in, so a cross-link resolves to a page only when
	// that page is part of this (possibly narrowed) output.
	index := map[string]string{}
	for _, n := range emitted {
		index[strings.Join(n.segments, "/")] = n.filePath()
	}

	var pages []Page
	if isWholeModel {
		pages = append(pages, Page{Path: "README.md", Content: renderRoot(top)})
	}
	for _, n := range emitted {
		pages = append(pages, Page{Path: n.filePath(), Content: renderPage(n, index)})
	}
	return pages, nil
}

// buildTopLevel builds the domains and, if the model has any context
// mappings, the context-mapping namespace that holds them.
func buildTopLevel(m *model.Model) []*node {
	var top []*node
	for _, domain := range sortedByName(m.Domains) {
		domainName := bareName(domain)
		top = append(top, newNode("domain", []string{domainName}, domain, buildDomainChildren(m, domainName)))
	}

	mappings := buildContextMappings(m)
	if len(mappings) > 0 {
		top = append(top, &node{
			kind:     "context-mappings",
			name:     "context-mapping",
			segments: []string{"context-mapping"},
			children: mappings,
		})
	}

	return top
}

// buildDomainChildren builds every element scoped to the given domain.
func buildDomainChildren(m *model.Model, domain string) []*node {
	var out []*node

	for _, view := range selectSorted(m.Subdomains, func(v model.SubdomainView) bool { return inScope(v, domain, "") }) {
		out = append(out, newNode("subdomain", []string{domain, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.EventHandlers, func(v model.EventHandlerView) bool { return inScope(v, domain, "") }) {
		out = append(out, newNode("event-handler", []string{domain, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.Policies, func(v model.PolicyView) bool { return inScope(v, domain, "") }) {
		out = append(out, newNode("policy", []string{domain, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.ExternalSystems, func(v model.ExternalSystemView) bool { return inScope(v, domain, "") }) {
		out = append(out, newNode("external-system", []string{domain, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.Extensions.DomainStorytelling.Stories, func(v model.DomainStoryView) bool { return inScope(v, domain, "") }) {
		out = append(out, newNode("domain-story", []string{domain, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.ProcessManagers, func(v model.ProcessManagerView) bool { return inScope(v, domain, "") }) {
		name := bareName(view)
		out = append(out, newNode("process-manager", []string{domain, name}, view, buildFeatures(m, []string{domain, name}, "process-manager", domain, "", name)))
	}
	for _, view := range selectSorted(m.BoundedContexts, func(v model.BoundedContextView) bool { return inScope(v, domain, "") }) {
		name := bareName(view)
		out = append(out, newNode("bounded-context", []string{domain, name}, view, buildBoundedContextChildren(m, domain, name)))
	}

	return out
}

// buildBoundedContextChildren builds every element scoped to the given
// bounded context, including its free-standing events.
func buildBoundedContextChildren(m *model.Model, domain, boundedContext string) []*node {
	var out []*node

	for _, view := range selectSorted(m.Aggregates, func(v model.AggregateView) bool { return inScope(v, domain, boundedContext) }) {
		name := bareName(view)
		out = append(out, newNode("aggregate", []string{domain, boundedContext, name}, view, buildAggregateChildren(m, domain, boundedContext, name)))
	}
	for _, view := range selectSorted(m.DynamicConsistencyBoundaries, func(v model.DynamicConsistencyBoundaryView) bool { return inScope(v, domain, boundedContext) }) {
		name := bareName(view)
		out = append(out, newNode("dynamic-consistency-boundary", []string{domain, boundedContext, name}, view, buildDcbChildren(m, domain, boundedContext, name)))
	}
	for _, view := range selectSorted(m.ReadModels, func(v model.ReadModelView) bool { return inScope(v, domain, boundedContext) }) {
		name := bareName(view)
		out = append(out, newNode("read-model", []string{domain, boundedContext, name}, view, buildFeatures(m, []string{domain, boundedContext, name}, "read-model", domain, boundedContext, name)))
	}
	for _, view := range selectSorted(m.Queries, func(v model.QueryView) bool { return inScope(v, domain, boundedContext) }) {
		out = append(out, newNode("query", []string{domain, boundedContext, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.Entities, func(v model.EntityView) bool { return inScope(v, domain, boundedContext) }) {
		out = append(out, newNode("entity", []string{domain, boundedContext, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.ValueObjects, func(v model.ValueObjectView) bool { return inScope(v, domain, boundedContext) }) {
		out = append(out, newNode("value-object", []string{domain, boundedContext, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.DomainServices, func(v model.DomainServiceView) bool { return inScope(v, domain, boundedContext) }) {
		out = append(out, newNode("domain-service", []string{domain, boundedContext, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.Actors, func(v model.ActorView) bool { return inScope(v, domain, boundedContext) }) {
		out = append(out, newNode("actor", []string{domain, boundedContext, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.Events, func(v model.EventView) bool {
		return inScope(v, domain, boundedContext) && model.ScopeText(v.Scope(), "aggregate") == ""
	}) {
		out = append(out, newNode("event", []string{domain, boundedContext, bareName(view)}, view, nil))
	}

	return out
}

// buildAggregateChildren builds the commands, events, and features of
// an aggregate.
func buildAggregateChildren(m *model.Model, domain, boundedContext, aggregate string) []*node {
	var out []*node

	for _, view := range selectSorted(m.Commands, func(v model.CommandView) bool {
		return inScope(v, domain, boundedContext) && model.ScopeText(v.Scope(), "aggregate") == aggregate
	}) {
		out = append(out, newNode("command", []string{domain, boundedContext, aggregate, bareName(view)}, view, nil))
	}
	for _, view := range selectSorted(m.Events, func(v model.EventView) bool {
		return inScope(v, domain, boundedContext) && model.ScopeText(v.Scope(), "aggregate") == aggregate
	}) {
		out = append(out, newNode("event", []string{domain, boundedContext, aggregate, bareName(view)}, view, nil))
	}
	out = append(out, buildFeatures(m, []string{domain, boundedContext, aggregate}, "aggregate", domain, boundedContext, aggregate)...)

	return out
}

// buildDcbChildren builds the commands and features of a dynamic
// consistency boundary.
func buildDcbChildren(m *model.Model, domain, boundedContext, dcb string) []*node {
	var out []*node

	for _, view := range selectSorted(m.Commands, func(v model.CommandView) bool {
		return inScope(v, domain, boundedContext) && model.ScopeText(v.Scope(), "dynamicConsistencyBoundary") == dcb
	}) {
		out = append(out, newNode("command", []string{domain, boundedContext, dcb, bareName(view)}, view, nil))
	}
	out = append(out, buildFeatures(m, []string{domain, boundedContext, dcb}, "dynamic-consistency-boundary", domain, boundedContext, dcb)...)

	return out
}

// buildFeatures builds the features that target the element at
// parentSegments, identified by its kind, domain, bounded context, and
// name.
func buildFeatures(m *model.Model, parentSegments []string, kind, domain, boundedContext, target string) []*node {
	var out []*node
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		featureKind, featureDomain, featureBC, featureTarget := featureParent(feature)
		if featureKind != kind || featureDomain != domain || featureBC != boundedContext || featureTarget != target {
			continue
		}
		segments := append(append([]string{}, parentSegments...), bareName(feature))
		out = append(out, newNode("feature", segments, feature, nil))
	}
	return out
}

// buildContextMappings builds the context-mapping leaves, addressed
// through the context-mapping namespace because they have no enclosing
// domain.
func buildContextMappings(m *model.Model) []*node {
	var out []*node
	for _, mapping := range sortedByName(m.ContextMappings) {
		out = append(out, newNode("context-mapping", []string{"context-mapping", bareName(mapping)}, mapping, nil))
	}
	return out
}

// featureParent reports the kind, domain, bounded context, and name of
// the element a feature targets. The bounded context is empty for a
// process-manager target, which is domain-scoped.
func featureParent(feature model.FeatureView) (kind, domain, boundedContext, target string) {
	scope := feature.Scope()
	domain = model.ScopeText(scope, "domain")

	switch {
	case scope.HasField("aggregate"):
		return "aggregate", domain, model.ScopeText(scope, "boundedContext"), model.ScopeText(scope, "aggregate")
	case scope.HasField("dynamicConsistencyBoundary"):
		return "dynamic-consistency-boundary", domain, model.ScopeText(scope, "boundedContext"), model.ScopeText(scope, "dynamicConsistencyBoundary")
	case scope.HasField("processManager"):
		return "process-manager", domain, "", model.ScopeText(scope, "processManager")
	case scope.HasField("readModel"):
		return "read-model", domain, model.ScopeText(scope, "boundedContext"), model.ScopeText(scope, "readModel")
	}

	return "", domain, "", ""
}

// narrow walks the top-level nodes segment by segment, returning the
// node the path identifies. Its error wording mirrors esdm view so
// both commands reject invalid input the same way.
func narrow(nodes []*node, segments []string) (*node, error) {
	var matched *node
	current := nodes

	for i, seg := range segments {
		var found *node
		for _, candidate := range current {
			if candidate.name == seg {
				found = candidate
				break
			}
		}
		if found == nil {
			if i == 0 {
				return nil, fmt.Errorf("no entity %q under model root", seg)
			}
			return nil, fmt.Errorf("no entity %q under %q", seg, strings.Join(segments[:i], "/"))
		}
		matched = found
		current = found.children
	}

	return matched, nil
}

// flatten appends every node in the trees to out, depth first.
func flatten(nodes []*node, out *[]*node) {
	for _, n := range nodes {
		*out = append(*out, n)
		flatten(n.children, out)
	}
}

// filePath is the node's relative slash path: a README.md index inside
// a directory when the node has children, otherwise a <name>.md leaf.
func (n *node) filePath() string {
	joined := strings.Join(n.segments, "/")
	if len(n.children) > 0 {
		return joined + "/README.md"
	}
	return joined + ".md"
}

// newNode builds a node from a view, reading its name and description
// and keeping the view for kind-specific detail rendering.
func newNode(kind string, segments []string, view documented, children []*node) *node {
	name, _ := view.Name().Text()
	description, _ := view.Description().Text()
	return &node{
		kind:        kind,
		name:        name,
		description: strings.TrimSpace(description),
		segments:    segments,
		view:        view,
		children:    children,
	}
}

// relLink is the relative link from the page at fromPath to the page
// at toPath, both relative slash paths within the tree.
func relLink(fromPath, toPath string) string {
	fromDir := path.Dir(fromPath)
	var fromParts []string
	if fromDir != "." && fromDir != "" {
		fromParts = strings.Split(fromDir, "/")
	}
	toParts := strings.Split(toPath, "/")

	shared := 0
	for shared < len(fromParts) && shared < len(toParts)-1 && fromParts[shared] == toParts[shared] {
		shared++
	}

	var rel []string
	for range fromParts[shared:] {
		rel = append(rel, "..")
	}
	rel = append(rel, toParts[shared:]...)
	return strings.Join(rel, "/")
}

// sortedByName returns the map's views ordered by their bare name.
func sortedByName[V documented](m map[string]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return bareName(out[i]) < bareName(out[j])
	})
	return out
}

// selectSorted returns the matching views ordered by their bare name.
func selectSorted[V scoped](m map[string]V, match func(V) bool) []V {
	var out []V
	for _, v := range m {
		if match(v) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return bareName(out[i]) < bareName(out[j])
	})
	return out
}

// bareName reads a view's bare name.
func bareName(view documented) string {
	name, _ := view.Name().Text()
	return name
}

// inScope reports whether a scoped view sits in the given domain and,
// when boundedContext is non-empty, that bounded context.
func inScope[V scoped](view V, domain, boundedContext string) bool {
	if model.ScopeText(view.Scope(), "domain") != domain {
		return false
	}
	if boundedContext != "" && model.ScopeText(view.Scope(), "boundedContext") != boundedContext {
		return false
	}
	return true
}
