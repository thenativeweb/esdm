package rules

import (
	"context"
	"fmt"
	"sort"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDAmbiguousName = "esdm/structure/ambiguous-name"

type ambiguousNameRule struct{}

func newAmbiguousNameRule() *ambiguousNameRule {
	return &ambiguousNameRule{}
}

func (*ambiguousNameRule) Meta() Meta {
	return Meta{
		ID:          ruleIDAmbiguousName,
		Severity:    diag.SeverityError,
		Description: "Within a bounded context, aggregates, dynamic consistency boundaries, entities, value objects, and domain services share one namespace: they become the types of one code module, with nothing composed into their names, so two of them cannot share a name.",
	}
}

// domainType is one element of the shared namespace: its
// kind and the node holding its name, which is also where
// the diagnostic points.
type domainType struct {
	kind           string
	boundedContext string
	name           ast.Node
}

// Check groups the five domain-type kinds by bounded context
// and name and reports every element beyond the first, in
// file order, so the first definition stays the reference
// point the note points at. Other kinds at the same position
// - read models, queries, actors, free-standing events - are
// deliberately not part of the namespace: an aggregate and a
// read model called `order` are the same concept in two
// roles, and CQRS code keeps them in separate modules.
func (*ambiguousNameRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	byContext := make(map[string]map[string][]domainType)
	add := func(kind string, scope ast.Node, name ast.Node) {
		boundedContext := scopeField(scope, "boundedContext")
		position := scopeField(scope, "domain") + "/" + boundedContext
		bare, ok := name.Text()
		if !ok {
			return
		}
		if byContext[position] == nil {
			byContext[position] = make(map[string][]domainType)
		}
		byContext[position][bare] = append(byContext[position][bare], domainType{kind: kind, boundedContext: boundedContext, name: name})
	}
	for _, agg := range m.Aggregates {
		add("aggregate", agg.Scope(), agg.Name())
	}
	for _, dcb := range m.DynamicConsistencyBoundaries {
		add("dynamic-consistency-boundary", dcb.Scope(), dcb.Name())
	}
	for _, entity := range m.Entities {
		add("entity", entity.Scope(), entity.Name())
	}
	for _, valueObject := range m.ValueObjects {
		add("value-object", valueObject.Scope(), valueObject.Name())
	}
	for _, domainService := range m.DomainServices {
		add("domain-service", domainService.Scope(), domainService.Name())
	}

	positions := make([]string, 0, len(byContext))
	for position := range byContext {
		positions = append(positions, position)
	}
	sort.Strings(positions)

	for _, position := range positions {
		byName := byContext[position]
		names := make([]string, 0, len(byName))
		for name := range byName {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			occurrences := byName[name]
			if len(occurrences) < 2 {
				continue
			}
			sort.Slice(occurrences, func(i, j int) bool {
				return isBefore(occurrences[i].name.Location(), occurrences[j].name.Location())
			})
			first := occurrences[0]
			for _, later := range occurrences[1:] {
				report.Report(diag.Diagnostic{
					Message:  fmt.Sprintf("%q names both %s and %s in bounded-context %q", name, withArticle(first.kind), withArticle(later.kind), first.boundedContext),
					Location: later.name.Location(),
					Related: []diag.Related{
						{
							Message:  fmt.Sprintf("%s %q defined here", first.kind, name),
							Location: first.name.Location(),
						},
					},
				})
			}
		}
	}
}

// isBefore orders two locations by file, then line, then
// column, which is the order a reader encounters them in.
func isBefore(a, b diag.Location) bool {
	if a.File != b.File {
		return a.File < b.File
	}
	if a.Line != b.Line {
		return a.Line < b.Line
	}
	return a.Column < b.Column
}

// withArticle prefixes a kind with its indefinite article,
// so the message reads as a sentence.
func withArticle(kind string) string {
	switch kind[0] {
	case 'a', 'e', 'i', 'o', 'u':
		return "an " + kind
	}
	return "a " + kind
}
