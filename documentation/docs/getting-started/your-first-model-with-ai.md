# Your First Model with AI

This guide shows how to **draft your first ESDM model with the help of an AI coding agent**, instead of writing every `.esdm.yaml` file by hand. By the end you'll have a small, complete model – typically a Domain, a Bounded Context, an Aggregate, an Event, a Command, a Read Model, and a Query – that **lints cleanly without errors or warnings**, produced through a short conversation with the agent.

!!! tip "Prefer to write every file by hand?"

    **[Your First Model by Hand](/getting-started/your-first-model.md)** covers the same scope artifact by artifact, walking through each `.esdm.yaml` file in turn. The two paths produce comparable results – pick whichever way you'd rather work.

The model deliberately stays small. A second Bounded Context, a second Aggregate, and a Context Mapping all belong in the **[Modeling Guide](/modeling-guides/overview.md)** rather than here. We focus on the write-and-read loop on a single Aggregate first.

Before you start, make sure you have **[installed ESDM](/getting-started/installing-esdm.md)** and that `./esdm version` (or `esdm version` on Windows) prints a version number.

## Why This Works

ESDM ships its schemas inside the binary, and **`esdm add-schema`** writes them into your project so an agent can read them locally. The schemas are more than a validation contract – they carry the **vocabulary**, the **file conventions**, and a **project-layout description** the agent should follow. Pointing the agent at those files is enough; you don't have to teach it ESDM in the prompt.

`esdm lint` then closes the loop. Whatever the agent produces, the linter says whether it's a valid model. Anything it flags becomes the next round of the conversation.

## Setting Up the Project

Create a directory for the model and put the `esdm` binary next to it. The materialized schemas, the YAML files the agent writes, and every linter run all live in this directory.

```shell
mkdir my-first-model
cd my-first-model
```

ESDM lints whatever it finds in the directory you point it at, recursively. The convention is **one artifact per `.esdm.yaml` file**; the agent picks that up from the schemas in the next step.

## Materializing the Schemas

From your project's root directory, run:

```shell
./esdm add-schema
```

The command writes the embedded core schema and the embedded extension schemas into a `schemas/` directory inside the project. Commit that directory so every contributor – and every agent run – reads the same schema revision. The **[Editor Support](/getting-started/editor-support.md)** page describes the materialized layout in detail.

If `schemas/` already exists from an earlier setup, refresh it with `./esdm update-schema` instead of re-running `add-schema`.

## Picking an Agent

Any **AI coding agent** that can read and write files in your workspace works for this. Open the agent inside the project directory so it can see `schemas/` and the model files it's going to create.

In practice this means tools like **Claude Code**, **Codex**, **GitHub Copilot**, and others. The starter prompt below is **tool-agnostic** – it doesn't depend on agent-specific features, just on the agent being able to read the schemas and write `.esdm.yaml` files alongside them.

## The Starter Prompt

Paste this prompt into the agent as the opening message:

```text
You are a Domain-Driven Design and Event Sourcing expert helping me model a
domain using ESDM (Event-Sourced Domain Modeling).

Read the ESDM schemas in the current working directory before you write
anything. They define the entire vocabulary, the file conventions, and the
project layout you must follow.

Before producing any YAML, interview me about the domain. Ask what we are
modeling, who the actors are, which events happen, where the consistency
boundaries sit, and how the things involved are identified. Ask one question
at a time and phrase the questions in the vocabulary from the schemas.

When you have enough context, propose the model following the conventions
from the schemas. After the files are written, ask me to run `esdm lint`
and we will work through any findings together.
```

The prompt has four moving parts. The **role** in the first paragraph sets the agent's perspective: a DDD and Event Sourcing expert, not a generic YAML scribe. The **schema instruction** in the second paragraph keeps the agent honest – it has to read the canonical source rather than rely on whatever it thinks ESDM looks like. The **interview style** in the third paragraph keeps the conversation grounded in your domain and uses the vocabulary you'll have to live with afterwards. The **lint loop** in the fourth paragraph turns the linter into a co-author rather than an afterthought.

!!! tip "Try it on the library example"

    If you'd like a concrete domain to anchor the conversation, tell the agent: *"I want to model a city library's cataloging side – acquiring books and listing them."* That's the same scope **[Your First Model by Hand](/getting-started/your-first-model.md)** covers, which makes the agent's output easy to compare against a hand-written reference.

## Verifying the Model

When the agent says it's done, run the linter from the project root:

```shell
./esdm lint
```

A clean model produces **no output and exits with status `0`**. If anything is off, the linter prints findings with locations and a one-line description of the rule that flagged them. Paste the output back to the agent and ask it to fix the issues – it has both the schemas and your conversation history, so it has everything it needs to address each finding in place.

For a richer view of what the agent produced, run `./esdm view`. The **[Running esdm view](/getting-started/running-esdm-view.md)** page covers it in detail.

## Where to Go Next

- **[Your First Model by Hand](/getting-started/your-first-model.md)** walks through the same kind of model artifact by artifact, so you can compare the agent's output against a written-out reference.
- **[Concepts](/concepts/overview.md)** is the full vocabulary, one term per page – useful to fact-check anything the agent proposes.
- **[Running esdm lint](/getting-started/running-esdm-lint.md)** explains the lint workflow in more depth, including CI integration.
- **[Editor Support](/getting-started/editor-support.md)** sets up your editor so it offers autocomplete and validation against the local schemas while you and the agent work.
- **[Modeling Guides](/modeling-guides/overview.md)** grow a small model into a multi-Aggregate one, with Context Mapping and a Process Manager – the natural next step.
