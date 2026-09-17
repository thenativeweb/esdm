### `esdm/gwt/feature-references-unknown-aggregate` { #gwt-feature-references-unknown-aggregate }

Severity: `error`

A Feature scoped to an Aggregate must name an Aggregate the model declares. When the scope does not resolve, every Scenario in the Feature describes the behavior of something that does not exist.

### `esdm/gwt/feature-references-unknown-dynamic-consistency-boundary` { #gwt-feature-references-unknown-dynamic-consistency-boundary }

Severity: `error`

A Feature scoped to a Dynamic Consistency Boundary must name one the model declares. When the scope does not resolve, every Scenario in the Feature describes the behavior of something that does not exist.

### `esdm/gwt/feature-references-unknown-process-manager` { #gwt-feature-references-unknown-process-manager }

Severity: `error`

A Feature scoped to a Process Manager must name one the model declares. When the scope does not resolve, every Scenario in the Feature describes the behavior of something that does not exist.

### `esdm/gwt/feature-references-unknown-read-model` { #gwt-feature-references-unknown-read-model }

Severity: `error`

A Feature scoped to a Read Model must name one the model declares. When the scope does not resolve, every Scenario in the Feature describes the behavior of something that does not exist.

### `esdm/gwt/feature-without-scenarios` { #gwt-feature-without-scenarios }

Severity: `warning`

Every Feature must declare at least one Scenario; a Feature without Scenarios specifies nothing. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/gwt/scenario-actor-not-permitted` { #gwt-scenario-actor-not-permitted }

Severity: `error`

When a Scenario names an Actor in its `when`, that Actor must be listed in the `actors` of the Command the Scenario sends. A Scenario that lets an unlisted Actor send the Command describes a path the model does not permit.

### `esdm/gwt/scenario-references-unknown-actor` { #gwt-scenario-references-unknown-actor }

Severity: `error`

An Actor named in the `when` of a Scenario must be declared in the Bounded Context of the Feature. A Scenario cannot exercise an Actor the model does not know.

### `esdm/gwt/scenario-references-unknown-command` { #gwt-scenario-references-unknown-command }

Severity: `error`

Every Command a Scenario names in its `when`, or in the `emits` of a Process Manager Scenario, must be declared in the model. A bare Command name resolves through the scope of the Feature; a scoped reference names its Bounded Context and consistency unit explicitly.

### `esdm/gwt/scenario-references-unknown-event` { #gwt-scenario-references-unknown-event }

Severity: `error`

Every Event a Scenario names in its `given` or in `then.events` must be declared in the model. A bare Event name resolves through the scope of the Feature; a scoped reference names its Bounded Context and, where the Event has one, its Aggregate explicitly.

### `esdm/gwt/scenario-references-unknown-query` { #gwt-scenario-references-unknown-query }

Severity: `error`

Every Query a Read Model Scenario names in its `when` must be declared in the Bounded Context of the Feature.

### `esdm/gwt/scenario-references-unknown-timer` { #gwt-scenario-references-unknown-timer }

Severity: `error`

Every timer a Process Manager Scenario names in its `when`, `then.setTimers`, or `then.cancelTimers` must be declared in the `timers` of the Process Manager the Feature targets.

### `esdm/gwt/scenario-rejection-references-unknown-invariant` { #gwt-scenario-rejection-references-unknown-invariant }

Severity: `error`

When `then.rejection` names an invariant, that invariant must be declared in the `invariants` of the consistency unit the Feature targets. A rejection by an invariant the unit does not have describes a guarantee the model does not make.

### `esdm/gwt/scenario-then-mismatched-feature-scope` { #gwt-scenario-then-mismatched-feature-scope }

Severity: `warning`

The `then` of a Scenario must match the variant of its Feature: Aggregate and Dynamic Consistency Boundary Features expect `events` or `rejection`, Process Manager Features expect `emits`, `setTimers`, `cancelTimers`, `state`, or `ended`, and Read Model Features expect `result` or `readModel`. The schema already binds these shapes; the rule keeps the binding in place independently of the schema.

### `esdm/gwt/scenario-when-mismatched-feature-scope` { #gwt-scenario-when-mismatched-feature-scope }

Severity: `warning`

The `when` of a Scenario must match the variant of its Feature: Aggregate and Dynamic Consistency Boundary Features take a Command, Process Manager Features take an Event or a timer, and Read Model Features take a Query. The schema already binds these shapes; the rule keeps the binding in place independently of the schema.

### `esdm/gwt/scenario-without-then` { #gwt-scenario-without-then }

Severity: `warning`

Every Scenario must declare a `then` outcome; a Scenario without one asserts nothing. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/gwt/scenario-without-when` { #gwt-scenario-without-when }

Severity: `warning`

Every Scenario must declare a `when` trigger; a Scenario without one exercises nothing. The schema already requires the field; the rule keeps the requirement in place independently of the schema.

### `esdm/gwt/uncovered-invariant` { #gwt-uncovered-invariant }

Severity: `warning`

On an Aggregate or Dynamic Consistency Boundary that has a Feature, every named invariant should be exercised by at least one Scenario. An invariant no Scenario covers is a guarantee nobody has written down how to test. Units without a Feature are exempt.
