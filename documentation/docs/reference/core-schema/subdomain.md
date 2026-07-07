# subdomain

Strategic classification of a portion of the Domain. See **[Concepts: Subdomain](/concepts/subdomain.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/subdomain.yaml"
```

## Anatomy

A Subdomain is anchored in its Domain through `scope`, which carries the single field `domain` (the Domain name). That positioning is not optional: a Subdomain that doesn't know which Domain it belongs to is meaningless.

The `type` field classifies the Subdomain and is restricted to one of three values. `core` marks the strategically differentiating part of the business; `supporting` marks parts that are necessary but not differentiating; `generic` marks parts that any business in the space would have. The schema enforces the closed set so the typology cannot drift.

The `boundedContexts` array binds the Subdomain to one or more Bounded Contexts by name. The list is required and non-empty: a Subdomain that has no Bounded Contexts under it has no presence in the model. Each entry follows the kebab-case `name` pattern.

`description` and `metadata` round out the document with the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `subdomain`, and `name` is the Subdomain's own kebab-case identifier.
