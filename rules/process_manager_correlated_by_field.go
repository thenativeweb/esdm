package rules

import (
	"context"
	"fmt"
	"sort"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDProcessManagerCorrelatedByField = "esdm/structure/process-manager-correlated-by-field"

type processManagerCorrelatedByFieldRule struct{}

func newProcessManagerCorrelatedByFieldRule() *processManagerCorrelatedByFieldRule {
	return &processManagerCorrelatedByFieldRule{}
}

func (*processManagerCorrelatedByFieldRule) Meta() Meta {
	return Meta{
		ID:          ruleIDProcessManagerCorrelatedByField,
		Severity:    diag.SeverityError,
		Description: "When process-manager.correlatedBy uses source: event-field, the named `field` must be declared in every consumed event's data.properties.",
	}
}

func (*processManagerCorrelatedByFieldRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, pm := range sortedByName(m.ProcessManagers) {
		pmName, _ := pm.Name().Text()
		pmDomain := scopeField(pm.Scope(), "domain")

		correlatedBy := pm.CorrelatedBy()
		source, _ := correlatedBy.Field("source").Text()
		if source != "event-field" {
			continue
		}
		fieldNode := correlatedBy.Field("field")
		fieldName, ok := fieldNode.Text()
		if !ok {
			continue
		}

		for _, ref := range consumedEventRefs(pm) {
			event, ok := lookupEventByRef(m, pmDomain, ref)
			if !ok {
				continue
			}
			if SchemaHasProperty(event.Data(), fieldName) {
				continue
			}
			eventName, _ := event.Name().Text()
			report.Report(diag.Diagnostic{
				Message: fmt.Sprintf(
					"process-manager %q correlatedBy.field %q is not declared in event %q data.properties",
					pmName, fieldName, eventName,
				),
				Location: fieldNode.Location(),
			})
		}
	}
}

// consumedEventRefs returns the eventReference nodes for
// every event the process manager consumes - startsWhen
// entries plus reactions whose `when` is an event
// reference rather than a timer reference. The triples
// allow scope-aware lookup of the referenced events.
// References are deduplicated and sorted by composite
// (boundedContext, aggregate, event) tuple for
// deterministic diagnostic order.
func consumedEventRefs(pm model.ProcessManagerView) []ast.Node {
	type tuple struct {
		bc, aggregate, event string
		ref                  ast.Node
	}
	seen := make(map[string]struct{})
	var collected []tuple
	push := func(ref ast.Node) {
		bc, _ := ref.Field("boundedContext").Text()
		event, _ := ref.Field("event").Text()
		var aggregate string
		if aggNode := ref.Field("aggregate"); aggNode.Exists() {
			aggregate, _ = aggNode.Text()
		}
		key := bc + "/" + aggregate + "/" + event
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		collected = append(collected, tuple{bc: bc, aggregate: aggregate, event: event, ref: ref})
	}
	for _, ref := range pm.StartsWhen().Seq() {
		push(ref)
	}
	for _, reaction := range pm.Reactions().Seq() {
		when := reaction.Field("when")
		if when.HasField("timer") {
			continue
		}
		push(when)
	}
	sort.Slice(collected, func(i, j int) bool {
		if collected[i].bc != collected[j].bc {
			return collected[i].bc < collected[j].bc
		}
		if collected[i].aggregate != collected[j].aggregate {
			return collected[i].aggregate < collected[j].aggregate
		}
		return collected[i].event < collected[j].event
	})
	out := make([]ast.Node, 0, len(collected))
	for _, entry := range collected {
		out = append(out, entry.ref)
	}
	return out
}

// lookupEventByRef resolves an eventReference node
// (carrying boundedContext, optional aggregate, and
// event) to its EventView within the given domain.
// Returns (zero, false) when the reference does not
// resolve.
func lookupEventByRef(m *model.Model, domain string, ref ast.Node) (model.EventView, bool) {
	bc, _ := ref.Field("boundedContext").Text()
	event, _ := ref.Field("event").Text()
	if bc == "" || event == "" {
		return model.EventView{}, false
	}
	var aggregate string
	if aggNode := ref.Field("aggregate"); aggNode.Exists() {
		aggregate, _ = aggNode.Text()
	}
	v, exists := m.LookupEvent(domain, bc, aggregate, event)
	return v, exists
}
