package view_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/view"
)

func TestRender(t *testing.T) {
	t.Run("renders a single node", func(t *testing.T) {
		root := &view.Node{Kind: "domain", Name: "sample"}
		out := view.Render(root, view.RenderOptions{})
		assert.Equal(t, "domain sample\n", out)
	})

	t.Run("renders the synthetic model root transparently", func(t *testing.T) {
		root := &view.Node{Kind: "model", Children: []*view.Node{
			{Kind: "domain", Name: "sample"},
		}}
		out := view.Render(root, view.RenderOptions{})
		assert.Equal(t, "domain sample\n", out)
	})

	t.Run("renders children with box-drawing connectors", func(t *testing.T) {
		root := &view.Node{Kind: "domain", Name: "sample", Children: []*view.Node{
			{Kind: "bounded-context", Name: "context-one"},
			{Kind: "bounded-context", Name: "context-two"},
		}}
		out := view.Render(root, view.RenderOptions{})
		expected := strings.Join([]string{
			"domain sample",
			"├─ bounded-context context-one",
			"└─ bounded-context context-two",
			"",
		}, "\n")
		assert.Equal(t, expected, out)
	})

	t.Run("renders nested children with continuation pipes", func(t *testing.T) {
		root := &view.Node{Kind: "domain", Name: "sample", Children: []*view.Node{
			{Kind: "bounded-context", Name: "context-one", Children: []*view.Node{
				{Kind: "aggregate", Name: "widget"},
			}},
			{Kind: "bounded-context", Name: "context-two", Children: []*view.Node{
				{Kind: "aggregate", Name: "gadget"},
			}},
		}}
		out := view.Render(root, view.RenderOptions{})
		expected := strings.Join([]string{
			"domain sample",
			"├─ bounded-context context-one",
			"│  └─ aggregate widget",
			"└─ bounded-context context-two",
			"   └─ aggregate gadget",
			"",
		}, "\n")
		assert.Equal(t, expected, out)
	})

	t.Run("renders tags and stats inline on the header", func(t *testing.T) {
		root := &view.Node{
			Kind:  "subdomain",
			Name:  "primary",
			Tags:  []string{"core"},
			Stats: []string{"1 BC"},
		}
		out := view.Render(root, view.RenderOptions{})
		assert.Equal(t, "subdomain primary (core)  1 BC\n", out)
	})

	t.Run("renders severity glyph in plain mode", func(t *testing.T) {
		root := &view.Node{Kind: "aggregate", Name: "widget", Severity: view.SeverityError}
		out := view.Render(root, view.RenderOptions{})
		assert.Contains(t, out, "✗")
	})

	t.Run("renders severity glyph with ANSI when colors are enabled", func(t *testing.T) {
		root := &view.Node{Kind: "aggregate", Name: "widget", Severity: view.SeverityWarning}
		out := view.Render(root, view.RenderOptions{Colors: true})
		assert.Contains(t, out, "\x1b[33m⚠\x1b[0m")
	})

	t.Run("renders the name in bold when colors are enabled", func(t *testing.T) {
		root := &view.Node{Kind: "aggregate", Name: "widget"}
		out := view.Render(root, view.RenderOptions{Colors: true})
		assert.Contains(t, out, "\x1b[1mwidget\x1b[0m")
	})

	t.Run("renders the kind in cyan when colors are enabled", func(t *testing.T) {
		root := &view.Node{Kind: "aggregate", Name: "widget"}
		out := view.Render(root, view.RenderOptions{Colors: true})
		assert.Contains(t, out, "\x1b[36maggregate\x1b[0m")
	})

	t.Run("renders tags and stats dim when colors are enabled", func(t *testing.T) {
		root := &view.Node{
			Kind:  "subdomain",
			Name:  "primary",
			Tags:  []string{"core"},
			Stats: []string{"1 BC"},
		}
		out := view.Render(root, view.RenderOptions{Colors: true})
		assert.Contains(t, out, "\x1b[2m(core)\x1b[0m")
		assert.Contains(t, out, "\x1b[2m1 BC\x1b[0m")
	})

	t.Run("does not emit ANSI bold, dim or cyan in plain mode", func(t *testing.T) {
		root := &view.Node{Kind: "aggregate", Name: "widget", Tags: []string{"alpha"}, Stats: []string{"1 thing"}}
		out := view.Render(root, view.RenderOptions{})
		assert.NotContains(t, out, "\x1b[1m")
		assert.NotContains(t, out, "\x1b[2m")
		assert.NotContains(t, out, "\x1b[36m")
	})

	t.Run("renders detail lines dim when colors are enabled", func(t *testing.T) {
		root := &view.Node{Kind: "aggregate", Name: "widget", Lines: []string{"identifiedBy: generated/uuid"}}
		out := view.Render(root, view.RenderOptions{Colors: true, ShowDetails: true})
		assert.Contains(t, out, "\x1b[2midentifiedBy: generated/uuid\x1b[0m")
	})

	t.Run("emits per-node lines only when ShowDetails is true", func(t *testing.T) {
		root := &view.Node{Kind: "aggregate", Name: "widget", Lines: []string{"identifiedBy: generated/uuid"}}
		without := view.Render(root, view.RenderOptions{})
		assert.NotContains(t, without, "identifiedBy")
		with := view.Render(root, view.RenderOptions{ShowDetails: true})
		assert.Contains(t, with, "identifiedBy: generated/uuid")
	})
}
