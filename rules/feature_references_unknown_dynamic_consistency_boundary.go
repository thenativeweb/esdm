package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDFeatureReferencesUnknownDynamicConsistencyBoundary = "esdm/gwt/feature-references-unknown-dynamic-consistency-boundary"

type featureReferencesUnknownDynamicConsistencyBoundaryRule struct{}

func newFeatureReferencesUnknownDynamicConsistencyBoundaryRule() *featureReferencesUnknownDynamicConsistencyBoundaryRule {
	return &featureReferencesUnknownDynamicConsistencyBoundaryRule{}
}

func (*featureReferencesUnknownDynamicConsistencyBoundaryRule) Meta() Meta {
	return Meta{
		ID:          ruleIDFeatureReferencesUnknownDynamicConsistencyBoundary,
		Extension:   "given-when-then",
		Severity:    diag.SeverityError,
		Description: "A Feature scoped to a Dynamic Consistency Boundary must name one the model declares. When the scope does not resolve, every Scenario in the Feature describes the behavior of something that does not exist.",
	}
}

func (*featureReferencesUnknownDynamicConsistencyBoundaryRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		if featureVariant(feature) != featureVariantDynamicConsistencyBoundary {
			continue
		}
		scope := feature.Scope()
		domain := scopeField(scope, "domain")
		boundedContext := scopeField(scope, "boundedContext")
		dcb := scopeField(scope, "dynamicConsistencyBoundary")
		if _, ok := m.LookupDynamicConsistencyBoundary(domain, boundedContext, dcb); ok {
			continue
		}
		featureName, _ := feature.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("feature %q scope references dynamic-consistency-boundary %q in bounded-context %q which is not declared in the model", featureName, dcb, boundedContext),
			Location: scope.Field("dynamicConsistencyBoundary").Location(),
		})
	}
}
