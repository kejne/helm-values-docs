# Helm Chart API Docs

Human-readable API documentation for Helm chart values schemas.

`helm-values-docs` turns a chart's `values.schema.json` into clear, navigable documentation for chart users. It is intended to make the chart values contract as discoverable and useful to humans as it is to Helm and editors.

## Status

Early development.

## Goals

- Render Helm `values.schema.json` as Markdown and HTML documentation.
- Present nested values, types, defaults, descriptions, constraints, enums, and examples clearly.
- Resolve schema references and support the JSON Schema features commonly used by Helm.
- Integrate chart metadata and practical defaults from `values.yaml` where useful.
- Be deterministic and friendly to CI, static sites, and GitHub repositories.

## Development

This project uses [mise](https://mise.jdx.dev/) for the Go toolchain and [hk](https://hk.jdx.dev/) for Git hooks.

```sh
mise install
hk install
```

## License

MIT
