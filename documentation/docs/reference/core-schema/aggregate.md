# aggregate

Container-based consistency unit inside a Bounded Context. See **[Concepts: Aggregate](/concepts/aggregate.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/aggregate.yaml"
```

## Anatomy

An Aggregate scopes to a Bounded Context. The `scope` object carries `domain` and `boundedContext`, both required, and no other fields are allowed. That positioning is what gives the Aggregate its consistency boundary: every Command, Event, and invariant attached to it lives inside that single Bounded Context.

The `identifiedBy` field is a discriminated `oneOf` on `source`, and the chosen value selects the sibling field that completes the strategy. When `source` is `state`, the entry pairs with a required `field`, naming a property of `state` that holds the identifier. When `source` is `static`, the entry pairs with `value`, a non-empty string that fixes the identifier for all instances – useful for singleton Aggregates. When `source` is `generated`, the entry pairs with `generator`, the kebab-case name of a generation strategy such as `uuid`, `ulid`, `snowflake`, `cuid`, or `nanoid`. The schema does not validate that the named generator exists; that's a downstream concern.

The `state` field is a JSON Schema object that describes the Aggregate's per-instance state. Its precise shape is up to the model; the schema only requires that it is an object, leaving the structure entirely to the team.

The optional `invariants` array holds named rules that must always hold over `state`. Each entry is `{ name, rule }`, where `name` follows the kebab-case `name` pattern and `rule` is the prose statement of the constraint. Invariants are descriptive, not executable: they tell the reader what the Aggregate guarantees, not how it enforces it.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `aggregate`, and `name` is the Aggregate's kebab-case identifier.
