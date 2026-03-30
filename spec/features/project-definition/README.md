# Feature: Project Definition

**Status:** Conceptual

## Summary

How SpecScore discovers and configures specification projects.

## Project File

Every SpecScore project has a `specscore-project.yaml` file as its entry point. This file lives at the root of the specification repository and declares project metadata and directory conventions.

### Mandatory fields

| Field   | Description                 |
|---------|-----------------------------|
| `title` | Human-readable project name |

### Optional fields

| Field                         | Default | Description                                                           |
|-------------------------------|---------|-----------------------------------------------------------------------|
| `description`                 | —       | Brief description of the project                                      |
| `repos`                       | —       | List of code repository URLs associated with this project             |
| `project_dirs.specifications` | `spec`  | Directory for technical specifications (features, architecture, etc.) |
| `project_dirs.documents`      | `docs`  | Directory for user-facing documentation                               |

Orchestration tools may extend this file with additional fields (see [Orchestration Tool Extensions](#orchestration-tool-extensions)).

## Repository Structure

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

The project entry point is `specscore-project.yaml` at the repository root. The `spec/` and `docs/` directories are configurable via `project_dirs` fields.

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

## Orchestration Tool Extensions

Orchestration tools (like Synchestra) may extend `specscore-project.yaml` with additional fields:

```yaml
# Example: Synchestra extension
state_repo: https://github.com/org/project-synchestra
planning:
  auto_create: false
```

SpecScore ignores unknown fields, allowing any tool to add its own configuration alongside the standard SpecScore fields.

## Outstanding Questions

- Should there be a schema version field in `specscore-project.yaml` to support future evolution?
- Should `project_dirs` support additional custom directories beyond `specifications` and `documents`?
- How should SpecScore handle a repository that has no `specscore-project.yaml` — infer defaults or refuse to operate?
