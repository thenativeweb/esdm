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
		Extension:   "given-when-then",
		Severity:    diag.SeverityError,
		Description: "A Feature scoped to an Aggregate must name an Aggregate the model declares. When the scope does not resolve, every Scenario in the Feature describes the behavior of something that does not exist.",
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
