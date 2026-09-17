package rules_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/rules"
	"github.com/thenativeweb/esdm/schema"
)

func TestMeta(t *testing.T) {
	meta := rules.Meta{ID: "esdm/structure/ambiguous-name"}

	t.Run("derives the category from the ID", func(t *testing.T) {
		assert.Equal(t, "structure", meta.Category())
	})

	t.Run("derives the name from the ID", func(t *testing.T) {
		assert.Equal(t, "ambiguous-name", meta.Name())
	})

	t.Run("derives the documentation anchor from category and name", func(t *testing.T) {
		assert.Equal(t, "structure-ambiguous-name", meta.Anchor())
	})
}

// TestCatalogExtensions guards the Extension field of every
// rule in the catalog. The field is what places a rule on the
// Linter Rules page of its extension instead of the core page,
// so a forgotten assignment would silently misfile the rule.
func TestCatalogExtensions(t *testing.T) {
	extensions, err := schema.Extensions()
	require.NoError(t, err)

	knownExtensions := map[string]bool{}
	for _, extension := range extensions {
		knownExtensions[extension.Name] = true
	}

	t.Run("declares only the empty string or a known extension", func(t *testing.T) {
		for _, rule := range rules.Catalog() {
			meta := rule.Meta()
			if meta.Extension == "" {
				continue
			}
			assert.Truef(t, knownExtensions[meta.Extension], "rule %s declares unknown extension %q", meta.ID, meta.Extension)
		}
	})

	// The gwt category exists only for rules over the
	// given-when-then extension, so category and extension
	// must agree in both directions.
	t.Run("assigns exactly the gwt rules to the given-when-then extension", func(t *testing.T) {
		for _, rule := range rules.Catalog() {
			meta := rule.Meta()
			isGwtCategory := meta.Category() == "gwt"
			isGwtExtension := meta.Extension == "given-when-then"
			assert.Equalf(t, isGwtCategory, isGwtExtension, "rule %s: category %q does not match extension %q", meta.ID, meta.Category(), meta.Extension)
		}
	})

	// The domain-storytelling rules share the structure and
	// modeling categories with the core rules; their story-
	// prefix is the only trace of their origin in the ID, so
	// the explicit field is what keeps them off the core page.
	t.Run("assigns exactly the story rules to the domain-storytelling extension", func(t *testing.T) {
		for _, rule := range rules.Catalog() {
			meta := rule.Meta()
			isStoryRule := strings.HasPrefix(meta.Name(), "story-")
			isStoryExtension := meta.Extension == "domain-storytelling"
			assert.Equalf(t, isStoryRule, isStoryExtension, "rule %s: name %q does not match extension %q", meta.ID, meta.Name(), meta.Extension)
		}
	})
}
