# Actor

An **Actor** is someone (or something) that does things in the domain. A customer, a sales clerk, a scheduler.

Actors live at **story scope**: one Actor icon appears once in the entire story, no matter how many Sentences reference it. That story-wide identity is what makes Actors readable across long flows – the customer in Sentence 1 is the same customer in Sentence 7.

## Why Actor Scope Is the Whole Story

The single-icon-per-story rule is a Domain Storytelling convention, not an arbitrary technical constraint. A story tells one flow; the cast that performs the flow stays the same. If a Sentence introduces a new Actor, that's a new icon in the legend, but the previously introduced Actors keep their identity.

This is one of the splits between Actors and **[Work Objects](/extensions/domain-storytelling/concepts/work-object.md)**: Work Objects exist per sentence, Actors exist per story.
