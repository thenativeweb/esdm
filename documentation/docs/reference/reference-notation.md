# Reference Notation

A **reference notation** is a way to point at an element of an ESDM model from *outside* it. A product spec that discusses a specific Aggregate, a ticket that traces a bug to one Command, an architecture note that cites an Event – each of them needs to name a model element precisely, and each of them lives in a document that is not part of the model. The reference notation gives them a single, stable string to do it with.

Inside a model, elements already reference each other: a Policy names the Events it reacts to, a Query names its Read Model. Those internal references are logical, not physical – they identify an element by what it *is* in the model, never by which file it happens to sit in. **The reference notation extends that same idea outward**, so a foreign document can hold a reference that keeps resolving even as the model is reorganized on disk.

## Why Not Just Point at a File

The obvious way to reference an element is to point at the file that defines it – `catalog/book.esdm.yaml`, line 40. It is also the wrong way. **A file path binds your reference to a decision that has nothing to do with the domain**: which document a modeler chose to put the element in, and where in that document it landed.

ESDM lets you split a model across as many `.esdm.yaml` files as you like, merge several into one, and move elements between them, all without changing the model itself. A file-based reference breaks on every one of those moves. A reference that names the element logically does not: it keeps resolving as long as the element still exists under the same name at the same place in the domain.

## The Shape of a Reference

A reference is a URI in the `esdm` scheme, followed by the path that locates the element inside its Domain:

```text
esdm:<domain>/<segment>/.../<name>
```

The segments are the names of the elements that contain your target, from the Domain inward, ending with the target's own name. Every segment is a bare name, written exactly as it appears in the model – lowercase and kebab-case, matching the rules for the `name` field, with no slugging or escaping of any kind.

Here is one element at each level of a model, addressed:

| Element | Reference |
| --- | --- |
| a Domain | `esdm:library` |
| a Bounded Context | `esdm:library/catalog` |
| a Policy | `esdm:library/notify-on-overdue` |
| an Aggregate | `esdm:library/catalog/book` |
| a Command | `esdm:library/catalog/book/register-book` |
| an Event | `esdm:library/catalog/book/book-registered` |

The path follows the model's containment, so its length says where an element sits. An Aggregate's Event carries the Aggregate in its path; a free-standing Event – one published by a Command on a Dynamic Consistency Boundary rather than by an Aggregate – has no Aggregate to carry, so it is one segment shorter: `esdm:library/lending/loan-extended` against an Aggregate's `esdm:library/catalog/book/book-registered`.

## A Name, Not a Location

**A reference is a name, not a location.** It does not say where the model's documentation is hosted, or on which server it might be browsed; it says which element it means, and nothing more. This is why the scheme is `esdm:` and not `esdm://`: the two slashes in `http://` introduce an *authority* – a host on a network – and an ESDM reference has no authority to name. It belongs with the schemes that name things, `urn:` and `tag:` and `mailto:`, not with the schemes that locate them.

Keeping the host out of the reference is deliberate. **A reference outlives any one place the model is published** – a specification written today should still point at the right element after the documentation moves to a new domain, or is generated fresh into a different repository. The host is not part of the element's identity, so it is not part of the reference.

That said, the reference is built so that a host turns it into a location. The path after `esdm:` is exactly the path an element occupies when the model is rendered as a documentation tree, one page per element along the containment hierarchy. Supply the base URL of such a rendering, and the reference resolves by concatenation:

```text
esdm:library/catalog/book/book-registered
+ https://docs.example.com/
= https://docs.example.com/library/catalog/book/book-registered
```

## One Name per Position

The path carries no kind – it does not say `aggregate` or `event` anywhere – so resolving it means walking the containment hierarchy by name, one segment at a time. For that walk to land on a single element, **a name has to identify exactly one element at each position, across all kinds that can sit there**, not merely within one kind.

This is a rule ESDM places on a well-formed model. A Bounded Context may not hold both an Aggregate and a Dynamic Consistency Boundary named `loan`; a Domain may not hold both a Bounded Context and a Policy named `orders`; an Aggregate may not hold both a Command and an Event named `place-order`. Each of these would make `esdm:library/lending/loan` or `esdm:library/orders` mean two things at once, and the model would be ambiguous with or without the notation.

**Unique names at each position are what make the notation total**: every addressable element has exactly one reference, and every well-formed reference names at most one element. The rule costs nothing in practice – two elements at the same place with the same name are confusing to a reader long before they confuse a tool – and in return the reference needs no kind, no disambiguator, and no escaping.

## What You Can Point At

Every named kind is addressable, at the level where it lives:

- **At the top of the model** – a **[Domain](/concepts/domain.md)**. A **[Context Mapping](/concepts/context-mapping.md)** sits here too, but it has no enclosing Domain – its endpoints may straddle domains – so it is the one kind named through a leading marker rather than a containment path: `esdm:context-mapping/catalog-to-lending`.
- **Within a Domain** – a **[Subdomain](/concepts/subdomain.md)**, a **[Bounded Context](/concepts/bounded-context.md)**, a **[Policy](/concepts/policy.md)**, an **[Event Handler](/concepts/event-handler.md)**, a **[Process Manager](/concepts/process-manager.md)**, and an **[External System](/concepts/external-system.md)**.
- **Within a Bounded Context** – an **[Aggregate](/concepts/aggregate.md)**, a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**, a **[Read Model](/concepts/read-model.md)**, a **[Query](/concepts/query.md)**, an **[Entity](/concepts/entity.md)**, a **[Value Object](/concepts/value-object.md)**, a **[Domain Service](/concepts/domain-service.md)**, and an **[Actor](/concepts/actor.md)**.
- **Within an Aggregate or a Dynamic Consistency Boundary** – a **[Command](/concepts/command.md)** and an **[Event](/concepts/event.md)**.

The **[Extensions](/extensions/overview.md)** add two more addressable kinds. A **[Domain Story](/extensions/domain-storytelling/concepts/overview.md)** sits at Domain level, like a Process Manager: `esdm:library/first-loan`. A **[Feature](/extensions/given-when-then/concepts/feature.md)** attaches to whatever it specifies, so its path mirrors that element's – `esdm:library/catalog/book/registering-a-book` for a Feature about an Aggregate, `esdm:library/overdue-escalation/escalating-an-overdue-loan` for one about a Process Manager.

References stop at the element. They do not reach into its schema fields, an Aggregate's invariants, or a single scenario inside a Feature. **The unit you point at is a modeling element, not a line inside one.**

## References and Renames

A reference is stable against everything physical – which file an element lives in, how the files are split or merged, where they sit in the repository. It is *not* stable against renaming the element or any of its containers. Rename the `book` Aggregate to `title`, and every `esdm:library/catalog/book` that pointed at it goes stale.

This is deliberate, and it matches how references behave *inside* a model, where renaming an element breaks every internal reference to it until those are updated too. **A rename is a refactoring, and a refactoring updates its references** – the ones in other model files and the ones in the specs, tickets, and notes that point in from outside. The notation makes those outside references easy to find, because they all share the `esdm:` prefix and carry the element's name.

The **[Concepts overview](/concepts/overview.md)** lists every kind a reference can name, and the paths you pass to **[esdm view](/getting-started/running-esdm-view.md)** are the same containment path in relative form.
