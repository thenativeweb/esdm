# domain-service

Stateless domain logic that does not belong to any single Aggregate. See **[Concepts: Domain Service](/concepts/domain-service.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/domain-service.yaml"
```

## Anatomy

A Domain Service scopes to a Bounded Context via `scope`, which carries `domain` and `boundedContext`. The service belongs to its Bounded Context's vocabulary even though it does not reside inside an Aggregate.

The `functions` array is required and non-empty. Each function declares a `name`, `arguments`, `returns`, and optionally a `description` and `rules`. The `name` follows the kebab-case `name` pattern. The `arguments` field is a JSON Schema object describing the input the function needs to produce its result; an empty object schema is valid but rare. The `returns` field is a JSON Schema object describing the function's return value, and is always required because a domain function always returns computed information.

The optional `rules` array on a function holds named rules expressing the computation in prose. Each entry is `{ name, rule }`. Rules sit on individual functions rather than on the Domain Service as a whole because the meaningful unit of behavior is the function, not the service.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments at the document level. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `domain-service`, and `name` is the Domain Service's kebab-case identifier.
