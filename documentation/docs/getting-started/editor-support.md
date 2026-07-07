# Editor Support

Most editors with a **YAML language server** can validate `.esdm.yaml` files against the ESDM schemas while you type, surface autocomplete on field names, and flag structural mistakes long before you reach `esdm lint`. ESDM gives you the two pieces an LSP needs: a local copy of the schemas, and a per-file pointer that says which schema applies. Wiring those into a specific editor is the editor's job; the convention itself lives here.

## Local Schemas

Run `esdm add-schema` once at the project root:

```shell
./esdm add-schema
```

The command writes the embedded schemas into a `schemas/` directory under the working directory, one subdirectory per schema:

```text
schemas/
├── core/
│   └── v1.yaml
├── given-when-then/
│   └── v1.yaml
└── domain-storytelling/
    └── v1.yaml
```

Commit the directory to your repository so every contributor (and CI, and your editor) reads the same schema revision.

## Schema Pointers

Each `.esdm.yaml` file declares which schema applies through a **modeline** at the top of the document:

```yaml
# yaml-language-server: $schema=./schemas/core/v1.yaml
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: book
# ...
```

The path is **relative to the file itself**, so a document deeper in the tree adjusts accordingly – `../schemas/core/v1.yaml` from one level down, `../../schemas/core/v1.yaml` from two. Extension documents point at their matching schema – `./schemas/given-when-then/v1.yaml` for a `feature` document, `./schemas/domain-storytelling/v1.yaml` for a `domain-story` document.

Multi-document files (`---`-separated) place one modeline per document. As long as one document equals one editor "view", the LSP picks the right schema for each.

## Refreshing the Schemas

After upgrading the `esdm` binary, run `esdm update-schema` to refresh the local copy:

```shell
./esdm update-schema
```

The command rejects downgrades. If your project's schemas are newer than the binary's, the upgrade is the binary, not the project. See **[CLI: esdm update-schema](/reference/cli.md#esdm-update-schema)** for the full semantics.

## Editor-Specific Setup

How you turn on the YAML language server in your editor is **out of scope for these docs**: the surface area is large, every editor handles it differently, and the answers move on a schedule we don't control. Consult your editor's YAML LSP documentation; once the LSP is active, the modeline above is all it needs.

## Where to Go Next

- **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** produces a small model through a short conversation with a coding agent that you can immediately validate against the local schemas.
- **[Your First Model by Hand](/getting-started/your-first-model.md)** is a comparable model you can paste these modelines onto artifact by artifact.
- **[Running esdm lint](/getting-started/running-esdm-lint.md)** is the structural counterpart to editor-time validation.
- **[CLI: esdm add-schema](/reference/cli.md#esdm-add-schema)** documents the schema-writing command in detail.
