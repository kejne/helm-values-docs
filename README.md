# Helm Chart API Docs

Human-readable API documentation for Helm chart values schemas.

`helm-values-docs` turns a chart's `values.schema.json` into clear, navigable documentation for chart users. It is intended to make the chart values contract as discoverable and useful to humans as it is to Helm and editors.

## Status

The Markdown renderer is available. It documents the common JSON Schema
object/property subset used by Helm values schemas.

## Render Markdown

```sh
hvd render \
  --schema values.schema.json \
  --values values.yaml \
  --values values-production.yaml \
  --output docs/values.md
```

`--values` is repeatable. Files are merged in command-line order: mappings are
merged recursively and later scalar or sequence values replace earlier values.
Displayed defaults come only from the merged values files; schema defaults are
not used. Properties are sorted lexically, so output is deterministic. The
renderer supports types,
descriptions, required properties, defaults, enums, examples, and common
string, number, and array constraints. Unsupported advanced schema constructs
are listed in the generated document rather than inferred.

The command requires `--schema`, at least one `--values`, and `--output`.
Invalid input and output failures return a nonzero status. Output is assembled
in memory and atomically replaced after successful rendering.

## Check generated Markdown

Use the same inputs with `check` in CI:

```sh
hvd check \
  --schema values.schema.json \
  --values values.yaml \
  --values values-production.yaml \
  --output docs/values.md
```

`check` compares bytes exactly and never writes the output. It exits 0 when
current, 1 when the output is missing or stale, and 2 for usage, input, or
output errors.

## Template composition

Both commands accept an optional `--template TEMPLATE` flag. The template is
executed with the stable data contract `TemplateData{API string}`, so
`{{.API}}` inserts the canonical generated API Markdown:

```markdown
# Chart values

{{.API}}
## Additional chart notes
```

Without `--template`, output remains the standalone API document. Template
read, parse, and execution failures return status 2 and identify the template
path.

## Goals

- Render Helm `values.schema.json` as Markdown documentation.
- Present nested values, types, defaults, descriptions, constraints, enums, and examples clearly.
- Be deterministic and friendly to CI, static sites, and GitHub repositories.

## Architecture

The implementation follows hexagonal architecture. `internal/domain` holds the
schema model; `internal/ports/in` defines use cases and `internal/ports/out`
defines the application's adapter interfaces. `internal/app` coordinates those
ports. CLI, JSON Schema, YAML, Markdown, and filesystem implementations live in
`internal/adapters/in` or `internal/adapters/out`. The executable is the
composition root in `cmd/hvd`.

`.go-arch-lint.yml` describes and enforces the allowed import direction. Run
`go-arch-lint check` to validate it; `hk check` runs it alongside the Go checks.

## Development

This project uses [mise](https://mise.jdx.dev/) for the Go toolchain and [hk](https://hk.jdx.dev/) for Git hooks.

```sh
mise install
hk install
```

## License

MIT
