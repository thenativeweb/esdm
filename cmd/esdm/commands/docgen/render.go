package docgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/tree"
)

// renderRoot renders the index page of the whole tree: the
// domains and the context mappings, which have no domain to
// sit under.
func (g *generator) renderRoot() string {
	var b strings.Builder
	b.WriteString("# Model\n")

	b.WriteString("\n## Domains\n\n")
	for _, domain := range g.root.Children {
		fmt.Fprintf(&b, "- %s\n", g.childLine("README.md", domain))
	}

	mappings := g.contextMappings()
	if len(mappings) > 0 {
		b.WriteString("\n## Context Mappings\n\n")
		for _, mapping := range mappings {
			fmt.Fprintf(&b, "- %s\n", g.childLine("README.md", mapping))
		}
	}
	return b.String()
}

// contextMappings collects every context mapping the tree
// lists under any domain, once per mapping, sorted by name.
func (g *generator) contextMappings() []*tree.Node {
	seen := make(map[string]*tree.Node)
	for _, domain := range g.root.Children {
		for _, child := range domain.Children {
			if child.Kind == "context-mapping" {
				seen[child.Name] = child
			}
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]*tree.Node, 0, len(names))
	for _, name := range names {
		out = append(out, seen[name])
	}
	return out
}

// renderPage renders one element's page: header, stats,
// description, the kind's own section, details, and the
// children grouped by kind.
func (g *generator) renderPage(n *tree.Node) string {
	page := g.pathOf[n]
	var b strings.Builder

	fmt.Fprintf(&b, "# %s\n\n%s `%s`\n\n", n.Name, singular(n.Kind), g.refOf[n])

	if stats := statsLine(n); stats != "" && !hasRelations(n.Kind) && n.Kind != "context-mapping" {
		b.WriteString(stats + "\n\n")
	}

	description, _ := n.Document.Description().Text()
	if description != "" {
		b.WriteString(description + "\n\n")
	}

	switch n.Kind {
	case "command", "event", "query", "read-model":
		if relations := g.relations(page, n); relations != "" {
			b.WriteString(capitalize(relations) + ".\n\n")
		}
	case "bounded-context":
		b.WriteString(g.renderUbiquitousLanguage(n))
	case "context-mapping":
		b.WriteString(g.renderContextMapping(page, n))
	}

	if details := detailLines(n, description); len(details) > 0 {
		b.WriteString("## Details\n\n")
		for _, line := range details {
			fmt.Fprintf(&b, "- %s\n", line)
		}
		b.WriteString("\n")
	}

	b.WriteString(g.renderChildren(page, n))
	return strings.TrimRight(b.String(), "\n") + "\n"
}

// statsLine joins a node's tags and stats the way the view
// shows them next to the name.
func statsLine(n *tree.Node) string {
	var parts []string
	if len(n.Tags) > 0 {
		parts = append(parts, "("+strings.Join(n.Tags, ", ")+")")
	}
	if len(n.Stats) > 0 {
		parts = append(parts, strings.Join(n.Stats, " · "))
	}
	return strings.Join(parts, " · ")
}

// hasRelations names the kinds whose stats in the view are
// relationship annotations; the page writes those out
// instead of repeating the arrows.
func hasRelations(kind string) bool {
	return kind == "command" || kind == "event" || kind == "query" || kind == "read-model"
}

// detailLines returns the node's detail lines without the
// ones the page already shows elsewhere: the description,
// and a bounded context's terms, which its language section
// renders in full.
func detailLines(n *tree.Node, description string) []string {
	if n.Kind == "bounded-context" {
		return nil
	}
	var out []string
	for _, line := range n.Lines {
		if line == description || strings.HasPrefix(line, "description: ") {
			continue
		}
		out = append(out, line)
	}
	return out
}

// renderChildren writes one section per kind of child, in
// the order the tree lists them, each child as a link with
// its relations or its stats.
func (g *generator) renderChildren(page string, n *tree.Node) string {
	var b strings.Builder
	var currentKind string
	for _, child := range n.Children {
		if child.Kind != currentKind {
			if currentKind != "" {
				b.WriteString("\n")
			}
			fmt.Fprintf(&b, "## %s\n\n", plural(child.Kind))
			currentKind = child.Kind
		}
		fmt.Fprintf(&b, "- %s\n", g.childLine(page, child))
	}
	return b.String()
}

