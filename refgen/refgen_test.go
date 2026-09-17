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
			if !strings.HasSuffix(key, ".yaml") {
				continue
			}

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

func TestRuleSnippets(t *testing.T) {
	t.Run("contains a Markdown snippet for the core rules and for every extension with rules", func(t *testing.T) {
		snippets := requireSnippets(t)

		assert.Contains(t, snippets, "reference/linter-rules/core.md")
		assert.Contains(t, snippets, "reference/linter-rules/given-when-then.md")
		assert.Contains(t, snippets, "reference/linter-rules/domain-storytelling.md")
	})

	t.Run("renders a rule as a heading with an explicit anchor, its severity, and its description", func(t *testing.T) {
		snippets := requireSnippets(t)
		core := string(snippets["reference/linter-rules/core.md"])

		assert.Contains(t, core, "### `esdm/structure/ambiguous-name` { #structure-ambiguous-name }\n\nSeverity: `error`\n\n")
		assert.Contains(t, core, "### `esdm/modeling/aggregate-without-commands` { #modeling-aggregate-without-commands }\n\nSeverity: `warning`\n\n")
	})

	t.Run("groups the core rules under one heading per category, structure first", func(t *testing.T) {
		snippets := requireSnippets(t)
		core := string(snippets["reference/linter-rules/core.md"])

		assert.True(t, strings.HasPrefix(core, "## Structure\n\n### "), "core snippet starts with the Structure heading")
		assert.Contains(t, core, "\n## Modeling\n\n### ")
		assert.Less(t, strings.Index(core, "## Structure"), strings.Index(core, "## Modeling"))
	})

	t.Run("sorts the rules of a category by ID", func(t *testing.T) {
		snippets := requireSnippets(t)
		core := string(snippets["reference/linter-rules/core.md"])

		var previous string
		for _, line := range strings.Split(core, "\n") {
			if strings.HasPrefix(line, "## ") {
				previous = ""
				continue
			}
			if !strings.HasPrefix(line, "### ") {
				continue
			}
			assert.Less(t, previous, line)
			previous = line
		}
	})

	t.Run("keeps the extension rules off the core snippet", func(t *testing.T) {
		snippets := requireSnippets(t)
		core := string(snippets["reference/linter-rules/core.md"])

		assert.NotContains(t, core, "esdm/gwt/")
		assert.NotContains(t, core, "/story-")
	})

	t.Run("places every extension rule on the snippet of its extension", func(t *testing.T) {
		snippets := requireSnippets(t)
		gwt := string(snippets["reference/linter-rules/given-when-then.md"])
		story := string(snippets["reference/linter-rules/domain-storytelling.md"])

		assert.Contains(t, gwt, "### `esdm/gwt/scenario-without-then` { #gwt-scenario-without-then }")
		assert.Contains(t, story, "### `esdm/structure/story-without-sentences` { #structure-story-without-sentences }")
		assert.Contains(t, story, "### `esdm/modeling/story-orphan-actor` { #modeling-story-orphan-actor }")
	})

	t.Run("omits the category heading when a snippet has a single category", func(t *testing.T) {
		snippets := requireSnippets(t)
		gwt := string(snippets["reference/linter-rules/given-when-then.md"])

		assert.True(t, strings.HasPrefix(gwt, "### `esdm/gwt/"), "gwt snippet starts with its first rule")
		assert.NotContains(t, gwt, "\n## ")
	})

	t.Run("ends every Markdown snippet with exactly one newline", func(t *testing.T) {
		snippets := requireSnippets(t)

		for key, data := range snippets {
			if !strings.HasSuffix(key, ".md") {
				continue
			}
			text := string(data)
			assert.True(t, strings.HasSuffix(text, "\n"), key)
			assert.False(t, strings.HasSuffix(text, "\n\n"), key)
		}
	})
}
