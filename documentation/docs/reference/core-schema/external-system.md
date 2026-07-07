# external-system

System outside the modeled Domain that the Domain talks to. See **[Concepts: External System](/concepts/external-system.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/external-system.yaml"
```

## Anatomy

An External System scopes to a Domain via `scope.domain`. Domain-scope is appropriate because an External System sits outside any single Bounded Context – it is something the whole Domain talks to, even when only one Bounded Context happens to use it today.

The required `direction` field carries `inbound`, `outbound`, or `bidirectional`. The value determines which connections the External System can participate in. An `inbound` system can be the source of an Event that lands in the Domain. An `outbound` system is what an Event Handler ultimately calls. A `bidirectional` system does both.

The optional `category` field is a free-form short tag for the kind of External System – `payment`, `mail`, `identity`, `geocoding`, and so on. It is not an enum, because new categories appear all the time and constraining the set would only force teams to misclassify.

The optional `capabilities` array carries free-form descriptive strings naming the operations or interactions the Domain uses – "create a charge", "receive payment-succeeded webhook". The list is documentation, not a callable surface; the linter does not link these strings to anything.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `external-system`, and `name` is the External System's kebab-case identifier.
