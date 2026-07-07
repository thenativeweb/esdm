package glossary_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/glossary"
)

const renderedGlossary = `# Glossary

## billing

### Invoice

A demand for payment.

## ordering

### Customer

A person who places orders.

### Order

A request to purchase.

_Avoid the term "Basket"._ Reserved for the cart.

_Avoid the term "Bag"._
`

func TestRender(t *testing.T) {
	t.Run("renders an empty glossary as just the top-level heading", func(t *testing.T) {
		out := glossary.Render(&glossary.Glossary{})
		assert.Equal(t, "# Glossary\n", out)
	})

	t.Run("renders sections, terms, and avoid hints in the agreed Markdown layout", func(t *testing.T) {
		g := &glossary.Glossary{
			Sections: []glossary.Section{
				{
					BoundedContext: "billing",
					Terms: []glossary.Term{
						{Term: "Invoice", Definition: "A demand for payment."},
					},
				},
				{
					BoundedContext: "ordering",
					Terms: []glossary.Term{
						{Term: "Customer", Definition: "A person who places orders."},
						{
							Term:       "Order",
							Definition: "A request to purchase.",
							Avoid: []glossary.Avoid{
								{Term: "Basket", Reason: "Reserved for the cart."},
								{Term: "Bag"},
							},
						},
					},
				},
			},
		}

		assert.Equal(t, renderedGlossary, glossary.Render(g))
	})

	t.Run("ends the output with exactly one trailing newline", func(t *testing.T) {
		g := &glossary.Glossary{
			Sections: []glossary.Section{
				{
					BoundedContext: "ordering",
					Terms:          []glossary.Term{{Term: "Order", Definition: "A request."}},
				},
			},
		}
		out := glossary.Render(g)
		assert.True(t, len(out) >= 2)
		assert.Equal(t, "\n", out[len(out)-1:])
		assert.NotEqual(t, "\n\n", out[len(out)-2:])
	})
}