// childLine renders a child as a link followed by what
// distinguishes it: written-out relations for commands,
// events and queries, tags and stats for everything else.
func (g *generator) childLine(page string, child *tree.Node) string {
	line := g.link(page, child)
	if hasRelations(child.Kind) {
		if relations := g.relations(page, child); relations != "" {
			line += " – " + relations
		}
		return line
	}
	if len(child.Tags) > 0 {
		line += " (" + strings.Join(child.Tags, ", ") + ")"
	}
	if len(child.Stats) > 0 {
		line += " · " + strings.Join(child.Stats, " · ")
	}
	return line
}

// link renders a Markdown link from page to the node's page,
// or the node's name with its reference when the node lies
// outside the written tree, so that no link is ever broken.
func (g *generator) link(page string, target *tree.Node) string {
	if g.selected[target] {
		return fmt.Sprintf("[%s](%s)", target.Name, relativeLink(page, g.pathOf[target]))
	}
	return fmt.Sprintf("%s (`%s`)", target.Name, g.refOf[target])
}

// linkDocument links to the page of a model document, or
// falls back to the bare name when the tree did not place it.
func (g *generator) linkDocument(page string, document model.DocumentViewBase, name string) string {
	if target, ok := g.nodeFor(document); ok {
		return g.link(page, target)
	}
	return name
}

// relations writes a command's, event's or query's
// relationships out in words, with each counterpart linked.
func (g *generator) relations(page string, n *tree.Node) string {
	switch n.Kind {
	case "command":
		cmd := model.CommandView{DocumentViewBase: n.Document}
		var parts []string
		if published := g.publishedEvents(page, cmd); len(published) > 0 {
			parts = append(parts, "publishes "+joinComma(published))
		}
		if actors := g.actors(page, cmd.Scope(), cmd.Actors()); len(actors) > 0 {
			parts = append(parts, "issued by "+joinComma(actors))
		}
		return strings.Join(parts, ", ")
	case "event":
		event := model.EventView{DocumentViewBase: n.Document}
		var publishers []string
		for _, cmd := range g.m.PublishersOf(event) {
			name, _ := cmd.Name().Text()
			publishers = append(publishers, g.linkDocument(page, cmd.DocumentViewBase, name))
		}
		if len(publishers) == 0 {
			return ""
		}
		return "published by " + joinComma(publishers)
	case "read-model":
		readModel := model.ReadModelView{DocumentViewBase: n.Document}
		domain := model.ScopeText(readModel.Scope(), "domain")
		var projected []string
		for _, projection := range readModel.Projections().Seq() {
			boundedContext, _ := projection.Field("boundedContext").Text()
			aggregate, _ := projection.Field("aggregate").Text()
			name, ok := projection.Field("event").Text()
			if !ok {
				continue
			}
			event, exists := g.m.LookupEvent(domain, boundedContext, aggregate, name)
			if !exists {
				projected = append(projected, name)
				continue
			}
			projected = append(projected, g.linkDocument(page, event.DocumentViewBase, name))
		}
		if len(projected) == 0 {
			return ""
		}
		return "projects " + joinComma(projected)
	case "query":
		query := model.QueryView{DocumentViewBase: n.Document}
		var parts []string
		if readModelName, ok := query.ReadModel().Text(); ok {
			domain := model.ScopeText(query.Scope(), "domain")
			boundedContext := model.ScopeText(query.Scope(), "boundedContext")
			text := readModelName
			if readModel, exists := g.m.LookupReadModel(domain, boundedContext, readModelName); exists {
				text = g.linkDocument(page, readModel.DocumentViewBase, readModelName)
			}
			parts = append(parts, "reads "+text)
		}
		if actors := g.actors(page, query.Scope(), query.Actors()); len(actors) > 0 {
			parts = append(parts, "issued by "+joinComma(actors))
		}
		return strings.Join(parts, ", ")
	}
	return ""
}

// publishedEvents links every event the command publishes,
// resolved the way the model defines the relation.
func (g *generator) publishedEvents(page string, cmd model.CommandView) []string {
	var out []string
	for _, item := range cmd.Publishes().Seq() {
		name, ok := item.Text()
		if !ok {
			continue
		}
		event, exists := g.m.PublishedEvent(cmd, name)
		if !exists {
			out = append(out, name)
			continue
		}
		out = append(out, g.linkDocument(page, event.DocumentViewBase, name))
	}
	return out
}

// actors links the actors a command or query names, looked
// up in the element's own bounded context.
func (g *generator) actors(page string, scope ast.Node, actors ast.Node) []string {
	domain := model.ScopeText(scope, "domain")
	boundedContext := model.ScopeText(scope, "boundedContext")
	var out []string
	for _, item := range actors.Seq() {
		name, ok := item.Text()
		if !ok {
			continue
		}
		actor, exists := g.m.LookupActor(domain, boundedContext, name)
		if !exists {
			out = append(out, name)
			continue
		}
		out = append(out, g.linkDocument(page, actor.DocumentViewBase, name))
	}
	return out
}

