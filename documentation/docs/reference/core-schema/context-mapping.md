# context-mapping

Relationship between two Bounded Contexts, or between a Bounded Context and an External System. See **[Concepts: Context Mapping](/concepts/context-mapping.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/context-mapping.yaml"
```

## Anatomy

A Context Mapping carries no `scope` – its endpoints may straddle Domains. The natural identifier of a mapping is its `type` plus its endpoints; `name` is a human-readable handle, not a referenced identifier from elsewhere in the model.

The `type` field selects the mapping pattern from a closed set: `customer-supplier`, `conformist`, `anti-corruption-layer`, `open-host-service`, `published-language`, `shared-kernel`, `partnership`, or `separate-ways`. The chosen value determines which sibling fields are required, because each pattern has its own role vocabulary.

The asymmetric patterns name two distinct sides. `customer-supplier` requires `customer` and `supplier`. `conformist` requires `conformist` and `upstream`. `anti-corruption-layer` requires `downstream` and `upstream`. `open-host-service` requires `host` and `consumer`. `published-language` requires `publisher` and `consumer`. Each role accepts either a Bounded Context reference (`{ domain, boundedContext }`) or an External System reference (`{ domain, externalSystem }`), so a mapping can describe an internal collaboration or a relationship with the outside world equally well.

The symmetric patterns – `shared-kernel`, `partnership`, `separate-ways` – use a single `participants` field that takes exactly two Bounded Contexts. External Systems are not allowed here on purpose: sharing a kernel or forming a partnership with a third-party system does not fit the model. If the relationship needs to acknowledge an outside system, use one of the asymmetric patterns instead.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `context-mapping`, and `name` is the mapping's handle in the kebab-case `name` pattern.
