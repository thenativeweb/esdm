package model_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/schema"
)

// TestFacadesMatchSchema verifies that every facade method
// maps to a real schema field (no stale accessors) and
// every schema field that applies to the facade's kind has
// a matching accessor (no missed fields).
//
// The check replaces the safety effect of code generation
// without introducing a generator. Drift between schema
// and facades is caught by this test alone - nowhere else
// in the pipeline depends on the two staying in sync.
func TestFacadesMatchSchema(t *testing.T) {
	coreSchema := parseSchemaBytes(t, schema.Core())
	extensionSchema := parseSchemaBytes(t, extensionBytes(t, "domain-storytelling"))
	givenWhenThenSchema := parseSchemaBytes(t, extensionBytes(t, "given-when-then"))

	cases := []struct {
		name   string
		facade any
		kind   string
		schema map[string]any
	}{
		{name: "DomainView", facade: model.DomainView{}, kind: "domain", schema: coreSchema},
		{name: "SubdomainView", facade: model.SubdomainView{}, kind: "subdomain", schema: coreSchema},
		{name: "BoundedContextView", facade: model.BoundedContextView{}, kind: "bounded-context", schema: coreSchema},
		{name: "ContextMappingView", facade: model.ContextMappingView{}, kind: "context-mapping", schema: coreSchema},
		{name: "AggregateView", facade: model.AggregateView{}, kind: "aggregate", schema: coreSchema},
		{name: "DynamicConsistencyBoundaryView", facade: model.DynamicConsistencyBoundaryView{}, kind: "dynamic-consistency-boundary", schema: coreSchema},
		{name: "CommandView", facade: model.CommandView{}, kind: "command", schema: coreSchema},
		{name: "EventView", facade: model.EventView{}, kind: "event", schema: coreSchema},
		{name: "EventHandlerView", facade: model.EventHandlerView{}, kind: "event-handler", schema: coreSchema},
		{name: "PolicyView", facade: model.PolicyView{}, kind: "policy", schema: coreSchema},
		{name: "ProcessManagerView", facade: model.ProcessManagerView{}, kind: "process-manager", schema: coreSchema},
		{name: "ReadModelView", facade: model.ReadModelView{}, kind: "read-model", schema: coreSchema},
		{name: "QueryView", facade: model.QueryView{}, kind: "query", schema: coreSchema},
		{name: "EntityView", facade: model.EntityView{}, kind: "entity", schema: coreSchema},
		{name: "ValueObjectView", facade: model.ValueObjectView{}, kind: "value-object", schema: coreSchema},
		{name: "DomainServiceView", facade: model.DomainServiceView{}, kind: "domain-service", schema: coreSchema},
		{name: "ActorView", facade: model.ActorView{}, kind: "actor", schema: coreSchema},
		{name: "ExternalSystemView", facade: model.ExternalSystemView{}, kind: "external-system", schema: coreSchema},
		{name: "DomainStoryView", facade: model.DomainStoryView{}, kind: "domain-story", schema: extensionSchema},
		{name: "FeatureView", facade: model.FeatureView{}, kind: "feature", schema: givenWhenThenSchema},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			topLevel := topLevelFields(t, c.schema)
			branch := branchFields(t, c.schema, c.kind)

			schemaFields := make(map[string]bool)
			for f := range topLevel {
				schemaFields[f] = true
			}
			for f := range branch {
				schemaFields[f] = true
			}

			facadeFields, unused := reflectFacadeFields(t, c.facade)

			for _, u := range unused {
				delete(schemaFields, u)
			}

			for f := range schemaFields {
				assert.True(t, facadeFields[f], "%s: schema field %q has no matching facade accessor", c.name, f)
			}

			for f := range facadeFields {
				_, isInSchema := schemaFields[f]
				isInUnused := slices.Contains(unused, f)
				assert.True(t, isInSchema || isInUnused, "%s: facade accessor for %q but schema has no such field", c.name, f)
			}
		})
	}
}

func parseSchemaBytes(t *testing.T, data []byte) map[string]any {
	t.Helper()

	var m map[string]any
	require.NoError(t, yaml.Unmarshal(data, &m))
	return m
}

