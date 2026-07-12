# Reference Notation

A **reference notation** is a way to point at an element of an ESDM model from *outside* it. A product spec that discusses a specific Aggregate, a ticket that traces a bug to one Command, an architecture note that cites an Event – each of them needs to name a model element precisely, and each of them lives in a document that is not part of the model. The reference notation gives them a single, stable string to do it with.

Inside a model, elements already reference each other: a Policy names the Events it reacts to, a Query names its Read Model. Those internal references are logical, not physical – they identify an element by what it *is* in the model, never by which file it happens to sit in. **The reference notation extends that same idea outward**, so a foreign document can hold a reference that keeps resolving even as the model is reorganized on disk.

## Why Not Just Point at a File

The obvious way to reference an element is to point at the file that defines it – `catalog/book.esdm.yaml`, line 40. It is also the wrong way. **A file path binds your reference to a decision that has nothing to do with the domain**: which document a modeler chose to put the element in, and where in that document it landed.

ESDM lets you split a model across as many `.esdm.yaml` files as you like, merge several into one, and move elements between them, all without changing the model itself. A file-based reference breaks on every one of those moves. A reference that names the element logically does not: it keeps resolving as long as the element still exists under the same name at the same place in the domain.

## The Shape of a Reference

A reference is a URI in the `esdm` scheme. It names the *kind* of element you are pointing at, followed by the path that locates it:

```text
esdm:<kind>/<segment>/<segment>/...
```

The first segment after the scheme is the kind – `aggregate`, `event`, `command`, and so on, spelled exactly as the schema spells it. The segments that follow are the names of the elements that contain your target, from the Domain inward, ending with the target's own name. Every segment is a bare name, written exactly as it appears in the model – lowercase and kebab-case, matching the rules for the `name` field, with no slugging or escaping of any kind.

Here is one element at each level of a model, addressed:

| Element | Reference |
| --- | --- |
| a Domain | `esdm:domain/library` |
| a Bounded Context | `esdm:bounded-context/library/catalog` |
| a Policy | `esdm:policy/library/notify-on-overdue` |
| an Aggregate | `esdm:aggregate/library/catalog/book` |
| a Command | `esdm:command/library/catalog/book/register-book` |
| an Event | `esdm:event/library/catalog/book/book-registered` |

**A reference is always absolute**: it starts at the Domain and spells out the full containment path. That is the one way it differs from the paths you type into **[esdm view](/getting-started/running-esdm-view.md)**, which are read relative to the model as a whole and leave the Domain implicit. A foreign document has no such context, so the notation carries everything it needs to resolve on its own.

## Why `esdm:` and Not `esdm://`

The two slashes in `http://` are not decoration. In a URI they introduce an *authority* – a host on a network, the thing that answers when you dereference the address. `esdm://event/...` would, to any correct URI parser, read `event` as a hostname.

**An ESDM reference has no authority, because it is a name, not a location.** It does not point at a server; it points at a place in a domain model. That is exactly the distinction that separates schemes like `urn:`, `tag:`, and `mailto:`, which name things, from schemes like `http://` and `file://`, which locate them. ESDM references name things, so they take the authority-less `esdm:` form.

## Naming the Kind Makes It Unambiguous

You might wonder why the kind is part of the reference at all, when the path already leads to the element. The reason is that **a name in ESDM is unique only within its kind at a given position, not across kinds**. Nothing stops a Bounded Context from holding both an Aggregate and a Dynamic Consistency Boundary named `loan`; nothing stops a Domain from holding both a Bounded Context and a Policy named `orders`. Naming the kind up front says which of them you mean.

Naming the kind also fixes the shape of the rest of the path, because every kind has a fixed place in the model's hierarchy. Once a reference says `event`, its path can only be a Domain, a Bounded Context, and – for an Aggregate's Event – the Aggregate, then the Event itself. A free-standing Event, published by a Command on a Dynamic Consistency Boundary rather than by an Aggregate, has no Aggregate to name, so its reference is one segment shorter: `esdm:event/library/lending/loan-extended`, against an Aggregate's `esdm:event/library/catalog/book/book-registered`. **The number of segments tells the two apart.**

One case the notation deliberately leaves to you: a Command names its parent, and that parent can be either an Aggregate or a Dynamic Consistency Boundary. `esdm:command/library/lending/loan/extend-loan` does not say which `loan` it means. **Keeping an Aggregate and a Dynamic Consistency Boundary in one Bounded Context from sharing a name is part of drawing the boundary cleanly** – if they share it, the reference to their Command is ambiguous, and so, arguably, is the model.

## What You Can Point At

Every named kind in the core vocabulary is addressable, at the level where it lives:

- **At the top of the model** – a **[Domain](/concepts/domain.md)** and a **[Context Mapping](/concepts/context-mapping.md)**. A Context Mapping has no enclosing Domain, because its endpoints may straddle domains, so it is named on its own: `esdm:context-mapping/catalog-to-lending`.
- **Within a Domain** – a **[Subdomain](/concepts/subdomain.md)**, a **[Bounded Context](/concepts/bounded-context.md)**, a **[Policy](/concepts/policy.md)**, an **[Event Handler](/concepts/event-handler.md)**, a **[Process Manager](/concepts/process-manager.md)**, and an **[External System](/concepts/external-system.md)**.
- **Within a Bounded Context** – an **[Aggregate](/concepts/aggregate.md)**, a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**, a **[Read Model](/concepts/read-model.md)**, a **[Query](/concepts/query.md)**, an **[Entity](/concepts/entity.md)**, a **[Value Object](/concepts/value-object.md)**, a **[Domain Service](/concepts/domain-service.md)**, and an **[Actor](/concepts/actor.md)**.
- **Within an Aggregate or a Dynamic Consistency Boundary** – a **[Command](/concepts/command.md)** and an **[Event](/concepts/event.md)**.

The **[Extensions](/extensions/overview.md)** add two more addressable kinds. A **[Domain Story](/extensions/domain-storytelling/concepts/overview.md)** sits at Domain level, like a Process Manager: `esdm:domain-story/library/first-loan`. A **[Feature](/extensions/given-when-then/concepts/feature.md)** attaches to whatever it specifies, so its path mirrors that element's – `esdm:feature/library/catalog/book/registering-a-book` for a Feature about an Aggregate, `esdm:feature/library/overdue-escalation/escalating-an-overdue-loan` for one about a Process Manager.

References stop at the element. They do not reach into its schema fields, an Aggregate's invariants, or a single scenario inside a Feature. **The unit you point at is a modeling element, not a line inside one.**

## References and Renames

A reference is stable against everything physical – which file an element lives in, how the files are split or merged, where they sit in the repository. It is *not* stable against renaming the element or any of its containers. Rename the `book` Aggregate to `title`, and every `esdm:aggregate/library/catalog/book` that pointed at it goes stale.

This is deliberate, and it matches how references behave *inside* a model, where renaming an element breaks every internal reference to it until those are updated too. **A rename is a refactoring, and a refactoring updates its references** – the ones in other model files and the ones in the specs, tickets, and notes that point in from outside. The notation makes those outside references easy to find, because they all share the `esdm:` prefix and carry the element's name.

The paths you pass to **[esdm view](/getting-started/running-esdm-view.md)** are the same idea in relative form, and the **[Concepts overview](/concepts/overview.md)** lists every kind a reference can name.
