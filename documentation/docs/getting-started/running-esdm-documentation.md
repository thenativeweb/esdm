# Running esdm documentation

This page walks through the **documentation workflow**: how to render an ESDM model as a tree of Markdown files, how to read a page it produces, how to scope the output to one region of the model, and how to regenerate it safely.

This page assumes you have a small model on disk. Either **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** or **[Your First Model by Hand](/getting-started/your-first-model.md)** produces one of the right size; the examples below build on the by-hand `library` model's `cataloging` context. The **[CLI reference](/reference/cli.md#esdm-documentation)** documents every flag in detail; this page focuses on the workflow.

## What esdm documentation Produces

Where **[esdm view](/getting-started/running-esdm-view.md)** prints a transient summary to the terminal and **[esdm glossary](/getting-started/running-esdm-glossary.md)** writes a single document, `esdm documentation` writes a whole **directory tree** to disk: one Markdown page per element, laid out along the model's containment hierarchy. A Domain is a directory, the Bounded Contexts inside it are directories, and so on down to the Commands and Events, which are files.

The layout is the point. **Each element's page sits at the path you would use to name it** – the `acquired` Event of the `book` Aggregate lives at `library/cataloging/book/acquired.md` – so the tree doubles as an addressing scheme, and a reader who knows an element's place in the domain knows where its page is.

## Generating the Tree

The output directory is required, because writing a whole tree into the wrong place should never happen by accident. Pass it with `-o` / `--output`:

```shell
./esdm documentation --output ./docs
```

From the `library` model, that writes:

```text
docs/
  README.md
  library/
    README.md
    cataloging/
      README.md
      book/
        README.md
        acquire.md
        acquired.md
      books.md
      list-books.md
```

An element that contains others – a Domain, a Bounded Context, an Aggregate – becomes a directory with a `README.md` index; a leaf element becomes `<name>.md`. **GitHub renders each `README.md` as the landing page of its directory**, so browsing the tree on GitHub reads top-down without any configuration.

## Reading a Page

Every page opens with the element's name and its reference, then renders the detail the model carries. The `book` Aggregate's index page looks like this:

```markdown
# book

Reference: `esdm:library/cataloging/book` (aggregate)

## Identity

By its `isbn` field, from the state.

## State

- `author`
- `isbn`
- `title`

## Commands

- [acquire](acquire.md)

## Events

- [acquired](acquired.md)
```

The detail sections mirror what `esdm view --with-details` shows for each kind: an Aggregate's identity and state, a Command's payload and the Events it publishes, a Read Model's projections, and so on. **Wherever one element names another, the page links to it** with a relative link, so the tree is navigable: the `acquire` Command links to the `acquired` Event, and the `list-books` Query links to the `books` Read Model. When a linked element falls outside the generated output, the page shows its reference instead, so a link is never broken.

## Filtering by Path

When the model spans several Domains and Bounded Contexts, you'll often want just one. Pass a **path** to scope the output to a subtree. The path follows the model hierarchy and uses the same shape as **[esdm view](/getting-started/running-esdm-view.md)**, so what you learn there carries over.

```shell
./esdm documentation --output ./docs library/cataloging
```

This writes only the `cataloging` subtree, and – importantly – **keeps the full path prefix**, so `book`'s page is still at `library/cataloging/book/README.md`. The pages line up with their references whether you render the whole model or one corner of it. An unknown or too-deep segment is rejected as invalid input, so a typo fails loudly instead of producing an empty tree.

## Regenerating Safely

A documentation tree should mirror the model exactly, with no pages left over from elements you have since renamed or removed. To protect against writing over unrelated files, `esdm documentation` **refuses to write into a directory that already has content**:

```text
output directory "./docs" is not empty; pass --force to clear and rewrite it
```

Passing `--force` clears the output directory first and then writes the fresh tree, so the result always reflects the current model and nothing else:

```shell
./esdm documentation --output ./docs --force
```

Because `--force` deletes the directory's existing contents, point `--output` at a directory that holds only generated documentation – not at a directory you also keep other work in.

## Publishing the Tree

The output is neutral Markdown with relative links and no tool-specific configuration – no `mkdocs.yml`, no theme files. That keeps it portable: GitHub renders it directly, and a static-site generator such as MkDocs can pick it up as its `docs` directory. Committing the tree next to the model gives everyone a browsable, always-current view of the domain, and regenerating it after a model change keeps it honest.

Linter findings do not block the output; as long as the model resolves, the command renders whatever it finds. Run **[esdm lint](/getting-started/running-esdm-lint.md)** first when you want the model checked for correctness before you publish from it.

## Where to Go Next

- **[Running esdm view](/getting-started/running-esdm-view.md)** renders the same structure to the terminal for a quick look; `esdm documentation` writes it to disk to keep and publish.
- **[Running esdm glossary](/getting-started/running-esdm-glossary.md)** writes just the ubiquitous language as a single document; the documentation tree includes it on each Bounded Context's page.
- **[Running esdm lint](/getting-started/running-esdm-lint.md)** checks the model for correctness before you publish a tree from it.
- **[CLI: esdm documentation](/reference/cli.md#esdm-documentation)** documents every flag and the path syntax in full.
