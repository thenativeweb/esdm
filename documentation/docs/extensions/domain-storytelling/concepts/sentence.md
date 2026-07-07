# Sentence

A **Sentence** is one numbered unit of the story. It carries a sequence number (its place in the story), an optional list of declared **[Work Objects](/extensions/domain-storytelling/concepts/work-object.md)**, and a list of edges.

A Sentence reads like a sentence in natural language: an **[Actor](/extensions/domain-storytelling/concepts/actor.md)** does something to or with a Work Object, perhaps passing it to another Actor. Drawn as arrows on a diagram, a Sentence shows one step in the flow.

## Edges

An **edge** is one drawn arrow inside a Sentence. It has a from-end and a to-end – each pointing at an Actor or a Work Object – and an optional label, the activity that connects them: a verb, a preposition, or a conjunction.

The list of edges in a Sentence is flat, not a tree. The sentence patterns of Domain Storytelling – linear flows, collaboration, conjunctions, tool use, parallel arrows – are all expressible as independent edges without special-case structures.

## Sequence Numbers

Sequence numbers may repeat. Duplicate sequence numbers denote parallel Sentences at the same position in the flow – two things that happen at the same point of the story, without one preceding the other.
