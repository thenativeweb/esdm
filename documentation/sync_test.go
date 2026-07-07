package documentation_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/refgen"
	"github.com/thenativeweb/esdm/schema"
)

// commonFields are present on every core and extension document
// regardless of kind. The reference Overview pages document them
// once, so the per-kind pages do not need to repeat the full
// definition - they just have to mention each field by name.
var commonFields = []string{
	"apiVersion",
	"kind",
	"name",
	"description",
	"metadata",
}

func TestCoreReferenceMentionsEveryKindField(t *testing.T) {
	t.Parallel()

	core, err := loadYAMLDocument(schema.Core())
	require.NoError(t, err)

	index := newSchemaIndex(core)
	kinds := collectCoreKinds(t, core, index)

	for _, kind := range kinds {
		t.Run(kind.name, func(t *testing.T) {
			t.Parallel()

			markdown := readReference(t, filepath.Join("docs", "reference", "core-schema", kind.name+".md"))
			where := "core-schema/" + kind.name + ".md"

			assertCommonFields(t, markdown, where)
			assertInventory(t, markdown, kind.inventory, where)
		})
	}
}

// TestGivenWhenThenReference checks the feature.md and scenario.md
// pair as a single surface. The schema is one document split across
// two reference pages for editorial reasons; the test treats the
// pair as the union it is.
func TestGivenWhenThenReference(t *testing.T) {
	t.Parallel()

	document := loadExtension(t, "given-when-then")
	index := newSchemaIndex(document)
	inventory := index.inventory(document)

	feature := readReference(t, filepath.Join("docs", "extensions", "given-when-then", "reference", "feature.md"))
	scenario := readReference(t, filepath.Join("docs", "extensions", "given-when-then", "reference", "scenario.md"))
	combined := feature + "\n" + scenario
	where := "extensions/given-when-then/reference (feature.md + scenario.md)"

	t.Run("mentions every common field", func(t *testing.T) {
		assertCommonFields(t, combined, where)
	})

	t.Run("mentions every schema field and enum value", func(t *testing.T) {
		assertInventory(t, combined, inventory, where)
	})
}

func TestDomainStorytellingReference(t *testing.T) {
	t.Parallel()

	document := loadExtension(t, "domain-storytelling")
	index := newSchemaIndex(document)
	inventory := index.inventory(document)

	markdown := readReference(t, filepath.Join("docs", "extensions", "domain-storytelling", "reference", "domain-story.md"))
	where := "extensions/domain-storytelling/reference/domain-story.md"

	t.Run("mentions every common field", func(t *testing.T) {
		assertCommonFields(t, markdown, where)
	})

	t.Run("mentions every schema field and enum value", func(t *testing.T) {
		assertInventory(t, markdown, inventory, where)
	})
}

// TestReferenceSnippetsAreUpToDate makes sure the YAML excerpts the
// reference pages embed via pymdownx.snippets are exactly what
// `go run ./cmd/refgen` would produce against the embedded schemas.
// If a schema field, enum value, or const changes without the
// snippets being regenerated, this test fails.
func TestReferenceSnippetsAreUpToDate(t *testing.T) {
	t.Parallel()

	expected, err := refgen.Snippets()
	require.NoError(t, err)

	for _, key := range refgen.SortedPaths(expected) {
		t.Run(key, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join("snippets", filepath.FromSlash(key))
			actual, err := os.ReadFile(path)
			require.NoErrorf(t, err, "snippet %s missing - run `make generate-reference-snippets`", key)
			require.Equalf(
				t,
				string(expected[key]),
				string(actual),
				"snippet %s is out of date - run `make generate-reference-snippets`",
				key,
			)
		})
	}
}

// schemaIndex parses a JSON Schema YAML document and lets the test
// walk it recursively. The walker descends through `properties`,
// `items`, the polymorphic combinators (`oneOf`, `anyOf`, `allOf`),
// the conditional triplets (`if`, `then`, `else`, `not`), and
// follows `$ref: "#/$defs/..."` references - with cycle protection.
type schemaIndex struct {
	root map[string]any
	defs map[string]any
}

// pageInventory captures every name a kind's reference page is
// expected to mention. `fields` are property names (substring match
// in the markdown is sufficient because field names are typically
// distinct camelCase identifiers). `values` are enum entries and
// string constants, which are checked in their backtick form so that
// prose occurrences of common words like `core` or `state` cannot
// satisfy the test by accident.
type pageInventory struct {
	fields []string
	values []string
}

type kindEntry struct {
	name      string
	inventory pageInventory
}

func newSchemaIndex(root map[string]any) *schemaIndex {
	defs, _ := root["$defs"].(map[string]any)
	return &schemaIndex{root: root, defs: defs}
}

func (index *schemaIndex) inventory(node any) pageInventory {
	fields := map[string]struct{}{}
	values := map[string]struct{}{}
	index.walk(node, fields, values, map[string]struct{}{})

	for _, commonField := range commonFields {
		delete(fields, commonField)
	}

	return pageInventory{
		fields: sortedSet(fields),
		values: sortedSet(values),
	}
}

