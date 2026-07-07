# Subdomain

A **Subdomain** is a strategic classification of part of the **[Domain](/concepts/domain.md)**. It groups **[Bounded Contexts](/concepts/bounded-context.md)** by their value to the business, so the model captures not only what the system does but also where investment matters.

Subdomains are orthogonal to Bounded Contexts, not a hierarchy on top of them. A Bounded Context can exist without a Subdomain, and Subdomains classify Bounded Contexts after the fact.

## Core, Supporting, and Generic

Domain-Driven Design distinguishes three kinds of Subdomain. **Core** Subdomains carry the strategic value of the business and deserve the most modeling care. **Supporting** Subdomains are necessary but not differentiating – specific enough to build in-house, but not strategic. **Generic** Subdomains are commodity – authentication, payment, geocoding – and are usually solved by a third-party product.

The classification is part of the model, not an aside. Recording it makes the strategic distinction visible to whoever reads the model later, and it shapes how teams reason about investment, ownership, and build-vs-buy decisions.
