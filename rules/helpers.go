package rules

import (
	"sort"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/model"
)

// scopeField reads a string field from a scope node,
// returning "" when the field is missing or not a string
// scalar.
func scopeField(scope ast.Node, field string) string {
	v, _ := scope.Field(field).Text()
	return v
}

// sortedByName returns the values of a Model map sorted
// by their bare entity name. Map keys in Model are
// composite (scope tuple plus name) so rules cannot rely
// on the key for human-readable diagnostic messages;
// pulling the bare name from view.Name() and sorting on
// it gives both the right value to print and a stable
// iteration order.
func sortedByName[V interface{ Name() ast.Node }](m map[string]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

// Feature-scope variant constants. These match the
// kebab-case core-schema kind names so diagnostic messages
// can speak the same vocabulary as the rest of the linter.
const (
	featureVariantAggregate                  = "aggregate"
	featureVariantDynamicConsistencyBoundary = "dynamic-consistency-boundary"
	featureVariantProcessManager             = "process-manager"
	featureVariantReadModel                  = "read-model"
)

// featureVariant returns the kind name of the consistency
// unit a feature targets, derived from which discriminator
// field is present in `scope`. Returns "" when the scope
// is missing or has no recognized discriminator - the
// schema validator catches that case separately, so rules
// can simply skip such features.
func featureVariant(feature model.FeatureView) string {
	scope := feature.Scope()
	switch {
	case scope.Field("aggregate").Exists():
		return featureVariantAggregate
	case scope.Field("dynamicConsistencyBoundary").Exists():
		return featureVariantDynamicConsistencyBoundary
	case scope.Field("processManager").Exists():
		return featureVariantProcessManager
	case scope.Field("readModel").Exists():
		return featureVariantReadModel
	}
	return ""
}

// scenariosOf returns the typed scenario views of a
// feature.
func scenariosOf(feature model.FeatureView) []model.ScenarioView {
	scenarios := feature.Scenarios().Seq()
	out := make([]model.ScenarioView, 0, len(scenarios))
	for _, s := range scenarios {
		out = append(out, model.ScenarioView{Node: s})
	}
	return out
}

// SchemaHasProperty reports whether the JSON-Schema-shaped
// node declares a top-level property called fieldName via
// its `properties` map. The check is intentionally local:
// nested $ref, oneOf/allOf branches, and patternProperties
// are not followed, because the cross-field rules that
// consume this helper want a clear "is this field declared
// here?" answer rather than a full JSON-Schema resolution.
//
// Authors that wrap their state/data behind a $ref or use
// composition should declare the referenced field
// explicitly at the top level (e.g. via allOf with an
// inline branch that lists `properties`); otherwise the
// rule will flag the reference as unresolved, which is a
// modeling signal in itself.
func SchemaHasProperty(schema ast.Node, fieldName string) bool {
	properties := schema.Field("properties")
	if !properties.Exists() {
		return false
	}
	return properties.HasField(fieldName)
}
