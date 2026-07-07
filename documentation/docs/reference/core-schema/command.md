# command

Expression of intent to change the model. See **[Concepts: Command](/concepts/command.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/command.yaml"
```

## Anatomy

A Command targets exactly one consistency unit. The `scope` field is a structural `oneOf`: an Aggregate-scoped Command carries `domain`, `boundedContext`, and `aggregate`, while a DCB-scoped Command carries `domain`, `boundedContext`, and `dynamicConsistencyBoundary`. The two shapes are mutually exclusive – the structural `oneOf` does not need a discriminator field because the presence of `aggregate` versus `dynamicConsistencyBoundary` is itself the discriminator.

The `data` field is a JSON Schema object that describes the Command's payload. There is no shortcut for "no payload": the model uses an empty object schema (`data: {}`) to express that case deliberately, so a missing `data` field never silently means "no payload".

The `publishes` array lists the bare Event names the Command may publish, with at least one entry. The names refer to Events declared in the same scope; the schema does not link them, but the linter does. Bareness is the rule here: a Command publishes Events into its own consistency unit, so naming the surrounding scope would be redundant.

The optional `actors` array lists the Actor names permitted to issue this Command, with at least one entry when present. Authorization is descriptive in ESDM: the model says who is allowed to issue a Command, but does not say how that permission is enforced.

The optional `constraints` array holds named rules that shape how `publishes` is emitted – ordering, conditional emission, mutual exclusivity. Each entry is `{ name, rule }`, where `rule` is prose. Constraints describe the relationship between the Command and its Events without prescribing implementation.

Commands deliberately carry **no** `deliveryGuarantee` and **no** `idempotency`. Retries and idempotency keys live in the API or transport layer; idempotency of effect typically falls out of the target's invariants.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `command`, and `name` is the Command's kebab-case identifier.
