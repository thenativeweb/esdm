# Running esdm lint

This page walks through the **lint workflow**: what `esdm lint` does, what the output looks like on a clean and a broken model, how the exit code and severity levels relate, and how to wire the linter into a continuous-integration pipeline.

This page assumes you have a small model on disk. Either **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** or **[Your First Model by Hand](/getting-started/your-first-model.md)** produces a fitting companion; the concrete examples below are taken from the by-hand model, since pinning the artifact and field names down makes the diagnostics easier to follow. The **[CLI reference](/reference/cli.md)** documents the exact flags and arguments; this page focuses on the workflow around them.

## Linting a Clean Model

From inside the directory that holds your `.esdm.yaml` files:

```shell
./esdm lint
```

A clean model produces **no output and exits with status `0`**. Silence is success: if there is nothing wrong, the linter has nothing to say. Pipe the exit code into a script or a CI step and act on it.

```shell
./esdm lint && echo "model is clean"
```

To lint a directory other than the current one, pass `-d` (or `--directory`). The linter walks the directory recursively, so a project with one `.esdm.yaml` per consistency unit is linted in one invocation.

## Reading a Finding

When the linter has something to say, it groups the output around three pieces of information per finding: the **severity**, the **location** (file, line, column), and a one-line **message** that explains what's wrong.

For example, if the `state` field on `book.esdm.yaml` is renamed to `stateOmitted`, the linter reports:

```text
error: missing required field "state"
  at book.esdm.yaml:1:1

error: unknown field "stateOmitted"
  at book.esdm.yaml:11:3
```

Multiple findings on the same file are reported in source order. The location format is `<file>:<line>:<column>`, which most editors and CI log viewers render as a clickable link.

## Severity and Exit Code

The linter classifies every finding as either an **error** or a **warning**.

An **error** is a structural or modeling defect serious enough to break the model: a missing required field, an unresolved reference, a forbidden combination of fields. A run with at least one error exits with a **non-zero status**.

A **warning** is a softer signal: a Read Model with no Query reading from it, an Event with no consumer, a redundant naming pattern. Warnings are reported, but they do not by themselves fail the run – a warning-only model still exits with `0`. If you want warnings to fail your build too, pass `--warnings-as-errors`:

```shell
./esdm lint --warnings-as-errors
```

The flag does not change what is printed; it only escalates the exit code so a warning-only run reports a non-zero status. That is the ESDM equivalent of treating compiler warnings as errors – a deliberate choice for CI pipelines, opt-in everywhere else.

## JSON Output

For tooling, editor integrations, or CI pipelines that want to parse findings rather than read them, run with `--format json`:

```shell
./esdm lint --format json
```

Each finding is emitted as a single JSON object on its own line. The shape carries the same severity, location, and message that the human format prints, plus a stable identifier for the rule that flagged the finding. Stream the output, group it, or post-process it as you prefer.

## Linting in CI

Lint on every push and on every pull request. The linter is fast – sub-second on a model with hundreds of artifacts – so there is no reason to gate it.

A minimal **GitHub Actions** step that downloads the binary and lints the model looks like this:

```yaml
- name: Install ESDM
  run: |
    curl -sSL -o esdm \
      https://esdm.s3.fr-par.scw.cloud/0.14.0/esdm-linux-amd64
    chmod a+x esdm

- name: Lint ESDM model
  run: ./esdm lint --directory ./model --color never --warnings-as-errors
```

`--color never` keeps ANSI escape codes out of CI logs. `--warnings-as-errors` is the recommended default for CI: warnings rarely belong on a clean main branch, and the flag turns the exit code into a single binary signal. The exit code carries the success signal automatically; no extra `if` is needed.

## Where to Go Next

- **[Running esdm view](/getting-started/running-esdm-view.md)** complements the linter with a structural view of the model.
- **[Editor Support](/getting-started/editor-support.md)** sets up your editor so problems surface while you type, not only on a manual `esdm lint` run.
- **[CLI: esdm lint](/reference/cli.md#esdm-lint)** documents every flag and the exact exit-code semantics.
