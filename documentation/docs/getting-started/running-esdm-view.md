# Running esdm view

This page walks through the **view workflow**: how to render an ESDM model as a hierarchical tree, how to read the connections it draws, how to focus on one region of the model, and how to surface the per-artifact details when you need them.

This page assumes you have a small model on disk. Either **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** or **[Your First Model by Hand](/getting-started/your-first-model.md)** produces one of the right size; the concrete examples in this page show the output for the by-hand model, since it lets us pin every artifact name and field down for reference. The **[CLI reference](/reference/cli.md)** documents every flag in detail; this page focuses on the workflow.

## Rendering the Model

From inside the directory that holds your `.esdm.yaml` files:

```shell
./esdm view
```

The output for the Your First Model by Hand example looks like this:

```text
domain library  1 bc
└─ bounded-context cataloging  1 agg · 1 rm · 1 qry
   ├─ aggregate book  1 cmd · 1 evt
   │  ├─ command acquire  → acquired
   │  └─ event acquired  ← acquire
   ├─ read-model books  ← 1 evt
   └─ query list-books  → books
```

The tree mirrors the model hierarchy: Domain at the top, then Bounded Contexts, then per-Bounded-Context the consistency units (Aggregates, Dynamic Consistency Boundaries) and the read-side artifacts (Read Models, Queries). Each row carries the kind, the name, and a short summary of what's attached to it.

## Reading the Connections

Two arrow glyphs encode the direction of a relationship. The **`→`** points outward, from the artifact on the left to the one on the right – `command acquire  → acquired` means "this Command publishes the `acquired` Event". The **`←`** points inward, toward the artifact on the left – `event acquired  ← acquire` means "this Event is published by the `acquire` Command", and `read-model books  ← 1 evt` means "this Read Model consumes one Event".

The same direction reading applies to Queries: `query list-books  → books` means "this Query reads from the `books` Read Model".

Counts like `1 cmd · 1 evt` summarize what would otherwise be repetitive sub-rows. They give a fast read of the size and shape of an Aggregate without having to scroll through every Command and Event.

## Filtering by Path

When the model grows, rendering everything at once becomes noise. Pass a **path** to focus on one region. The path follows the hierarchy: Domain, Bounded Context, consistency unit, separated by slashes.

```shell
./esdm view library/cataloging/book
```

The output trims the tree to the requested subtree:

```text
aggregate book  1 cmd · 1 evt
├─ command acquire  → acquired
└─ event acquired  ← acquire
```

Use shorter paths to scope at a coarser level – `library` shows just the Domain summary plus its Bounded Contexts; `library/cataloging` shows one Bounded Context and everything inside it.

## Showing Details

The default view stays at the structural level – names and connections. To inspect what each artifact actually holds, add `--with-details`:

```shell
./esdm view --with-details
```

The output gains an indented detail line per artifact:

```text
domain library  1 bc
└─ bounded-context cataloging  1 agg · 1 rm · 1 qry
   ├─ aggregate book  1 cmd · 1 evt
   │     identifiedBy: state.isbn
   │  ├─ command acquire  → acquired
   │  │     data: {author, isbn, title}
   │  └─ event acquired  ← acquire
   │        data: {author, isbn, title}
   ├─ read-model books  ← 1 evt
   │     projects cataloging/book/acquired - Append a row with `title`, `author`, and `isbn`.
   │     paradigm: tabular
   └─ query list-books  → books
```

The details are deliberately compact: the identifier strategy on the Aggregate, the field set on Command and Event payloads, the projection rule on the Read Model. The view is meant to be skimmed, not read line by line.

`--with-details` and a path filter combine. `./esdm view library/cataloging/book --with-details` zooms into the Aggregate and shows its details in one step.

## Where to Go Next

- **[Running esdm lint](/getting-started/running-esdm-lint.md)** is the structural counterpart – correctness, not exploration.
- **[Editor Support](/getting-started/editor-support.md)** sets up autocomplete and validation while you write.
- **[CLI: esdm view](/reference/cli.md#esdm-view)** documents every flag and the path syntax in full.
