package parser_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/parser"
	"github.com/thenativeweb/esdm/schema"
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

// boundedContextWithTerminology is a bounded context whose
// ubiquitous language is written in English and carries one
// German translation, including a rejected alternative for
// that translation.
const boundedContextWithTerminology = `apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: billing
scope:
  domain: commerce
language: en
ubiquitousLanguage:
  - term: Invoice
    definition: A request for payment issued to a customer.
    avoid:
      - term: Bill
        reason: Colloquial; conflicts with the accounting sense.
    translations:
      - language: de
        term: Rechnung
        definition: Eine an den Kunden gerichtete Zahlungsaufforderung.
        avoid:
          - term: Faktura
            reason: Austrian usage, not spoken in this team.
`

func TestBoundedContextLanguage(t *testing.T) {
	t.Run("accepts a ubiquitous language with a declared language and translations", func(t *testing.T) {
		path := writeTempFile(t, boundedContextWithTerminology)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		assert.Empty(t, diagnostics)
	})

	t.Run("requires language once a bounded context declares a ubiquitous language", func(t *testing.T) {
		document := strings.Replace(boundedContextWithTerminology, "language: en\n", "", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/missing-required-field"),
			"expected missing-required-field in %+v", diagnostics)
	})

	t.Run("does not require language on a bounded context without ubiquitous language", func(t *testing.T) {
		document := `apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: billing
scope:
  domain: commerce
`
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		assert.Empty(t, diagnostics)
	})

	t.Run("rejects a language that is not a BCP 47 tag", func(t *testing.T) {
		document := strings.Replace(boundedContextWithTerminology, "language: en\n", "language: English\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/constraint-violation"),
			"expected constraint-violation in %+v", diagnostics)
	})

	t.Run("requires a definition on every translation", func(t *testing.T) {
		document := strings.Replace(boundedContextWithTerminology, "        definition: Eine an den Kunden gerichtete Zahlungsaufforderung.\n", "", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/missing-required-field"),
			"expected missing-required-field in %+v", diagnostics)
	})
}

const customerSupplierWithTerms = `apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: sales-billing
type: customer-supplier
customer:
  domain: commerce
  boundedContext: sales
supplier:
  domain: commerce
  boundedContext: billing
terms:
  - customer: Customer
    supplier: Account
`

func TestContextMappingTerms(t *testing.T) {
	t.Run("accepts term pairs named after the roles of an asymmetric mapping", func(t *testing.T) {
		path := writeTempFile(t, customerSupplierWithTerms)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		assert.Empty(t, diagnostics)
	})

	t.Run("requires both roles on every term pair", func(t *testing.T) {
		document := strings.Replace(customerSupplierWithTerms, "    supplier: Account\n", "", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.True(t, hasRuleID(diagnostics, "esdm/structure/missing-required-field"),
			"expected missing-required-field in %+v", diagnostics)
	})

	t.Run("rejects term pairs on a symmetric mapping", func(t *testing.T) {
		document := `apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: sales-billing
type: shared-kernel
participants:
  - domain: commerce
    boundedContext: sales
  - domain: commerce
    boundedContext: billing
terms:
  - customer: Customer
    supplier: Account
`
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)
	})
}

// camelCaseIdentifier references a state property that is
// not kebab-case. JSON Schema places no constraint on
// property names, so a reference to one must not either.
const camelCaseIdentifier = `apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: account
scope:
  domain: commerce
  boundedContext: billing
identifiedBy:
  source: state
  field: accountId
state:
  type: object
  properties:
    accountId:
      type: string
`

