package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDOrphanContextMapping = "esdm/modeling/orphan-context-mapping"

type orphanContextMappingRule struct{}

func newOrphanContextMappingRule() *orphanContextMappingRule {
	return &orphanContextMappingRule{}
}

func (*orphanContextMappingRule) Meta() Meta {
	return Meta{
		ID:          ruleIDOrphanContextMapping,
		Severity:    diag.SeverityWarning,
		Description: "A context-mapping linking two bounded contexts should be backed by an actual cross-BC reference somewhere in the model; otherwise it documents a relationship that nothing in the model uses.",
	}
}

func (*orphanContextMappingRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	traversed := traversedBCPairs(m)

	for _, cm := range sortedByName(m.ContextMappings) {
		bcs := contextMappingBCs(cm)
		// Mapping has no BC pair (e.g. external-system to
		// external-system) - not in this rule's scope.
		if !hasDistinctPair(bcs) {
			continue
		}
		if anyPairIsTraversed(bcs, traversed) {
			continue
		}

		name, _ := cm.Name().Text()
		pair := firstDistinctPair(bcs)
		report.Report(diag.Diagnostic{
			Message: fmt.Sprintf(
				"context-mapping %q links bounded contexts %q and %q, but no consumer in the model crosses them",
				name, pair[0], pair[1],
			),
			Location: cm.Name().Location(),
		})
	}
}

// traversedBCPairs collects every (BC, BC) pair that some
// consumer (event-handler, policy, process-manager,
// read-model) actually spans, so context-mapping
// orphanage can be decided in a single map lookup.
func traversedBCPairs(m *model.Model) map[[2]string]bool {
	out := make(map[[2]string]bool)

	addPairs := func(bcs []string) {
		for i := 0; i < len(bcs); i++ {
			for j := i + 1; j < len(bcs); j++ {
				if bcs[i] == bcs[j] {
					continue
				}
				out[canonicalPair(bcs[i], bcs[j])] = true
			}
		}
	}

	for _, eh := range m.EventHandlers {
		addPairs(collectEventHandlerBCs(eh))
	}
	for _, p := range m.Policies {
		addPairs(collectPolicyBCs(p))
	}
	for _, pm := range m.ProcessManagers {
		addPairs(collectProcessManagerBCs(pm))
	}
	for _, rm := range m.ReadModels {
		addPairs(collectReadModelBCs(rm))
	}

	return out
}

func hasDistinctPair(bcs []string) bool {
	for i := 0; i < len(bcs); i++ {
		for j := i + 1; j < len(bcs); j++ {
			if bcs[i] != bcs[j] {
				return true
			}
		}
	}
	return false
}

func anyPairIsTraversed(bcs []string, traversed map[[2]string]bool) bool {
	for i := 0; i < len(bcs); i++ {
		for j := i + 1; j < len(bcs); j++ {
			if bcs[i] == bcs[j] {
				continue
			}
			if traversed[canonicalPair(bcs[i], bcs[j])] {
				return true
			}
		}
	}
	return false
}

func firstDistinctPair(bcs []string) [2]string {
	for i := 0; i < len(bcs); i++ {
		for j := i + 1; j < len(bcs); j++ {
			if bcs[i] != bcs[j] {
				return canonicalPair(bcs[i], bcs[j])
			}
		}
	}
	return [2]string{}
}
