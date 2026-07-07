package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/parser"
)

// validEvent returns a minimal, schema-valid event
// document in YAML form. Tests mutate one aspect and
// assert that a specific diagnostic appears.
const validEvent = `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "doc.esdm.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func hasRuleID(diagnostics []diag.Diagnostic, ruleID string) bool {
	for _, d := range diagnostics {
		if d.RuleID == ruleID {
			return true
		}
	}
	return false
}

func TestParse(t *testing.T) {
	t.Run("returns an error when the file does not exist", func(t *testing.T) {
		_, _, err := parser.Parse("/nonexistent/path.esdm.yaml")
		assert.Error(t, err)
	})

	t.Run("produces no diagnostics for a valid event document", func(t *testing.T) {
		path := writeTempFile(t, validEvent)

		parsed, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotNil(t, parsed)

		assert.Empty(t, diagnostics)
		assert.Len(t, parsed.Documents, 1)
	})

	t.Run("splits multi-document YAML into one AST node per document", func(t *testing.T) {
		content := validEvent + "---\n" + strings.Replace(validEvent, "order-placed", "order-shipped", 1)
		path := writeTempFile(t, content)

		parsed, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotNil(t, parsed)

		assert.Empty(t, diagnostics)
		assert.Len(t, parsed.Documents, 2)
	})

	t.Run("reports a YAML syntax error as a structure/yaml-syntax-error diagnostic", func(t *testing.T) {
		path := writeTempFile(t, "name: [unterminated\n")

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.Len(t, diagnostics, 1)

		assert.Equal(t, "esdm/structure/yaml-syntax-error", diagnostics[0].RuleID)
		assert.Equal(t, path, diagnostics[0].Location.File)
	})

	t.Run("reports a missing required field as missing-required-field", func(t *testing.T) {
		document := strings.Replace(validEvent, "name: order-placed\n", "", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/missing-required-field"),
			"expected missing-required-field in %+v", diagnostics)
	})

	t.Run("reports a type mismatch as type-mismatch", func(t *testing.T) {
		document := strings.Replace(validEvent, "kind: event\n", "kind: 42\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/type-mismatch"),
			"expected type-mismatch in %+v", diagnostics)
	})

	t.Run("reports an unknown field as unknown-field", func(t *testing.T) {
		document := validEvent + "bogusField: true\n"
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/unknown-field"),
			"expected unknown-field in %+v", diagnostics)
	})

	t.Run("reports a pattern violation as constraint-violation", func(t *testing.T) {
		document := strings.Replace(validEvent, "name: order-placed\n", "name: Order Placed\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/constraint-violation"),
			"expected constraint-violation in %+v", diagnostics)
	})

	t.Run("attaches a did-you-mean hint when an enum value is close to a valid one", func(t *testing.T) {
		document := strings.Replace(validEvent, "kind: event\n", "kind: evant\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)

		var enum *diag.Diagnostic
		for i := range diagnostics {
			if diagnostics[i].RuleID == "esdm/structure/constraint-violation" {
				enum = &diagnostics[i]
				break
			}
		}
		require.NotNil(t, enum, "expected a constraint-violation diagnostic, got %+v", diagnostics)
		require.NotEmpty(t, enum.Related, "expected a did-you-mean hint on the enum violation, got %+v", enum)
		assert.Contains(t, enum.Related[0].Message, "event")
	})

	t.Run("accepts a valid domain-storytelling extension document without diagnostics", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: place-order
scope:
  domain: commerce
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`)

		parsed, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotNil(t, parsed)
		assert.Empty(t, diagnostics)
	})

	t.Run("accepts a valid given-when-then aggregate feature without diagnostics", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: order-cancellation
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
scenarios:
  - name: cancels-an-open-order
    given:
      - event: order-placed
        data: {}
    when:
      command: cancel-order
      data: {}
    then:
      events:
        - event: order-canceled
          data: {}
  - name: rejects-cancellation-of-shipped-order
    given:
      - event: order-placed
        data: {}
      - event: order-shipped
        data: {}
    when:
      command: cancel-order
      data: {}
    then:
      rejection:
        invariant: shipping-blocks-cancellation
  - name: idempotent-re-cancel-emits-no-events
    given:
      - event: order-placed
        data: {}
      - event: order-canceled
        data: {}
    when:
      command: cancel-order
      data: {}
    then:
      events: []
`)

		parsed, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotNil(t, parsed)
		assert.Empty(t, diagnostics, "expected no diagnostics, got %+v", diagnostics)
	})

	t.Run("rejects a feature whose when shape does not match its scope variant", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: order-cancellation
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
scenarios:
  - name: cancels-an-open-order
    given: []
    when:
      query: list-orders
      parameters: {}
    then:
      events: []
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		assert.NotEmpty(t, diagnostics, "expected schema violations for query-shaped when in aggregate feature")
	})

	t.Run("accepts a valid given-when-then process-manager feature with timer when and emits then", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: order-fulfillment
scope:
  domain: commerce
  processManager: order-fulfillment
scenarios:
  - name: timer-fires-and-cancels-pending-order
    given:
      - boundedContext: ordering
        aggregate: order
        event: order-placed
        data: {}
    when:
      timer: order-acceptance-deadline
    then:
      emits:
        - boundedContext: ordering
          aggregate: order
          command: cancel-order
          data: {}
      ended: true
`)

		parsed, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotNil(t, parsed)
		assert.Empty(t, diagnostics, "expected no diagnostics, got %+v", diagnostics)
	})

	t.Run("accepts a valid given-when-then read-model feature with query when and result then", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: order-overview
scope:
  domain: commerce
  boundedContext: ordering
  readModel: order-overview
scenarios:
  - name: lists-only-open-orders
    given:
      - boundedContext: ordering
        aggregate: order
        event: order-placed
        data: {}
    when:
      query: list-open-orders
      parameters: {}
    then:
      result: []
`)

		parsed, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotNil(t, parsed)
		assert.Empty(t, diagnostics, "expected no diagnostics, got %+v", diagnostics)
	})

	t.Run("emits unknown-api-version when a document's apiVersion is not compiled in", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: example.com/schema/imaginary/v1
kind: event
name: whatever
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/unknown-api-version"),
			"expected unknown-api-version, got %+v", diagnostics)
	})

	t.Run("emits unknown-api-version when a document has no apiVersion field", func(t *testing.T) {
		path := writeTempFile(t, "kind: event\nname: whatever\n")

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/unknown-api-version"),
			"expected unknown-api-version, got %+v", diagnostics)
	})

	t.Run("points diagnostics at a precise source location", func(t *testing.T) {
		document := strings.Replace(validEvent, "kind: event\n", "kind: 42\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		var typeMismatch *diag.Diagnostic
		for i := range diagnostics {
			if diagnostics[i].RuleID == "esdm/structure/type-mismatch" {
				typeMismatch = &diagnostics[i]
				break
			}
		}
		require.NotNil(t, typeMismatch, "did not find a type-mismatch diagnostic")

		assert.Equal(t, path, typeMismatch.Location.File)
		assert.Positive(t, typeMismatch.Location.Line)
		assert.Positive(t, typeMismatch.Location.Column)
	})
}