// renderUbiquitousLanguage renders a bounded context's terms
// sorted by name, each with its rejected alternatives and
// its translations indented beneath it, one line per
// language.
func (g *generator) renderUbiquitousLanguage(n *tree.Node) string {
	boundedContext := model.BoundedContextView{DocumentViewBase: n.Document}
	entries := boundedContext.UbiquitousLanguage().Seq()
	if len(entries) == 0 {
		return ""
	}
	sort.Slice(entries, func(i, j int) bool {
		ti, _ := entries[i].Field("term").Text()
		tj, _ := entries[j].Field("term").Text()
		return ti < tj
	})

	var b strings.Builder
	b.WriteString("## Ubiquitous Language\n\n")
	if language, ok := boundedContext.Language().Text(); ok {
		fmt.Fprintf(&b, "Written in `%s`.\n\n", language)
	}
	for _, entry := range entries {
		term, _ := entry.Field("term").Text()
		definition, _ := entry.Field("definition").Text()
		fmt.Fprintf(&b, "- **%s** – %s\n", term, definition)
		for _, avoid := range entry.Field("avoid").Seq() {
			fmt.Fprintf(&b, "  - %s\n", avoidNote(avoid))
		}
		for _, translation := range entry.Field("translations").Seq() {
			language, _ := translation.Field("language").Text()
			translatedTerm, _ := translation.Field("term").Text()
			translatedDefinition, _ := translation.Field("definition").Text()
			fmt.Fprintf(&b, "  - `%s` **%s** – %s", language, translatedTerm, translatedDefinition)
			for _, avoid := range translation.Field("avoid").Seq() {
				b.WriteString(" " + avoidNote(avoid))
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
	return b.String()
}

// avoidNote renders one rejected alternative the way the
// glossary does: an italic sentence, then the reason.
func avoidNote(avoid ast.Node) string {
	term, _ := avoid.Field("term").Text()
	note := fmt.Sprintf("_Avoid the term %q._", term)
	if reason, ok := avoid.Field("reason").Text(); ok && reason != "" {
		note += " " + reason
	}
	return note
}

// renderContextMapping renders the mapping's type with its
// endpoints and, for the asymmetric types, the term pairs.
func (g *generator) renderContextMapping(page string, n *tree.Node) string {
	mapping := model.ContextMappingView{DocumentViewBase: n.Document}
	mappingType, _ := mapping.Type().Text()

	var b strings.Builder
	roles, isAsymmetric := model.MappingRoles(mappingType)
	if isAsymmetric {
		fmt.Fprintf(&b, "%s: %s %s, %s %s\n\n", mappingType, roles[0], g.endpointLink(page, mapping.Field(roles[0])), roles[1], g.endpointLink(page, mapping.Field(roles[1])))
	} else {
		var participants []string
		for _, participant := range mapping.Participants().Seq() {
			participants = append(participants, g.endpointLink(page, participant))
		}
		fmt.Fprintf(&b, "%s: participants %s\n\n", mappingType, joinComma(participants))
	}

	terms := mapping.Terms().Seq()
	if isAsymmetric && len(terms) > 0 {
		b.WriteString("## Terms\n\n")
		for _, pair := range terms {
			first, _ := pair.Field(roles[0]).Text()
			second, _ := pair.Field(roles[1]).Text()
			fmt.Fprintf(&b, "- %s in %s corresponds to %s in %s\n", first, g.endpointLink(page, mapping.Field(roles[0])), second, g.endpointLink(page, mapping.Field(roles[1])))
		}
		b.WriteString("\n")
	}
	return b.String()
}

// endpointLink links a mapping endpoint, a bounded context
// or an external system, to its page.
func (g *generator) endpointLink(page string, endpoint ast.Node) string {
	domain := model.ScopeText(endpoint, "domain")
	if name, ok := endpoint.Field("boundedContext").Text(); ok {
		if boundedContext, exists := g.m.LookupBoundedContext(domain, name); exists {
			return g.linkDocument(page, boundedContext.DocumentViewBase, name)
		}
		return name
	}
	if name, ok := endpoint.Field("externalSystem").Text(); ok {
		if externalSystem, exists := g.m.LookupExternalSystem(domain, name); exists {
			return g.linkDocument(page, externalSystem.DocumentViewBase, name)
		}
		return name
	}
	return ""
}

func joinComma(values []string) string {
	return strings.Join(values, ", ")
}

func capitalize(text string) string {
	if text == "" {
		return text
	}
	return strings.ToUpper(text[:1]) + text[1:]
}
