package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDUncoveredInvariant = "esdm/gwt/uncovered-invariant"

type uncoveredInvariantRule struct{}

func newUncoveredInvariantRule() *uncoveredInvariantRule {
	return &uncoveredInvariantRule{}
}

func (*uncoveredInvariantRule) Meta() Meta {
	return Meta{
		ID:          ruleIDUncoveredInvariant,
		Severity:    diag.SeverityWarning,
		Description: "On an aggregate or dynamic-consistency-boundary that has a Given-When-Then feature, every named invariant should be covered by at least one scenario; an invariant that no scenario exercises is untested. Units without a Given-When-Then feature are exempt.",
	}
}

func (*uncoveredInvariantRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	// Collect, per Given-When-Then-tested unit, which named
	// invariants a scenario rejection references. Only aggregate and
	// DCB features carry a then.rejection, so only those two unit
	// kinds can be covered this way - process managers, value objects
	// and entities have no rejection mechanism and are out of scope,
	// exactly like the sibling rule
	// scenario-rejection-references-unknown-invariant.
	covered := make(map[string]map[string]bool)
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		variant := featureVariant(feature)
		scope := feature.Scope()
		var key string
		switch variant {
		case featureVariantAggregate:
			key = unitKey(variant, scopeField(scope, "domain"), scopeField(scope, "boundedContext"), scopeField(scope, "aggregate"))
		case featureVariantDynamicConsistencyBoundary:
			key = unitKey(variant, scopeField(scope, "domain"), scopeField(scope, "boundedContext"), scopeField(scope, "dynamicConsistencyBoundary"))
		default:
			continue
		}
		refs, ok := covered[key]
		if !ok {
			refs = make(map[string]bool)
			covered[key] = refs
		}
		for _, scenario := range scenariosOf(feature) {
			if invariant, ok := scenario.Then().Field("rejection").Field("invariant").Text(); ok {
				refs[invariant] = true
			}
		}
	}

	for _, agg := range sortedByName(m.Aggregates) {
		scope := agg.Scope()
		name, _ := agg.Name().Text()
		key := unitKey(featureVariantAggregate, scopeField(scope, "domain"), scopeField(scope, "boundedContext"), name)
		reportUncoveredInvariants(agg.Field("invariants"), featureVariantAggregate, name, covered[key], report)
	}
	for _, dcb := range sortedByName(m.DynamicConsistencyBoundaries) {
		scope := dcb.Scope()
		name, _ := dcb.Name().Text()
		key := unitKey(featureVariantDynamicConsistencyBoundary, scopeField(scope, "domain"), scopeField(scope, "boundedContext"), name)
		reportUncoveredInvariants(dcb.Field("invariants"), featureVariantDynamicConsistencyBoundary, name, covered[key], report)
	}
}

// unitKey builds a lookup key that keeps aggregates and DCBs of the
// same domain/bounded-context/name apart by tagging the variant.
func unitKey(variant, domain, boundedContext, name string) string {
	return variant + "\x00" + domain + "/" + boundedContext + "/" + name
}

// reportUncoveredInvariants flags each named invariant of a unit that
// no scenario rejection references. A nil coveredRefs means the unit
// has no Given-When-Then feature at all; the rule then stays silent,
// because Given-When-Then is optional.
func reportUncoveredInvariants(invariants ast.Node, variant, unitName string, coveredRefs map[string]bool, report diag.Reporter) {
	if coveredRefs == nil {
		return
	}
	for _, invariant := range invariants.Seq() {
		name, ok := invariant.Field("name").Text()
		if !ok {
			continue
		}
		if coveredRefs[name] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("invariant %q of %s %q is not covered by any scenario", name, variant, unitName),
			Location: invariant.Field("name").Location(),
		})
	}
}
