# feature

Top-level Given-When-Then document. A Feature carries one or more Scenarios about one consistency unit.

## Schema

```yaml
--8<-- "reference/extensions/given-when-then/feature.yaml"
```

## Anatomy

A Feature targets exactly one consistency unit, and the chosen target fixes the shape of every Scenario inside the document. The `scope` field is a structural `oneOf` over four variants. The Aggregate variant carries `domain`, `boundedContext`, and `aggregate`. The DCB variant carries `domain`, `boundedContext`, and `dynamicConsistencyBoundary`. The Process Manager variant carries `domain` and `processManager` – Process Managers are domain-scoped, so no Bounded Context is involved. The Read Model variant carries `domain`, `boundedContext`, and `readModel`. The presence of `aggregate`, `dynamicConsistencyBoundary`, `processManager`, or `readModel` is itself the discriminator. Mixing variants inside one Feature is forbidden by the schema's per-variant `if`/`then` rules.

The `scenarios` array is required and non-empty. Each entry carries `name`, `given`, `when`, and `then`, plus an optional `description`. The shape of `given`, `when`, and `then` depends on the Feature's variant; the per-scenario fields are documented on **[scenario](/extensions/given-when-then/reference/scenario.md)**.

The common fields round out the document: `apiVersion` is `schema.esdm.io/given-when-then/v1`, `kind` is `feature`, `name` is the Feature's kebab-case identifier, `description` carries free-form prose, and `metadata` holds non-semantic `labels` and `annotations`.
