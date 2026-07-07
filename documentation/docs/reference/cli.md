# CLI

The `esdm` binary exposes six subcommands. This page is the canonical place to look up the exact invocation, flags, and behavior of each one; subcommands are listed alphabetically so the chapter behaves as a lookup surface. The **[Getting Started](/getting-started/installing-esdm.md)** chapter walks through the most common invocations in context.

Run `esdm --help` for the in-binary version of this overview, or `esdm <subcommand> --help` for a single subcommand.

## `esdm add-schema`

### Invocation

```shell
esdm add-schema
```

### Anatomy

`esdm add-schema` writes the schemas embedded in the running binary into a `schemas/` directory in the current working directory, materializing the layout that the schema's `x-esdm-project-layout` describes – `schemas/core/v1.yaml`, `schemas/<extension>/v1.yaml`. Use this once at project setup so editors with a YAML language server can offer autocomplete and validation against the schemas the linter actually uses.

The command refuses to run if a `schemas/` directory already exists. Refreshing an existing one is the job of `esdm update-schema`; refusing here keeps the two paths cleanly separated and prevents an accidental overwrite of edited or pinned schema files.

## `esdm glossary`

### Invocation

```shell
esdm glossary [path] [flags]
```

### Anatomy

`esdm glossary` reads the ubiquitous language declared on the Bounded Contexts of an ESDM model and writes it to stdout as a Markdown glossary. Each Bounded Context with a `ubiquitousLanguage` block becomes a section, each term an entry, and every discouraged alternative is called out with a short *Avoid* note. The output is plain Markdown, so redirecting it into a file – `esdm glossary > glossary.md` – is the intended way to keep it around.

The optional `[path]` argument narrows the output to a sub-region of the model. The path follows the model hierarchy – Domain, then Bounded Context – separated by a slash. A single segment selects a Domain and emits the glossary for every Bounded Context inside it; two segments select a single Bounded Context. A bare `esdm glossary` with no path covers the whole model. A segment that names no such Domain or Bounded Context, or that reaches below the Bounded Context level, is rejected as invalid input.

The directory holding the model is selected with `-d` / `--directory`, defaulting to the current working directory. There is no `--color` flag – the output is Markdown meant for files and renderers, not a terminal. Linter findings do not block the glossary; as long as the model resolves, the command emits whatever ubiquitous language it finds. When no Bounded Context in scope declares any, the output is just the `# Glossary` heading.

The exit code is `0` on success, the empty-glossary case included. It is non-zero only when the path argument is invalid or the model cannot be resolved at all – run `esdm lint` to find out why in the latter case.

Typical invocations look like this:

```shell
esdm glossary
esdm glossary <domain>
esdm glossary <domain>/<bounded-context> > glossary.md
```

## `esdm lint`

### Invocation

```shell
esdm lint [flags]
```

### Anatomy

`esdm lint` walks a directory and lints every `.esdm.yaml` file it finds. The directory is selected with `-d` / `--directory`, defaulting to the current working directory. The walk is recursive; the linter does not stop at sub-directory boundaries.

The output format is selected with `--format`. The default is `human`, which produces a readable report grouped by file, with severity-colored headers and source excerpts. The alternative is `json`, which produces a machine-readable stream of one diagnostic per line, suitable for CI pipelines and editor integrations. The `--color` flag controls human-output coloring: `auto` detects the terminal capability, `always` forces colors on, and `never` suppresses them. The flag is ignored when `--format` is `json`.

The optional `--warnings-as-errors` flag escalates warning-severity findings to errors for exit-code purposes only. The output and the formatted findings are unchanged; the flag affects nothing other than whether a warning-only run exits with `0` or with a non-zero status. Default is off.

**The exit code reflects the model's correctness rather than the run's success.** `esdm lint` exits with `0` when the model is clean and with a non-zero status when at least one finding has severity `error`. Findings of severity `warning` are reported but do not by themselves cause a non-zero exit, unless `--warnings-as-errors` is set – then a single warning is enough.

## `esdm update-schema`

### Invocation

```shell
esdm update-schema
```

### Anatomy

`esdm update-schema` refreshes the local `schemas/` directory to match the schemas embedded in the running binary. The schema's `x-esdm-schema-revision` field is the comparison key: the local files are overwritten when the binary's revision is newer.

The command rejects downgrades. If the local schema revision is newer than the binary's – for example, because a project has been updated to a fresh schema while the local `esdm` binary has not – the command fails rather than silently overwriting newer files with older ones. Update the binary first, then re-run the command.

## `esdm version`

### Invocation

```shell
esdm version
```

### Anatomy

`esdm version` prints the binary's release version followed by the git commit it was built from. Released binaries print a **[Semantic Versioning](https://semver.org/)** (SemVer) string (e.g. `1.4.0`); unreleased builds print `(version unavailable)` in place of the version. The commit is always present, so a build from an unreleased tree is still uniquely identifiable.

## `esdm view`

### Invocation

```shell
esdm view [path] [flags]
```

### Anatomy

`esdm view` renders a hierarchical summary of an ESDM model. The optional `[path]` argument filters the rendered tree to a sub-region of the model, with the path following the model hierarchy – Domain, Bounded Context, Consistency Unit – separated by slashes. A bare `esdm view` with no path renders the full model.

The directory holding the model is selected with `-d` / `--directory`, defaulting to the current working directory. The optional `--with-details` flag (default `false`) includes node-level details such as schemas, invariants, and rule prose alongside the skeleton; without it, the output is just the structural tree. The `--color` flag controls coloring with the same `auto` / `always` / `never` semantics as on `esdm lint`.

Typical invocations look like this:

```shell
esdm view
esdm view <domain>/<bounded-context>/<aggregate>
esdm view --with-details
```
