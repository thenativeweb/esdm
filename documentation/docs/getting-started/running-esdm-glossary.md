# Running esdm glossary

This page walks through the **glossary workflow**: how to turn the ubiquitous language declared in an ESDM model into a human-readable Markdown glossary, how to read what it produces, how to scope it to one region of the model, and how to keep the result around as a file.

This page assumes you have a small model on disk. Either **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** or **[Your First Model by Hand](/getting-started/your-first-model.md)** produces one of the right size; the examples below build on the by-hand `library` model's `cataloging` context. The **[CLI reference](/reference/cli.md#esdm-glossary)** documents every flag in detail; this page focuses on the workflow.

## What the Glossary Reads

A glossary is only as good as the terminology the model actually records. ESDM keeps that terminology on each Bounded Context, in its `ubiquitousLanguage` block – the shared, agreed vocabulary of that context, the heart of how Domain-Driven Design keeps language and model aligned. `esdm glossary` reads exactly that block and nothing else; it does not invent definitions from artifact names.

Suppose the `cataloging` Bounded Context of the `library` model declares its language like this:

```yaml
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: cataloging
scope:
  domain: library
language: en
ubiquitousLanguage:
  - term: Acquisition
    definition: The process of adding a book to the catalog, whether bought or donated.
    avoid:
      - term: Purchase
        reason: Not every acquisition is bought – donations are acquisitions too.
  - term: Catalog
    definition: The authoritative list of every book the library holds.
```

The `language` field names the language the terms are written in – here English. It is required as soon as a Bounded Context declares a `ubiquitousLanguage`, because the glossary renders per language and every context has to say which one it speaks.

## Generating the Glossary

From inside the directory that holds your `.esdm.yaml` files:

```shell
./esdm glossary
```

The command writes Markdown to stdout:

```text
# Glossary

## cataloging

### Acquisition

The process of adding a book to the catalog, whether bought or donated.

_Avoid the term "Purchase"._ Not every acquisition is bought – donations are acquisitions too.

### Catalog

The authoritative list of every book the library holds.
```

The structure mirrors the model: a single `# Glossary` heading, one `##` section per Bounded Context, and one `###` entry per term. Terms are sorted alphabetically within each section, and sections are sorted by Bounded Context name, so the output is stable across runs and diffs cleanly.

## Reading the Avoid Hints

Each discouraged alternative recorded under `avoid` becomes its own short paragraph beneath the term: an italicized *Avoid the term …* sentence, followed by the reason as a plain sentence when the model gives one. Several discouraged alternatives simply stack as separate paragraphs. These notes are the part teams reach for most in review – they catch the synonyms that quietly pull a model out of alignment.

!!! info

    A Bounded Context without a `ubiquitousLanguage` block contributes no section, and a model without any ubiquitous language at all produces just the `# Glossary` heading. That's the expected result, not an error – the exit code is still `0`.

## Rendering Another Language

Domain experts often speak a different language than the code. A term can carry **translations** – the same concept in another language, with its own definition and its own rejected alternatives. Suppose the `Acquisition` entry above gains a German one:

```yaml
  - term: Acquisition
    definition: The process of adding a book to the catalog, whether bought or donated.
    avoid:
      - term: Purchase
        reason: Not every acquisition is bought – donations are acquisitions too.
    translations:
      - language: de
        term: Erwerb
        definition: Die Aufnahme eines Buchs in den Katalog, ob gekauft oder gespendet.
        avoid:
          - term: Kauf
            reason: Auch Spenden sind ein Erwerb.
```

Pass `--language` with a BCP 47 tag to render the glossary in that language:

```shell
./esdm glossary --language de
```

```text
# Glossary

## cataloging

### Catalog

The authoritative list of every book the library holds.

_No translation into de._

### Erwerb

Die Aufnahme eines Buchs in den Katalog, ob gekauft oder gespendet.

_Avoid the term "Kauf"._ Auch Spenden sind ein Erwerb.
```

A Bounded Context whose own `language` is the requested one is rendered as is. Everywhere else the translations take over, including their rejected alternatives. A term without a translation keeps its primary form and is marked with an italic note, so the glossary shows where the language is still incomplete instead of quietly dropping the term – the same idea as the avoid hints: the gaps are the useful part.

## Filtering by Path

When the model spans several Domains and Bounded Contexts, you'll often want just one. Pass a **path** to scope the output. The path follows the model hierarchy – a Domain, then a Bounded Context inside it – with the two segments separated by a slash. The **[esdm view](/getting-started/running-esdm-view.md)** command uses the same path shape, so what you learn here carries over.

```shell
./esdm glossary library/cataloging
```

A single segment stays at the Domain level and emits every Bounded Context inside it; two segments narrow to one Bounded Context. An unknown or too-deep segment is rejected as invalid input, so a typo fails loudly instead of silently producing an empty glossary.

## Keeping the Glossary as a File

The output is plain Markdown with no terminal coloring, which makes it safe to redirect straight into a file:

```shell
./esdm glossary > glossary.md
```

Committing that file next to the model gives non-technical stakeholders – product owners, domain experts, new team members – a readable artifact they can review without the ESDM toolchain, and regenerating it after a model change keeps it honest.

## Where to Go Next

- **[Running esdm view](/getting-started/running-esdm-view.md)** renders the structure of the model; the glossary renders its language.
- **[Running esdm lint](/getting-started/running-esdm-lint.md)** checks the model for correctness before you publish a glossary from it.
- **[Bounded Context](/concepts/bounded-context.md)** explains where the `ubiquitousLanguage` block lives and why it is scoped per context.
- **[CLI: esdm glossary](/reference/cli.md#esdm-glossary)** documents every flag and the path syntax in full.
