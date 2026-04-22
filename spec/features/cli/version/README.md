# Feature: Version

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Fversion) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore version` and `specscore --version` (with its `-v` short alias) report the CLI's build identity. The subcommand prints the version, commit, and build date on a single line for humans and bug reports. The flag prints only the bare semver so scripts, installers, and CI gates can consume it without parsing.

## Synopsis

```
specscore version
specscore --version
specscore -v
```

## Problem

Users, install scripts, and support workflows need a reliable way to identify which `specscore` binary is running. Two things typically go wrong when this is not specified:

- **Inconsistent surfaces.** A CLI that only ships `--version` forces humans to read a terse, context-free string. A CLI that only ships a `version` subcommand forces scripts to parse a multi-field line. Shipping both without pinning their formats leads to drift when one is changed and the other is not.
- **Unstable output format.** Users write regexes against version output. When the format changes — extra prose, a `v` prefix, a reordered field — those regexes break. Without a spec, every release is free to change the shape of the output.

Pinning both surfaces keeps humans and scripts from stepping on each other.

## Behavior

### Two output surfaces

The CLI exposes version information through two distinct surfaces that serve different audiences but read from the same underlying values.

| Surface | Audience | Output shape |
|---|---|---|
| `specscore version` | Humans, bug reports, support | `specscore <version> (<commit>) <date>` |
| `specscore --version` / `-v` | Scripts, installers, CI | `<version>` |

#### REQ: subcommand-output

`specscore version` MUST print a single line of the form:

```
specscore <version> (<commit>) <date>
```

`<version>` is the bare semver (see [REQ: no-v-prefix](#req-no-v-prefix)). `<commit>` is the full git commit SHA the binary was built from. `<date>` is the build timestamp in RFC 3339 / ISO 8601 form. The line MUST end with a single trailing newline.

#### REQ: flag-output

`specscore --version` MUST print only the bare semver on a single line, terminated by a newline. The output MUST NOT include the program name, the commit, the build date, or any other decoration. A caller MUST be able to consume the output with `$(specscore --version)` and receive exactly the version string.

#### REQ: short-flag

`-v` MUST be accepted as a short alias for `--version` and MUST produce identical output.

#### REQ: no-v-prefix

The `<version>` field MUST NOT carry a leading `v`. `0.11.0` is correct; `v0.11.0` is not. This holds on both surfaces. The `v` prefix remains on git tags, release filenames, and the `SPECSCORE_VERSION` install variable — those are outside the scope of CLI output.

The convention matches the broader Go ecosystem: `go version`, `gh --version`, `docker version`, `hugo version`, and `kubectl` client version all print bare semver even though their tags carry `v`.

### Build-time value injection

The version, commit, and date fields are populated at build time and baked into the binary.

#### REQ: ldflag-injection

The three values MUST be injected via Go linker flags against package-level `var` symbols in `internal/cli`:

```
-X github.com/synchestra-io/specscore/internal/cli.version=<semver>
-X github.com/synchestra-io/specscore/internal/cli.commit=<full-sha>
-X github.com/synchestra-io/specscore/internal/cli.date=<rfc3339>
```

A release build MUST supply all three. The release workflow (goreleaser / equivalent) is the canonical producer of these values.

#### REQ: default-placeholders

When the binary is built without `-ldflags` (typical of `go run`, `go build` during development, or `go install` from a consumer), the three fields MUST fall back to literal placeholders: `version="dev"`, `commit="none"`, `date="unknown"`. The CLI MUST NOT error on missing version information — placeholders are a valid state.

A `dev` binary therefore prints:

- `specscore --version` → `dev`
- `specscore version` → `specscore dev (none) unknown`

### No public API

Version information is read-only and confined to the CLI.

#### REQ: no-public-api

The `pkg/` tree MUST NOT expose version, commit, or date values. Consumers that import `specscore` as a library depend on the Go module version, not the CLI's build-time metadata.

## Parameters

None. Neither the subcommand nor the flag accepts arguments.

## Exit codes

| Exit code | Meaning |
|---|---|
| `0` | Success (always, for both surfaces) |

The version command cannot fail under normal operation — missing build metadata falls back to placeholders rather than erroring. Exits other than `0` indicate an unexpected runtime fault in the CLI itself, handled by the shared error path defined in the [CLI parent feature](../README.md).

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [CLI](../README.md) | Parent feature. `version` inherits any shared CLI conventions introduced there. |
| [Project Definition](../../project-definition/README.md) | Unrelated. The CLI version is a property of the binary, not of the spec repo it operates on. A spec repo does not pin a CLI version. |

The installation documentation in `docs/installation.md` shows users how to verify a successful install by matching the shape defined in [REQ: subcommand-output](#req-subcommand-output). Changes to that format MUST be reflected in the installation docs.

## Acceptance Criteria

### AC: surfaces-agree

**Requirements:** cli/version#req:subcommand-output, cli/version#req:flag-output, cli/version#req:short-flag

`specscore version`, `specscore --version`, and `specscore -v` all print the same version string (with different surrounding context). The flag surfaces print only the bare semver; the subcommand adds commit and date in parentheses.

### AC: scripting-friendly-flag

**Requirements:** cli/version#req:flag-output, cli/version#req:no-v-prefix

`$(specscore --version)` yields a single bare semver with no prefix, no program name, no commit, no trailing whitespace beyond a single newline — safe to feed directly into semver comparators, installer version checks, and CI gates.

### AC: dev-build-works

**Requirements:** cli/version#req:default-placeholders

A `specscore` binary built without `-ldflags` (e.g., `go run ./cmd/specscore --version`) exits `0` and prints `dev`. `specscore version` prints `specscore dev (none) unknown`. The CLI never errors or panics because version metadata is missing.

### AC: go-idiom-format

**Requirements:** cli/version#req:no-v-prefix

Version output follows Go-ecosystem convention: no `v` prefix on the printed number, even though the underlying git tag is `v`-prefixed. Any CLI output containing `v0.` or `v1.` at the start of the version field is a regression.

## Outstanding Questions

- Should `specscore version` gain a `--json` (or `--format`) flag for machine consumption that includes the commit and date, now that the bare `--version` flag is reserved for the minimal scripting form?
- Should the build-date field have a normative time-zone requirement (UTC only) to make output reproducible across release machines, or is any valid RFC 3339 value acceptable?

---
*This document follows the https://specscore.md/feature-specification*
