package docgen_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/docgen"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/modelpath"
	"github.com/thenativeweb/esdm/runner"
)

const modelYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: library
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: cataloging
scope:
  domain: library
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: book
scope:
  domain: library
  boundedContext: cataloging
identifiedBy:
  source: state
  field: isbn
state:
  type: object
  properties:
    title: { type: string }
    isbn: { type: string }
  required: [title, isbn]
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: acquire
scope:
  domain: library
  boundedContext: cataloging
  aggregate: book
data:
  type: object
  properties:
    title: { type: string }
    isbn: { type: string }
  required: [title, isbn]
publishes:
  - acquired
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: acquired
scope:
  domain: library
  boundedContext: cataloging
  aggregate: book
data:
  type: object
  properties:
    title: { type: string }
    isbn: { type: string }
  required: [title, isbn]
---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: books
scope:
  domain: library
  boundedContext: cataloging
paradigm: tabular
schema:
  type: array
  items:
    type: object
    properties:
      title: { type: string }
      isbn: { type: string }
    required: [title, isbn]
projections:
  - boundedContext: cataloging
    aggregate: book
    event: acquired
    rule: Append a row.
---
apiVersion: schema.esdm.io/core/v1
kind: query
name: list-books
scope:
  domain: library
  boundedContext: cataloging
readModel: books
result:
  type: array
  items:
    type: object
    properties:
      title: { type: string }
      isbn: { type: string }
    required: [title, isbn]
`

func loadModel(t *testing.T, yaml string) *model.Model {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "model.esdm.yaml"), []byte(yaml), 0o644))
	_, m, err := runner.RunWithModel(context.Background(), dir)
	require.NoError(t, err)
	require.NotNil(t, m)
	return m
}

func paths(pages []docgen.Page) []string {
	out := make([]string, 0, len(pages))
	for _, page := range pages {
		out = append(out, page.Path)
	}
	sort.Strings(out)
	return out
}

func pageByPath(t *testing.T, pages []docgen.Page, path string) docgen.Page {
	t.Helper()
	for _, page := range pages {
		if page.Path == path {
			return page
		}
	}
	require.Failf(t, "page not found", "no page at %q", path)
	return docgen.Page{}
}

func TestBuild(t *testing.T) {
	t.Run("places every element at its containment path", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{
			"README.md",
			"library/README.md",
			"library/cataloging/README.md",
			"library/cataloging/book/README.md",
			"library/cataloging/book/acquire.md",
			"library/cataloging/book/acquired.md",
			"library/cataloging/books.md",
			"library/cataloging/list-books.md",
		}, paths(pages))
	})

	t.Run("writes a container as a README index and a leaf as a named file", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		names := paths(pages)
		assert.Contains(t, names, "library/cataloging/book/README.md")
		assert.Contains(t, names, "library/cataloging/book/acquired.md")
	})

	t.Run("narrows to a single element and its subtree, keeping full paths", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{Segments: []string{"library", "cataloging", "book"}})
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{
			"library/cataloging/book/README.md",
			"library/cataloging/book/acquire.md",
			"library/cataloging/book/acquired.md",
		}, paths(pages))
	})

	t.Run("renders the reference and contents on a container page", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		book := pageByPath(t, pages, "library/cataloging/book/README.md")
		assert.Contains(t, book.Content, "# book")
		assert.Contains(t, book.Content, "Reference: `esdm:library/cataloging/book` (aggregate)")
		assert.Contains(t, book.Content, "## Commands")
		assert.Contains(t, book.Content, "[acquire](acquire.md)")
		assert.Contains(t, book.Content, "## Events")
		assert.Contains(t, book.Content, "[acquired](acquired.md)")
	})

	t.Run("renders the reference on a leaf page", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		acquired := pageByPath(t, pages, "library/cataloging/book/acquired.md")
		assert.Contains(t, acquired.Content, "# acquired")
		assert.Contains(t, acquired.Content, "Reference: `esdm:library/cataloging/book/acquired` (event)")
	})

	t.Run("indexes the domains from the root page", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		root := pageByPath(t, pages, "README.md")
		assert.Contains(t, root.Content, "# Documentation")
		assert.Contains(t, root.Content, "[library](library/README.md)")
	})

	t.Run("links a bounded context to its bounded-context-scoped elements", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		cataloging := pageByPath(t, pages, "library/cataloging/README.md")
		assert.Contains(t, cataloging.Content, "[book](book/README.md)")
		assert.Contains(t, cataloging.Content, "[books](books.md)")
		assert.Contains(t, cataloging.Content, "[list-books](list-books.md)")
	})

	t.Run("rejects an unknown top-level segment", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		_, err := docgen.Build(m, modelpath.Path{Segments: []string{"nonexistent"}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "under model root")
	})

	t.Run("rejects an unknown segment under a known one", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		_, err := docgen.Build(m, modelpath.Path{Segments: []string{"library", "nonexistent"}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), `under "library"`)
	})

	t.Run("rejects a path that reaches below a leaf", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		_, err := docgen.Build(m, modelpath.Path{Segments: []string{"library", "cataloging", "book", "acquired", "deeper"}})
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "acquired"))
	})
}

func TestDetails(t *testing.T) {
	t.Run("renders an aggregate's identity and state fields", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		book := pageByPath(t, pages, "library/cataloging/book/README.md")
		assert.Contains(t, book.Content, "## Identity")
		assert.Contains(t, book.Content, "`isbn` field, from the state")
		assert.Contains(t, book.Content, "## State")
		assert.Contains(t, book.Content, "- `title`")
	})

	t.Run("renders a command's payload and links its published events", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		acquire := pageByPath(t, pages, "library/cataloging/book/acquire.md")
		assert.Contains(t, acquire.Content, "## Payload")
		assert.Contains(t, acquire.Content, "## Publishes")
		assert.Contains(t, acquire.Content, "- [acquired](acquired.md)")
	})

	t.Run("links a read model's projections to the projected events", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		books := pageByPath(t, pages, "library/cataloging/books.md")
		assert.Contains(t, books.Content, "## Paradigm")
		assert.Contains(t, books.Content, "## Projections")
		assert.Contains(t, books.Content, "- [acquired](book/acquired.md):")
	})

	t.Run("links a query to its read model", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{})
		require.NoError(t, err)

		listBooks := pageByPath(t, pages, "library/cataloging/list-books.md")
		assert.Contains(t, listBooks.Content, "## Read Model")
		assert.Contains(t, listBooks.Content, "[books](books.md)")
	})

	t.Run("falls back to the reference when a linked target is outside the output", func(t *testing.T) {
		m := loadModel(t, modelYAML)

		pages, err := docgen.Build(m, modelpath.Path{Segments: []string{"library", "cataloging", "books"}})
		require.NoError(t, err)

		books := pageByPath(t, pages, "library/cataloging/books.md")
		assert.Contains(t, books.Content, "`esdm:library/cataloging/book/acquired`")
	})
}
