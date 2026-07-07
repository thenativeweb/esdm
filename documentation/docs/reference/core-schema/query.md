# query

Read operation served by a Read Model. See **[Concepts: Query](/concepts/query.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/query.yaml"
```

## Anatomy

A Query scopes to a Bounded Context via `scope`, which carries `domain` and `boundedContext`. The Query lives next to the Read Model it serves.

The `readModel` field is required and names the Read Model that backs the Query. The name is bare – the Bounded Context is fixed by the Query's own scope, so naming it again here would be redundant. The linter resolves the name against the Read Models in the same Bounded Context.

The `result` field is a required JSON Schema object describing the Query's response shape. Whether the response is a single object, a paginated list, or a paradigm-specific structure is up to the model; the schema only requires that it is an object.

The optional `parameters` field is a JSON Schema object describing the Query's input parameters. Queries without parameters simply omit the field. The optional `paradigm` field hints at the paradigm used to describe `result`, choosing from `tabular`, `document`, `graph`, `stream`, `key-value`, `search-index`, `time-series`, `column`, or `vector`. As on Read Models, the hint is non-binding for the linter.

The optional `actors` array lists the Actor names permitted to run the Query, with at least one entry when present. Authorization is descriptive: the model states who is allowed to run a Query, but does not enforce it.

The optional `constraints` array holds named rules that shape the Query result – filtering, sorting, limiting, paradigm-specific traversal. Each entry is `{ name, rule }`, where `rule` is prose. A Query has no side effects and produces no Events; constraints describe how the result is computed, not what the Query changes.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `query`, and `name` is the Query's kebab-case identifier.
