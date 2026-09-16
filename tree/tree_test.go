package tree_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/modelpath"
	"github.com/thenativeweb/esdm/runner"
	"github.com/thenativeweb/esdm/tree"
)

const twoDomainsYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shop
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: shop
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: order
scope:
  domain: shop
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: order
scope:
  domain: shop
  boundedContext: ordering
projections: []
schema:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: domain
name: warehouse
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

func kindsAndNames(nodes []*tree.Node) []string {
	var out []string
	for _, n := range nodes {
		out = append(out, n.Kind+" "+n.Name)
	}
	return out
}

func TestBuild(t *testing.T) {
	t.Run("returns a synthetic root with one child per domain, sorted by name", func(t *testing.T) {
		root := tree.Build(loadModel(t, twoDomainsYAML), false)

		assert.Equal(t, "model", root.Kind)
		assert.Equal(t, []string{"domain shop", "domain warehouse"}, kindsAndNames(root.Children))
	})

	t.Run("places elements at the position their scope names", func(t *testing.T) {
		root := tree.Build(loadModel(t, twoDomainsYAML), false)

		shop := root.Children[0]
		require.Equal(t, []string{"bounded-context ordering"}, kindsAndNames(shop.Children))
		assert.Equal(t, []string{"aggregate order", "read-model order"}, kindsAndNames(shop.Children[0].Children))
	})
}

func TestNarrow(t *testing.T) {
	t.Run("returns the root unchanged for an empty path", func(t *testing.T) {
		root := tree.Build(loadModel(t, twoDomainsYAML), false)

		narrowed, err := tree.Narrow(root, modelpath.Path{})
		require.NoError(t, err)
		assert.Same(t, root, narrowed)
	})

	t.Run("keeps every sibling that matches the last segment under a synthetic root", func(t *testing.T) {
		root := tree.Build(loadModel(t, twoDomainsYAML), false)

		narrowed, err := tree.Narrow(root, modelpath.Path{Segments: []string{"shop", "ordering", "order"}})
		require.NoError(t, err)
		assert.Equal(t, "model", narrowed.Kind)
		assert.Equal(t, []string{"aggregate order", "read-model order"}, kindsAndNames(narrowed.Children))
	})

	t.Run("rejects an unknown segment", func(t *testing.T) {
		root := tree.Build(loadModel(t, twoDomainsYAML), false)

		_, err := tree.Narrow(root, modelpath.Path{Segments: []string{"shop", "nowhere"}})
		require.Error(t, err)
		assert.Equal(t, `no element "nowhere" under "shop"`, err.Error())
	})
}
