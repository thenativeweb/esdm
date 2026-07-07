# ESDM

ESDM is the Event-Sourced Domain Modeling language. This repository holds the `esdm` CLI written in Go – a linter and toolset for `.esdm.yaml` models –, the embedded schemas that define the language, and the documentation site published at https://www.esdm.io.

## Workflow

- Lint and test with `make qa` (add `WITH_RACE_DETECTION=true` for the race detector; CI runs it that way).
- Build the cross-platform binaries with `make build`.
- Regenerate the Reference schema snippets with `make generate-reference-snippets` whenever the embedded schemas change; a sync test fails if the committed snippets drift.
- Preview the documentation with `make dev-documentation`.
- The example domains under `documentation/guides/` are lint-gated; `make qa` runs the linter over them.

## Layout

**Linting pipeline**

- `loader/`: finds and reads `.esdm.yaml` files.
- `parser/`: parses documents and validates them against the schemas.
- `resolver/`: resolves references between documents.
- `model/`: the typed facade over parsed documents.
- `modelpath/`: parses the slash-separated model-path strings the commands accept to scope their output.
- `rules/`: the lint rules, one file per rule.
- `runner/`: wires loading, parsing, resolving, and rules into the pipeline.

**Diagnostics and output**

- `diag/`: diagnostic types and locations.
- `reporter/`: output formatting for diagnostics.
- `hint/`: did-you-mean suggestions.

**Application and infrastructure**

- `cmd/`: entry points (the `esdm` CLI and the `refgen` snippet writer).
- `schema/`: the embedded schema YAML files and the schema host.
- `update/`: the daily update notification check.
- `logging/`: structured logging.
- `refgen/`: extracts the per-kind schema excerpts the Reference pages embed.

**Docs and conventions**

- `documentation/`: the published documentation site, its lint-gated example guides, and its sync tests.
- `.claude/rules/`: coding, collaboration, Git, and language conventions.

For the concept behind a package, read its `documentation.go`.
