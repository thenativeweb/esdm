package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEntityIdentifiedByField = "esdm/structure/entity-identified-by-field"

type entityIdentifiedByFieldRule struct{}

func newEntityIdentifiedByFieldRule() *entityIdentifiedByFieldRule {
	return &entityIdentifiedByFieldRule{}
}

func (*entityIdentifiedByFieldRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEntityIdentifiedByField,
		Severity:    diag.SeverityError,
		Description: "When entity.identifiedBy uses source: schema, the named `field` must be declared in the entity's schema.properties.",
	}
}

func (*entityIdentifiedByFieldRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, e := range sortedByName(m.Entities) {
		name, _ := e.Name().Text()
		identifiedBy := e.IdentifiedBy()

		source, _ := identifiedBy.Field("source").Text()
		if source != "schema" {
			continue
		}
		fieldNode := identifiedBy.Field("field")
		fieldName, ok := fieldNode.Text()
		if !ok {
			continue
		}

		if SchemaHasProperty(e.Schema(), fieldName) {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("entity %q identifiedBy.field %q is not declared in schema.properties", name, fieldName),
			Location: fieldNode.Location(),
		})
	}
}
