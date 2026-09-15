package docgen

import (
	"fmt"
	"strings"
)

// childKindOrder fixes the order in which a page's child sections
// appear, so a page always lists its contents the same way.
var childKindOrder = []string{
	"domain",
	"subdomain",
	"bounded-context",
	"aggregate",
	"dynamic-consistency-boundary",
	"command",
	"event",
	"read-model",
	"query",
	"entity",
	"value-object",
	"domain-service",
	"actor",
	"process-manager",
	"event-handler",
	"policy",
	"external-system",
	"feature",
	"domain-story",
	"context-mappings",
}

// sectionHeadings maps a child kind to the heading of its section.
var sectionHeadings = map[string]string{
	"domain":                       "Domains",
	"subdomain":                    "Subdomains",
	"bounded-context":              "Bounded Contexts",
	"aggregate":                    "Aggregates",
	"dynamic-consistency-boundary": "Dynamic Consistency Boundaries",
	"command":                      "Commands",
	"event":                        "Events",
	"read-model":                   "Read Models",
	"query":                        "Queries",
	"entity":                       "Entities",
	"value-object":                 "Value Objects",
	"domain-service":               "Domain Services",
	"actor":                        "Actors",
	"process-manager":              "Process Managers",
	"event-handler":                "Event Handlers",
	"policy":                       "Policies",
	"external-system":              "External Systems",
	"feature":                      "Features",
	"domain-story":                 "Domain Stories",
	"context-mappings":             "Context Mappings",
}

// renderRoot renders the tree's entry page, indexing the domains and
// the context-mapping namespace.
func renderRoot(top []*node) string {
	var b strings.Builder
	b.WriteString("# Documentation\n")
	renderContents(&b, "README.md", top)
	return b.String()
}

// renderPage renders one element's page: its name, reference,
// description, kind-specific details, and an index of its children.
func renderPage(n *node, index map[string]string) string {
	var b strings.Builder

	if n.kind == "context-mappings" {
		b.WriteString("# Context Mappings\n")
		renderContents(&b, n.filePath(), n.children)
		return b.String()
	}

	fmt.Fprintf(&b, "# %s\n\n", n.name)
	fmt.Fprintf(&b, "Reference: `esdm:%s` (%s)\n", strings.Join(n.segments, "/"), n.kind)
	if n.description != "" {
		fmt.Fprintf(&b, "\n%s\n", n.description)
	}
	renderDetails(&b, n, index)
	renderContents(&b, n.filePath(), n.children)

	return b.String()
}

// renderContents writes one section per child kind, in a fixed order,
// linking each child relative to fromPath.
func renderContents(b *strings.Builder, fromPath string, children []*node) {
	if len(children) == 0 {
		return
	}

	groups := map[string][]*node{}
	for _, child := range children {
		groups[child.kind] = append(groups[child.kind], child)
	}

	for _, kind := range childKindOrder {
		group := groups[kind]
		if len(group) == 0 {
			continue
		}
		fmt.Fprintf(b, "\n## %s\n\n", sectionHeadings[kind])
		for _, child := range group {
			fmt.Fprintf(b, "- [%s](%s)\n", child.name, relLink(fromPath, child.filePath()))
		}
	}
}