func TestFieldReferences(t *testing.T) {
	t.Run("accepts a field reference in any form the referenced schema allows", func(t *testing.T) {
		path := writeTempFile(t, camelCaseIdentifier)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		assert.Empty(t, diagnostics)
	})

	t.Run("accepts a timer that names a camelCase state field", func(t *testing.T) {
		document := `apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: reminder
scope:
  domain: commerce
deliveryGuarantee: at-most-once
correlatedBy:
  source: event-field
  field: orderId
state:
  type: object
  properties:
    dueAt:
      type: string
timers:
  - name: due
    at: dueAt
startsWhen:
  - boundedContext: ordering
    aggregate: order
    event: placed
endsWhen:
  - name: done
    condition: state.completed is true
reactions:
  - when:
      timer: due
    rule: remind
`
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		assert.Empty(t, diagnostics)
	})

	t.Run("still rejects an empty field reference", func(t *testing.T) {
		document := strings.Replace(camelCaseIdentifier, "  field: accountId\n", "  field: \"\"\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.NotEmpty(t, diagnostics)
	})
}

// messagesOf renders diagnostics as "<ruleID>: <message>"
// so a test can assert the exact set of findings one
// defect produces, not just that a certain rule ID is
// among them.
func messagesOf(diagnostics []diag.Diagnostic) []string {
	out := make([]string, 0, len(diagnostics))
	for _, d := range diagnostics {
		out = append(out, d.RuleID+": "+d.Message)
	}
	return out
}

// TestFollowOnErrors pins down that one defect yields one
// finding. Schema validation drops the annotations of a
// failed branch, which makes every field that branch
// would have evaluated look unknown, and it reports the
// causes of all alternatives of a oneOf. Both turn a
// single defect into a list of findings that sends the
// reader after the wrong one.
func TestFollowOnErrors(t *testing.T) {
	t.Run("reports only the missing language when a bounded context declares a ubiquitous language without one", func(t *testing.T) {
		document := strings.Replace(boundedContextWithTerminology, "language: en\n", "", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.Equal(t, []string{
			`esdm/structure/missing-required-field: missing required field "language"`,
		}, messagesOf(diagnostics))
	})

	t.Run("reports only the rejection branch when a scenario's rejection carries an unknown field", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: order-cancellation
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
scenarios:
  - name: rejects-cancellation-of-shipped-order
    given:
      - event: order-placed
        data: {}
    when:
      command: cancel-order
      data: {}
    then:
      rejection:
        invariant: shipping-blocks-cancellation
        reason: A shipped order cannot be canceled.
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.Equal(t, []string{
			`esdm/structure/unknown-field: unknown field "reason"`,
		}, messagesOf(diagnostics))
	})

	t.Run("keeps reporting a field that no branch of the document's kind declares", func(t *testing.T) {
		document := strings.Replace(boundedContextWithTerminology, "language: en\n", "bogusField: true\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{
			`esdm/structure/missing-required-field: missing required field "language"`,
			`esdm/structure/unknown-field: unknown field "bogusField"`,
		}, messagesOf(diagnostics))
	})

	t.Run("keeps reporting a misspelled field next to the required field it was meant to be", func(t *testing.T) {
		document := strings.Replace(validEvent, "data:\n", "dat:\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{
			`esdm/structure/missing-required-field: missing required field "data"`,
			`esdm/structure/unknown-field: unknown field "dat"`,
		}, messagesOf(diagnostics))
	})

	t.Run("reports every branch of a oneOf when the document resembles none of them", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: order-cancellation
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
scenarios:
  - name: rejects-cancellation-of-shipped-order
    given: []
    when:
      command: cancel-order
      data: {}
    then: {}
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{
			`esdm/structure/missing-required-field: missing required field "events"`,
			`esdm/structure/missing-required-field: missing required field "rejection"`,
		}, messagesOf(diagnostics))
	})

	t.Run("reports the branch a discriminating const identifies", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: notify-customer
scope:
  domain: commerce
deliveryGuarantee: at-most-once
handles:
  - boundedContext: ordering
    aggregate: order
    event: order-placed
sideEffects:
  - type: external-call
    rule: send the order-placed notification email
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.Equal(t, []string{
			`esdm/structure/missing-required-field: missing required field "externalSystem"`,
		}, messagesOf(diagnostics))
	})
}

// sweepDocuments are schema-valid documents covering the
// kinds whose schemas combine branches: `allOf` with
// `if`/`then` for the kind, and `oneOf` for the variants
// of a field. Those are the shapes that produce follow-on
// errors, so they are the ones the sweep below mutates.
var sweepDocuments = map[string]string{
	"event":                        validEvent,
	"bounded-context":              boundedContextWithTerminology,
	"aggregate":                    camelCaseIdentifier,
	"context-mapping":              customerSupplierWithTerms,
	"dynamic-consistency-boundary": validDynamicConsistencyBoundary,
	"feature":                      validAggregateFeature,
}

