package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDFeatureReferencesUnknownAggregate = "esdm/gwt/feature-references-unknown-aggregate"

type featureReferencesUnknownAggregateRule struct{}

func newFeatureReferencesUnknownAggregateRule() *featureReferencesUnknownAggregateRule {
	return &featureReferencesUnknownAggregateRule{}
}

func (*featureReferencesUnknownAggregateRule) Meta() Meta {
	return Meta{
		ID:          ruleIDFeatureReferencesUnknownAggregate,
		Severity:    diag.SeverityError,
		Description: "An aggregate-scoped feature must point at a declared aggregate; an unresolved scope means every scenario underneath is testing nothing.",
	}
}

func (*featureReferencesUnknownAggregateRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		if featureVariant(feature) != featureVariantAggregate {
			continue
		}
		scope := feature.Scope()
		domain := scopeField(scope, "domain")
		boundedContext := scopeField(scope, "boundedContext")
		aggregate := scopeField(scope, "aggregate")
		if _, ok := m.LookupAggregate(domain, boundedContext, aggregate); ok {
			continue
		}
		featureName, _ := feature.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("feature %q scope references aggregate %q in bounded-context %q which is not declared in the model", featureName, aggregate, boundedContext),
			Location: scope.Field("aggregate").Location(),
		})
	}
}
