package rules_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/parser"
	"github.com/thenativeweb/esdm/reporter"
	"github.com/thenativeweb/esdm/resolver"
	"github.com/thenativeweb/esdm/rules"
)

// buildModel parses and resolves the given YAML document
// set and returns the fully populated Model. Tests use
// this for realistic fixtures when a rule depends on
// cross-entity relationships in the Model.
func buildModel(t *testing.T, yamlDocs string) *model.Model {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "model.esdm.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yamlDocs), 0o644))

	parsed, _, err := parser.Parse(path)
	require.NoError(t, err)

	m, _ := resolver.Resolve([]*parser.ParsedFile{parsed})
	return m
}

// buildModelTolerantOfSchemaErrors builds the Model from a
// YAML document set whose content the schema would reject.
// The schema-mirror rules use it to guard against an
// accidental relaxation of a schema constraint. It behaves
// exactly like buildModel - both ignore parse and resolve
// diagnostics - and exists only to signal at the call site
// that the input intentionally violates the schema.
func buildModelTolerantOfSchemaErrors(t *testing.T, yamlDocs string) *model.Model {
	t.Helper()

	return buildModel(t, yamlDocs)
}

// runRule invokes the rule's Check against the model and
// returns the diagnostics it reported.
func runRule(t *testing.T, rule rules.Rule, m *model.Model) []diag.Diagnostic {
	t.Helper()

	c := reporter.NewCollector()
	rule.Check(context.Background(), m, c)
	return c.All()
}

// minimalParents is a small but complete prefix that
// many modeling-rule tests extend with the entity under
// test: a domain, a bounded-context, an aggregate, a
// command in that aggregate, plus an event the command
// publishes. It exists so rules that rely on cross-
// references (publisher, consumer, aggregate presence)
// see a coherent model around the entity being checked.
const minimalParents = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: bc
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: agg
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: user
scope:
  domain: d
  boundedContext: bc
type: human
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: do-it
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
publishes:
  - agg-done
actors:
  - user
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: agg-done
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
`

// findCatalogRule returns the rule matching the given ID
// from Catalog(); it fails the test when no rule matches.
func findCatalogRule(t *testing.T, id string) rules.Rule {
	t.Helper()

	for _, r := range rules.Catalog() {
		if r.Meta().ID == id {
			return r
		}
	}
	t.Fatalf("rule %q not found in catalog", id)
	return nil
}
