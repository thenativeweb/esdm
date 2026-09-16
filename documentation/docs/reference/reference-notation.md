# Reference Notation

A **reference notation** is a way to point at an element of an ESDM model from *outside* it. A product spec that discusses a specific Aggregate, a ticket that traces a bug to one Command, an architecture note that cites an Event – each of them needs to name a model element precisely, and each of them lives in a document that is not part of the model. The reference notation gives them a single, stable string to do it with.

Inside a model, elements already reference each other: a Policy names the Events it reacts to, a Query names its Read Model. Those internal references are logical, not physical – they identify an element by what it *is* in the model, never by which file it happens to sit in. **The reference notation extends that same idea outward**, so a foreign document can hold a reference that keeps resolving even as the model is reorganized on disk.

## Why Not Just Point at a File

The obvious way to reference an element is to point at the file that defines it – `catalog/book.esdm.yaml`, line 40. It is also the wrong way. **A file path binds your reference to a decision that has nothing to do with the domain**: which document a modeler chose to put the element in, and where in that document it landed.

ESDM lets you split a model across as many `.esdm.yaml` files as you like, merge several into one, and move elements between them, all without changing the model itself. A file-based reference breaks on every one of those moves. A reference that names the element logically does not: it keeps resolving as long as the element still exists under the same name at the same place in the domain.

## The Shape of a Reference

A reference is a URI in the `esdm` scheme. Its path walks the model's containment from the Domain inward, one segment per element, and **every segment names the element's kind and its name**, joined by an equals sign:

```text
esdm:<kind>=<name>/<kind>=<name>/.../<kind>=<name>
```

The kinds are the values the `kind` field takes in the model – `domain`, `bounded-context`, `aggregate`, `command`, and so on. The names are written exactly as they appear in the model, lowercase and kebab-case, with no escaping of any kind. Here is one element at each level of a model, addressed:

| Element | Reference |
| --- | --- |
| a Domain | `esdm:domain=library` |
| a Bounded Context | `esdm:domain=library/bounded-context=catalog` |
| a Policy | `esdm:domain=library/policy=notify-on-overdue` |
| an Aggregate | `esdm:domain=library/bounded-context=catalog/aggregate=book` |
| a Command | `esdm:domain=library/bounded-context=catalog/aggregate=book/command=register` |
| an Event | `esdm:domain=library/bounded-context=catalog/aggregate=book/event=registered` |

