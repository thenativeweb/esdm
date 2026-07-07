# scenario

A Scenario is a nested entry inside a `feature.scenarios[]` list. It is not a top-level document on its own; it lives inside its Feature and inherits the Feature's variant. See **[feature](/extensions/given-when-then/reference/feature.md)** for the surrounding document and the variant choice it fixes.

## Schema

The shape of a Scenario depends on its surrounding Feature's variant. The four variants resolve to the following schemas.

### Aggregate variant

```yaml
--8<-- "reference/extensions/given-when-then/scenario-aggregate.yaml"
```

### DCB variant

```yaml
--8<-- "reference/extensions/given-when-then/scenario-dynamic-consistency-boundary.yaml"
```

### Process Manager variant

```yaml
--8<-- "reference/extensions/given-when-then/scenario-process-manager.yaml"
```

### Read Model variant

```yaml
--8<-- "reference/extensions/given-when-then/scenario-read-model.yaml"
```

## Anatomy

Every Scenario carries four fields. The required `name` is a kebab-case identifier that uniquely names the Scenario within its Feature. The optional `description` carries free-form prose. The required `given`, `when`, and `then` describe the Scenario itself, and their shape depends on which variant the surrounding Feature has chosen.

`given` is always a list and is always required, but the list may be empty. An empty `given` means "no preceding history" – the consistency unit is in its initial state when the Scenario starts. The shape of each `given` entry depends on the variant: bare for Aggregate Features, scoped for the rest.

The Aggregate variant uses bare-name Events in `given`. Each entry carries `event` (the bare Event name) and `data` (the concrete payload value, an empty object for an Event without a payload). The `when` of an Aggregate Scenario carries `command` (the bare Command name), `data` (the concrete payload value), and an optional `actor` (the Actor name, useful for permission-style Scenarios). The `then` of an Aggregate Scenario is one of two shapes: `events`, listing `{ event, data }` entries with bare Event names – the empty list is legal and expresses idempotency – or `rejection`, an expected refusal expressed as either `{ invariant: <name> }` or `{ reason: <prose> }`.

The DCB variant uses scoped Event references in `given`, mirroring the structure of `consults` on the DCB itself. Each entry carries `boundedContext`, `event`, `data`, and (when the Event is Aggregate-bound) `aggregate`. The `when` and `then` shapes match the Aggregate variant, because the DCB still produces Events from a single triggering Command.

The Process Manager variant uses scoped Event references in `given`, with the same shape as the DCB variant. The `when` of a Process Manager Scenario is a `oneOf` over three shapes: an Aggregate-owned Event delivered to the instance (carrying `boundedContext`, `aggregate`, `event`, `data`), a free-standing Event delivered to the instance (carrying `boundedContext`, `event`, `data`), or a tick of a named timer (carrying `timer`). The `then` of a Process Manager Scenario is an object that may carry any of `emits` (Command references with their `data`, Aggregate-bound or DCB-bound), `setTimers` (timer names the reaction arms), `cancelTimers` (timer names the reaction cancels), `state` (the concrete state value of the instance after the reaction), or `ended` (`true` when the reaction is expected to retire the instance; `false` is not meaningful and should be omitted).

The Read Model variant uses scoped Event references in `given` so the projection state can be set up. The `when` of a Read Model Scenario carries `query` (the bare Query name) and `parameters` (the concrete parameter value). The `then` is one of two shapes: `result`, which is the expected query result and is intentionally free-form because result shapes are paradigm-specific, or `readModel`, which is the expected materialized Read Model content.

`description` is the only optional top-level field on a Scenario. The common document-level fields – `apiVersion`, `kind`, `name`, `metadata` – are not present on a Scenario; they live on the surrounding Feature.
