# Actor

An **Actor** is **whoever or whatever issues a [Command](/concepts/command.md)**. It can be a human – a customer, an administrator, a clerk – or it can be an automated piece of the system itself.

Actors are how the model captures **agency**. Every Command in ESDM has at least one issuer, and every issuer is named, scoped, and typed. That makes the question *who can do this?* a model-level question rather than a deployment-time one.

## Human and Non-Human Actors

ESDM distinguishes Actor types. A human Actor represents a real person or a role a real person takes on. A non-human Actor – such as a scheduler, a watchdog, or a system service – represents an automated trigger. The distinction matters because the two come with different concerns: human Actors usually have permissions and interfaces, non-human Actors usually have schedules and runtimes.

A human Actor cannot also be backed by a system component – the model keeps the distinction sharp so a person is never quietly conflated with a process. When a system component issues Commands on behalf of a human, model both: the human as the Actor of intent, the system component as the **[External System](/concepts/external-system.md)** that delivers the Command.

## Actor and Permissions

Actors are also how the model captures **authorization** – but the constraint lives on the Command, not on the Actor. A Command that should only be issued by an administrator carries that constraint by listing only the administrative Actor among its permitted issuers (`command.actors`); the Actor itself holds no permission list. Business-level access policies (roles, scopes, tokens) belong to the runtime; what the model captures is which kinds of Actor are permitted in the first place – a useful baseline that runtime checks then build on.

## An Actor Is Not an Invariant

An Actor answers *who or what issues a Command*; an invariant is a rule about a consistency unit's state that must always hold, checked when a Command is handled. The two are easy to conflate, because both can sound like rules – but **"only an administrator may cancel an order" is authorization, not an invariant**. Authorization is expressed through the Actors a Command permits (`command.actors`), not as an invariant of the unit. An invariant constrains what the unit's state may become; an Actor constrains who may ask.