func (index *schemaIndex) walk(node any, fields, values map[string]struct{}, visitedRefs map[string]struct{}) {
	switch node := node.(type) {
	case map[string]any:
		if rawRef, ok := node["$ref"].(string); ok {
			const prefix = "#/$defs/"
			if strings.HasPrefix(rawRef, prefix) {
				name := strings.TrimPrefix(rawRef, prefix)
				if _, isVisited := visitedRefs[name]; isVisited {
					return
				}
				if index.defs != nil {
					if target, ok := index.defs[name]; ok {
						next := make(map[string]struct{}, len(visitedRefs)+1)
						for k := range visitedRefs {
							next[k] = struct{}{}
						}
						next[name] = struct{}{}
						index.walk(target, fields, values, next)
					}
				}
				return
			}
		}

		if properties, ok := node["properties"].(map[string]any); ok {
			for key, value := range properties {
				fields[key] = struct{}{}
				index.walk(value, fields, values, visitedRefs)
			}
		}

		if items, ok := node["items"]; ok {
			index.walk(items, fields, values, visitedRefs)
		}

		for _, key := range []string{"oneOf", "anyOf", "allOf"} {
			if branches, ok := node[key].([]any); ok {
				for _, branch := range branches {
					index.walk(branch, fields, values, visitedRefs)
				}
			}
		}

		for _, key := range []string{"if", "then", "else", "not"} {
			if sub, ok := node[key]; ok {
				index.walk(sub, fields, values, visitedRefs)
			}
		}

		if enum, ok := node["enum"].([]any); ok {
			for _, v := range enum {
				if s, ok := v.(string); ok {
					values[s] = struct{}{}
				}
			}
		}

		if constVal, ok := node["const"]; ok {
			if s, ok := constVal.(string); ok {
				values[s] = struct{}{}
			}
		}
	case []any:
		for _, item := range node {
			index.walk(item, fields, values, visitedRefs)
		}
	}
}

func collectCoreKinds(t *testing.T, core map[string]any, index *schemaIndex) []kindEntry {
	t.Helper()

	allOf, ok := core["allOf"].([]any)
	require.Truef(t, ok, "core schema is missing top-level allOf")

	byName := map[string]any{}
	for _, raw := range allOf {
		entry, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		ifBlock, ok := entry["if"].(map[string]any)
		if !ok {
			continue
		}
		ifProperties, ok := ifBlock["properties"].(map[string]any)
		if !ok {
			continue
		}
		kindNode, ok := ifProperties["kind"].(map[string]any)
		if !ok {
			continue
		}
		kindName, _ := kindNode["const"].(string)
		if kindName == "" {
			continue
		}
		thenBlock, ok := entry["then"]
		if !ok {
			continue
		}
		byName[kindName] = thenBlock
	}

	out := make([]kindEntry, 0, len(byName))
	for name, then := range byName {
		out = append(out, kindEntry{name: name, inventory: index.inventory(then)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

func sortedSet(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for v := range set {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func loadYAMLDocument(data []byte) (map[string]any, error) {
	var document map[string]any
	err := yaml.Unmarshal(data, &document)
	if err != nil {
		return nil, err
	}
	return document, nil
}

func loadExtension(t *testing.T, name string) map[string]any {
	t.Helper()

	extensions, err := schema.Extensions()
	require.NoError(t, err)

	for _, extension := range extensions {
		if extension.Name != name {
			continue
		}
		document, err := loadYAMLDocument(extension.Bytes)
		require.NoError(t, err)
		return document
	}

	t.Fatalf("extension %q not found", name)
	return nil
}

func readReference(t *testing.T, relativePath string) string {
	t.Helper()

	absolutePath, err := filepath.Abs(relativePath)
	require.NoError(t, err)

	contents, err := os.ReadFile(absolutePath)
	require.NoErrorf(t, err, "reference file %s does not exist", absolutePath)

	return string(contents)
}

func assertCommonFields(t *testing.T, contents, where string) {
	t.Helper()

	for _, field := range commonFields {
		assertFieldMentioned(t, contents, field, where)
	}
}

func assertInventory(t *testing.T, contents string, inventory pageInventory, where string) {
	t.Helper()

	for _, field := range inventory.fields {
		assertFieldMentioned(t, contents, field, where)
	}
	for _, value := range inventory.values {
		assertEnumMentioned(t, contents, value, where)
	}
}

func assertFieldMentioned(t *testing.T, contents, field, where string) {
	t.Helper()

	// A raw t.Errorf keeps the failure message precise. assert.Contains
	// would append the full haystack on failure, and here that haystack
	// is an entire reference page - several kilobytes of Markdown that
	// would bury the one field name that is actually missing.
	if !strings.Contains(contents, field) {
		t.Errorf("%s does not mention schema field %q", where, field)
	}
}

func assertEnumMentioned(t *testing.T, contents, value, where string) {
	t.Helper()

	// Enum values often coincide with common English words (`core`,
	// `state`, `static`, `pure`, ...). Requiring the backtick form
	// keeps prose occurrences from satisfying the check by accident.
	// The raw t.Errorf keeps the failure message precise instead of
	// dumping the whole reference page, as in the field check above.
	needle := "`" + value + "`"
	if !strings.Contains(contents, needle) {
		t.Errorf("%s does not mention schema value %q in backticks", where, value)
	}
}
