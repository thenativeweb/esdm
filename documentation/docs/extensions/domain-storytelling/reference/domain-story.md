# domain-story

Top-level Domain Storytelling document. A story is a numbered sequence of Sentences that describe one narrative flow through the domain.

## Schema

```yaml
--8<-- "reference/extensions/domain-storytelling/domain-story.yaml"
```

## Anatomy

A Domain Story scopes to a Domain via `scope.domain`. Stories are not bound to a single Bounded Context because a narrative often crosses contexts – an invitation drafted in one and accepted in another, for instance.

Three optional classifiers position the story in the broader landscape. `pointInTime` carries `as-is` or `to-be`, distinguishing a description of current behavior from a target picture. `granularity` carries `coarse-grained` or `fine-grained`, signaling how zoomed-in the narrative is. `domainPurity` carries `pure` or `digitalized`, separating a story about humans-and-paper from a story about how the digital system supports them. All three are non-binding for the linter; they document intent for the reader.

The optional `groups` array declares groups that nodes and edges can belong to. Each group carries a required `name`, plus optional `description` and `annotation`. The `description` is prose; the `annotation` is in-diagram text shown alongside the group when rendered. The list must be non-empty when present, with at least one group.

The optional `actors` array declares story-global Actors. Bare Actors can be referenced inline from edges without a declaration, so this list is required only for Actors that need to carry an `annotation` or `groups` membership – metadata that has no place to live on a bare inline reference. Each entry carries a required `name`, an optional `annotation`, and an optional `groups` list naming the groups the Actor belongs to. The `groups` list, when present, must be non-empty.

The required `sentences` array carries the narrative itself, with at least one entry. Each Sentence has a required `sequenceNumber` (an integer of at least one), a required `edges` array, and an optional `workObjects` array. Duplicate `sequenceNumber` values denote parallel Sentences at the same position – two things happening simultaneously in the narrative.

The per-sentence `workObjects` list declares Work Objects that need an `annotation` or `groups` membership. As with Actors, bare Work Objects can be referenced inline from edges without declaration, so the list is required only when there is something to say beyond the name. Each entry carries a required `name`, an optional `annotation`, and an optional non-empty `groups` list. The same name in a different Sentence is a new instance with its own annotation, because Domain Storytelling redraws Work Objects per sentence.

The per-sentence `edges` list is required and non-empty. Each edge carries a required `from` node reference, a required `to` node reference, an optional `label` (the Activity text – verb, preposition, or conjunction), an optional `annotation`, and an optional non-empty `groups` list. A node reference is `{ actor: <name> }` or `{ workObject: <name> }`; mixing both shapes within one reference is rejected.

The `annotation` field on Actors, Work Objects, and edges is distinct from `metadata.annotations`: it carries in-diagram text, while `metadata.annotations` is for tooling and provenance.

The common document-level fields complete the picture: `apiVersion` is `schema.esdm.io/domain-storytelling/v1`, `kind` is `domain-story`, `name` is the story's kebab-case identifier, `description` carries free-form prose, and `metadata` holds non-semantic `labels` and `annotations`.
