# Bounded Context

A **Bounded Context** is the largest unit inside which a single, consistent vocabulary applies. The same word can mean different things in different Bounded Contexts; inside a Bounded Context, every term has exactly one meaning.

Bounded Contexts live inside the **[Domain](/concepts/domain.md)** (and may be classified by a **[Subdomain](/concepts/subdomain.md)**). They hold the **consistency units** that actually do the work – **[Aggregates](/concepts/aggregate.md)**, **[Dynamic Consistency Boundaries](/concepts/dynamic-consistency-boundary.md)**, **[Process Managers](/concepts/process-manager.md)**, and **[Read Models](/concepts/read-model.md)**.

## Why the Boundary Matters

The boundary is a translation point. When a concept crosses from one Bounded Context to another, you don't carry the meaning along; you re-interpret. A `customer` in a sales Bounded Context might be a *Lead*, and the same `customer` in a billing Bounded Context might be an *AccountHolder*. Both are real, both are right, and the Bounded Context is what tells you which one you're in.

That is why a Bounded Context is also the place where **[Context Mappings](/concepts/context-mapping.md)** become meaningful. A Context Mapping describes how two Bounded Contexts translate between each other – which parts of one are exposed, which parts of the other are consumed, and what shape they take in the middle.

## One Namespace for the Domain Types

Inside a Bounded Context, the **Aggregates, Dynamic Consistency Boundaries, Entities, Value Objects, and Domain Services share one namespace**. They become the types of one module in code, and nothing gets composed into their names the way an Aggregate is composed into its Events, so two of them cannot share a name: an Aggregate and a Dynamic Consistency Boundary both called `loan`, or an Entity and a Value Object both called `money`, is an error the linter reports.

Everything else may share a name, because it plays a different role or lives in a different module. An Aggregate and a Read Model both called `order` are the same concept on the write and the read side; an Entity and an Actor both called `applicant` are the data and the role; a Subdomain and a Bounded Context both called `billing` are the classification and the context. The model keeps them apart by kind, and `esdm view` shows every element that matches a path.

## The Ubiquitous Language

A Bounded Context's most valuable artifact is its **ubiquitous language** – the canonical terms that apply inside it, each with its definition. The language is one entry per term, deliberately without aliases: the whole point of choosing one term is that there is one term. Where helpful, an entry can also list rejected alternatives, so a reader who looks up a synonym still lands on the canonical word.

The vocabulary is written in **one language**, and the Bounded Context names it in its `language` field – `en`, `de`, `de-AT` – as soon as it declares a `ubiquitousLanguage`. That language is the one the model and the code speak. Domain experts often speak another one, and a term can carry **translations** into it: the same concept, with its own term, its own definition, and its own rejected alternatives that explain why it was translated this way and not otherwise. A translation is not a second term. Within each language there is still exactly one, so the one-term principle holds per language rather than being given up. The **[glossary](/getting-started/running-esdm-glossary.md)** renders any of these languages on request.
