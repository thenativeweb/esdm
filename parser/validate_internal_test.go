package parser

import (
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/ast"
)

// overlappingSchema declares two variants that a single
// document can satisfy at once. No embedded schema is
// shaped this way - every `oneOf` in them separates its
// branches by a contradicting constant or by a closed
// object - so this case is only reachable through a
// schema written for the test. It is covered all the same
// because the failure mode is silence, and a schema
// change could introduce it without anyone noticing.
const overlappingSchema = `$id: https://schema.esdm.io/overlapping/v1
$schema: https://json-schema.org/draft/2020-12/schema
type: object
properties:
  outcome:
    oneOf:
      - type: object
        properties:
          accepted:
            type: string
        required: [accepted]
      - type: object
        properties:
          rejected:
            type: string
        required: [rejected]
`

func TestAmbiguousVariant(t *testing.T) {
	t.Run("reports a document that satisfies more than one variant of a oneOf", func(t *testing.T) {
		schemas := make(map[string]*compiledSchema)
		require.NoError(t, compileSchemaInto(schemas, []byte(overlappingSchema)))

		schema := schemas["schema.esdm.io/overlapping/v1"]
		require.NotNil(t, schema)

		source := "outcome:\n  accepted: yes please\n  rejected: no thanks\n"

		var raw yaml.Node
		require.NoError(t, yaml.Unmarshal([]byte(source), &raw))

		var generic any
		require.NoError(t, yaml.Unmarshal([]byte(source), &generic))

		err := schema.Validator.Validate(generic)
		require.Error(t, err, "the test schema is meant to reject this document")

		validationError, isValidationError := err.(*jsonschema.ValidationError)
		require.True(t, isValidationError)

		diagnostics := validationDiagnostics(validationError, ast.NewNode("doc.esdm.yaml", &raw), schema)

		require.Len(t, diagnostics, 1)
		assert.Equal(t, "esdm/structure/constraint-violation", diagnostics[0].RuleID)
		assert.Equal(t, "matches more than one of the forms allowed here; exactly one must apply", diagnostics[0].Message)
	})
}
