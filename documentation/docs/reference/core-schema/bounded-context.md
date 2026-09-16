# bounded-context

Largest unit inside which a single, consistent vocabulary applies. See **[Concepts: Bounded Context](/concepts/bounded-context.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/bounded-context.yaml"
```

## Anatomy

A Bounded Context scopes to a Domain via `scope.domain`. That's the only structural anchor it needs; everything else inside is about vocabulary.

The optional `ubiquitousLanguage` array fixes the canonical terms used inside this Bounded Context. Each entry carries a required `term`, a required `definition`, and an optional `avoid` list of rejected alternatives. The `avoid` shape is itself a list of `{ term, reason? }` entries. Naming `reason` as optional is deliberate: forcing a justification for every rejection tends to produce filler text rather than insight, so the schema lets a team simply state that a term is wrong without explaining why.

A vocabulary is written in one language, and `language` names it as a BCP 47 tag such as `en` or `de-AT`. The field is required as soon as the Bounded Context declares a `ubiquitousLanguage`, because translations only make sense relative to a declared source language. An entry's optional `translations` list carries the same concept in other languages: each translation has the shape of the main entry – its own `language`, a `term`, a required `definition`, and an optional `avoid` list for translations that were considered and rejected. Within each language there is still exactly one term; a translation is not an alias.

`description` and `metadata` are the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `bounded-context`, and `name` is the Bounded Context's kebab-case identifier.
