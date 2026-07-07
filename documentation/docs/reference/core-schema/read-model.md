# read-model

Query-optimized projection of Events. See **[Concepts: Read Model](/concepts/read-model.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/read-model.yaml"
```

## Anatomy

A Read Model scopes to a Bounded Context. The `scope` object carries `domain` and `boundedContext`, both required, with no other fields permitted.

The `schema` field is a required JSON Schema object describing the entire materialized Read Model. For tabular and document paradigms this is commonly a collection shape – `type: array` with `items`, or a keyed map. The schema describes the result the Read Model exposes to its Queries, not the internal storage representation.

The optional `paradigm` field is a hint naming the paradigm used to describe `schema`. It carries one of `tabular`, `document`, `graph`, `stream`, `key-value`, `search-index`, `time-series`, `column`, or `vector`. The hint is non-binding for the linter; it's documentation for the reader and a flag for tooling that wants to render the Read Model differently per paradigm.

The `projections` array is required and non-empty. Each entry is a flat object combining an Event reference with a prose `rule`. An Aggregate-bound entry carries `boundedContext`, `aggregate`, `event`, and `rule`; a Bounded-Context-scoped entry carries `boundedContext`, `event`, and `rule`. `event` is the bare Event name; the surrounding fields fix the producer.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `read-model`, and `name` is the Read Model's kebab-case identifier.
