# Domain

A **Domain** is the top-level area of activity that the model describes. It is the answer to the question *what business are we in?* – not a single feature, not a product, but the whole field within which the system operates.

A model has exactly one Domain. It sits above **[Subdomains](/concepts/subdomain.md)** and everything below them, and serves as the natural boundary: when you ask "is this in scope?", the answer is "yes if it falls under this Domain".

## When to Use the Term

Use *Domain* when you want to talk about the whole. A logistics company has a logistics Domain. A bank has a banking Domain. A municipal library has a library Domain. The Domain is the umbrella under which all the other terms make sense.

Resist the temptation to split a Domain into more than one. If you find yourself wanting to, what you're really discovering is that you have **multiple Subdomains**, not multiple Domains.
