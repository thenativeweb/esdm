## Structure

### `esdm/structure/aggregate-identified-by-field` { #structure-aggregate-identified-by-field }

Severity: `error`

When the `identifiedBy` of an Aggregate uses `source: state`, the named `field` must be a property of the `state` of the Aggregate. Otherwise the identifier points at a value that does not exist.

### `esdm/structure/ambiguous-name` { #structure-ambiguous-name }

Severity: `error`

Within a Bounded Context, Aggregates, Dynamic Consistency Boundaries, Entities, Value Objects, and Domain Services share one namespace. They become the types of one code module, with nothing composed into their names, so two of them cannot share a name.

### `esdm/structure/duplicate-translation-language` { #structure-duplicate-translation-language }

Severity: `error`

A term has exactly one translation per language. Two translations into the same language contradict the one-term principle the ubiquitous language rests on.

### `esdm/structure/dynamic-consistency-boundary-identified-by-field` { #structure-dynamic-consistency-boundary-identified-by-field }

Severity: `error`

When an `identifiedBy` entry of a Dynamic Consistency Boundary uses `source: command-payload`, the named `field` must be a property of the `data` of every Command that triggers the boundary. Otherwise a decision cannot be tied to its instance.

### `esdm/structure/entity-identified-by-field` { #structure-entity-identified-by-field }

Severity: `error`

When the `identifiedBy` of an Entity uses `source: schema`, the named `field` must be a property of the `schema` of the Entity. Otherwise the identifier points at a value that does not exist.

### `esdm/structure/process-manager-correlated-by-field` { #structure-process-manager-correlated-by-field }

Severity: `error`

When the `correlatedBy` of a Process Manager uses `source: event-field`, the named `field` must be a property of the `data` of every Event the Process Manager consumes. Otherwise an Event cannot be routed to its process instance.

### `esdm/structure/process-manager-timer-at-field` { #structure-process-manager-timer-at-field }

Severity: `error`

When a timer of a Process Manager uses the absolute `at` shape, the named field must be a property of the `state` of the Process Manager. Otherwise the timer has no point in time to fire at.

### `esdm/structure/translation-into-own-language` { #structure-translation-into-own-language }

Severity: `error`

A translation into the language the Bounded Context itself is written in is a second term in that language, which the ubiquitous language does not allow.

## Modeling

### `esdm/modeling/actor-without-type` { #modeling-actor-without-type }

Severity: `warning`

