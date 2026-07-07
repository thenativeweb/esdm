---
paths:
  - "**/*.go"
---

# Go

## Language

- Stick to idiomatic Go, except where the rules below deliberately deviate.

## Error Handling

- Assign and check errors on separate lines. Do not use an `if` init statement for the error check.
  - Write `result, err := doSomething()` followed by a separate `if err != nil { ... }`.
  - Never write `if result, err := doSomething(); err != nil`.
- Never write an `if` statement on a single line, not even `if err != nil { return err }`.
- Pass errors through unwrapped (`return err`) when you do not handle them.
  - Only wrap with `fmt.Errorf("...: %w", err)` when you add real, meaningful context.
  - Never wrap an error without adding context.
- Fail fast: when you hit an unexpected, unrecoverable state, stop early and loudly instead of continuing with a corrupted state.

## Package Structure

- Do not use an `internal/` or `pkg/` segregation layer. Packages start at the repository root and are organized by feature; nesting subpackages is fine.
- Every package has a `documentation.go` file containing only the package comment and the `package` declaration.
  - The package comment starts with `// Package <name> ...`.
  - For load-bearing packages, let the package comment explain the package's concept and design decisions in prose, not just name its purpose.
  - Keep concept documentation reference-free: no line numbers, no file or function names. Package names are the only allowed reference, because in Go the directory equals the package equals the import path, which is the stable architectural boundary.
  - Local invariants (why one specific spot is the way it is) belong in an inline comment at that spot, not in the package comment.
  - Before changing a package, read its `documentation.go` as the entry point to its concept.
  - When a change touches the concept or a design decision, update the `documentation.go` along with it. This is what keeps the concept documentation living with the code.
  - New load-bearing packages include their concept documentation from the start; it is part of "done", not a later addition.
  - When working on a load-bearing package whose concept documentation is missing or stale, fill it in within the scope of the task (boy scout rule). Trivial packages stay deliberately brief; do not force concept prose where there is no concept to explain.

## Testing

- Use `testify/assert` and `testify/require` for assertions.
  - Use `require` for preconditions that should abort the test immediately.
  - Use `assert` for the actual checks.
- Write tests as `t.Run` subtests.
  - Every assertion lives inside a `t.Run` block with a descriptive name.
  - No bare assertions sit directly in the top-level `Test...` function body.
  - Table-driven tests dispatch each row through `t.Run(c.name, ...)`.
