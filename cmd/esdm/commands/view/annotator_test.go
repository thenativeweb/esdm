package view_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/view"
	"github.com/thenativeweb/esdm/diag"
)

func TestAnnotate(t *testing.T) {
	t.Run("attaches matching severity to a node by exact location", func(t *testing.T) {
		root := &view.Node{
			Kind:     "domain",
			Name:     "sample",
			Location: view.SourceLocation{File: "domain.esdm.yaml", Line: 3, Column: 7},
		}
		view.Annotate(root, []diag.Diagnostic{
			{Severity: diag.SeverityWarning, Location: diag.Location{File: "domain.esdm.yaml", Line: 3, Column: 7}},
		})
		assert.Equal(t, view.SeverityWarning, root.Severity)
	})

	t.Run("bubbles a child's severity up to its ancestors", func(t *testing.T) {
		child := &view.Node{
			Kind:     "aggregate",
			Name:     "widget",
			Location: view.SourceLocation{File: "widget.esdm.yaml", Line: 5, Column: 7},
		}
		root := &view.Node{
			Kind:     "domain",
			Name:     "sample",
			Location: view.SourceLocation{File: "domain.esdm.yaml", Line: 3, Column: 7},
			Children: []*view.Node{child},
		}
		view.Annotate(root, []diag.Diagnostic{
			{Severity: diag.SeverityError, Location: diag.Location{File: "widget.esdm.yaml", Line: 5, Column: 7}},
		})
		assert.Equal(t, view.SeverityError, child.Severity)
		assert.Equal(t, view.SeverityError, root.Severity)
	})

	t.Run("error severity wins over warning when both target the same node", func(t *testing.T) {
		root := &view.Node{
			Kind:     "domain",
			Name:     "sample",
			Location: view.SourceLocation{File: "f.esdm.yaml", Line: 1, Column: 1},
		}
		view.Annotate(root, []diag.Diagnostic{
			{Severity: diag.SeverityWarning, Location: diag.Location{File: "f.esdm.yaml", Line: 1, Column: 1}},
			{Severity: diag.SeverityError, Location: diag.Location{File: "f.esdm.yaml", Line: 1, Column: 1}},
		})
		assert.Equal(t, view.SeverityError, root.Severity)
	})

	t.Run("falls back to the deepest enclosing node when no exact match exists", func(t *testing.T) {
		// A diagnostic anchored inside the aggregate's
		// invariants block (line 20) should land on the
		// aggregate (whose name is on line 5), not the
		// domain.
		domain := &view.Node{
			Kind:     "domain",
			Name:     "sample",
			Location: view.SourceLocation{File: "domain.esdm.yaml", Line: 3, Column: 7},
		}
		agg := &view.Node{
			Kind:     "aggregate",
			Name:     "widget",
			Location: view.SourceLocation{File: "agg.esdm.yaml", Line: 5, Column: 7},
		}
		domain.Children = []*view.Node{agg}

		view.Annotate(domain, []diag.Diagnostic{
			{Severity: diag.SeverityError, Location: diag.Location{File: "agg.esdm.yaml", Line: 20, Column: 5}},
		})
		assert.Equal(t, view.SeverityError, agg.Severity)
		assert.Equal(t, view.SeverityError, domain.Severity)
	})
}
