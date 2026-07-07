package glossary_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/glossary"
)

const boundedContextWithoutLanguageYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shop
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: shop
`

func runGlossaryCommand(t *testing.T, dir, content string, args []string) (string, error) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "model.esdm.yaml"), []byte(content), 0o644))

	var buf bytes.Buffer
	glossary.Command.SetOut(&buf)
	glossary.Command.SetErr(&buf)
	glossary.Command.SetArgs(append([]string{"--directory", dir}, args...))
	err := glossary.Command.Execute()
	return buf.String(), err
}

func TestGlossaryCommand(t *testing.T) {
	t.Run("writes the whole-model glossary as Markdown without arguments", func(t *testing.T) {
		dir := t.TempDir()

		out, err := runGlossaryCommand(t, dir, glossaryModelYAML, nil)
		require.NoError(t, err)

		assert.Contains(t, out, "# Glossary\n")
		assert.Contains(t, out, "## billing")
		assert.Contains(t, out, "## inventory")
		assert.Contains(t, out, "## ordering")
		assert.Contains(t, out, "### Order")
		assert.Contains(t, out, `_Avoid the term "Basket"._ Reserved for the pre-checkout cart.`)
		assert.Contains(t, out, `_Avoid the term "Bag"._`)
	})

	t.Run("narrows to a single bounded context when given a path", func(t *testing.T) {
		dir := t.TempDir()

		out, err := runGlossaryCommand(t, dir, glossaryModelYAML, []string{"shop/ordering"})
		require.NoError(t, err)

		assert.Contains(t, out, "## ordering")
		assert.NotContains(t, out, "## billing")
		assert.NotContains(t, out, "## inventory")
	})

	t.Run("returns an error for an unknown path segment", func(t *testing.T) {
		dir := t.TempDir()

		_, err := runGlossaryCommand(t, dir, glossaryModelYAML, []string{"shop/nonexistent"})
		assert.Error(t, err)
	})

	t.Run("emits just the heading when no bounded context has ubiquitous language", func(t *testing.T) {
		dir := t.TempDir()

		out, err := runGlossaryCommand(t, dir, boundedContextWithoutLanguageYAML, nil)
		require.NoError(t, err)
		assert.Equal(t, "# Glossary\n", out)
	})
}
