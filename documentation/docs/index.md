# Event-Sourced Domain Modeling

Welcome to the official documentation for **ESDM** – the Event-Sourced Domain Modeling language.

**ESDM** describes event-sourced domains as YAML and ships with the tools to manage them. The language captures the building blocks of **Domain-Driven Design**, **CQRS**, and **Event Sourcing** – Aggregates, Events, Commands, Process Managers, Read Models, Context Mappings, and the rest – along with the artifacts that surround modeling work, such as **Domain Storytelling** discoveries and **Given-When-Then** specifications.

Whether you model by hand, build tools that consume domain models, lean on AI to model or to analyze code, or simply want a written record of an event-sourced system you already built, this documentation is the starting point.

!!! tip "Get ESDM"

    ESDM ships as pre-built binaries for **macOS**, **Linux**, and **Windows**. **[Download and install ESDM](/getting-started/installing-esdm.md)**. New to ESDM? Start with **[What is ESDM](/introduction/what-is-esdm.md)**.

## Pick your path

Different visitors want different things. Pick the one that matches what you're here for.

### Modeling for the first time

You're learning Domain-Driven Design or Event Sourcing, or you want to capture a model from scratch. Start with the basics and walk through a guided example.

- **[Getting Started](/getting-started/installing-esdm.md)**

    Install ESDM and write your first model from scratch.

- **[Modeling Guides](/modeling-guides/overview.md)**

    Follow a worked example end to end.

### Modeling with AI

You want an LLM to help you draft a model, or to extract one from existing code. ESDM's YAML is plain enough that LLMs can read and write it directly, and the Concepts pages give the model exactly the vocabulary it needs.

- **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)**

    A short conversation with a coding agent that produces a lint-clean model.

- **[Concepts](/concepts/overview.md)**

    The full vocabulary, one term per page.

- **[Recipes](/recipes/overview.md)**

    Focused answers to specific modeling questions.

### Documenting an existing system

You already have an event-sourced system and want to capture it as a model. Start with the vocabulary and the schema reference – they describe what every kind of artifact looks like in ESDM.

- **[Concepts](/concepts/overview.md)**

    Define the parts of the language in your own words first.

- **[Reference](/reference/cli.md)**

    Look up CLI commands, schema fields, and exact field names.

### Building tools that consume ESDM

You're building tooling – validators, generators, transformers, IDE plugins – that interoperates with ESDM. The schema reference is the contract you build against, and the extensions show how the format scales beyond the core.

- **[Reference: Core Schema](/reference/core-schema/overview.md)**

    The canonical description of every kind in the core schema.

- **[Extensions](/extensions/overview.md)**

    Given-When-Then and Domain Storytelling, each with its own schema.

### Already know ESDM, just need a lookup

Skip the prose, jump straight to the answer.

- **[CLI Reference](/reference/cli.md)**

- **[Core Schema Reference](/reference/core-schema/overview.md)**

- **[Extensions](/extensions/overview.md)**

## Licensing

ESDM is **open source** under the **MIT license**.

- **[Licensing Details](/introduction/license.md)**

    Learn how ESDM is provided and what that means in practice.

## Need Support?

If you or your team need help designing, integrating, or scaling an event-sourced system, we're happy to assist. Just reach out to **[hello@thenativeweb.io](mailto:hello@thenativeweb.io)**.
