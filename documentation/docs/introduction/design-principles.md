# Design Principles

ESDM is shaped by a small number of decisions that we made deliberately and that we expect to keep, even as the language grows. They are listed here so you can recognize the shape of the tool and predict how it will behave in situations the documentation has not yet covered.

## The Model Is Files, Not a Service

ESDM models live as `.esdm.yaml` files in your repository. There is **no central server**, no registry, no online lookup. When you run the linter, it reads your files, parses them against an embedded schema, and reports findings. That is the entire pipeline.

This keeps ESDM **completely offline**. It works in air-gapped environments, in CI runners without outbound network access, and on the train. It also makes the tool trivial to reason about – nothing happens that you cannot reproduce from the source.

## The Schema Is the Contract

Every ESDM document declares an `apiVersion`, and that version pins the document to a specific schema. The schema is **embedded in the binary** – the very same schema every tool in the toolchain validates against. There is no version negotiation, no schema migration, no surprise.

When the schema changes in a non-breaking way, the schema revision goes up but the `apiVersion` stays the same. When a breaking change ever happens, the `apiVersion` moves to a new major, and old documents continue to validate against the old schema. This is the same versioning discipline you find in Kubernetes API groups, and it has the same property: **stability is the default**, change is opt-in.

## Extensions Sit Alongside the Core, Not Inside It

Beyond the core vocabulary, ESDM defines **extensions** – independent schemas that describe artifacts the core leaves out. The Given-When-Then extension models behavioral scenarios. The Domain Storytelling extension captures discovery stories.

Crucially, extensions do not inject kinds into the core. A core document validates against the core schema and only the core schema; an extension document validates against its own. That asymmetry is what lets us add new extensions over time without ever touching the core.

## Rules Are Fixed, Not Configurable

ESDM has a single, fixed catalog of linter rules. You cannot disable them, you cannot change their severity, and there is no project-level configuration file that overrides any of this. A rule either applies or it doesn't, and that decision is made by us, not by a YAML file in your repository.

This is opinionated by design. A linter that ships with knobs ends up describing **a hundred different dialects** of the same language, and the value of a shared model collapses with it. We pick the rules carefully, we remove rules that turn out to be wrong, and we trust the result.

The trade-off is that some rules will occasionally feel too strict for a specific situation. When that happens, the right move is usually to revisit the model – the rule is almost always pointing at something real.

## Diagnostics Are Locations, Not Stack Traces

Every diagnostic ESDM emits points at a **file, a line, and a column**. There are no stack traces, no internal error frames, no error codes that require a lookup table. The message reads as a sentence, the location is a place you can navigate to, and the fix is in your editor.

The principle is that the toolchain reports problems in the user's vocabulary, not in its own.

## Where to Go Next

If you want to see how these principles play out in practice, **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** drafts a hands-on example through a short conversation with a coding agent, or **[Your First Model by Hand](/getting-started/your-first-model.md)** walks through the same scope artifact by artifact. If you want to understand the parts of the language, the **[Concepts](/concepts/overview.md)** section is the next step. The **[License](/introduction/license.md)** page describes how ESDM is provided.
