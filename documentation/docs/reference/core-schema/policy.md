# policy

Stateless reaction that emits Commands when an Event occurs. See **[Concepts: Policy](/concepts/policy.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/policy.yaml"
```

## Anatomy

A Policy scopes to a Domain via `scope.domain`. Domain-scope reflects the Policy's role as an in-domain reaction that connects Events from one consistency unit to Commands targeting another – the binding crosses Bounded Contexts, so anchoring it inside one would misrepresent its reach.

The required `deliveryGuarantee` field carries `at-least-once` or `at-most-once`. As with Event Handlers, the schema turns `at-least-once` into a constraint: if the Policy is at-least-once, it must also carry an `idempotency` object. With `at-most-once`, `idempotency` is optional.

The `idempotency` object names who tolerates duplicate delivery and, optionally, how. `idempotency.owner` is one of `self`, `downstream`, `infrastructure`, `none`, or `not-required`, with the same meanings as on Event Handlers. The optional `idempotency.strategy` is free-form prose describing the concrete mechanism.

The `handles` array lists Event references with at least one entry. Each reference is a discriminated `oneOf`: Aggregate-bound entries carry `boundedContext`, `aggregate`, and `event`, while Bounded-Context-scoped entries carry `boundedContext` and `event`.

The `emits` array lists Command references with at least one entry. Each reference is also a discriminated `oneOf`: an Aggregate-bound Command reference carries `boundedContext`, `aggregate`, and `command`, while a DCB-bound Command reference carries `boundedContext`, `dynamicConsistencyBoundary`, and `command`. The two shapes mirror the two `scope` shapes a Command itself can carry.

The optional `constraints` array holds named rules under which the Policy runs (or doesn't). Each entry is `{ name, rule }`. Constraints sit outside `handles` and `emits`; they describe activation conditions for the Policy as a whole.

A Policy is **stateless**. For stateful, multi-Event coordination with timers and lifecycle, use a Process Manager.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `policy`, and `name` is the Policy's kebab-case identifier.
