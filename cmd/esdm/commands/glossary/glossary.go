package glossary

import (
	"fmt"
	"sort"
	"strings"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/modelpath"
)

// Avoid is a discouraged alternative for a term, optionally
// carrying the reason it should not be used.
type Avoid struct {
	Term   string
	Reason string
}

// Term is a single glossary entry: the term, its definition,
// and any discouraged alternatives.
type Term struct {
	Term       string
	Definition string
	Avoid      []Avoid
}

// Section groups the glossary entries of one bounded
// context. A bounded context without ubiquitous language
// produces no section at all.
type Section struct {
	BoundedContext string
	Terms          []Term
}

// Glossary is the resolved, ordered glossary of a model:
// sections sorted by bounded-context name, terms sorted
// alphabetically within each section.
type Glossary struct {
	Sections []Section
}

// Build extracts the glossary from the resolved model,
// narrowed to the region identified by p. The empty path
// selects every bounded context; a one-segment path selects
// a domain; a two-segment path selects a single bounded
// context. An unknown or too-deep path segment is rejected
// as invalid input.
func Build(m *model.Model, p modelpath.Path) (*Glossary, error) {
	boundedContexts, err := selectBoundedContexts(m, p.Segments)
	if err != nil {
		return nil, err
	}

	g := &Glossary{}
	for _, boundedContext := range boundedContexts {
		terms := CollectTerms(boundedContext.UbiquitousLanguage())
		if len(terms) == 0 {
			continue
		}
		name, _ := boundedContext.Name().Text()
		g.Sections = append(g.Sections, Section{
			BoundedContext: name,
			Terms:          terms,
		})
	}
	return g, nil
}

// CollectTerms turns the ubiquitousLanguage sequence node
// into the sorted, typed term list. Entries missing a term
// or definition are skipped defensively, even though the
// schema requires both. It is exported so other commands (the
// documentation tree) can render the same terms without
// duplicating the extraction.
func CollectTerms(ubiquitousLanguage ast.Node) []Term {
	var terms []Term
	for _, entry := range ubiquitousLanguage.Seq() {
		term := strings.TrimSpace(textOf(entry, "term"))
		definition := strings.TrimSpace(textOf(entry, "definition"))
		if term == "" || definition == "" {
			continue
		}

		var avoid []Avoid
		for _, a := range entry.Field("avoid").Seq() {
			avoidTerm := strings.TrimSpace(textOf(a, "term"))
			if avoidTerm == "" {
				continue
			}
			avoid = append(avoid, Avoid{
				Term:   avoidTerm,
				Reason: strings.TrimSpace(textOf(a, "reason")),
			})
		}

		terms = append(terms, Term{
			Term:       term,
			Definition: definition,
			Avoid:      avoid,
		})
	}

	sort.Slice(terms, func(i, j int) bool {
		return terms[i].Term < terms[j].Term
	})
	return terms
}

// textOf reads a scalar string field, returning "" when the
// field is absent or not a scalar.
func textOf(n ast.Node, field string) string {
	v, _ := n.Field(field).Text()
	return v
}

// selectBoundedContexts resolves the path segments against
// the model's domain -> bounded-context hierarchy and returns
// the matching bounded contexts, sorted by name. The error
// wording mirrors the `esdm view` path narrowing so both
// commands reject invalid input the same way.
func selectBoundedContexts(m *model.Model, segments []string) ([]model.BoundedContextView, error) {
	if len(segments) == 0 {
		return sortedBoundedContexts(boundedContextsInDomain(m, "")), nil
	}

	domain := segments[0]
	if !domainExists(m, domain) {
		return nil, fmt.Errorf("no entity %q under model root", domain)
	}

	inDomain := sortedBoundedContexts(boundedContextsInDomain(m, domain))
	if len(segments) == 1 {
		return inDomain, nil
	}

	boundedContextName := segments[1]
	var match *model.BoundedContextView
	for i := range inDomain {
		name, _ := inDomain[i].Name().Text()
		if name == boundedContextName {
			match = &inDomain[i]
			break
		}
	}
	if match == nil {
		return nil, fmt.Errorf("no entity %q under %q", boundedContextName, domain)
	}

	if len(segments) > 2 {
		return nil, fmt.Errorf("no entity %q under %q", segments[2], strings.Join(segments[:2], "/"))
	}
	return []model.BoundedContextView{*match}, nil
}

// boundedContextsInDomain returns every bounded context
// scoped to the given domain. An empty domain matches every
// bounded context, which is how the empty-path case selects
// the whole model.
func boundedContextsInDomain(m *model.Model, domain string) []model.BoundedContextView {
	var out []model.BoundedContextView
	for _, boundedContext := range m.BoundedContexts {
		if domain == "" || textOf(boundedContext.Scope(), "domain") == domain {
			out = append(out, boundedContext)
		}
	}
	return out
}

// sortedBoundedContexts orders bounded contexts by their
// bare name so the glossary is deterministic.
func sortedBoundedContexts(boundedContexts []model.BoundedContextView) []model.BoundedContextView {
	sort.Slice(boundedContexts, func(i, j int) bool {
		nameI, _ := boundedContexts[i].Name().Text()
		nameJ, _ := boundedContexts[j].Name().Text()
		return nameI < nameJ
	})
	return boundedContexts
}

// domainExists reports whether the model contains a domain
// with the given name.
func domainExists(m *model.Model, domain string) bool {
	for _, d := range m.Domains {
		if name, _ := d.Name().Text(); name == domain {
			return true
		}
	}
	return false
}
