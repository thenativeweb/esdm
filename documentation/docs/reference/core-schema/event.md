# event

Immutable record of a fact that happened. See **[Concepts: Event](/concepts/event.md)**.

## Schema

```yaml
--8<-- "reference/core-schema/event.yaml"
```

## Anatomy

An Event has one of two shapes of `scope`, expressed as a structural `oneOf`. The container-owned variant carries `domain`, `boundedContext`, and `aggregate` – this is the common case, an Event that belongs to an Aggregate. The free-standing variant carries `domain` and `boundedContext` only, and is reserved for Events that have no Aggregate of their own; these are typically published by DCB-bound Commands. The presence of `aggregate` is the discriminator.

The `data` field is a JSON Schema object that describes the immutable payload. As with Commands, the model uses an empty object schema (`data: {}`) to express "no payload" deliberately, rather than allowing a missing field to imply it.

ESDM does not model Event versioning as a first-class field. An evolved Event is a new Event type with its own name (e.g. `invoice-paid`, `invoice-paid-v2`). Migration and upcasting are implementation concerns outside this schema.

`description` and `metadata` carry the usual free-form prose and non-semantic attachments. `apiVersion` is `schema.esdm.io/core/v1`, `kind` is `event`, and `name` is the Event's kebab-case identifier in the past tense.