The Command is `register` and the Event is `registered`, not `register-book` and `book-registered`: names in the model are bare, and the Aggregate is already in the path. The **[Naming](/concepts/event.md#naming)** section of the Event concept explains that convention.

The path follows the model's containment, so its length says where an element sits. A free-standing Event – one published by a Command on a Dynamic Consistency Boundary rather than by an Aggregate – has no Aggregate to carry, so it is one segment shorter: `esdm:domain=library/bounded-context=lending/event=loan-extended`. A **[Context Mapping](/concepts/context-mapping.md)** has no Domain at all, because its endpoints may straddle Domains, so its reference is a single segment: `esdm:context-mapping=catalog-to-lending`.

Writing a parameter into a path segment as `key=value` is the convention the URI specification itself describes for segment parameters, so every URI parser accepts the form, and splitting at `/` and `=` recovers the pairs.

## Why Every Segment Carries Its Kind

A reference has to work for a reader as much as for a tool, and a reader who sees `book/register` does not know whether `register` is a Command, an Event, or a Feature. **The kind is part of what the reference says**, so it is written down – for the element you point at, and for every container above it.

The containers need it as much as the target does. A Command sits inside an Aggregate or inside a Dynamic Consistency Boundary; naming the container's kind tells the reader which. And ESDM allows the same name to be used by different kinds at the same place, where that is the idiom rather than a mistake: an Aggregate and a Read Model both called `order` are one concept on the write and the read side, an Entity and an Actor both called `applicant` are the data and the role, a Subdomain and a Bounded Context both called `billing` are the classification and the context. A reference that carried only names would be ambiguous for every one of them. With the kind in each segment, `aggregate=order` and `read-model=order` are two references, as they should be. The one group of kinds that may *not* share a name – the domain types of a Bounded Context, which become the types of one code module – is described on the **[Bounded Context](/concepts/bounded-context.md#one-namespace-for-the-domain-types)** page and enforced by the linter.

There is a pleasant consequence. Every ESDM document already states its kind, its name, and its `scope` – the Domain, the Bounded Context, the Aggregate it belongs to. **A reference is that document head, serialized**: the scope fields in order, each written as the kind it names, then the element's own kind and name. You can form a reference from any document without looking anything up, and a tool can check one against the model by walking exactly those fields.

## A Name, Not a Location

**A reference is a name, not a location.** It does not say where the model's documentation is hosted, or on which server it might be browsed; it says which element it means, and nothing more. This is why the scheme is `esdm:` and not `esdm://`: the two slashes in `http://` introduce an *authority* – a host on a network – and an ESDM reference has no authority to name. It belongs with the schemes that name things, `urn:` and `tag:` and `mailto:`, not with the schemes that locate them.

Keeping the host out of the reference is deliberate. **A reference outlives any one place the model is published** – a specification written today should still point at the right element after the documentation moves to a new domain, or is generated fresh into a different repository. The host is not part of the element's identity, so it is not part of the reference.

That said, a reference is built so that a rendering of the model can turn it into a location. Its segments are the containment path, and a documentation tree rendered from the model – one page per element along that path – maps every reference to exactly one page by a fixed rule. Supply the base URL of such a rendering, and a tool resolves the reference to a link. The **[esdm documentation](https://github.com/thenativeweb/esdm/issues/5)** command that produces such a tree is in the works.

## What You Can Point At

Every named kind is addressable, at the level where it lives:

- **At the top of the model** – a **[Domain](/concepts/domain.md)**, and a **[Context Mapping](/concepts/context-mapping.md)**, which stands alone because its endpoints may straddle Domains.
- **Within a Domain** – a **[Subdomain](/concepts/subdomain.md)**, a **[Bounded Context](/concepts/bounded-context.md)**, a **[Policy](/concepts/policy.md)**, an **[Event Handler](/concepts/event-handler.md)**, a **[Process Manager](/concepts/process-manager.md)**, and an **[External System](/concepts/external-system.md)**.
- **Within a Bounded Context** – an **[Aggregate](/concepts/aggregate.md)**, a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**, a free-standing **[Event](/concepts/event.md)**, a **[Read Model](/concepts/read-model.md)**, a **[Query](/concepts/query.md)**, an **[Entity](/concepts/entity.md)**, a **[Value Object](/concepts/value-object.md)**, a **[Domain Service](/concepts/domain-service.md)**, and an **[Actor](/concepts/actor.md)**.
- **Within an Aggregate or a Dynamic Consistency Boundary** – a **[Command](/concepts/command.md)** and an **[Event](/concepts/event.md)**.

The **[Extensions](/extensions/overview.md)** add two more addressable kinds. A **[Domain Story](/extensions/domain-storytelling/concepts/overview.md)** sits at Domain level, like a Process Manager: `esdm:domain=library/domain-story=first-loan`. A **[Feature](/extensions/given-when-then/concepts/feature.md)** sits under the unit it specifies, so its path continues that unit's – `esdm:domain=library/bounded-context=catalog/aggregate=book/feature=registering-a-book` for a Feature about an Aggregate, `esdm:domain=library/process-manager=overdue-escalation/feature=escalating-an-overdue-loan` for one about a Process Manager.

References stop at the element. They do not reach into its schema fields, an Aggregate's invariants, a term of a Bounded Context's ubiquitous language, or a single scenario inside a Feature. **The unit you point at is a modeling element, not a line inside one.**

## References and Renames

A reference is stable against everything physical – which file an element lives in, how the files are split or merged, where they sit in the repository. It is *not* stable against renaming the element or any of its containers. Rename the `book` Aggregate to `title`, and every reference that carried `aggregate=book` goes stale.

This is deliberate, and it matches how references behave *inside* a model, where renaming an element breaks every internal reference to it until those are updated too. **A rename is a refactoring, and a refactoring updates its references** – the ones in other model files and the ones in the specs, tickets, and notes that point in from outside. The notation makes those outside references easy to find, because they all share the `esdm:` prefix and carry the element's kind and name.

The **[Concepts overview](/concepts/overview.md)** lists every kind a reference can name. The paths you pass to **[esdm view](/getting-started/running-esdm-view.md)** are the same containment walk in a relative, names-only form; where a name is shared by several kinds at one position, `esdm view` renders every match, and a reference names one.
