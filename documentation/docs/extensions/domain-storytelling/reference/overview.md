# Domain Storytelling Reference

This chapter is the canonical reference for the Domain Storytelling extension. The extension defines a single document kind, `domain-story`, that captures a numbered sequence of Sentences describing one narrative flow through the domain.

## Schema Identity

The machine-readable schema is identified by `https://schema.esdm.io/domain-storytelling/v1`. The `apiVersion` field of every domain-storytelling document carries the same identifier without the `https://` scheme – `apiVersion: schema.esdm.io/domain-storytelling/v1`. Hosting under `schema.esdm.io` is in preparation; the URL is currently an identifier, not a fetchable endpoint. The extension's version moves independently of the core schema.

## Kinds

- **[domain-story](/extensions/domain-storytelling/reference/domain-story.md)** – Top-level Domain Storytelling document. A story is a numbered sequence of Sentences that describe one narrative flow through the domain.

## Common Fields

A Domain Storytelling document carries the same top-level shape as a core document: `apiVersion`, `kind`, `name`, an optional `description`, and an optional `metadata` block holding non-semantic `labels` and `annotations`. The required set is `apiVersion`, `kind`, and `name`; unknown top-level keys are rejected via `unevaluatedProperties: false`. The reference page restates these in its anatomy section.

## Annotations Versus `metadata.annotations`

A Domain Story carries an in-diagram `annotation` field on Actors, Work Objects, and edges. That field is distinct from `metadata.annotations`: `annotation` (singular) carries text that is part of the rendered diagram, while `metadata.annotations` (plural, on the surrounding document) is for tooling and provenance. The two never overlap in purpose, and the linter does not treat them as equivalent.

## File Convention

The file suffix and document separator are the same as for the core schema: documents end in `.esdm.yaml` and multi-document files use `---` as the separator. A project that uses Domain Storytelling alongside its core model typically keeps stories in a sibling directory next to the Bounded Contexts they originated in, but the schema imposes no specific layout.
