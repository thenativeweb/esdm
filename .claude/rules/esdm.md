# ESDM

Repository-specific rules that complement the shared rule set.

## Architectural Negatives

These are deliberate product decisions, not gaps. Do not "fix" them.

- No code generation for the model facades; `model/schema_sync_test.go` is the safety net that keeps every facade method mapped 1:1 to a schema field.
- No end-user configurability for rules: no toggling, no severity overrides.
- No project-level configuration file (no `.esdm-lint.yaml` or similar).
- No stack traces in diagnostics; diagnostics describe the model, not the linter's internals.
- No `--verbose` flag.

## Linter Rules

- Rules do not set `RuleID` or `Severity`; the runner stamps both from `Meta()`.
- RuleID format: `esdm/<category>/<name>`. Allowed categories: `structure`, `naming`, `modeling`, `linguistic`, `system`, `gwt` (rules over the given-when-then extension).
- Severity is fixed per rule, never per finding.
- A new rule is a new file in `rules/` with a `newXxxRule()` constructor and an entry in `Catalog()`.
- Describe a rule producing a diagnostic as the rule "throwing", not "firing". This applies to test names, assertion messages, comments, docs, and commit messages – anywhere the linter's own mechanics are described. Domain terms in user content (for example policies or timers in the ESDM schema) are unaffected.
- The `rules` package uses established short forms for the kind nouns throughout its identifiers, because they recur so densely that the full words would drown the logic: `eh` (event-handler), `pm` (process-manager), `rm` (read-model), `dcb` (dynamic-consistency-boundary), `bc` (bounded-context), `cm` (context-mapping), `ep` (external system's endpoint), `agg` (aggregate), `cmd` (command), `ev` (event), `q` (query). This is a deliberate, package-wide convention; keep it consistent rather than expanding individual occurrences.

## Schema

- Schemas live at `schema/<name>/v<N>.yaml`. A new extension requires an additional `//go:embed` directive in `schema/embed.go`.
- The schema-facade sync test (`model/schema_sync_test.go`) must stay green.
- ESDM documents end in `.esdm.yaml`; multiple documents per file are separated by `---`.

## Typography

- Go identifiers and comments follow the ASCII-only rule from `general.md`. The one exception is string literals that carry user-facing terminal output: the `view` renderer's box-drawing connectors and severity glyphs, the update notification's symbols, and the occasional glyph inside a diagnostic message may use Unicode, because that output is deliberate visual design. Do not rewrite those glyphs as `\u` escapes; keep them readable in the literal.
- The prose inside the schema YAML files (comments and description strings) follows the writing rules instead: it uses the en-dash (U+2013) as the punctuation dash, never the em-dash (U+2014). The reason: the Reference pages of the documentation embed schema excerpts verbatim via snippets, so schema prose is documentation prose.
- Plain ASCII hyphens stay in compound words, kebab-case identifiers, CLI flags, YAML sequences, and YAML's `---` document separator; never replace those with en-dashes.
