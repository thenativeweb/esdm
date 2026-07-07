# What is ESDM

ESDM is the **Event-Sourced Domain Modeling** language. It describes event-sourced domains – Domains, Bounded Contexts, Aggregates, Events, Commands, Process Managers, Read Models, and the relationships between them – as a set of YAML files.

ESDM is **a YAML language plus a built-in toolchain**. The schema defines what an ESDM document may contain. The toolchain reads, validates, and renders those documents. Today the binary ships with a linter that catches structural and modeling errors, a command that renders a hierarchical summary of a model, and helpers that materialize the schemas locally for editor support. More tools may follow.

## Why a Standard Format

Event Sourcing and Domain-Driven Design come with a lot of vocabulary – Aggregates, Bounded Contexts, Process Managers, Dynamic Consistency Boundaries, Read Models, Context Mappings. In practice, the language is shared, but the **artifacts are not**. Teams capture their models in slide decks, whiteboard photos, `README.md` files, or scattered code comments, and the model drifts the moment it leaves the room.

ESDM gives the model a **first-class home**. The model is files, the files live next to the code, the files are reviewed in pull requests, and a shared format means every tool – yours, ours, third-party – speaks the same language about your domain.

The **standardization** is the point. When the format is fixed, you can build tooling against it: linters, validators, generators, transformers, IDE integrations, AI-assisted modelers. ESDM provides the schema and a starting toolchain; the format is open enough that more tools, internal or external, can grow around it.

## Extensions

Beyond the core vocabulary, ESDM ships **extension schemas** for artifacts that surround the modeling work. **Domain Storytelling** captures discovery stories – Actors, Work Objects, and the activities that connect them – as their own kind of document. **Given-When-Then** captures behavioral specifications – preceding Events, a triggering action, expected outcomes – on top of the core kinds.

Each extension follows the same format conventions as the core and is validated by the same toolchain. Extensions never inject kinds into the core, and core documents never validate against an extension – the asymmetry keeps the core lean while letting the surface grow.

## What ESDM Is Not

ESDM is **not** an event store, **not** a runtime, **not** a code generator, **not** a framework. It does not execute your domain. It does not know about your database, your message bus, or your deployment pipeline. It is a static, descriptive layer that lives **alongside** your code and helps you keep the model honest.

That separation is deliberate. A modeling language that tries to also be a framework forces you into specific runtime choices, and a runtime that tries to also be a modeling language buries the model under implementation detail. ESDM keeps the two apart so you can mix it freely with whatever you actually run in production.

## Where to Go Next

If you want to see ESDM in action, head to **[Installing ESDM](/getting-started/installing-esdm.md)** and write your first model. If you want to understand the vocabulary first, the **[Concepts](/concepts/overview.md)** section walks through every kind the core schema defines. The **[Reference](/reference/cli.md)** section is the canonical place to look up a CLI command or a schema field once you're past the basics.
