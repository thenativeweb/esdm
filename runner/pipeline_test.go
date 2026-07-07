package runner_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/runner"
	"github.com/thenativeweb/esdm/schema"
)

// parentsYAML provides the chain of parent entities every
// scope-bearing test fixture needs (domain, bounded
// context, aggregate) plus an actor, a command and an
// event-handler. This makes the baseline a fully
// consistent model, so no esdm/modeling/* rule throws,
// and individual tests can focus on the specific signal
// they want to exercise.
const parentsYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: commerce
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: commerce
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: order
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: customer
scope:
  domain: commerce
  boundedContext: ordering
type: human
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: place-order
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
publishes:
  - placed
actors:
  - customer
---
apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: notify-customer
scope:
  domain: commerce
deliveryGuarantee: at-most-once
handles:
  - boundedContext: ordering
    aggregate: order
    event: placed
sideEffects:
  - type: other
    rule: send confirmation
`

const validEventYAML = `apiVersion: schema.esdm.io/core/v1
kind: event
name: placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`

func writeESDM(t *testing.T, dir, fileName, content string) {
	t.Helper()

	path := filepath.Join(dir, fileName)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func writeParents(t *testing.T, dir string) {
	t.Helper()
	writeESDM(t, dir, "parents.esdm.yaml", parentsYAML)
}

func TestRun(t *testing.T) {
	t.Run("returns a Walk error when the directory does not exist", func(t *testing.T) {
		_, err := runner.Run(context.Background(), "/nonexistent/path")
		assert.Error(t, err)
	})

	t.Run("produces no diagnostics for a single valid event model", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "event.esdm.yaml", validEventYAML)

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)
		assert.Empty(t, diagnostics)
	})

	t.Run("runs rules when parser and resolver produce no errors and the naming rule throws", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "event.esdm.yaml", validEventYAML)
		// A second event whose name starts with its
		// aggregate's name makes the aggregate-prefix rule
		// throw while leaving the baseline model intact.
		writeESDM(t, dir, "extra-event.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-shipped
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)

		hasFoundDiagnostic := false
		for _, d := range diagnostics {
			if d.RuleID == "esdm/modeling/event-name-with-aggregate-prefix" {
				hasFoundDiagnostic = true
				assert.Equal(t, diag.SeverityWarning, d.Severity)
			}
		}
		assert.True(t, hasFoundDiagnostic, "expected event-name-with-aggregate-prefix diagnostic, got %+v", diagnostics)
	})

	t.Run("skips the rule engine when parser reports schema errors", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "event.esdm.yaml", strings.Replace(validEventYAML, "name: placed\n", "", 1))

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)

		for _, d := range diagnostics {
			assert.NotEqual(t, "esdm/modeling/event-name-with-aggregate-prefix", d.RuleID,
				"rule engine should have been skipped when schema errors exist")
		}
	})

	t.Run("detects duplicate event names across files", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "a.esdm.yaml", validEventYAML)
		writeESDM(t, dir, "b.esdm.yaml", validEventYAML)

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)

		hasFoundDiagnostic := false
		for _, d := range diagnostics {
			if d.RuleID == "esdm/structure/duplicate-name" {
				hasFoundDiagnostic = true
			}
		}
		assert.True(t, hasFoundDiagnostic, "expected duplicate-name diagnostic, got %+v", diagnostics)
	})

	t.Run("flags an unresolved reference when an event names an unknown aggregate", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "event.esdm.yaml", strings.Replace(validEventYAML, "aggregate: order", "aggregate: ordr", 1))

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)

		var unresolvedReference *diag.Diagnostic
		for i := range diagnostics {
			if diagnostics[i].RuleID == "esdm/structure/unresolved-reference" {
				unresolvedReference = &diagnostics[i]
				break
			}
		}
		require.NotNil(t, unresolvedReference, "expected unresolved-reference diagnostic, got %+v", diagnostics)
		assert.Equal(t, diag.SeverityError, unresolvedReference.Severity)
		assert.Contains(t, unresolvedReference.Message, "ordr")

		require.Len(t, unresolvedReference.Related, 1)
		assert.Contains(t, unresolvedReference.Related[0].Message, "order")
	})

	t.Run("passes when the local schemas directory matches the embedded set", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "event.esdm.yaml", validEventYAML)
		require.NoError(t, schema.Write(filepath.Join(dir, "schemas")))

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)
		assert.Empty(t, diagnostics)
	})

	t.Run("aborts with a drift diagnostic when the local schemas directory has an extra file", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "event.esdm.yaml", validEventYAML)
		require.NoError(t, schema.Write(filepath.Join(dir, "schemas")))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "schemas", "stray.yaml"), []byte("x"), 0o644))

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)

		require.Len(t, diagnostics, 1)
		assert.Equal(t, "esdm/system/schemas-directory-drift", diagnostics[0].RuleID)
		assert.Equal(t, diag.SeverityError, diagnostics[0].Severity)
		assert.Contains(t, diagnostics[0].Message, "update-schema")
	})

	t.Run("aborts with a drift diagnostic when a schemas file has been hand-edited", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "event.esdm.yaml", validEventYAML)
		require.NoError(t, schema.Write(filepath.Join(dir, "schemas")))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "schemas", "core", "v1.yaml"), []byte("hand-edited"), 0o644))

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)

		require.Len(t, diagnostics, 1)
		assert.Equal(t, "esdm/system/schemas-directory-drift", diagnostics[0].RuleID)
	})

	t.Run("ignores a missing local schemas directory", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "event.esdm.yaml", validEventYAML)

		diagnostics, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)
		assert.Empty(t, diagnostics)
	})

	t.Run("returns diagnostics in deterministic order", func(t *testing.T) {
		dir := t.TempDir()
		writeParents(t, dir)
		writeESDM(t, dir, "a.esdm.yaml", strings.Replace(validEventYAML, "name: placed", "name: initiate", 1))
		writeESDM(t, dir, "b.esdm.yaml", strings.Replace(validEventYAML, "name: placed", "name: destroy", 1))

		first, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)

		second, err := runner.Run(context.Background(), dir)
		require.NoError(t, err)

		require.Equal(t, len(first), len(second))
		for i := range first {
			assert.Equal(t, first[i].RuleID, second[i].RuleID)
			assert.Equal(t, first[i].Location, second[i].Location)
		}
	})
}
