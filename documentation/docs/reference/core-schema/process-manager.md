# process-manager

Stateful coordinator that drives a long-running flow. See **[Concepts: Process Manager](/concepts/process-manager.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/process-manager.yaml"
```

## Anatomy

A Process Manager scopes to a Domain via `scope.domain`. Domain-scope is necessary because a Process Manager often coordinates across Bounded Contexts; binding it to a single one would misrepresent its reach.

The required `deliveryGuarantee` field carries `at-least-once` or `at-most-once`, and the schema enforces the same conditional pairing as on Event Handlers and Policies: **an at-least-once Process Manager must also carry an `idempotency` object**. The shape of `idempotency` is the same – `idempotency.owner` from the closed set `self`, `downstream`, `infrastructure`, `none`, `not-required`, plus an optional `strategy` carrying prose.

The `correlatedBy` field is a discriminated `oneOf` on `source`. Today the only variant is `event-field`, which pairs `source` with a required `field`. The named field must exist on every Event in `startsWhen` and in any `reactions[].when` Event – the linter checks this; the schema only enforces the structural shape. The instance key for a running Process Manager is the value found at that field path on each incoming Event.

The `state` field is a required JSON Schema object describing per-instance state. Unlike on an Aggregate, a Process Manager's state is bookkeeping for the flow itself: which Events have arrived, which decisions have been made, which timers are armed.

The `startsWhen` field is a single Event reference – Aggregate-bound (`boundedContext`, `aggregate`, `event`) or Bounded-Context-scoped (`boundedContext`, `event`) – that triggers a new instance. The `endsWhen` array holds the conditions under which an existing instance terminates, with at least one entry. Each entry is `{ name, rule }`, where `rule` is prose.

The `reactions` array binds triggers to per-trigger logic, with at least one entry. Each reaction has a required `when` and `rule`, and may carry any of `emits`, `setTimers`, `cancelTimers`. The `when` field is itself a `oneOf`: it accepts an Event reference (Aggregate-bound or Bounded-Context-scoped) or a timer trigger (`{ timer: <name> }`) where `<name>` is one of the timers declared on the Process Manager. The `rule` is prose. The `emits` array carries Command references, each Aggregate-bound (`boundedContext`, `aggregate`, `command`) or DCB-bound (`boundedContext`, `dynamicConsistencyBoundary`, `command`). The `setTimers` and `cancelTimers` arrays carry timer names.

The optional `timers` array declares timer definitions the Process Manager may arm. Each timer requires a `name` and an optional `description`, plus exactly one of `after` or `at` (a `oneOf`). The `after` shape is `{ value, unit }`, with `value` a positive integer and `unit` one of `seconds`, `minutes`, `hours`, `days`, `weeks`, `months`, `years`. The `at` shape names a state field whose value carries the absolute target timestamp.

The optional `invariants` array holds named rules every instance's state must satisfy; the optional `constraints` array holds entity-wide activation conditions. Both follow the `{ name, rule }` shape with prose rules.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `process-manager`, and `name` is the Process Manager's kebab-case identifier.
