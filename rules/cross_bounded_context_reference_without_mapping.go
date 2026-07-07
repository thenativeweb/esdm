package rules

import (
	"context"
	"fmt"
	"sort"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDCrossBoundedContextReferenceWithoutMapping = "esdm/modeling/cross-bounded-context-reference-without-mapping"

type crossBoundedContextReferenceWithoutMappingRule struct{}

func newCrossBoundedContextReferenceWithoutMappingRule() *crossBoundedContextReferenceWithoutMappingRule {
	return &crossBoundedContextReferenceWithoutMappingRule{}
}

func (*crossBoundedContextReferenceWithoutMappingRule) Meta() Meta {
	return Meta{
		ID:          ruleIDCrossBoundedContextReferenceWithoutMapping,
		Severity:    diag.SeverityWarning,
		Description: "When a consumer touches more than one bounded context (read-model projection across BCs, domain-scoped event-handler/policy/process-manager spanning BCs), each pair of touched BCs should be linked by a context-mapping.",
	}
}

func (*crossBoundedContextReferenceWithoutMappingRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	mapped := mappedBCPairs(m)

	for _, eh := range sortedByName(m.EventHandlers) {
		name, _ := eh.Name().Text()
		reportMissingPairs(name, "event-handler", collectEventHandlerBCs(eh), mapped, eh.Name().Location(), report)
	}
	for _, p := range sortedByName(m.Policies) {
		name, _ := p.Name().Text()
		reportMissingPairs(name, "policy", collectPolicyBCs(p), mapped, p.Name().Location(), report)
	}
	for _, pm := range sortedByName(m.ProcessManagers) {
		name, _ := pm.Name().Text()
		reportMissingPairs(name, "process-manager", collectProcessManagerBCs(pm), mapped, pm.Name().Location(), report)
	}
	for _, rm := range sortedByName(m.ReadModels) {
		name, _ := rm.Name().Text()
		reportMissingPairs(name, "read-model", collectReadModelBCs(rm), mapped, rm.Name().Location(), report)
	}
}

// mappedBCPairs returns the set of unordered (BC, BC)
// pairs covered by any context-mapping. Endpoints that
// are external systems are ignored, since this rule cares
// about BC-to-BC integration only.
func mappedBCPairs(m *model.Model) map[[2]string]bool {
	out := make(map[[2]string]bool)
	for _, cm := range m.ContextMappings {
		bcs := contextMappingBCs(cm)
		for i := 0; i < len(bcs); i++ {
			for j := i + 1; j < len(bcs); j++ {
				if bcs[i] == bcs[j] {
					continue
				}
				out[canonicalPair(bcs[i], bcs[j])] = true
			}
		}
	}
	return out
}

func contextMappingBCs(cm model.ContextMappingView) []string {
	var bcs []string
	for _, ep := range []ast.Node{
		cm.Customer(), cm.Supplier(), cm.Conformist(),
		cm.Upstream(), cm.Downstream(), cm.Host(),
		cm.Consumer(), cm.Publisher(),
	} {
		if name, ok := ep.Field("boundedContext").Text(); ok {
			bcs = append(bcs, name)
		}
	}
	for _, p := range cm.Participants().Seq() {
		if name, ok := p.Field("boundedContext").Text(); ok {
			bcs = append(bcs, name)
		}
	}
	return bcs
}

func canonicalPair(a, b string) [2]string {
	if a < b {
		return [2]string{a, b}
	}
	return [2]string{b, a}
}

func collectEventHandlerBCs(eh model.EventHandlerView) []string {
	bcs := make(map[string]bool)
	for _, ref := range eh.Handles().Seq() {
		if name, ok := ref.Field("boundedContext").Text(); ok {
			bcs[name] = true
		}
	}
	return sortedSet(bcs)
}

func collectPolicyBCs(p model.PolicyView) []string {
	bcs := make(map[string]bool)
	for _, ref := range p.Handles().Seq() {
		if name, ok := ref.Field("boundedContext").Text(); ok {
			bcs[name] = true
		}
	}
	for _, ref := range p.Emits().Seq() {
		if name, ok := ref.Field("boundedContext").Text(); ok {
			bcs[name] = true
		}
	}
	return sortedSet(bcs)
}

func collectProcessManagerBCs(pm model.ProcessManagerView) []string {
	bcs := make(map[string]bool)
	for _, ref := range pm.StartsWhen().Seq() {
		if name, ok := ref.Field("boundedContext").Text(); ok {
			bcs[name] = true
		}
	}
	for _, reaction := range pm.Reactions().Seq() {
		when := reaction.Field("when")
		if !when.HasField("timer") {
			if name, ok := when.Field("boundedContext").Text(); ok {
				bcs[name] = true
			}
		}
		for _, ref := range reaction.Field("emits").Seq() {
			if name, ok := ref.Field("boundedContext").Text(); ok {
				bcs[name] = true
			}
		}
	}
	return sortedSet(bcs)
}

func collectReadModelBCs(rm model.ReadModelView) []string {
	bcs := make(map[string]bool)
	if name := scopeField(rm.Scope(), "boundedContext"); name != "" {
		bcs[name] = true
	}
	for _, proj := range rm.Projections().Seq() {
		if name, ok := proj.Field("boundedContext").Text(); ok {
			bcs[name] = true
		}
	}
	return sortedSet(bcs)
}

func sortedSet(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func reportMissingPairs(consumerName, consumerKind string, touchedBCs []string, mapped map[[2]string]bool, loc diag.Location, report diag.Reporter) {
	for i := 0; i < len(touchedBCs); i++ {
		for j := i + 1; j < len(touchedBCs); j++ {
			pair := canonicalPair(touchedBCs[i], touchedBCs[j])
			if mapped[pair] {
				continue
			}
			report.Report(diag.Diagnostic{
				Message: fmt.Sprintf(
					"%s %q crosses bounded contexts %q and %q without a context-mapping linking them",
					consumerKind, consumerName, pair[0], pair[1],
				),
				Location: loc,
			})
		}
	}
}
