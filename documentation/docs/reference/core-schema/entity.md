# entity

Identity-bearing modeling element without a consistency container of its own. See **[Concepts: Entity](/concepts/entity.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/entity.yaml"
```

## Anatomy

An Entity scopes to a Bounded Context via `scope`, which carries `domain` and `boundedContext`. Like a Value Object, an Entity is a structural definition that lives inside one Bounded Context; cross-context reuse is a separate modeling decision and not expressed here.

The `schema` field is a required JSON Schema object describing the shape of a single instance. The naming mirrors `value-object.schema` rather than `aggregate.state`: an Entity has identity but no projected state of its own, so the schema describes what one instance looks like, not the evolving state of a container.

The `identifiedBy` field is a discriminated `oneOf` on `source`, and the chosen value selects the sibling field that completes the strategy. When `source` is `schema`, the entry pairs with a required `field`, naming a property of `schema` that holds the identifier. When `source` is `static`, the entry pairs with `value`, a non-empty string that fixes the identifier for all occurrences – useful for Entities that exist exactly once in the model. The `state` variant from `aggregate.identifiedBy` is intentionally absent, because an Entity has no container state to draw from. The `generated` variant is absent as well: an Entity describes *what* identifies it, not how a new identifier is minted; minting belongs on the Aggregate or DCB-bound Command that brings the Entity into being.

The optional `invariants` array holds named rules over the Entity's data. Each entry is `{ name, rule }`, where `rule` is prose. Invariants on an Entity express constraints that any instance of this Entity must satisfy, regardless of where it is referenced. They are purely structural and value-level – Entities have no lifecycle, so there are no lifecycle invariants here.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `entity`, and `name` is the Entity's kebab-case identifier.
