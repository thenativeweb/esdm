# Overview

The ESDM core schema defines a fixed vocabulary for describing event-sourced domains. Each section in this chapter introduces one of those terms, says what it is in one or two sentences, and then explains it in enough detail that you can use it in conversation and recognize it when it appears in someone else's model.

This chapter doubles as a glossary. If you've landed here from a search, look up the term you came for; the first paragraph of every page is the short definition.

## How the Vocabulary Fits Together

A model starts with a **[Domain](/concepts/domain.md)**, which is classified by **[Subdomains](/concepts/subdomain.md)** and contains **[Bounded Contexts](/concepts/bounded-context.md)**. Inside a Bounded Context live the **consistency units** that own behavior: **[Aggregates](/concepts/aggregate.md)**, **[Dynamic Consistency Boundaries](/concepts/dynamic-consistency-boundary.md)**, **[Process Managers](/concepts/process-manager.md)**, and **[Read Models](/concepts/read-model.md)**.

Consistency units exchange messages. **[Commands](/concepts/command.md)** express intent and target a single consistency unit. **[Events](/concepts/event.md)** record what happened and are published by exactly one consistency unit. **[Queries](/concepts/query.md)** read from a Read Model.

The remaining vocabulary fills in around those primitives. **[Entities](/concepts/entity.md)** are identity-bearing modeling elements that carry no consistency container of their own – the right shape for things a Dynamic Consistency Boundary refers to. **[Value Objects](/concepts/value-object.md)** are the typed building blocks that appear inside Command and Event payloads. **[Policies](/concepts/policy.md)**, **[Event Handlers](/concepts/event-handler.md)**, and **[Domain Services](/concepts/domain-service.md)** are different shapes of behavior that react to Events or sit beside Aggregates. **[Actors](/concepts/actor.md)** and **[External Systems](/concepts/external-system.md)** describe who or what initiates Commands and produces Events from outside the domain. **[Context Mappings](/concepts/context-mapping.md)** describe the relationship between Bounded Contexts.

## Where to Go Next

If you're new to Event Sourcing or Domain-Driven Design, read this chapter linearly. If you're looking for a specific term, navigate to it directly from the sidebar. To put the vocabulary into practice, the **[Modeling Guides](/modeling-guides/overview.md)** walk through example domains end to end. To see how this vocabulary extends with behavioral specifications and discovery stories, head to the **[Extensions](/extensions/overview.md)**. To look up the schema fields each kind has, the **[Reference](/reference/core-schema/overview.md)** is the canonical place. And to point at any of these elements from a document outside the model, the **[Reference Notation](/reference/reference-notation.md)** gives you a stable way to name them.
