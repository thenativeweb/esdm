# Running esdm documentation

This page walks through the **documentation workflow**: how to render an ESDM model as a tree of Markdown pages on disk, how the pages and their paths relate to the model, how to narrow the output to one region, and how to keep the result next to your code.

This page assumes you have a small model on disk. Either **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** or **[Your First Model by Hand](/getting-started/your-first-model.md)** produces one of the right size; the examples below show the output for the by-hand `library` model. The **[CLI reference](/reference/cli.md#esdm-documentation)** documents every flag in detail.

## Where the Command Sits

Two commands already present the model. **[esdm view](/getting-started/running-esdm-view.md)** renders it as a tree in the terminal, transient and compact; **[esdm glossary](/getting-started/running-esdm-glossary.md)** writes the ubiquitous language as one Markdown document. `esdm documentation` writes the **whole model as a directory of Markdown pages** – one page per element, cross-linked – ready to be committed next to the code, browsed on GitHub, or fed to a static-site generator. It renders the same tree `esdm view` renders, so the two never disagree about where an element sits.

## Generating the Tree

From inside the directory that holds your `.esdm.yaml` files, name an output directory:

```shell
./esdm documentation --output docs
```

The directory is required and has no default, so a tree is never written by accident. The command writes these files for the library model:

```text
docs/
  README.md
  domain_library/
    README.md
    bounded-context_cataloging/
      README.md
      aggregate_book/
        README.md
        command_acquire.md
        event_acquired.md
      read-model_books/
        README.md
      query_list-books.md
```

Every element becomes one page. An element that contains others – a Domain, a Bounded Context, an Aggregate, a Dynamic Consistency Boundary, a Process Manager, a Read Model – becomes a **directory with a `README.md`**, which GitHub renders as the folder's landing page. Everything else becomes a **file**. The root `README.md` lists the Domains and the Context Mappings.

## Reading a Page

Here is the page of the `book` Aggregate, `domain_library/bounded-context_cataloging/aggregate_book/README.md`:

```markdown
# book

Aggregate `esdm:domain=library/bounded-context=cataloging/aggregate=book`

1 cmd · 1 evt

## Details

- identifiedBy: state.isbn

## Commands

- [acquire](command_acquire.md) – publishes [acquired](event_acquired.md)

## Events

- [acquired](event_acquired.md) – published by [acquire](command_acquire.md)
```

The heading is the element's bare name. The line below states its kind and its **[reference](/reference/reference-notation.md)**, so a reader can cite the element from anywhere. Then come the same stats `esdm view` shows, the element's description if it has one, its details as `esdm view --with-details` lists them, and its children grouped by kind. **Relationships are written out and linked**: a Command publishes Events and is issued by Actors, an Event is published by Commands, a Query reads a Read Model, a Read Model projects Events. A Bounded Context's page renders its ubiquitous language with each term's translations, a Context Mapping's page its endpoints and term pairs.

## Paths and References

A page's path follows the element's reference, segment by segment: `esdm:domain=library/bounded-context=cataloging/aggregate=book` becomes `domain_library/bounded-context_cataloging/aggregate_book/README.md`. The rule is small enough to apply by hand – replace the equals sign of each segment with an underscore, then append `.md` for a leaf or `/README.md` for a container – and it means that two elements of different kinds sharing a name, an Aggregate and a Read Model both called `order`, say, get two pages that cannot collide.

## Focusing on One Region

Pass a **path** to write only one region of the model. The path follows the hierarchy, Domain, Bounded Context, consistency unit, separated by slashes, exactly as for **[esdm view](/getting-started/running-esdm-view.md)**:

```shell
./esdm documentation --output docs library/cataloging/book
```

The pages keep their full paths, so `aggregate_book/README.md` still lands under `domain_library/bounded-context_cataloging/`, and references keep resolving. A link to an element outside the written region falls back to that element's reference in plain text, so no link is ever broken.

## Writing Again

The command refuses to write into a directory that is not empty, so a stale tree is never silently mixed with a fresh one. Pass `--force` to clear the directory first:

```shell
./esdm documentation --output docs --force
```

Because the whole tree is regenerated, no page from an earlier run survives, and the directory mirrors the model exactly. Committing the output next to the model gives everyone who reads the repository a browsable picture of the domain, and regenerating it after a model change keeps that picture true.

## Where to Go Next

- **[Reference Notation](/reference/reference-notation.md)** explains the `esdm:` references every page carries.
- **[Running esdm view](/getting-started/running-esdm-view.md)** renders the same tree in the terminal.
- **[CLI: esdm documentation](/reference/cli.md#esdm-documentation)** documents every flag and the path syntax in full.
