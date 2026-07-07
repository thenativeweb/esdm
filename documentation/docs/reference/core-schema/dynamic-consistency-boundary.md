# dynamic-consistency-boundary

Selector-based consistency unit whose scope is determined per Command. See **[Concepts: Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/dynamic-consistency-boundary.yaml"
```

## Anatomy

A Dynamic Consistency Boundary scopes to a Bounded Context. The `scope` object carries `domain` and `boundedContext`, and no other fields are allowed. The DCB is the alternative to an Aggregate when the consistency unit cannot be expressed as a fixed container.

The `identifiedBy` array holds one or more identifier components. Each entry has a required `name` and a discriminated `source`. When `source` is `command-payload`, the entry pairs with `field`, naming a field on the triggering Command's `data` whose value provides the identifier component. When `source` is `static`, the entry pairs with `value`, a non-empty fixed string. When `source` is `generated`, the entry pairs with `generator`, the kebab-case name of a generation strategy such as `uuid` or `ulid`. The `state` source available on `aggregate.identifiedBy` is intentionally absent here – DCBs have no container-level state to draw from.

The `consults` array is the heart of a DCB. Each entry references an Event the DCB takes into account when forming its decision state, plus a prose `criteria` describing which occurrences are relevant. The Event reference is a discriminated `oneOf`: an Aggregate-bound entry carries `boundedContext`, `aggregate`, `event`, and `criteria`, while a free-standing entry omits `aggregate`. Both shapes require all four (or three) fields explicitly; partial entries are rejected.

The optional `invariants` array holds named rules over the consulted state. Each entry is `{ name, rule }`, where the rule is prose and the name follows the kebab-case `name` pattern. The DCB itself carries no `state` object; its decision surface is the union of consulted Events filtered by the criteria.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `dynamic-consistency-boundary`, and `name` is the DCB's kebab-case identifier.