Every Actor must declare whether it is `human` or `system`; the type decides which other fields make sense on it. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/aggregate-without-commands` { #modeling-aggregate-without-commands }

Severity: `warning`

An Aggregate should expose at least one Command. Without one, nothing can drive it from the outside, so it never changes state and never publishes an Event.

### `esdm/modeling/aggregate-without-events` { #modeling-aggregate-without-events }

Severity: `warning`

An Aggregate should have at least one Event. An Aggregate that records nothing is either unfinished or belongs to a different kind, such as a Read Model or a Value Object.

### `esdm/modeling/aggregate-without-identified-by` { #modeling-aggregate-without-identified-by }

Severity: `warning`

Every Aggregate must declare how its instances are identified, via `identifiedBy`. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/aggregate-without-state` { #modeling-aggregate-without-state }

Severity: `warning`

Every Aggregate must declare a `state` schema. An empty schema (`type: object`) is fine and says explicitly that the Aggregate has no observable state. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/bounded-context-without-consistency-unit` { #modeling-bounded-context-without-consistency-unit }

Severity: `warning`

A Bounded Context should host at least one Aggregate or Dynamic Consistency Boundary. Without a consistency unit it holds no behavior and is a placeholder.

### `esdm/modeling/command-without-data` { #modeling-command-without-data }

Severity: `warning`

Every Command must declare a `data` schema. An empty schema is fine and says explicitly that the Command carries no payload. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/command-without-publishes` { #modeling-command-without-publishes }

Severity: `warning`

Every Command must publish at least one Event. A Command that publishes nothing expresses an intent without any consequence the model can see. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/consumer-at-least-once-without-idempotency` { #modeling-consumer-at-least-once-without-idempotency }

Severity: `warning`

An Event Handler, Policy, or Process Manager with `deliveryGuarantee: at-least-once` must declare an idempotency strategy, because it will see the same Event more than once. The schema already requires this combination; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/cross-bounded-context-reference-without-mapping` { #modeling-cross-bounded-context-reference-without-mapping }

Severity: `warning`

When a consumer reaches into more than one Bounded Context, for example a Read Model projecting Events from two contexts or a domain-scoped Policy handling Events across contexts, every pair of contexts it touches should be linked by a Context Mapping. The mapping is where the relationship between the two contexts is made explicit.

### `esdm/modeling/domain-service-without-functions` { #modeling-domain-service-without-functions }

Severity: `warning`

Every Domain Service must declare at least one function; a Domain Service exists to offer operations. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/domain-without-bounded-context` { #modeling-domain-without-bounded-context }

Severity: `warning`

A Domain should own at least one Bounded Context. Without one, nothing in the model belongs to it and its presence is decorative.

### `esdm/modeling/dynamic-consistency-boundary-without-consults` { #modeling-dynamic-consistency-boundary-without-consults }

Severity: `warning`

Every Dynamic Consistency Boundary must consult at least one Event; the consulted Events are what its decisions are based on. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/dynamic-consistency-boundary-without-identified-by` { #modeling-dynamic-consistency-boundary-without-identified-by }

Severity: `warning`

Every Dynamic Consistency Boundary must declare at least one `identifiedBy` entry that says which Events belong to one decision. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/entity-without-identified-by` { #modeling-entity-without-identified-by }

Severity: `warning`

Every Entity must declare how its instances are identified, via `identifiedBy`; identity is what distinguishes an Entity from a Value Object. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/entity-without-schema` { #modeling-entity-without-schema }

Severity: `warning`

Every Entity must declare a `schema` field describing the shape of one instance. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/event-handler-without-handles` { #modeling-event-handler-without-handles }

Severity: `warning`

Every Event Handler must declare at least one Event it handles. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/event-handler-without-side-effects` { #modeling-event-handler-without-side-effects }

Severity: `warning`

Every Event Handler must declare at least one side effect; reacting to Events with further state changes instead is the job of a Process Manager. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/event-name-with-aggregate-prefix` { #modeling-event-name-with-aggregate-prefix }

Severity: `warning`

An Event bound to an Aggregate should not repeat the name of the Aggregate at the start of its own name. The scope already conveys the Aggregate, and ESDM composes the full name from Aggregate and Event wherever it is needed, so `book-registered` on the Aggregate `book` would read as `BookBookRegistered`. Events scoped to a Bounded Context are exempt, because no enclosing Aggregate provides the context.

### `esdm/modeling/event-without-consumer` { #modeling-event-without-consumer }

Severity: `warning`

An Event should be consumed somewhere: by an Event Handler, a Policy, a Process Manager, a Read Model, or a Dynamic Consistency Boundary. An Event nobody consumes is usually a leftover from an earlier version of the model.

### `esdm/modeling/event-without-data` { #modeling-event-without-data }

Severity: `warning`

Every Event must declare a `data` schema. An empty schema is fine and says explicitly that the Event carries no payload beyond the fact that it happened. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/event-without-publisher` { #modeling-event-without-publisher }

Severity: `warning`

Every Event should be published by at least one Command. An Event without a publisher has no way to come into existence in a running system.

### `esdm/modeling/external-system-without-direction` { #modeling-external-system-without-direction }

Severity: `warning`

Every External System must declare its `direction`: `inbound`, `outbound`, or `bidirectional`. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/human-actor-with-backed-by` { #modeling-human-actor-with-backed-by }

Severity: `warning`

A `human` Actor must not declare `backedBy`. The field names the External Systems that implement the channel of a `system` Actor and has no meaning for a person. The schema already forbids the combination; the rule keeps the restriction in place independently of the schema.

### `esdm/modeling/orphan-actor` { #modeling-orphan-actor }

Severity: `warning`

Every Actor should be named in the `actors` of at least one Command or Query. An Actor nothing references does nothing in the model.

### `esdm/modeling/orphan-context-mapping` { #modeling-orphan-context-mapping }

Severity: `warning`

A Context Mapping between two Bounded Contexts should be backed by an actual reference across that boundary somewhere in the model. Otherwise it documents a relationship nothing in the model uses.

### `esdm/modeling/orphan-external-system` { #modeling-orphan-external-system }

Severity: `warning`

Every External System should be referenced somewhere: from the `backedBy` of an Actor, from an external-call side effect of an Event Handler, or as the endpoint of a Context Mapping. Otherwise it has no role in the model.

### `esdm/modeling/policy-without-emits` { #modeling-policy-without-emits }

Severity: `warning`

Every Policy must emit at least one Command; a Policy exists to turn Events into Commands, and a Policy that only observes is an Event Handler. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/policy-without-handles` { #modeling-policy-without-handles }

Severity: `warning`

Every Policy must declare at least one Event it handles. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/process-manager-without-ends-when` { #modeling-process-manager-without-ends-when }

Severity: `warning`

Every Process Manager must declare at least one termination condition in `endsWhen`; a process that never ends is a design smell the model should not hide. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/process-manager-without-event-reactions` { #modeling-process-manager-without-event-reactions }

Severity: `warning`

A Process Manager whose reactions are all timer-based has no event-driven behavior beyond its starting Event. That is usually a gap: a process that only waits is not coordinating anything.

### `esdm/modeling/process-manager-without-starts-when` { #modeling-process-manager-without-starts-when }

Severity: `warning`

Every Process Manager must declare at least one starting Event in `startsWhen`. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/query-without-read-model` { #modeling-query-without-read-model }

Severity: `warning`

Every Query must name the Read Model it reads from. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/query-without-result` { #modeling-query-without-result }

Severity: `warning`

Every Query must declare a `result` schema. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/read-model-without-projections` { #modeling-read-model-without-projections }

Severity: `warning`

A Read Model should have at least one projection; the projections are what fill it. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/read-model-without-query` { #modeling-read-model-without-query }

Severity: `warning`

A Read Model should have at least one Query reading from it. A Read Model nobody reads has no purpose.

### `esdm/modeling/read-model-without-schema` { #modeling-read-model-without-schema }

Severity: `warning`

Every Read Model must declare a `schema` field describing what it materializes. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/subdomain-without-bounded-context` { #modeling-subdomain-without-bounded-context }

Severity: `warning`

A subdomain should list at least one Bounded Context; a subdomain without one names a part of the Domain that nothing in the model fills. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/modeling/value-object-without-schema` { #modeling-value-object-without-schema }

Severity: `warning`

Every Value Object must declare a `schema` field describing the shape of one instance. The schema already requires the field; the rule keeps the requirement in place independently of the schema.
