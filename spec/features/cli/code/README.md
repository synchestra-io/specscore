# Feature: Code

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Fcode) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore code` commands query the relationship between source code and the SpecScore resources that govern it. They read `specscore:` annotations and bare `https://specscore.md/...` URLs from source-file comments and report the features, plans, or docs those files depend on.

## Problem

Code-to-spec traceability is only useful if it is easy to query. Authors want to answer "what spec does this file implement?" without reading through every comment by hand. CI pipelines want to detect drift between code and the specs it references. This command group is the read side of the [source-references](../../source-references/README.md) contract.

## Contents

| Directory | Description |
|---|---|
| [deps/](deps/README.md) | Show SpecScore resources that given source files depend on |

### deps

Scans a set of source files for `specscore:` annotations and URL references in comments and lists the resources (features, plans, docs) those files point to. The inverse query — "what source files reference this feature?" — is not part of this group today; use `specscore feature refs` to walk spec-to-spec edges, or grep for direct source references.

## Behavior

### Scope of this group

Commands under `specscore code` are **read-only**. They MUST NOT mutate source files, the spec tree, or any project state. They operate on the working tree as it currently exists.

#### REQ: read-only

No `specscore code` subcommand may write to disk, make network calls, or require writable state. The commands are safe to run in any environment.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [source-references](../../source-references/README.md) | Defines the `specscore:` annotation and URL-in-comment conventions. This command group reads them. |
| [CLI](../README.md) | Inherits shared exit-code contract, project autodetection, and output-format conventions. |

## Outstanding Questions

- Should there be a `specscore code refs <feature_id>` command that returns the inverse (source files → feature), complementing `feature refs`? Today that query is only available via grep.

---
*This document follows the https://specscore.md/feature-specification*
