# Use Cases

ESDM was built for teams who are serious about Event Sourcing and Domain-Driven Design and want their model to stay coherent as the system grows. The use cases below describe the situations in which ESDM most clearly earns its keep.

## Capturing a Model During Discovery

When a team runs **Event Storming**, **Domain Storytelling**, or any similar discovery format, the artifacts at the end of the day are usually photos and sticky notes. ESDM gives you a place to put those discoveries that survives the workshop.

You can write the Events you discovered, the Commands that produce them, the Actors that issue those Commands, and the Bounded Contexts that hold them, and you get **immediate feedback** on whether the model is internally consistent. An Aggregate that owns no Events, a Command that nobody issues, or a Process Manager without a starting condition – ESDM flags those before they harden into assumptions.

For Domain Storytelling specifically, ESDM ships an **[extension](/extensions/domain-storytelling/introduction.md)** that captures stories as their own kind of document and validates them under the same toolchain.

## Documenting an Existing System

You already have code. The model exists, but only in the heads of the people who built it, in scattered comments, and in the names of types and methods. Writing the model down – in ESDM YAML, on disk, next to the code – turns that implicit knowledge into a **first-class artifact** the rest of the team can read.

Start with the consistency units: every Aggregate, Dynamic Consistency Boundary (DCB), Process Manager, and Read Model the system has. List the Events each one publishes, the Commands each one accepts. The exercise alone surfaces gaps: an Event nobody listens to, a Command without an Actor, a Read Model fed by no projection. Once written, the model can grow with the code instead of trailing behind it.

## Keeping the Model and the Code in Sync

The most common failure mode for a domain model is not that the model is wrong, but that the **model and the code drift apart**. The team makes a refactoring decision in a sprint, the diagram on the wiki stays untouched, and six months later nobody trusts either source.

When the model lives in version-controlled YAML next to the code, drift becomes visible. A pull request that renames an Aggregate has to update the model in the same change set, or the linter fails the check. A new Event has to declare its publisher and its consumers, or the linter fails the check. The pressure of code review keeps the model honest.

This is the use case that benefits most from running the linter in CI. The check is fast, the failures are precise, and the cost of fixing a drift inline is much smaller than the cost of recovering a wrong model later.

## Communicating with Stakeholders

A model expressed as YAML is **also a model expressed in domain language**. The names of Aggregates, Events, and Commands are the same names the domain experts use. When you sit down with a product owner or a subject-matter expert, the model is a document you can read together, paragraph by paragraph.

The view command renders an ESDM model as a hierarchical summary – Domain, Subdomains, Bounded Contexts, Consistency Units, the Events and Commands they own. It is the closest thing to a model overview that you can produce on demand and trust to match the source. For more depth, the **[Given-When-Then extension](/extensions/given-when-then/introduction.md)** lets you capture concrete behavioral scenarios alongside the model, in the format domain experts already recognize from BDD.

## Governing a Multi-Team System

In a single-team system, the rules of Event Sourcing tend to be implicit – everyone knows them, and they don't need to be written down. In a system with multiple teams and shared Bounded Contexts, **what one team assumes about another's Events** becomes load-bearing, and assumptions diverge fast.

ESDM makes those contracts explicit. A Context Mapping declares the relationship between two Bounded Contexts. A cross-context Event reference is a fact, not a guess. When two teams agree on a contract and write it down in the model, the toolchain holds them to it.

## Modeling with AI Assistance

Large language models can help with both ends of modeling work: drafting a fresh model from a domain conversation, and extracting an implicit model out of an existing code base. ESDM's YAML is plain enough that LLMs read and write it directly, and the **fixed schema** plus **named vocabulary** give them precisely the constraints they need to produce something coherent.

Point an LLM at the **[Concepts](/concepts/overview.md)** chapter and it has the full ESDM vocabulary in its context. Feed it source code, and it can extract candidate Aggregates, Events, and Commands. Feed it a transcript of a domain expert interview, and it can sketch the first draft of a model. The output is editable YAML that you check into version control and refine like any other model.

## Building Tools on Top of ESDM

ESDM is a **format**, and a format is something tools can be built against. Validators, generators, transformers, IDE plugins, dashboards, AI assistants, code generators that target a specific runtime – all of them can read and write ESDM YAML, and all of them speak the same vocabulary about your domain because the schema is fixed.

If you build CQRS, Event Sourcing, or DDD tooling, ESDM gives you an **interoperable substrate** to build on. Your tool emits ESDM, another tool consumes it, a third tool transforms it. The format is the contract; the toolchain is open.

## Where to Go Next

To start using ESDM in any of the modeling scenarios above, **[install it](/getting-started/installing-esdm.md)** and run through **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** to have a coding agent draft it, or **[Your First Model by Hand](/getting-started/your-first-model.md)** to write the YAML yourself. To understand the parts of the language before you write anything, the **[Concepts](/concepts/overview.md)** chapter walks through every kind. If you're building tooling against the format, head straight to the **[Reference](/reference/cli.md)** and the **[Extensions](/extensions/overview.md)**.
