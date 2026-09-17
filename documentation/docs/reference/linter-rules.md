# Linter Rules

`esdm lint` checks a model in two stages. First, the parser and the resolver make sure the documents are well-formed: valid YAML, a known `apiVersion`, the shape the schema prescribes, and references that point at declared elements. Only when that stage reports no error do the rules run. A rule looks at the resolved model as a whole and throws when it finds something the schema cannot express: an Event that nothing publishes, a Feature that names a Command the model does not have, two domain types that share a name.

Every finding carries the ID of the rule that threw it, in the form `esdm/<category>/<name>`, and the address of its entry on this page, so you can follow it straight from the terminal or from the JSON output. This page has one entry per rule over the core schema: the severity it throws at, and what it checks and why. The rules of the extensions live next to the extension they belong to, on **[Given-When-Then: Linter Rules](/extensions/given-when-then/reference/linter-rules.md)** and **[Domain Storytelling: Linter Rules](/extensions/domain-storytelling/reference/linter-rules.md)**.

## Severity

A rule throws either an `error` or a `warning`, and the severity is fixed per rule. An error means the model is inconsistent: something refers to what does not exist, or two declarations contradict each other. A run with at least one error exits with a non-zero status. A warning points at a gap the model can technically live with, such as an Aggregate without Commands. Warnings do not affect the exit code unless you pass `--warnings-as-errors`; the **[CLI reference](/reference/cli.md)** has the details.

There is no configuration to disable a rule or to change its severity. Every model is held to the same standard, so a finding means the same thing in every project.

## Categories

The category is the middle segment of the ID. `structure` rules throw when the model contradicts itself, for example when an `identifiedBy` names a field the state does not declare. `modeling` rules throw when the model is consistent but incomplete or misleading, for example when a Read Model has no Query reading from it.

--8<-- "reference/linter-rules/core.md"

## Pipeline Diagnostics

The parser, the resolver, and the runner report their findings in the same shape as the rules, but they are not rules of the catalog. They throw before any rule runs, and an error from the parser or the resolver stops the rules from running at all: a model that does not parse or resolve cannot be checked for its meaning. Their IDs use the `structure` category for problems in the model and the `system` category for problems around it.

### `esdm/structure/yaml-syntax-error` { #structure-yaml-syntax-error }

Severity: `error`

The file is not valid YAML. The message carries the YAML parser's own description of the problem. Parsing stops at the first syntax error, so fix it and lint again to see what follows.

### `esdm/structure/unknown-api-version` { #structure-unknown-api-version }

Severity: `error`

The document has no `apiVersion`, or names one this version of ESDM does not know. The message lists the versions it does know. Without a known `apiVersion` there is no schema to validate against, so the rest of the document is skipped.

### `esdm/structure/missing-required-field` { #structure-missing-required-field }

Severity: `error`

The document leaves out a field the schema requires, for example an Aggregate without `identifiedBy`. The message names the field, and the finding points at the object that should hold it.

### `esdm/structure/unknown-field` { #structure-unknown-field }

Severity: `error`

The document carries a field the schema does not know at that position. Unknown fields are rejected rather than ignored, because a typo in a field name would otherwise silently drop the field's meaning. A common cause is a field that belongs to a different kind, or one that sits a level too high or too low.

### `esdm/structure/type-mismatch` { #structure-type-mismatch }

Severity: `error`

A value has a different type than the schema expects, for example a string where a list is required. The message names the expected and the actual type.

### `esdm/structure/constraint-violation` { #structure-constraint-violation }

Severity: `error`

A value breaks a constraint of the schema that none of the findings above covers: a value outside an enumeration, a name that does not match the kebab-case pattern, a list shorter than its minimum length, a combination of fields the schema forbids. When the value is close to an allowed enumeration entry, the finding suggests it.

### `esdm/structure/duplicate-name` { #structure-duplicate-name }

Severity: `error`

Two documents of the same kind declare the same name in the same scope, for example two Commands named `place` on the same Aggregate. The finding points at the second declaration and notes where the first one is. The same name in different scopes is fine; a name shared between kinds is the concern of **[`esdm/structure/ambiguous-name`](#structure-ambiguous-name)**.

### `esdm/structure/unresolved-reference` { #structure-unresolved-reference }

Severity: `error`

A document refers to an element the model does not declare: a Command publishes an Event that does not exist, a scope names an Aggregate its Bounded Context does not have, a Context Mapping names an unknown endpoint. The message names the missing element, and when a declared name is close, the finding suggests it.

### `esdm/system/rule-panic` { #system-rule-panic }

Severity: `error`

A rule crashed while checking the model. This is a defect in ESDM, not in your model. The message names the rule and the reason; please **[report it](https://github.com/thenativeweb/esdm/issues)** together with the model that triggered it. The other rules keep running, so the rest of the output is still valid.

### `esdm/system/schemas-directory-drift` { #system-schemas-directory-drift }

Severity: `error`

The project's local `schemas/` directory, created by `esdm add-schema` for editor support, does not match the schemas embedded in this version of ESDM. The linter stops before reading any model file, because its results would disagree with what your editor validates. Run `esdm update-schema` to refresh the directory, or update ESDM if the directory holds the newer revision.
