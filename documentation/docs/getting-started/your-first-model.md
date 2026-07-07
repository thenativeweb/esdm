# Your First Model by Hand

This guide walks you through writing your **first ESDM model from scratch** and verifying it with the linter. By the end you will have a small, complete model of a city library's cataloging side – a Domain, a Bounded Context, an Aggregate, an Event, a Command, a Read Model, and a Query – that **lints cleanly without errors or warnings**.

!!! tip "Prefer to draft the model with an AI agent?"

    **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** covers the same scope through a short conversation with a coding agent that reads the ESDM schemas and writes the YAML for you. The two paths produce comparable results – pick whichever way you'd rather work.

The model deliberately stays small. Lending – borrowing a book, returning it – is the natural next step, and lives in the **[Modeling Guide](/modeling-guides/overview.md)** rather than here, because it requires a second Bounded Context, a second Aggregate, and a Context Mapping. We focus on the write-and-read loop on a single Aggregate first.

Before you start, make sure you have **[installed ESDM](/getting-started/installing-esdm.md)** and that `./esdm version` (or `esdm version` on Windows) prints a version number.

## Setting Up the Project

Create a directory for the model and put the `esdm` binary next to it. Each artifact will live in its own `.esdm.yaml` file inside that directory.

```shell
mkdir my-first-model
cd my-first-model
```

ESDM lints whatever it finds under the directory you point it at, recursively, so the layout below is just a convention – every artifact in its own file, all of them in one directory. Real projects group documents by Bounded Context; with one Bounded Context, a flat layout reads more clearly.

## The Domain

Start with the Domain – the top-level area the model describes. Domains carry no kind-specific fields; the document is the common shape only.

```yaml
--8<-- "your-first-model/model/domain.esdm.yaml"
```

Save this as `domain.esdm.yaml`.

## The Bounded Context

A Bounded Context is the largest unit inside which a single, consistent vocabulary applies. We pick **`cataloging`**: everything that is about getting books into the library and keeping them findable. Lending lives elsewhere; that's a different vocabulary, a different consistency story.

```yaml
--8<-- "your-first-model/model/bounded-context.esdm.yaml"
```

Save this as `bounded-context.esdm.yaml`.

## The Aggregate

An Aggregate is a container-based consistency unit. Every Command and Event in this model targets the **`book`** Aggregate, which carries the data for one library book.

```yaml
--8<-- "your-first-model/model/book.esdm.yaml"
```

Three things to notice:

- The `state` is a JSON Schema that describes one book: `title`, `author`, `isbn`. The fields are required because every book in the catalog must have all three.
- `identifiedBy` says how a book is identified. We use `source: state` with `field: isbn` – the ISBN already lives in `state` and is unique per book, so it's the natural identifier. There's no need for a generated UUID.
- The Aggregate carries no behavior. Behavior lives in the Commands and Events that target it.

Save this as `book.esdm.yaml`.

## The Event

An Event is the immutable record of a fact that happened. The Aggregate is `book`, so the Event's scope is `book`-bound. The Event itself is named **`acquired`**, not `book-acquired`: the surrounding scope already conveys the Aggregate, so a `book-` prefix would be redundant – ESDM's linter actively discourages it.

```yaml
--8<-- "your-first-model/model/acquired.esdm.yaml"
```

The `metadata.annotations` block carries the **CloudEvents type** for this Event – `io.eventsourcingdb.library.book-acquired` – following the convention documented in **[EventSourcingDB's event-types guide](https://docs.eventsourcingdb.io/fundamentals/event-types/)**. The CloudEvents type is non-semantic for ESDM; it's a tooling hint that downstream consumers (event stores, integrators) can read off the model.

Save this as `acquired.esdm.yaml`.

## The Command

A Command expresses intent to change the model. **`acquire`** carries the same payload shape as the Event it produces, and `publishes: [acquired]` ties them together.

```yaml
--8<-- "your-first-model/model/acquire.esdm.yaml"
```

Two things to notice:

- The Command's `data` is required, even when – as here – it happens to mirror the Event's `data`. ESDM keeps the two schemas separate so commands and events can evolve independently.
- `publishes` lists Events by their bare names. The surrounding scope (`book`) fixes the producer, so naming the scope again would be redundant.

Save this as `acquire.esdm.yaml`.

## The Read Model

The write side is now complete: a Command is issued, an Event is produced. The read side projects Events into a query-optimized shape. Our **`books`** Read Model is a tabular list, one row per book.

```yaml
--8<-- "your-first-model/model/books.esdm.yaml"
```

The `projections` array carries one entry per Event the Read Model consumes. Each entry is a flat object combining an Event reference (`boundedContext`, `aggregate`, `event`) with a prose `rule` describing how the Event updates the Read Model. The rule is descriptive, not executable; ESDM models the *what*, the implementation provides the *how*.

The optional `paradigm: tabular` is a hint for tooling. The schema describes what the Read Model exposes to its Queries, not the storage layout underneath.

Save this as `books.esdm.yaml`.

## The Query

A Query is a read operation served by a Read Model. **`list-books`** returns the entire list.

```yaml
--8<-- "your-first-model/model/list-books.esdm.yaml"
```

The `readModel: books` field binds the Query to the Read Model in the same Bounded Context. The `result` schema mirrors the Read Model's `schema` because we're returning the whole table, but they're two separate schemas in principle – a Query may filter, sort, or shape its result independently of the underlying Read Model.

Save this as `list-books.esdm.yaml`.

## Linting the Model

With all seven files in place, run the linter from inside the directory:

```shell
./esdm lint
```

A clean model produces **no output and exits with status `0`**. If anything is off, the linter prints findings with locations and a one-line description of the rule that flagged them. Try editing one of the files – delete a required field, rename a reference – and re-run; the diagnostic points you at exactly what changed.

For a richer view of the model, run `./esdm view` from the same directory; the **[Running esdm view](/getting-started/running-esdm-view.md)** page covers it in detail.

## Where to Go Next

- **[Running esdm lint](/getting-started/running-esdm-lint.md)** explains the lint workflow in more depth, including CI integration.
- **[Running esdm view](/getting-started/running-esdm-view.md)** shows how to render and explore the model you just wrote.
- **[Editor Support](/getting-started/editor-support.md)** sets up your editor so it offers autocomplete and validation against the ESDM schemas while you type.
- **[Modeling Guides](/modeling-guides/overview.md)** picks the same library and grows it into the lending side, with multiple Aggregates, Context Mapping, and a Process Manager.
