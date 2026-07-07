# Core Schema

This chapter is the canonical reference for the ESDM core schema. Each kind has its own page – Schema excerpt at the top, anatomy underneath – and the kinds are listed alphabetically so the chapter behaves as a lookup surface rather than a curriculum.

## Schema Identity

The machine-readable schema is identified by `https://schema.esdm.io/core/v1`. The `apiVersion` field of every core document carries the same identifier without the `https://` scheme – `apiVersion: schema.esdm.io/core/v1`. The schema identifier is stable; non-breaking edits move `x-esdm-schema-revision` instead. Hosting under `schema.esdm.io` is in preparation; the URL is currently an identifier, not a fetchable endpoint.

## Kinds

- **[actor](/reference/core-schema/actor.md)** – Initiator of a Command.
- **[aggregate](/reference/core-schema/aggregate.md)** – Container-based consistency unit inside a Bounded Context.
- **[bounded-context](/reference/core-schema/bounded-context.md)** – Largest unit inside which a single, consistent vocabulary applies.
- **[command](/reference/core-schema/command.md)** – Expression of intent to change the model.
- **[context-mapping](/reference/core-schema/context-mapping.md)** – Relationship between two Bounded Contexts, or between a Bounded Context and an External System.
- **[domain](/reference/core-schema/domain.md)** – Top-level area of activity that the model describes.
- **[domain-service](/reference/core-schema/domain-service.md)** – Stateless domain logic that does not belong to any single Aggregate.
- **[dynamic-consistency-boundary](/reference/core-schema/dynamic-consistency-boundary.md)** – Selector-based consistency unit whose scope is determined per Command.
- **[entity](/reference/core-schema/entity.md)** – Identity-bearing modeling element without a consistency container of its own.
- **[event](/reference/core-schema/event.md)** – Immutable record of a fact that happened.
- **[event-handler](/reference/core-schema/event-handler.md)** – Reaction to an Event that produces an externally observable side effect.
- **[external-system](/reference/core-schema/external-system.md)** – System outside the modeled Domain that the Domain talks to.
- **[policy](/reference/core-schema/policy.md)** – Stateless reaction that emits Commands when an Event occurs.
- **[process-manager](/reference/core-schema/process-manager.md)** – Stateful coordinator that drives a long-running flow.
- **[query](/reference/core-schema/query.md)** – Read operation served by a Read Model.
- **[read-model](/reference/core-schema/read-model.md)** – Query-optimized projection of Events.
- **[subdomain](/reference/core-schema/subdomain.md)** – Strategic classification of a portion of the Domain.
- **[value-object](/reference/core-schema/value-object.md)** – Typed, named structural definition without identity.

## Common Fields

Every core document carries the same top-level shape: `apiVersion`, `kind`, `name`, an optional `description`, and an optional `metadata` block holding non-semantic `labels` and `annotations`. The required set is `apiVersion`, `kind`, and `name`; unknown top-level keys are rejected via `unevaluatedProperties: false`. Each kind's reference page restates these fields in its anatomy section, so a reader landing on a single page sees the complete picture without having to bounce back here. The `domain` page additionally shows the full top-level shape in its Schema excerpt, because Domain has no kind-specific fields of its own.

## Scope

Most kinds carry a `scope` field that places the artifact inside the model hierarchy. The shape depends on the kind. Subdomain, Bounded Context, Event Handler, Policy, Process Manager, and External System scope to a Domain via `scope.domain`. Aggregate, Dynamic Consistency Boundary, Read Model, Query, Entity, Value Object, Domain Service, and Actor scope to a Bounded Context via `scope.boundedContext`, which carries `domain` and `boundedContext`. Commands and Events that belong to an Aggregate scope to it via `scope.aggregate`, which carries `domain`, `boundedContext`, and `aggregate`; Commands targeting a Dynamic Consistency Boundary scope to it via `scope.dynamicConsistencyBoundary`. Domain documents carry no `scope` (they are top-level), and Context Mapping documents carry no `scope` either – their endpoints can straddle Domains.

## File Convention

Documents that follow this schema are stored as `.esdm.yaml` files. A single file may contain multiple documents separated by `---` (standard YAML multi-document syntax). The schema exposes both conventions through `x-esdm-file-suffix` and `x-esdm-document-separator`, so generators can pick them up mechanically.

## Project Layout

An ESDM project keeps schemas under a single `schemas/` directory at the project root, one subdirectory per schema, and one file per version inside – `schemas/core/v1.yaml`, `schemas/<extension>/v1.yaml`. The user's own `.esdm.yaml` model files live alongside `schemas/`. The `x-esdm-project-layout` field on the schema captures this convention so generators materialize it without guessing.
