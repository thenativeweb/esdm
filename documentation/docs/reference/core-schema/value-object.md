# value-object

Typed, named structural definition without identity. See **[Concepts: Value Object](/concepts/value-object.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/value-object.yaml"
```

## Anatomy

A Value Object scopes to a Bounded Context via `scope`, which carries `domain` and `boundedContext`. A Value Object is a structural definition shared inside one Bounded Context; reusing the same shape across Bounded Contexts is a separate decision and is not represented here.

The `schema` field is a required JSON Schema object describing the Value Object's shape. The field is named `schema` (a type definition), in contrast to `event.data` and `command.data` (a payload value). That naming distinction is deliberate: a Value Object has no instance identity, so its schema is the type itself, not a payload of a particular occurrence.

The optional `invariants` array holds named rules over the Value Object's data. Each entry is `{ name, rule }`, where `rule` is prose. Invariants on a Value Object express constraints that any value of this type must satisfy, regardless of where it is used.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `value-object`, and `name` is the Value Object's kebab-case identifier.
