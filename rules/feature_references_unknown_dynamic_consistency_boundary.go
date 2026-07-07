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
		Severity:    diag.SeverityError,
		Description: "A DCB-scoped feature must point at a declared dynamic-consistency-boundary; an unresolved scope means every scenario underneath is testing nothing.",
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
