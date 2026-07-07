# Given-When-Then Reference

This chapter is the canonical reference for the Given-When-Then extension. The extension defines a single document kind, `feature`, plus a nested entry shape, `scenario`, that lives inside Features.

## Schema Identity

The machine-readable schema is identified by `https://schema.esdm.io/given-when-then/v1`. The `apiVersion` field of every given-when-then document carries the same identifier without the `https://` scheme – `apiVersion: schema.esdm.io/given-when-then/v1`. Hosting under `schema.esdm.io` is in preparation; the URL is currently an identifier, not a fetchable endpoint. The extension's version moves independently of the core schema, mirroring how Kubernetes API groups version independently of `core/v1`.

## Kinds

- **[feature](/extensions/given-when-then/reference/feature.md)** – Top-level Given-When-Then document. A Feature carries one or more Scenarios about one consistency unit.
- **[scenario](/extensions/given-when-then/reference/scenario.md)** – Nested entry inside `feature.scenarios[]`. Not a top-level document on its own; it lives inside its Feature and inherits the Feature's variant.

## Common Fields

Every given-when-then top-level document carries the same shape as a core document: `apiVersion`, `kind`, `name`, an optional `description`, and an optional `metadata` block holding non-semantic `labels` and `annotations`. The required set is `apiVersion`, `kind`, and `name`; unknown top-level keys are rejected via `unevaluatedProperties: false`. The `feature` reference page restates these in its anatomy section. Scenarios are nested entries and carry only `name`, `description`, `given`, `when`, and `then`; the document-level common fields live on the surrounding Feature.

## Variants

A Feature targets exactly one consistency unit, and the chosen target fixes the shape of every Scenario inside the document. Four variants are admitted by the schema, discriminated structurally on which sibling field is present alongside `domain` in `scope`. The Aggregate variant carries `aggregate`. The DCB variant carries `dynamicConsistencyBoundary`. The Process Manager variant carries `processManager` (Process Managers are domain-scoped, so no Bounded Context is involved). The Read Model variant carries `readModel`. Mixing variants inside one Feature is forbidden by the schema's per-variant `if`/`then` rules.

## File Convention

The file suffix and document separator are the same as for the core schema: documents end in `.esdm.yaml` and multi-document files use `---` as the separator. Given-When-Then documents typically live next to the consistency unit they describe rather than in a separate `tests/` directory, so the Feature stays close to the model it talks about.