const validDynamicConsistencyBoundary = `apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: course-enrollment
scope:
  domain: education
  boundedContext: enrollment
identifiedBy:
  - name: course
    source: command-payload
    field: courseId
  - name: tenant
    source: static
    value: acme
consults:
  - boundedContext: enrollment
    aggregate: course
    event: enrollment-recorded
    criteria: all enrollments with the same courseId
invariants:
  - name: at-most-thirty-enrollments
    rule: a course holds at most 30 enrollments
`

const validAggregateFeature = `apiVersion: schema.esdm.io/given-when-then/v1
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
  - name: rejects-cancellation-of-a-shipped-order
    given:
      - event: order-shipped
        data: {}
    when:
      command: cancel-order
      data: {}
    then:
      rejection:
        invariant: shipping-blocks-cancellation
`

// referenceValidators compiles the embedded schemas
// independently of the parser. The sweep needs a second
// opinion on whether a mutated document is valid at all;
// asking the parser itself would make the assertion
// circular, since the parser is what is under test.
func referenceValidators(t *testing.T) map[string]*jsonschema.Schema {
	t.Helper()

	validators := make(map[string]*jsonschema.Schema)

	compile := func(source []byte) {
		var decoded any
		require.NoError(t, yaml.Unmarshal(source, &decoded))

		id, isString := decoded.(map[string]any)["$id"].(string)
		require.True(t, isString, "schema without a string $id")

		compiler := jsonschema.NewCompiler()
		require.NoError(t, compiler.AddResource(id, decoded))

		compiled, err := compiler.Compile(id)
		require.NoError(t, err)

		validators[strings.TrimPrefix(id, "https://")] = compiled
	}

	compile(schema.Core())

	extensions, err := schema.Extensions()
	require.NoError(t, err)
	for _, extension := range extensions {
		compile(extension.Bytes)
	}

	return validators
}

// fieldPaths lists every mapping key in a decoded
// document, each as the sequence of segments leading to
// it.
func fieldPaths(value any, path []string) [][]string {
	var paths [][]string

	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			here := append(append([]string{}, path...), key)
			paths = append(paths, here)
			paths = append(paths, fieldPaths(child, here)...)
		}
	case []any:
		for index, item := range typed {
			here := append(append([]string{}, path...), strconv.Itoa(index))
			paths = append(paths, fieldPaths(item, here)...)
		}
	}

	return paths
}

// renameField returns a copy of the document with the
// field at path renamed. Renaming is the mutation the
// sweep uses because it produces both halves of the
// problem at once: the original field goes missing and an
// unknown one takes its place.
func renameField(t *testing.T, document string, path []string, renamed string) (any, string) {
	t.Helper()

	var mutated any
	require.NoError(t, yaml.Unmarshal([]byte(document), &mutated))

	parent := mutated
	for _, segment := range path[:len(path)-1] {
		switch typed := parent.(type) {
		case map[string]any:
			parent = typed[segment]
		case []any:
			index, err := strconv.Atoi(segment)
			require.NoError(t, err)
			parent = typed[index]
		}
	}

	mapping, isMapping := parent.(map[string]any)
	require.True(t, isMapping)

	key := path[len(path)-1]
	mapping[renamed] = mapping[key]
	delete(mapping, key)

	encoded, err := yaml.Marshal(mutated)
	require.NoError(t, err)

	return mutated, string(encoded)
}

