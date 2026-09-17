## Structure

### `esdm/structure/story-duplicate-actor-name` { #structure-story-duplicate-actor-name }

Severity: `error`

Within the `actors` of a Domain Story, every name appears at most once. Domain Storytelling draws one icon per Actor in a story; a second declaration with the same name confuses the diagram and every reference to it.

### `esdm/structure/story-duplicate-group-name` { #structure-story-duplicate-group-name }

Severity: `error`

Within the `groups` of a Domain Story, every name appears at most once. Memberships refer to Groups by name, so a duplicate makes them ambiguous.

### `esdm/structure/story-duplicate-work-object-in-sentence` { #structure-story-duplicate-work-object-in-sentence }

Severity: `error`

Within the `workObjects` of one Sentence, every name appears at most once. Domain Storytelling redraws Work Objects per Sentence, so the same name in another Sentence is a fresh instance and fine; a second declaration inside the same Sentence is a mistake.

### `esdm/structure/story-sentence-without-edges` { #structure-story-sentence-without-edges }

Severity: `warning`

Every Sentence of a Domain Story must draw at least one edge; a Sentence without edges tells nothing. The schema already requires this; the rule keeps the requirement in place independently of the schema.

### `esdm/structure/story-unknown-group-membership` { #structure-story-unknown-group-membership }

Severity: `error`

Every Group an Actor, Work Object, or edge claims membership in via `groups` must be declared in the `groups` of the story. Membership in an undeclared Group is a stale reference.

### `esdm/structure/story-without-sentences` { #structure-story-without-sentences }

Severity: `warning`

Every Domain Story must have at least one Sentence; a story without Sentences tells nothing. The schema already requires this; the rule keeps the requirement in place independently of the schema.

## Modeling

### `esdm/modeling/story-orphan-actor` { #modeling-story-orphan-actor }

Severity: `warning`

An Actor declared in the `actors` of a Domain Story should be drawn by at least one edge of the story; an Actor nothing draws is decorative. Actors that appear only in edges, without a declaration, are intentional and not flagged.

### `esdm/modeling/story-orphan-group` { #modeling-story-orphan-group }

Severity: `warning`

Every Group declared in the `groups` of a Domain Story should have at least one member: an Actor, a Work Object, or an edge. An empty Group is a frame around nothing.

### `esdm/modeling/story-orphan-work-object` { #modeling-story-orphan-work-object }

Severity: `warning`

A Work Object declared in the `workObjects` of a Sentence should be drawn by at least one edge of that Sentence; a Work Object nothing draws is decorative. Work Objects that appear only in edges, without a declaration, are intentional and not flagged.
