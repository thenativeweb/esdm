package lint_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/lint"
)

// parentsYAML is a fully consistent minimal model that
// does not make any esdm/modeling/* rule throw. Tests that
// want to exercise a specific rule add focused YAML on
// top of this baseline.
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

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
}

func writeValidModel(t *testing.T, dir string) {
	t.Helper()
	writeFile(t, dir, "parents.esdm.yaml", parentsYAML)
	writeFile(t, dir, "event.esdm.yaml", validEventYAML)
}

func runLintCommand(t *testing.T, args []string) (string, error) {
	t.Helper()

	var buf bytes.Buffer
	lint.Command.SetOut(&buf)
	lint.Command.SetErr(&buf)
	lint.Command.SetArgs(args)
	err := lint.Command.Execute()
	return buf.String(), err
}

func TestLintCommand(t *testing.T) {
	t.Run("produces empty output for a valid event model in human format", func(t *testing.T) {
		dir := t.TempDir()
		writeValidModel(t, dir)

		output, err := runLintCommand(t, []string{"--directory", dir})
		require.NoError(t, err)
		assert.Empty(t, output)
	})

	t.Run("emits JSON output when a modeling rule throws", func(t *testing.T) {
		dir := t.TempDir()
		writeValidModel(t, dir)
		// A second event whose name starts with its
		// aggregate's name makes the aggregate-prefix rule
		// throw without introducing unresolved references
		// in the baseline.
		writeFile(t, dir, "extra-event.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-shipped
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)

		output, err := runLintCommand(t, []string{"--directory", dir, "--format", "json"})
		require.NoError(t, err)

		var diagnostics []map[string]any
		require.NoError(t, json.Unmarshal([]byte(output), &diagnostics))
		assert.NotEmpty(t, diagnostics)

		didFindPrefix := false
		for _, d := range diagnostics {
			if d["ruleId"] == "esdm/modeling/event-name-with-aggregate-prefix" {
				didFindPrefix = true
			}
		}
		assert.True(t, didFindPrefix, "expected event-name-with-aggregate-prefix rule to throw, got %+v", diagnostics)
	})

	t.Run("rejects unknown --format values", func(t *testing.T) {
		dir := t.TempDir()
		writeValidModel(t, dir)

		_, err := runLintCommand(t, []string{"--directory", dir, "--format", "yaml"})
		assert.Error(t, err)
	})

	t.Run("returns an error when the directory does not exist", func(t *testing.T) {
		_, err := runLintCommand(t, []string{"--directory", filepath.Join(t.TempDir(), "does-not-exist")})
		assert.Error(t, err)
	})

	t.Run("reports exit code 0 after a successful run with no error diagnostics", func(t *testing.T) {
		dir := t.TempDir()
		writeValidModel(t, dir)

		_, err := runLintCommand(t, []string{"--directory", dir, "--format", "human"})
		require.NoError(t, err)
		assert.Equal(t, 0, lint.ExitCode())
	})

	t.Run("reports exit code 1 when an error-severity diagnostic is produced", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "parents.esdm.yaml", parentsYAML)
		// Reference an aggregate that does not exist - the
		// resolver flags this as an error-severity
		// `esdm/structure/unresolved-reference` diagnostic.
		writeFile(t, dir, "event.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: event
name: placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: nonexistent
data:
  type: object
`)

		_, err := runLintCommand(t, []string{"--directory", dir, "--format", "human"})
		require.NoError(t, err)
		assert.Equal(t, lint.ExitCodeHasErrors, lint.ExitCode())
	})

	t.Run("reports exit code 0 for a warning-only model when --warnings-as-errors is unset", func(t *testing.T) {
		dir := t.TempDir()
		writeValidModel(t, dir)
		// Adds a warning-severity diagnostic (event name
		// starts with the aggregate's name) without
		// introducing any error-severity diagnostic. The
		// run is therefore "warning-only".
		writeFile(t, dir, "extra-event.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-shipped
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)

		_, err := runLintCommand(t, []string{"--directory", dir, "--format", "human", "--warnings-as-errors=false"})
		require.NoError(t, err)
		assert.Equal(t, 0, lint.ExitCode())
	})

	t.Run("reports exit code 1 for a warning-only model when --warnings-as-errors is set", func(t *testing.T) {
		dir := t.TempDir()
		writeValidModel(t, dir)
		writeFile(t, dir, "extra-event.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-shipped
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)

		_, err := runLintCommand(t, []string{"--directory", dir, "--format", "human", "--warnings-as-errors"})
		require.NoError(t, err)
		assert.Equal(t, lint.ExitCodeHasErrors, lint.ExitCode())
	})

	t.Run("resets the exit code between runs so a clean run after a dirty one reports 0", func(t *testing.T) {
		dirtyDir := t.TempDir()
		writeFile(t, dirtyDir, "parents.esdm.yaml", parentsYAML)
		writeFile(t, dirtyDir, "event.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: event
name: placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: nonexistent
data:
  type: object
`)
		_, err := runLintCommand(t, []string{"--directory", dirtyDir, "--format", "human"})
		require.NoError(t, err)
		require.Equal(t, lint.ExitCodeHasErrors, lint.ExitCode())

		cleanDir := t.TempDir()
		writeValidModel(t, cleanDir)
		_, err = runLintCommand(t, []string{"--directory", cleanDir, "--format", "human"})
		require.NoError(t, err)
		assert.Equal(t, 0, lint.ExitCode())
	})
}
