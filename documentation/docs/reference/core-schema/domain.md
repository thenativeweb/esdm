# domain

Top-level area of activity that the model describes. See **[Concepts: Domain](/concepts/domain.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/domain.yaml"
```

## Anatomy

A Domain document is the common top-level shape with no kind-specific fields added. `apiVersion` is exactly `schema.esdm.io/core/v1`, `kind` is `domain`, and `name` is the Domain's natural identifier in the kebab-case `name` pattern. The required set is `apiVersion`, `kind`, and `name`. `description` and `metadata` are optional in the usual way: free-form prose for the former, non-semantic `labels` and `annotations` for the latter. `unevaluatedProperties: false` keeps unknown top-level keys from sneaking in.

There is no `scope`, no `state`, no behavioral fields. Subdomains and other artifacts scope back to a Domain by name; the Domain itself is a label that other kinds anchor to.
