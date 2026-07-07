# event-handler

Reaction to an Event that produces an externally observable side effect. See **[Concepts: Event Handler](/concepts/event-handler.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/event-handler.yaml"
```

## Anatomy

An Event Handler scopes to a Domain via `scope.domain`. Domain-scope is appropriate because handlers often coordinate with External Systems that sit at the edge of the Domain rather than inside any single Bounded Context.

The required `deliveryGuarantee` field carries one of two values: `at-least-once` or `at-most-once`. The two are not symmetric in their consequences. With `at-least-once`, **the schema enforces that `idempotency` is also present**, because the model has to say who absorbs the duplicates that the delivery guarantee admits. With `at-most-once`, `idempotency` is optional.

The `idempotency` object documents who tolerates duplicate delivery and, optionally, how. `idempotency.owner` selects from a closed set of values: `self` means the handler deduplicates internally, `downstream` means a system the handler calls handles it, `infrastructure` means the delivery layer deduplicates, `none` means no mechanism is in place – a deliberate, dangerous choice – and `not-required` means the handler is naturally idempotent. The optional `idempotency.strategy` carries free-form prose where `owner` alone is not self-explanatory.

The `handles` array lists Event references with at least one entry. Each reference is a discriminated `oneOf`: Aggregate-bound entries carry `boundedContext`, `aggregate`, and `event`, while Bounded-Context-scoped entries carry `boundedContext` and `event`. Both shapes are required to be complete – partial references are rejected.

The `sideEffects` array lists what the handler causes in the world, with at least one entry. Each entry is a discriminated `oneOf` on `type`. When `type` is `external-call`, the entry carries the target External System name in `externalSystem` plus a prose `rule` describing the behavior. When `type` is `other`, the entry carries only `rule` and is meant for any observable effect that is not a third-party invocation – audit logs, metrics, cache invalidation, and the like.

The optional `constraints` array holds named rules under which the handler runs. Each entry is `{ name, rule }`. Constraints sit outside `idempotency` and `sideEffects`; they describe activation conditions, not the side effect itself.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `event-handler`, and `name` is the Event Handler's kebab-case identifier.
