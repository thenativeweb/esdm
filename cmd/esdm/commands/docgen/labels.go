package docgen

// label is how a kind reads on a page: as a heading for a
// group of children and as the noun in an element's header.
type label struct {
	singular string
	plural   string
}

// labels maps every kind the tree places to its labels. The
// kinds are the schema's kind enum values plus the extension
// kinds; a kind without an entry falls back to its raw name.
var labels = map[string]label{
	"domain":                       {"Domain", "Domains"},
	"subdomain":                    {"Subdomain", "Subdomains"},
	"bounded-context":              {"Bounded Context", "Bounded Contexts"},
	"context-mapping":              {"Context Mapping", "Context Mappings"},
	"aggregate":                    {"Aggregate", "Aggregates"},
	"dynamic-consistency-boundary": {"Dynamic Consistency Boundary", "Dynamic Consistency Boundaries"},
	"command":                      {"Command", "Commands"},
	"event":                        {"Event", "Events"},
	"event-handler":                {"Event Handler", "Event Handlers"},
	"policy":                       {"Policy", "Policies"},
	"process-manager":              {"Process Manager", "Process Managers"},
	"read-model":                   {"Read Model", "Read Models"},
	"query":                        {"Query", "Queries"},
	"entity":                       {"Entity", "Entities"},
	"value-object":                 {"Value Object", "Value Objects"},
	"domain-service":               {"Domain Service", "Domain Services"},
	"actor":                        {"Actor", "Actors"},
	"external-system":              {"External System", "External Systems"},
	"domain-story":                 {"Domain Story", "Domain Stories"},
	"feature":                      {"Feature", "Features"},
}

func singular(kind string) string {
	if l, ok := labels[kind]; ok {
		return l.singular
	}
	return kind
}

func plural(kind string) string {
	if l, ok := labels[kind]; ok {
		return l.plural
	}
	return kind
}

// containerKinds are rendered as a directory with a
// README.md, because elements sit beneath them. The set is
// fixed per kind rather than derived from whether a node
// happens to have children, so an element's page path does
// not move when its first child appears.
var containerKinds = map[string]bool{
	"domain":                       true,
	"bounded-context":              true,
	"aggregate":                    true,
	"dynamic-consistency-boundary": true,
	"process-manager":              true,
	"read-model":                   true,
}
