# Feature: Feature (CLI)

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ffeature) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore feature` commands query and scaffold features. They provide structured access to the feature tree (list, info, tree), dependency traversal (deps, refs), and creation (new). Together they let authors, agents, and tools navigate the spec without parsing Markdown by hand.

## Problem

The feature tree is the primary organizing structure of a SpecScore repository. Without dedicated query commands, callers either parse `README.md` files themselves (fragile, re-implements structural rules) or shell out to `grep` (loses the hierarchy). A structured query surface — with stable YAML/JSON output — lets scripts, LSP servers, and agents integrate without every consumer re-inventing the parser.

## Contents

| Directory | Description |
|---|---|
| [deps/](deps/README.md) | Features that a given feature depends on |
| [info/](info/README.md) | Feature metadata and section TOC |
| [list/](list/README.md) | Flat list of all feature IDs |
| [new/](new/README.md) | Scaffold a new feature directory |
| [refs/](refs/README.md) | Features that reference a given feature |
| [tree/](tree/README.md) | Feature hierarchy as an indented tree |

### deps

Reads the `## Dependencies` section of a feature's README and lists the features it depends on. `--transitive` follows the chain recursively.

### info

Returns structured metadata (status, parent, children, dependency counts) plus a section table-of-contents with line ranges for the feature's README. Default output is YAML.

### list

Flat listing of every feature in the project as feature IDs, one per line, sorted alphabetically. `--fields` adds metadata columns.

### new

Scaffolds a new feature directory with a README containing every required section. Supports `--parent` for sub-features, `--depends-on` for dependency wiring, and `--commit` / `--push` for atomic git integration.

### refs

Inverse of `deps`: walks all features and reports which ones reference the given feature in their `## Dependencies` section. `--transitive` follows the chain.

### tree

Renders the feature hierarchy as an indented tree. Without an argument, shows the full tree. With a `feature_id`, shows that feature in context (ancestors + subtree); `--direction up` narrows to ancestors only, `--direction down` to the subtree only.

## Behavior

### Shared flags

Every command in this group accepts the shared flags defined in the [CLI parent](../README.md#shared-flags): `--project`, `--format`, `-h/--help`.

Commands that return structured data (deps, info, list, refs, tree) additionally share:

| Flag | Meaning |
|---|---|
| `--fields` | Comma-separated metadata fields to include per entry (e.g., `status,oq`). |
| `--format` | `yaml` (default), `json`, or `text`. `--fields` forces YAML when `text` is incompatible. |

#### REQ: fields-shape

`--fields` MUST accept a comma-separated list of field names. Unknown field names MUST exit `2` (InvalidArgs) with a message naming the offending field. Recognized field names are a stable contract across patch releases.

#### REQ: format-selection

`--format` MUST accept `yaml`, `json`, and `text`. Default is `yaml`. Text output MUST be suitable for shell piping (one entry per line, no ornamental characters). Any other value MUST exit `2`.

### Feature-ID resolution

Commands that take a `<feature_id>` argument MUST resolve it relative to `spec/features/` in the project root. IDs use the path-identification rules from the [feature](../../feature/README.md) spec (e.g., `cli/version`, `authentication`, `billing/payments`).

#### REQ: not-found-exit-code

When a feature ID does not resolve to an existing feature directory, commands MUST exit `3` (NotFound) with a message that names the requested ID.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [feature](../../feature/README.md) | Defines what a feature is. This command group provides the CLI surface for querying and creating them. |
| [CLI](../README.md) | Inherits shared exit-code contract, flag conventions, and project autodetection. |

## Outstanding Questions

- The `--fields` flag overlaps between `deps`, `info`, `list`, `refs`, and `tree`, but the exact set of valid field names is not enumerated in a single place. Should field names move to a shared registry (e.g., a REQ in this parent) or stay documented per-command?
- `feature new --parent <id>` and `feature new --depends-on <ids>` overlap with Idea-derived scaffolding from `new idea`. Should there be a consolidated `new feature` path that handles both flows, or does the current split (general "new" group for pre-spec artifacts, `feature new` inside this group for features) remain clearer?

---
*This document follows the https://specscore.md/feature-specification*
