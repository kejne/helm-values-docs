# AGENTS.md

## Project

`helm-values-docs` is a Go project that will turn Helm `values.schema.json` files into human-readable documentation. Keep output deterministic and suitable for CI and static sites.

## Tooling and validation

For the normal agent loop, run one command after making changes:

```sh
hk check
```

`hk check` runs the inexpensive checks against modified files. It includes formatting, vetting, `golangci-lint`, Go tests, module-tidy validation, and applicable GitHub Actions checks. It is the canonical agent validation command; do not replace it with a collection of ad hoc linter commands.

Other scopes:

- `hk check --all` runs the inexpensive checks against the whole repository.
- `hk run pre-push` runs the full quality suite for the push range, including `gosec` and `govulncheck`.
- Review any files changed by a fixer before committing.

The pre-push suite intentionally keeps the slower/network-dependent security checks out of the normal agent loop. Run it before pushing or when changes affect security-sensitive code. Do not bypass a failing check; fix the cause or explain an intentional exception explicitly.

The commit-message hook requires conventional commits. Install hooks locally with `hk install` when setting up a clone.

## Go layout

Use the conventional Go layout:

```text
cmd/<program>/main.go       # executable entry points
internal/<area>/             # implementation packages private to this module
internal/<area>/*_test.go    # package tests
<package>/*_test.go          # tests for an intentional root package
testdata/                    # fixtures and approval files, adjacent to their package when local
```

- Put commands under `cmd/` and reusable implementation under `internal/`.
- Prefer `internal/` over a new public package. Add `pkg/` only when a package is deliberately part of the external API.
- Keep package boundaries clear; do not put application logic in `cmd` beyond wiring dependencies.
- Keep generated files and test fixtures deterministic. Do not commit build artifacts.

## TDD

Use test-driven development for behavior changes:

1. Add or update a focused test that fails for the missing behavior.
2. Implement the smallest change that makes it pass.
3. Refactor with the tests still passing.
4. Run `hk check` after each meaningful change and before handing work off.

Prefer table-driven tests and test behavior through stable package APIs. Keep tests isolated, deterministic, and independent of network services, wall-clock time, and machine-specific paths.

## Strict approval testing

Use approval/golden tests for substantial rendered or serialized output. Approval tests must be strict:

- Compare the complete canonical output byte-for-byte, including meaningful whitespace and trailing newlines.
- Store approved results in reviewed, checked-in `testdata` files near the package under test.
- A normal test run must never rewrite approvals or silently accept a mismatch.
- Approval updates require explicit opt-in, for example `UPDATE_APPROVALS=1`, followed by reviewing the diff and rerunning tests without the flag.
- Make output deterministic before approving it: sort maps and collections, normalize paths, and inject clocks or IDs rather than reading them directly.
- Treat approval-file changes as implementation changes. Do not update an approval merely to make a failing test green; verify the behavioral change first.
- Fail on missing, unexpected, or stale approval files unless their removal is intentional and reviewed.

Keep approvals focused enough to review. Use ordinary assertions for small, local invariants and approval tests for user-visible documents, schemas, and other large structured output.
