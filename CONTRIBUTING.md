# Contributing to ESDM

Thank you for your interest in ESDM. Everything here – the `esdm` CLI, the schemas, and the documentation at **[www.esdm.io](https://www.esdm.io)** – is open source under the MIT license. ESDM is maintained by the native web, and we're glad about every contribution: a bug report, a field report from a real model, a fix to the documentation, or code. This document explains how we work together, so that the time you invest ends up in the project.

## Talk to Us Before You Write Code

The one thing we'd ask you to do: **open an issue before you start on code**, and wait until a maintainer has confirmed the direction. If an issue for your topic already exists, comment there instead.

Here is why we insist on this. A pull request represents hours of your time, and turning one down is something we dislike doing – it is one of the least pleasant parts of maintaining a project. When it happens, the reason is rarely the code. It's direction: the change collides with a product decision, the design needs a different shape, or someone else is already working on the same thing. A short conversation in an issue catches nearly all of that before a single line is written, and it gives you a clear target: once we've agreed on the what and the how, the review can focus on the code itself.

Confirmed means that a maintainer has replied in the issue and agreed to both the goal and the approach. If you'd like to take on an existing issue, say so in a comment and wait for that reply as well – it keeps two people from working on the same thing, and it lets us settle open design questions before you invest in them.

A pull request that arrives without an agreed issue still gets a look, but it may be closed without a longer discussion. That is precisely the outcome the issue is meant to prevent.

Agreeing on the how includes everything a user will see. Schema fields and kinds, command and flag names, default values, output layout, and the wording of diagnostics are decided in the issue, not in the review of a pull request. If such a detail comes up while you're implementing, pause and ask in the issue – one more question costs far less than reworking an interface afterwards.

Please also skim the **[design principles](https://www.esdm.io/introduction/design-principles/)** first. Some things are absent on purpose – configurable rules, a project-level configuration file, stack traces in diagnostics – and a proposal that reintroduces one of them will be declined regardless of its quality.

Use one of the issue templates – **Bug**, **Feature**, or **Task** – and keep its structure: what the issue is about, and a checklist of what needs to be done. If you only want to fix a typo or a broken link in the documentation, skip the issue and open a pull request right away.

We're a small team, so a first reply can take a while. If you haven't heard from us after two weeks, a friendly ping in the issue is welcome – silence is not a no.

## Ways to Contribute

You don't have to write code to help. Field reports are among the most valuable contributions we get: tell us how ESDM behaves on your real model, where a rule got in your way, or which diagnostic sent you looking in the wrong place. If something misbehaves, file it as a **Bug**; everything else – an experience, a rule that doesn't fit your domain, a gap you noticed – is a **Task**.

The documentation lives in this repository under `documentation/docs/`, so corrections and improvements arrive as ordinary pull requests, and we welcome them. It is written in American English; follow the tone and structure of the existing pages. To preview your changes locally, run `make dev-documentation`, which needs Docker.

If you have built a tool that reads ESDM models, the **[Ecosystem page](documentation/docs/ecosystem.md)** of the documentation is where it belongs. Propose an entry with a pull request against that page, following the criteria it states; the entry itself is the pull request, no issue needed.

## Working on Code

ESDM is written in Go; the required version is pinned in `go.mod`. Branch off `main` and name the branch in kebab-case after what it does, for example `skip-update-check-on-redirected-stdout`.

We develop test-first, in the red-green-refactor cycle, and we'd encourage you to do the same. What we ask for in any case: every change comes with tests, and a bug fix comes with a test that reproduces the bug and fails without the fix. A new linter rule is one file in `rules/` with a constructor and an entry in the catalog. Every package explains the concept behind it in its `documentation.go`, which is the best place to start reading.

Write readable code: full words rather than abbreviations, comments that explain *why* rather than *what*, American English, and ASCII only in source files. Formatting follows the `.editorconfig`. The detailed conventions live under `.claude/rules/`; coding agents that read such rule files pick them up automatically, and they are worth a read for humans, too. If you change a schema, regenerate the reference snippets with `make generate-reference-snippets`; a sync test fails otherwise.

Using an AI assistant to write code is fine – we do it ourselves. What matters is that you understand what you submit, have tested it, and can explain it in the review. Output that was passed along without that scrutiny is not a contribution we can accept.

Before you open a pull request, run

```shell
$ make qa
```

It runs the static analysis, lints the example models, and executes the tests. CI runs the same on Linux, macOS, and Windows, plus a build of the documentation site.

## Commits and Pull Requests

A pull request addresses exactly one issue. If you discover something else along the way, open a new issue for it rather than folding it in – small pull requests get reviewed faster, and each one becomes a single, clean commit on `main`.

We merge by squashing, and the pull request title becomes the commit message. The title therefore follows the conventional-commit format with exactly three prefixes – `feat`, `fix`, and `chore` – starts with a capital letter, and ends with a period, for example `fix: Skip the update check when stdout is not a terminal.` CI checks the title. The commits inside your branch disappear in the squash; we'd still suggest the same format for them, since tidy commits make a larger pull request easier to review.

Describe in the pull request what you changed and why, and link the issue with `Closes #<number>`. Once the branch is pushed, keep its history intact: no force pushes and no rebasing. When `main` moves on, merge it into your branch. This keeps review comments attached to the lines they refer to, and lets reviewers look at what changed since their last pass instead of the whole pull request again.

The maintainers review every pull request and have the final say. Expect questions and requests for changes – that's the review doing its job, not a verdict on your work. We keep the tone in issues and reviews respectful and to the point, and we expect the same from everyone who takes part.

## License

By contributing you agree that your contribution is licensed under the **[MIT license](LICENSE.md)**, like the rest of the project.
