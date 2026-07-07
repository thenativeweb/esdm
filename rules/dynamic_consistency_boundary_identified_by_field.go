package rules

import (
	"context"
	"fmt"
	"sort"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDDynamicConsistencyBoundaryIdentifiedByField = "esdm/structure/dynamic-consistency-boundary-identified-by-field"

type dynamicConsistencyBoundaryIdentifiedByFieldRule struct{}

func newDynamicConsistencyBoundaryIdentifiedByFieldRule() *dynamicConsistencyBoundaryIdentifiedByFieldRule {
	return &dynamicConsistencyBoundaryIdentifiedByFieldRule{}
}

func (*dynamicConsistencyBoundaryIdentifiedByFieldRule) Meta() Meta {
	return Meta{
		ID:          ruleIDDynamicConsistencyBoundaryIdentifiedByField,
		Severity:    diag.SeverityError,
		Description: "When a DCB.identifiedBy entry uses source: command-payload, the named `field` must be declared in every triggering command's data.properties.",
	}
}

func (*dynamicConsistencyBoundaryIdentifiedByFieldRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, dcb := range sortedByName(m.DynamicConsistencyBoundaries) {
		dcbName, _ := dcb.Name().Text()
		dcbDomain := scopeField(dcb.Scope(), "domain")
		dcbBC := scopeField(dcb.Scope(), "boundedContext")

		commands := triggeringCommands(m, dcbDomain, dcbBC, dcbName)

		for _, entry := range dcb.IdentifiedBy().Seq() {
			source, _ := entry.Field("source").Text()
			if source != "command-payload" {
				continue
			}
			fieldNode := entry.Field("field")
			fieldName, ok := fieldNode.Text()
			if !ok {
				continue
			}

			for _, entry := range commands {
				if SchemaHasProperty(entry.view.Data(), fieldName) {
					continue
				}
				report.Report(diag.Diagnostic{
					Message: fmt.Sprintf(
						"DCB %q identifiedBy.field %q is not declared in command %q data.properties",
						dcbName, fieldName, entry.name,
					),
					Location: fieldNode.Location(),
				})
			}
		}
	}
}

type triggeringCommand struct {
	name string
	view model.CommandView
}

// triggeringCommands returns every command whose scope is
// (domain, boundedContext, dynamicConsistencyBoundary)
// matching the given DCB, paired with its bare name.
// Sorted by bare name for deterministic diagnostic order.
func triggeringCommands(m *model.Model, domain, boundedContext, dcb string) []triggeringCommand {
	var out []triggeringCommand
	for _, cmd := range m.Commands {
		scope := cmd.Scope()
		if scopeField(scope, "domain") != domain {
			continue
		}
		if scopeField(scope, "boundedContext") != boundedContext {
			continue
		}
		if scopeField(scope, "dynamicConsistencyBoundary") != dcb {
			continue
		}
		name, _ := cmd.Name().Text()
		out = append(out, triggeringCommand{name: name, view: cmd})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}
