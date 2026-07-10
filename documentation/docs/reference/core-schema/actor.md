# actor

Initiator of a Command. See **[Concepts: Actor](/concepts/actor.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/actor.yaml"
```

## Anatomy

An Actor scopes to a Bounded Context via `scope`, which carries `domain` and `boundedContext`. Actors live inside the vocabulary of the Bounded Context that recognizes them.

The required `type` field carries `human` or `system`. The two are mutually exclusive and have different downstream consequences. A `system` Actor may name the External Systems that implement its channel through the optional `backedBy` array, listing one or more External System names. A `human` Actor may not – the schema enforces this through a conditional rule that forbids `backedBy` whenever `type` is `human`. The intent is to keep the model from conflating a person with a process.

The optional `responsibilities` array carries free-form prose statements describing what the Actor does in the domain. The list reads like a glossary entry rather than a permission list, and is descriptive rather than referenced from elsewhere in the model. Permission, in ESDM, is expressed by listing Actor names in `command.actors`, not by a separate ACL on the Actor itself, and not as an invariant of the target unit.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `actor`, and `name` is the Actor's kebab-case identifier.
