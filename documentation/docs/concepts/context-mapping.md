# Context Mapping

A **Context Mapping** describes the relationship between two **[Bounded Contexts](/concepts/bounded-context.md)** – or between a Bounded Context and an **[External System](/concepts/external-system.md)**. It is the model's record of *how these two parts talk to each other*.

Context Mappings are a primitive of strategic Domain-Driven Design. The patterns they capture have names – customer-supplier, conformist, anti-corruption-layer, shared-kernel, partnership, separate-ways, open-host-service, published-language – and each pattern describes a different power dynamic, a different kind of translation, and a different set of risks.

## Why Map Contexts Explicitly

A Bounded Context that depends on another without saying so is a piece of architecture without a name. *We use the customer service somehow* is a sentence that hides decisions, and the decisions tend to come back as bugs.

When you write down that the relationship is customer-supplier, you commit to a contract direction. When you write down that it is anti-corruption-layer, you commit to a translation point. When you write down that it is separate-ways, you commit to *not* integrating – and the model holds you to it.

## Asymmetric and Symmetric

Context Mappings come in two shapes. **Asymmetric** mappings (customer-supplier, conformist, anti-corruption-layer, open-host-service, published-language) name the two sides distinctly – a customer and a supplier, a conformist and an upstream, a downstream and an upstream. The role names are part of the mapping, so direction is not a comment but a fact.

**Symmetric** mappings (shared-kernel, partnership, separate-ways) treat the two sides as peers. They name two participating Bounded Contexts. Sharing a kernel or forming a partnership with a third-party External System doesn't fit the model – those mappings live between Bounded Contexts only.

Each mapping represents exactly two endpoints. Richer topologies are expressed as multiple mappings – two customer-supplier arrows, not one three-party document.
