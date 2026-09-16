package glossary_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
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

	// glossary.Command is a package-level cobra command, so a
	// flag set by one run (for example --language) would
	// otherwise stay in effect for every run after it.
	glossary.Command.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
		flag.Changed = false
	})

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

	t.Run("renders the requested language with --language", func(t *testing.T) {
		dir := t.TempDir()

		out, err := runGlossaryCommand(t, dir, translatedModelYAML, []string{"--language", "de"})
		require.NoError(t, err)

		assert.Contains(t, out, "### Bestellung")
		assert.Contains(t, out, `_Avoid the term "Auftrag"._ Used for production orders elsewhere.`)
		assert.Contains(t, out, "### Customer\n\nA person who places orders.\n\n_No translation into de._")
		assert.NotContains(t, out, "### Order")
	})

	t.Run("rejects a --language value that is not a BCP 47 tag", func(t *testing.T) {
		dir := t.TempDir()

		_, err := runGlossaryCommand(t, dir, translatedModelYAML, []string{"--language", "German"})
		require.Error(t, err)
		assert.Equal(t, `invalid language "German": expected a BCP 47 language tag such as "de" or "de-AT"`, err.Error())
	})

	t.Run("emits just the heading when no bounded context has ubiquitous language", func(t *testing.T) {
		dir := t.TempDir()

		out, err := runGlossaryCommand(t, dir, boundedContextWithoutLanguageYAML, nil)
		require.NoError(t, err)
		assert.Equal(t, "# Glossary\n", out)
	})
}
