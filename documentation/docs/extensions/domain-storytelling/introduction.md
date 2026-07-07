# Domain Storytelling

The **Domain Storytelling** extension captures discovery stories – narrative flows through a domain – as ESDM documents. The technique was developed by Stefan Hofer and Henning Schwentner, and a Domain Story document follows their canon.

A story is a sequence of **Sentences**. Each Sentence is a numbered set of arrows between **Actors** and **Work Objects**, with the arrows labeled by the activities that connect them. Read in order, the Sentences tell what happens in the domain, in the language the domain experts already use.

## When You'd Reach for It

Domain Storytelling is a **discovery format**, not a modeling format. You reach for it when you sit down with subject-matter experts and want to understand how the domain actually works. The story is what comes out of those conversations: who talks to whom, who passes what to whom, in which order.

Capturing the story as an ESDM document means it survives the workshop. The names of Actors and Work Objects stay consistent across sessions, the model notices when a Work Object is referenced without ever being introduced, and the file lives in version control next to the rest of the code. Stories that informed the design stay reachable, and you can come back to them when the design needs to change.

A story precedes and informs the core ESDM model – Aggregates, Events, and Commands are typically derived from stories later. Stories themselves do not depend on the core: a Domain Story document validates in isolation against this extension's schema.

## The Three Classification Dimensions

Stories are classified along three dimensions that match the Domain Storytelling canon:

- **Point in time** – as-is for the current state of the domain as observed, to-be for the target state the story envisions.
- **Granularity** – coarse-grained for high-level stories with few Sentences per major flow, fine-grained for detailed stories with one Sentence per atomic step.
- **Domain purity** – pure for the domain as experienced by humans (independent of software support), digitalized for stories that include the software systems that participate in the domain.

You pick one combination per story. Discovery typically starts with as-is, coarse-grained, pure; the story refines as the conversation deepens.
