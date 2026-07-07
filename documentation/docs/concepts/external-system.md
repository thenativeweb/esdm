# External System

An **External System** is a system **outside the modeled domain** that the domain talks to. A payment provider, a mailer, an identity provider, a third-party search service – anything that produces or consumes data on the boundary of what we model.

External Systems give the model a way to be honest about its edges. The world is bigger than the part you're modeling, and pretending otherwise pushes the integrations into footnotes that drift. By naming External Systems explicitly, the model makes its dependencies visible and amenable to the same checks as everything else.

## Direction

An External System has a **direction**. It can be inbound (the External System sends Commands or Events into the domain), outbound (the domain sends Commands or Events to the External System), or bidirectional (both).

The direction determines which connections the External System can participate in. An inbound system can be the source of an Event that lands in the domain; an outbound system is what an **[Event Handler](/concepts/event-handler.md)** ultimately calls; a bidirectional system does both.

## External System and Context Mappings

External Systems also appear in **[Context Mappings](/concepts/context-mapping.md)** – the relationships between Bounded Contexts and the outside world. A customer-supplier mapping where the supplier is an External System is the explicit way to say *this Bounded Context depends on a third party for some of its data*. Modeling that explicitly turns "we use the payment provider somehow" into a fact the model can check.