// TestEveryDefectKeepsAFinding sweeps one renamed field
// at a time through every document above and insists that
// the parser still describes the defect. Dropping
// follow-on errors is only correct as long as the finding
// that names the cause survives; a filter that silenced a
// defect outright would be worse than the noise it
// removes.
func TestEveryDefectKeepsAFinding(t *testing.T) {
	validators := referenceValidators(t)

	for kind, document := range sweepDocuments {
		t.Run(kind, func(t *testing.T) {
			var decoded any
			require.NoError(t, yaml.Unmarshal([]byte(document), &decoded))

			for _, path := range fieldPaths(decoded, nil) {
				field := path[len(path)-1]
				renamed := field + "Renamed"

				mutated, encoded := renameField(t, document, path, renamed)

				apiVersion, _ := mutated.(map[string]any)["apiVersion"].(string)
				validator, isKnown := validators[apiVersion]
				if !isKnown {
					continue
				}

				// A rename inside a free-form object - a
				// payload schema, a metadata map - is still
				// valid, and then there is nothing to report.
				if validator.Validate(mutated) == nil {
					continue
				}

				_, diagnostics, err := parser.Parse(writeTempFile(t, encoded))
				require.NoError(t, err)

				found := messagesOf(diagnostics)
				namesTheDefect := false
				for _, message := range found {
					if strings.Contains(message, `"`+field+`"`) || strings.Contains(message, `"`+renamed+`"`) {
						namesTheDefect = true
					}
				}

				assert.True(t, namesTheDefect,
					"renaming %s to %q left no finding naming it, got %v",
					strings.Join(path, "/"), renamed, found)
			}
		})
	}
}

// TestSilentSchemaViolations pins down that a document
// the schema rejects always produces a finding. Some
// schema keywords fail without nested causes, and a
// translation that only ever walks causes turns such a
// failure into silence - the worst answer a linter can
// give.
func TestSilentSchemaViolations(t *testing.T) {
	t.Run("reports a field that is forbidden by the value of another one", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/core/v1
kind: actor
name: customer
scope:
  domain: commerce
  boundedContext: ordering
type: human
backedBy:
  - crm-gateway
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.Equal(t, []string{
			`esdm/structure/constraint-violation: field "backedBy" is not allowed when "type" is "human"`,
		}, messagesOf(diagnostics))
	})

	t.Run("points the finding at the forbidden field", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/core/v1
kind: actor
name: customer
scope:
  domain: commerce
  boundedContext: ordering
type: human
backedBy:
  - crm-gateway
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)
		require.Len(t, diagnostics, 1)

		assert.Equal(t, 8, diagnostics[0].Location.Line,
			"expected the finding on the backedBy line, got %+v", diagnostics[0].Location)
	})

	t.Run("accepts the same field on an actor whose type allows it", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/core/v1
kind: actor
name: billing-robot
scope:
  domain: commerce
  boundedContext: ordering
type: system
backedBy:
  - crm-gateway
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.Empty(t, diagnostics)
	})
}

// TestRepeatedAndNestedVariants covers the two ways a
// single defect could still draw more than one finding
// after the follow-on errors are gone: alternatives that
// fail on the same field report it once each, and an
// alternative that is itself a set of alternatives states
// its requirements a level down, where the comparison
// against its siblings could not see them.
func TestRepeatedAndNestedVariants(t *testing.T) {
	t.Run("reports a field that every alternative rejects only once", func(t *testing.T) {
		document := strings.Replace(validEvent, "  aggregate: order\n", "  agregate: order\n", 1)
		path := writeTempFile(t, document)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{
			`esdm/structure/missing-required-field: missing required field "aggregate"`,
			`esdm/structure/unknown-field: unknown field "agregate"`,
		}, messagesOf(diagnostics))
	})

	t.Run("identifies an alternative whose requirements sit one level down", func(t *testing.T) {
		path := writeTempFile(t, `apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: order-fulfillment
scope:
  domain: commerce
deliveryGuarantee: at-most-once
correlatedBy:
  source: event-field
  field: orderId
state:
  type: object
startsWhen:
  - boundedContext: ordering
    aggregate: order
    event: order-placed
endsWhen:
  - name: fulfilled
    condition: state.completed is true
reactions:
  - when:
      boundedContext: ordering
      aggregate: order
      evnt: order-placed
    rule: start the fulfillment
`)

		_, diagnostics, err := parser.Parse(path)
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{
			`esdm/structure/missing-required-field: missing required field "event"`,
			`esdm/structure/unknown-field: unknown field "evnt"`,
		}, messagesOf(diagnostics))
	})
}