func extensionBytes(t *testing.T, name string) []byte {
	t.Helper()

	extensions, err := schema.Extensions()
	require.NoError(t, err)
	for _, e := range extensions {
		if e.Name == name {
			return e.Bytes
		}
	}
	t.Fatalf("extension %q not found in embedded schemas", name)
	return nil
}

func topLevelFields(t *testing.T, m map[string]any) map[string]bool {
	t.Helper()

	properties, ok := m["properties"].(map[string]any)
	require.True(t, ok, "schema is missing top-level properties")

	out := make(map[string]bool, len(properties))
	for name := range properties {
		out[name] = true
	}
	return out
}

// branchFields finds the allOf branch whose
// `if.properties.kind.const` matches the given kind and
// returns the union of every property name reachable at
// the same nesting level as `then` itself: directly via
// `then.properties`, and indirectly via `then.oneOf`,
// `then.anyOf`, and `then.allOf`. Nested schemas inside
// individual property values are deliberately NOT
// followed - only the same-level combinators contribute
// top-level field names.
func branchFields(t *testing.T, m map[string]any, kind string) map[string]bool {
	t.Helper()

	out := make(map[string]bool)

	branches, ok := m["allOf"].([]any)
	if !ok {
		return out
	}

	for _, b := range branches {
		branch, ok := b.(map[string]any)
		if !ok {
			continue
		}

		ifClause, _ := branch["if"].(map[string]any)
		ifProperties, _ := ifClause["properties"].(map[string]any)
		kindSpec, _ := ifProperties["kind"].(map[string]any)
		constVal, _ := kindSpec["const"].(string)

		if constVal != kind {
			continue
		}

		thenClause, _ := branch["then"].(map[string]any)
		collectSiblingProperties(thenClause, out)
	}

	return out
}

// collectSiblingProperties accumulates property names
// from `schema.properties` and from each branch of
// `schema.oneOf`, `schema.anyOf` and `schema.allOf`. It
// does not recurse into individual property values, so a
// nested `oneOf` inside `properties.identifiedBy` does
// not contribute its inner field names to the
// outer-level set.
func collectSiblingProperties(schema map[string]any, into map[string]bool) {
	if schema == nil {
		return
	}

	if properties, ok := schema["properties"].(map[string]any); ok {
		for name := range properties {
			into[name] = true
		}
	}

	for _, key := range []string{"oneOf", "anyOf", "allOf"} {
		branches, ok := schema[key].([]any)
		if !ok {
			continue
		}
		for _, b := range branches {
			branch, ok := b.(map[string]any)
			if !ok {
				continue
			}
			collectSiblingProperties(branch, into)
		}
	}
}

// reflectFacadeFields enumerates the zero-argument methods
// on the facade that return ast.Node, applies any
// FacadeOverrides renames, and returns the resolved schema
// field names together with the list of fields that are
// explicitly declared unused.
func reflectFacadeFields(t *testing.T, facade any) (map[string]bool, []string) {
	t.Helper()

	typ := reflect.TypeOf(facade)
	nodeType := reflect.TypeOf(ast.Node{})

	overrides := loadOverrides(facade)
	rename := overrides.Rename

	fields := make(map[string]bool)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)

		// Exclude meta methods and anything that does not
		// match the facade method signature: one output,
		// no arguments beyond the receiver, return type
		// ast.Node.
		if m.Name == "FacadeOverrides" {
			continue
		}
		if m.Type.NumIn() != 1 {
			continue
		}
		if m.Type.NumOut() != 1 {
			continue
		}
		if m.Type.Out(0) != nodeType {
			continue
		}

		fieldName := methodToField(m.Name)
		if renamed, ok := rename[m.Name]; ok {
			fieldName = renamed
		}
		fields[fieldName] = true
	}

	return fields, overrides.UnusedSchemaFields
}

func loadOverrides(facade any) model.FacadeOverrides {
	method := reflect.ValueOf(facade).MethodByName("FacadeOverrides")
	if !method.IsValid() {
		return model.FacadeOverrides{}
	}

	result := method.Call(nil)
	if len(result) != 1 {
		return model.FacadeOverrides{}
	}

	value, ok := result[0].Interface().(model.FacadeOverrides)
	if !ok {
		return model.FacadeOverrides{}
	}
	return value
}

func methodToField(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToLower(name[:1]) + name[1:]
}
