# Bounded Context

A **Bounded Context** is the largest unit inside which a single, consistent vocabulary applies. The same word can mean different things in different Bounded Contexts; inside a Bounded Context, every term has exactly one meaning.

Bounded Contexts live inside the **[Domain](/concepts/domain.md)** (and may be classified by a **[Subdomain](/concepts/subdomain.md)**). They hold the **consistency units** that actually do the work – **[Aggregates](/concepts/aggregate.md)**, **[Dynamic Consistency Boundaries](/concepts/dynamic-consistency-boundary.md)**, **[Process Managers](/concepts/process-manager.md)**, and **[Read Models](/concepts/read-model.md)**.

## Why the Boundary Matters

The boundary is a translation point. When a concept crosses from one Bounded Context to another, you don't carry the meaning along; you re-interpret. A `customer` in a sales Bounded Context might be a *Lead*, and the same `customer` in a billing Bounded Context might be an *AccountHolder*. Both are real, both are right, and the Bounded Context is what tells you which one you're in.

That is why a Bounded Context is also the place where **[Context Mappings](/concepts/context-mapping.md)** become meaningful. A Context Mapping describes how two Bounded Contexts translate between each other – which parts of one are exposed, which parts of the other are consumed, and what shape they take in the middle.

## The Ubiquitous Language

A Bounded Context's most valuable artifact is its **ubiquitous language** – the canonical terms that apply inside it, each with its definition. The language is one entry per term, deliberately without aliases: the whole point of choosing one term is that there is one term. Where helpful, an entry can also list rejected alternatives, so a reader who looks up a synonym still lands on the canonical word.
