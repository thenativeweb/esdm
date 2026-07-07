package refgen

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func requireSnippets(t *testing.T) map[string][]byte {
	t.Helper()

	snippets, err := Snippets()
	require.NoError(t, err)

	return snippets
}

func TestSnippets(t *testing.T) {
	t.Run("contains a snippet for every core kind that carries fields", func(t *testing.T) {
		snippets := requireSnippets(t)

		for _, kind := range []string{"aggregate", "event", "command", "domain"} {
			assert.Contains(t, snippets, "reference/core-schema/"+kind+".yaml")
		}
	})

	t.Run("contains a snippet for every extension kind", func(t *testing.T) {
		snippets := requireSnippets(t)

		assert.Contains(t, snippets, "reference/extensions/given-when-then/feature.yaml")
		assert.Contains(t, snippets, "reference/extensions/domain-storytelling/domain-story.yaml")
	})

	t.Run("contains the scenario variant snippets", func(t *testing.T) {
		snippets := requireSnippets(t)

		variants := 0

		for key := range snippets {
			if strings.HasPrefix(key, "reference/extensions/given-when-then/scenario-") {
				variants++
			}
		}

		assert.Equal(t, 4, variants)
	})

	t.Run("emits only valid, non-empty YAML", func(t *testing.T) {
		snippets := requireSnippets(t)

		for key, data := range snippets {
			assert.NotEmpty(t, data, key)

			var node yaml.Node
			assert.NoError(t, yaml.Unmarshal(data, &node), key)
		}
	})

	t.Run("resolves every internal reference", func(t *testing.T) {
		snippets := requireSnippets(t)

		for key, data := range snippets {
			assert.NotContains(t, string(data), "$ref", key)
		}
	})

	t.Run("strips the schema bookkeeping fields", func(t *testing.T) {
		snippets := requireSnippets(t)

		for key, data := range snippets {
			text := string(data)

			assert.NotContains(t, text, "$id:", key)
			assert.NotContains(t, text, "$schema:", key)
			assert.NotContains(t, text, "$defs:", key)
			assert.NotContains(t, text, "x-esdm-schema-revision:", key)
		}
	})
}

func TestSortedPaths(t *testing.T) {
	t.Run("returns every key in lexicographic order", func(t *testing.T) {
		snippets := map[string][]byte{
			"b.yaml": []byte("b"),
			"a.yaml": []byte("a"),
			"c.yaml": []byte("c"),
		}

		paths := SortedPaths(snippets)

		require.Len(t, paths, 3)
		assert.True(t, sort.StringsAreSorted(paths))
	})
}
