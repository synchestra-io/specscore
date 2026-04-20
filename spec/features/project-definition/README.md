# Feature: Project Definition

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fproject-definition) — graph, discussions, approvals

**Status:** Stable

## Summary

How SpecScore discovers and configures specification projects. Every SpecScore project has a `specscore-project.yaml` file as its entry point, declaring project metadata and directory conventions.

## Contents

| Directory | Description |
|-----------|-------------|
| [_tests/](_tests/README.md) | Test scenarios for project definition requirements |

## Problem

SpecScore needs a consistent way to discover where specifications and documents live within a repository. Without a project definition file, tools must guess directory layouts, project names, and repository associations. A standard project file eliminates ambiguity for both human contributors and spec-aware tooling.

## Behavior

### Project file location

The project configuration file is `specscore-project.yaml`. It lives at the root of the specification repository.

#### REQ: project-file-location

Every SpecScore project MUST have a `specscore-project.yaml` file at the repository root. This file is the entry point for all SpecScore tooling.

### Mandatory fields

The project file has one mandatory field:

| Field   | Description                 |
|---------|-----------------------------|
| `title` | Human-readable project name |

#### REQ: title-required

The `specscore-project.yaml` file MUST contain a `title` field. A project file without a `title` is invalid.

### Optional fields

| Field                         | Default | Description                                                           |
|-------------------------------|---------|-----------------------------------------------------------------------|
| `description`                 | —       | Brief description of the project                                      |
| `repos`                       | —       | List of code repository URLs associated with this project             |
| `project_dirs.specifications` | `spec`  | Directory for technical specifications (features, architecture, etc.) |
| `project_dirs.documents`      | `docs`  | Directory for user-facing documentation                               |

#### REQ: optional-field-defaults

When optional `project_dirs` fields are omitted, SpecScore MUST use the default values: `spec` for `project_dirs.specifications` and `docs` for `project_dirs.documents`.

### Repository structure

SpecScore specification projects follow a standard directory layout:

```
{repo}/
  specscore-project.yaml        # Project configuration
  README.md                     # Project overview
  spec/                         # Specifications (configurable)
    features/
      ...
    plans/
      ...
  docs/                         # Documentation (configurable)
    ...
```

#### REQ: directory-layout

The `spec/` and `docs/` directories MUST be configurable via `project_dirs` fields. Tools MUST resolve specification and document paths using the configured values rather than hardcoded defaults.

### Example

```yaml
title: My Service
description: Backend API and web frontend for Acme Corp
repos:
  - https://github.com/org/my-service-api
  - https://github.com/org/my-service-web
project_dirs:
  specifications: spec
  documents: docs
```

### Orchestration tool extensions

Orchestration tools (like Synchestra) may extend `specscore-project.yaml` with additional fields:

```yaml
# Example: Synchestra extension
state_repo: https://github.com/org/project-synchestra
planning:
  auto_create: false
```

#### REQ: unknown-fields-ignored

SpecScore MUST ignore unknown fields in `specscore-project.yaml`. This allows any orchestration tool to add its own configuration alongside the standard SpecScore fields without causing validation errors.

### Adherence footer

#### REQ: adherence-footer

Every project-definition document MUST end with an adherence footer per the [Adherence Footer feature](../adherence-footer/README.md). The footer URL MUST be `https://specscore.md/project-definition-specification`.

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Should there be a schema version field in `specscore-project.yaml` to support future evolution?
- Should `project_dirs` support additional custom directories beyond `specifications` and `documents`?
- How should SpecScore handle a repository that has no `specscore-project.yaml` — infer defaults or refuse to operate?
- Acceptance criteria not yet defined for this feature.

---
*This document follows the https://specscore.md/feature-specification*
