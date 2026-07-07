package glossary_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/glossary"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/modelpath"
	"github.com/thenativeweb/esdm/runner"
)

const glossaryModelYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shop
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: shop
ubiquitousLanguage:
  - term: Order
    definition: A customer's request to purchase one or more products.
    avoid:
      - term: Basket
        reason: Reserved for the pre-checkout cart.
      - term: Bag
  - term: Customer
    definition: A person who places orders.
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: billing
scope:
  domain: shop
ubiquitousLanguage:
  - term: Invoice
    definition: A demand for payment for an order.
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: shipping
scope:
  domain: shop
---
apiVersion: schema.esdm.io/core/v1
kind: domain
name: warehouse
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: inventory
scope:
  domain: warehouse
ubiquitousLanguage:
  - term: SKU
    definition: A stock keeping unit.
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

func TestBuild(t *testing.T) {
	t.Run("collects every bounded context with ubiquitous language, sorted by name", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		g, err := glossary.Build(m, modelpath.Path{})
		require.NoError(t, err)

		names := make([]string, 0, len(g.Sections))
		for _, s := range g.Sections {
			names = append(names, s.BoundedContext)
		}
		assert.Equal(t, []string{"billing", "inventory", "ordering"}, names)
	})

	t.Run("omits a bounded context that has no ubiquitous language", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		g, err := glossary.Build(m, modelpath.Path{})
		require.NoError(t, err)

		for _, s := range g.Sections {
			assert.NotEqual(t, "shipping", s.BoundedContext)
		}
	})

	t.Run("sorts terms alphabetically within a section", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		g, err := glossary.Build(m, modelpath.Path{})
		require.NoError(t, err)

		var ordering glossary.Section
		for _, s := range g.Sections {
			if s.BoundedContext == "ordering" {
				ordering = s
			}
		}
		require.Len(t, ordering.Terms, 2)
		assert.Equal(t, "Customer", ordering.Terms[0].Term)
		assert.Equal(t, "Order", ordering.Terms[1].Term)
		assert.Equal(t, "A customer's request to purchase one or more products.", ordering.Terms[1].Definition)
	})

	t.Run("captures avoid entries with and without a reason in document order", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		g, err := glossary.Build(m, modelpath.Path{})
		require.NoError(t, err)

		var order glossary.Term
		for _, s := range g.Sections {
			for _, term := range s.Terms {
				if term.Term == "Order" {
					order = term
				}
			}
		}
		require.Len(t, order.Avoid, 2)
		assert.Equal(t, glossary.Avoid{Term: "Basket", Reason: "Reserved for the pre-checkout cart."}, order.Avoid[0])
		assert.Equal(t, glossary.Avoid{Term: "Bag", Reason: ""}, order.Avoid[1])
	})

	t.Run("narrows to a single domain when given a one-segment path", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		g, err := glossary.Build(m, modelpath.Path{Segments: []string{"shop"}})
		require.NoError(t, err)

		names := make([]string, 0, len(g.Sections))
		for _, s := range g.Sections {
			names = append(names, s.BoundedContext)
		}
		assert.Equal(t, []string{"billing", "ordering"}, names)
	})

	t.Run("narrows to a single bounded context when given a two-segment path", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		g, err := glossary.Build(m, modelpath.Path{Segments: []string{"shop", "ordering"}})
		require.NoError(t, err)

		require.Len(t, g.Sections, 1)
		assert.Equal(t, "ordering", g.Sections[0].BoundedContext)
	})

	t.Run("returns an empty glossary for an existing bounded context without ubiquitous language", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		g, err := glossary.Build(m, modelpath.Path{Segments: []string{"shop", "shipping"}})
		require.NoError(t, err)
		assert.Empty(t, g.Sections)
	})

	t.Run("rejects an unknown domain", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		_, err := glossary.Build(m, modelpath.Path{Segments: []string{"nonexistent"}})
		assert.Error(t, err)
	})

	t.Run("rejects an unknown bounded context under a known domain", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		_, err := glossary.Build(m, modelpath.Path{Segments: []string{"shop", "nonexistent"}})
		assert.Error(t, err)
	})

	t.Run("rejects a path that reaches below the bounded-context level", func(t *testing.T) {
		m := loadModel(t, glossaryModelYAML)

		_, err := glossary.Build(m, modelpath.Path{Segments: []string{"shop", "ordering", "order"}})
		assert.Error(t, err)
	})
}
