---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Repo Config

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/specscore/specscore/spec/features/repo-config?op=explore) | [Edit](https://specscore.studio/app/github.com/specscore/specscore/spec/features/repo-config?op=edit) | [Ask question](https://specscore.studio/app/github.com/specscore/specscore/spec/features/repo-config?op=ask) | [Request change](https://specscore.studio/app/github.com/specscore/specscore/spec/features/repo-config?op=request-change) |

**Status:** Draft
**Source Ideas:** —

## Summary

Defines `specscore.yaml`, the single mandatory repository-level config file for SpecScore projects. Specifies the file name, the mandatory schema-pointer header comment, the optional `project` identity block, related-project navigation hints, explicit Plan-repository routing, top-level dir-name overrides, modules with code roots, studio configuration, publication policy, and inference defaults.

## Contents

| Directory | Description |
|-----------|-------------|
| [_tests/](_tests/README.md) | Test scenarios for repo-config requirements |

## Problem

SpecScore needs a single, predictable, repo-level config file. The legacy `specscore-spec-repo.yaml` carried Synchestra-flavored naming (the spec-repo / state-repo / code-repo split) that does not apply to SpecScore proper, and its schema did not support modular layouts (e.g., a frontend and a backend module sharing one repo). Tools that render the studio toolbar in artifact documents (the four canonical verbs: Explore, Edit, Ask question, Request change) need a name-agnostic place to read studio settings; the legacy `viewer:` block (with its single-link "View in …" model) is out of date and is superseded by `studio:`.

This feature replaces `project-definition` with a simpler, more flexible schema.

## Behavior

### Lesson classifications

#### REQ: lessons-classifications

A repository that uses Lessons MUST declare `lessons.classifications` in `specscore.yaml` as a non-empty list of unique lowercase, hyphen-separated strings. The list is the controlled vocabulary for [Lesson](../lesson/README.md) `**Classifications:**`; removing a value while any Lesson still uses it is a validation error. Repositories with no Lessons MAY omit the `lessons:` block until their first Lesson is created.

### File name and location

The config file is `specscore.yaml` and lives at the repository root.

#### REQ: file-name

Every SpecScore project MUST have a file named exactly `specscore.yaml` at the repository root. Any other name (including `specscore-project.yaml`, `specscore-spec-repo.yaml`, `.specscore.yaml`) is invalid.

### Schema-pointer header comment

Every `specscore.yaml` is self-identifying. Line 1 of the file is a fixed comment that names the schema. The comment is required even when the rest of the file is empty.

#### REQ: schema-header-comment

Line 1 of `specscore.yaml` MUST be exactly:

```yaml
# SpecScore Repo Config Schema: https://specscore.md/repo-config
```

A file whose line 1 is anything else — including a blank line, a different comment, a YAML document marker (`---`), or a leading BOM followed by other content — is invalid. Lint MUST emit a hard error.

#### REQ: empty-config-valid

A `specscore.yaml` containing only the schema-header comment (with an optional trailing newline) is valid. All field defaults defined by this feature apply.

### Project identity

The optional `project:` block declares the project's identity. Every field within it is also optional; missing fields are inferred from the git remote or from the working directory.

| Field | Default | Description |
|-------|---------|-------------|
| `title` | git repository name | Human-readable project name |
| `host` | parsed from `origin` git remote | Repository host (e.g., `github.com`) |
| `org` | parsed from `origin` git remote | Repository org / owner |
| `repo` | parsed from `origin` git remote | Repository name |
| `repositories` | — | List of code repository URLs associated with this project |

#### REQ: project-block-optional

The `project:` block MAY be omitted entirely. When omitted, every field below is inferred from the git remote and working directory; the file remains valid.

#### REQ: project-title-default

When `project.title` is omitted, the effective title is the basename of the repository's working directory. When `project.title` is present, its value is used as-is.

#### REQ: source-reference-overrides

`project.host`, `project.org`, and `project.repo` override `origin` git-remote inference for source-reference URL building. Tools resolving same-repo references MUST prefer explicit values when present. This handles forks, mirrors, and non-standard remote names. See the [Source References feature](../source-references/README.md).

#### REQ: project-repositories

`project.repositories` is an OPTIONAL list of repositories associated with this project. The list MUST round-trip on read/write. SpecScore validates the structural shape (entry shape, role enum, required `url` and non-empty `roles`) but does not act on the role values — interpretation is left to downstream orchestration tools (e.g., Synchestra). Every entry MUST be a role-tagged object per `repositories-entry-shape`; flat URL strings are not accepted.

#### REQ: repositories-entry-shape

Each entry in `project.repositories` is a YAML mapping (object) with the following fields. Flat scalar entries (URL strings, single tokens) are a hard error.

| Field | Required | Description |
|-------|----------|-------------|
| `url` | yes | Repository URL or short-path (same forms accepted everywhere `project.repositories` entries are written) |
| `roles` | yes | Non-empty list of role values; see `repositories-roles-list` and `repositories-roles-enum` |
| `title` | no | Human-readable repository name |
| `comment` | no | Free-form annotation |

A mapping entry that is missing `url`, has an empty `roles` list, or has a missing `roles` field is a hard error. Unknown fields inside a role-tagged entry MUST round-trip unchanged per `unknown-fields-preserved`.

#### REQ: repositories-roles-list

`roles` is a mandatory, non-empty YAML sequence (list) of one or more role values. Omitting `roles`, providing it as an empty list (`roles: []`), or providing a scalar value (e.g., `roles: code` instead of `roles: [code]`) is each a hard error — the field is always a non-empty list. Duplicate role values within a single entry's `roles` list are de-duplicated by tools and SHOULD emit a lint advisory; they are not a hard error.

#### REQ: repositories-roles-enum

Each value in a `roles` list MUST be drawn from the canonical role enum:

| Role | Meaning |
|------|---------|
| `code` | Source-code repository — scanned by code tooling for `specscore:` annotations |
| `specification` | SpecScore-managed spec repository — has its own `specscore.yaml` and `spec/` tree |
| `state` | Orchestrator state repository — managed by tools like Synchestra; not specscore-managed |
| `docs` | Documentation-only repository |
| `runner` | Remote runner / executor repository — referenced by orchestrator runner features |

A repository entry MAY combine multiple roles (e.g., `roles: [specification, code]`) — a single repo often plays several roles. Unknown role values are a hard error. The enum is closed in v1; future additions require a SpecScore Feature revision.

### Related projects

The optional root-level `projects:` field is a list of related SpecScore projects worth navigating to from this one. The relationship direction (parent, child, sibling, dependency) is unspecified; entries are navigation aids.

```yaml
projects:
  - https://github.com/org/parent-project          # external repo URL
  - ./packages/api                                  # path to dir with its own specscore.yaml
```

#### REQ: projects-list

`projects:` is an OPTIONAL root-level list. Each entry MUST be either:

1. A URL to an external repository (any scheme parseable as a URI), or
2. A path (relative or starting with `./`) to a directory in the current repository.

#### REQ: projects-local-path-must-resolve

When a `projects:` entry is a local path, that path MUST resolve to a directory containing its own `specscore.yaml` file. Lint MUST emit an error for local-path entries pointing at directories that do not exist or do not contain `specscore.yaml`.

### Plan repository routing

Plan storage is resolved by SpecScore itself. It does not depend on WB or another orchestrator. Routing uses portable repository identities, while checkout paths are kept in machine-local configuration.

Repository identities use either `owner/repo` shorthand (the source repository's host is implied) or the full `host/owner/repo` form. Normalization MUST preserve all three logical components and MUST reject absolute paths, `.` or `..` segments, and malformed identities.

The committed project config and its local override select one destination with a scalar:

```yaml
plans_repo: specscore/specscore
```

User and organization config may route several source projects to a destination:

```yaml
plan_repos:
  sneat-co/workbench:
    - datatug/datatug
    - github.com/datatug/dashboardius

repo_checkouts:
  sneat-co/workbench: /work/checkouts/sneat-co/workbench
```

#### REQ: plans-repo-project-selection

`plans_repo` in `specscore.local.yaml` or the repository's committed `specscore.yaml` MUST be a non-empty scalar repository identity. It selects the Plan destination for that source project. The local value overrides the committed value. Setting it to the source repository's own identity explicitly selects same-repository Plan storage; omission does not imply same-repository storage.

#### REQ: plan-repos-aggregate-routing

`plan_repos` in organization or user config MUST map destination repository identities to lists of source repository identities. A source MUST resolve to at most one destination within a layer; duplicate mappings to different normalized destinations are a hard error. Repository/local `plans_repo` overrides organization routing, which overrides user routing.

#### REQ: plan-route-required

Every Plan operation MUST resolve an explicit route from `plans_repo` or `plan_repos`. If none matches the source project, the Plan operation MUST fail before reading or writing Plan artifacts and explain where routing can be configured. Feature, Idea, Proposal, Lesson, and other non-Plan operations MUST remain independent of Plan routing.

#### REQ: repo-checkouts-machine-local

`repo_checkouts` maps a repository identity to an absolute path whose git top-level and origin identity match that repository. It MAY appear in user, organization, or `specscore.local.yaml` configuration and MUST NOT appear in committed `specscore.yaml`. Missing, relative, nested-subdirectory, non-repository, or identity-mismatched paths are hard errors when an external Plan checkout is required. Symlinks are resolved before these checks.

#### REQ: plan-config-precedence

Plan routing precedence is `specscore.local.yaml` then committed `specscore.yaml`, then organization `.specscore.yaml`, then user `~/.specscore.yaml`. Checkout precedence is local, organization, then user; committed project config is forbidden from contributing machine paths. The organization file is the `.specscore.yaml` inside the organization directory containing the canonical source clone, so linked worktrees share the same organization policy.

### Directory names

Top-level overrides for the spec and docs directory names. These apply uniformly across all modules. Per-module, per-artifact-kind overrides are provided separately by `path_overrides` (see [Path overrides](#path-overrides) below), which currently covers only the ideas directory.

| Field | Default | Description |
|-------|---------|-------------|
| `specs_dir_name` | `specs` | Name of the specs directory inside each module |
| `docs_dir_name` | `docs` | Name of the docs directory inside each module |

#### REQ: specs-dir-name-default

When `specs_dir_name` is omitted, tools MUST treat the value as `specs`. The same name applies to every module.

#### REQ: docs-dir-name-default

When `docs_dir_name` is omitted, tools MUST treat the value as `docs`. The same name applies to every module.

### Path overrides

Per-module overrides for where individual artifact kinds live, declared under a module's optional `path_overrides` mapping. Unlike `specs_dir_name`/`docs_dir_name` (repo-wide, name-only), `path_overrides` is per-module and takes a full module-relative path, so a kind can be relocated out of the spec tree entirely. In v1 the only interpreted key is `ideas_path`; the resolution behavior is owned by the [Configurable Ideas Path](../configurable-ideas-path/README.md) Feature.

```yaml
modules:
  - path: ""
    path_overrides:
      ideas_path: ideas        # ideas live at ./ideas instead of ./spec/ideas
```

#### REQ: path-overrides-optional

A module's `path_overrides` is an OPTIONAL mapping. When absent, every artifact kind uses its default location. Keys other than `ideas_path` MUST round-trip on read/write unchanged (per `unknown-fields-preserved`) but are not interpreted in v1.

#### REQ: ideas-path-field

`path_overrides.ideas_path`, when present, MUST be a module-relative path string giving the directory for that module's Ideas; when absent the default is `spec/ideas`. An absolute path (leading `/`) or a value escaping the module root via `../` is a hard error naming the offending module. The path's resolution semantics are defined by [Configurable Ideas Path](../configurable-ideas-path/README.md).

### Modules

A module is a code area within the repository, optionally with its own specs and docs subdirectories. The feature graph is unified across modules; modules partition where files live, not the graph itself.

```yaml
modules:
  - name: Highlevel        # specs at ./specs/, docs at ./docs/
  - name: Backend
    path: backend          # specs at ./backend/specs/, docs at ./backend/docs/
    code:
      - pkg/               # ./backend/pkg/
      - go.mod             # ./backend/go.mod
      - go.sum             # ./backend/go.sum
  - path: frontend         # name deduced as "frontend"
```

#### REQ: modules-default

When `modules:` is omitted, tools MUST behave as if a single module `{}` were declared — one module at the repository root with no `code:` paths.

#### REQ: module-name-deduction

When a module's `name` field is omitted, the effective name is the basename of `path`. When both `name` and `path` are omitted, the effective name is `default`.

#### REQ: module-path-resolution

A module's specs directory is `{module.path}/{specs_dir_name}/`. Its docs directory is `{module.path}/{docs_dir_name}/`. Each entry in `code:` is resolved as a path (file or directory) relative to `{module.path}/`. When `module.path` is omitted, all resolution is relative to the repository root.

#### REQ: module-code-list

`code:` is an OPTIONAL per-module list of paths (files or directories). Tools scanning code for `specscore:` annotations MUST use these paths as scan roots, scoped to the module's effective `path`. A module without a `code:` list is spec-only — no code is attributed to it.

#### REQ: module-paths-unique

Two modules MUST NOT have the same effective `path` (where omitted `path` is treated as `.`). Duplicate module paths are a hard error.

#### REQ: module-paths-non-nested

Two modules with explicit `path:` values MUST NOT have an ancestor-descendant relationship (i.e., one path is a prefix of the other interpreted as filesystem paths). A module with no `path:` (the implicit-root module) MAY coexist with any number of explicit-path modules. Nested explicit-path modules are a hard error.

### Studio

The optional `studio:` block names the upstream studio (Spec Studio web app) that renders SpecScore artifacts in a browser and serves the toolbar. Tools that emit the toolbar line in artifact documents (e.g., feature READMEs) read this block.

```yaml
studio:
  name: SpecScore.Studio
  url: https://specscore.studio/
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes (when `studio:` is present as a mapping) | Display name shown in the toolbar brand attribution |
| `url` | yes (when `studio:` is present as a mapping) | Base URL of the studio; toolbar links point under this URL |

#### REQ: studio-default-when-omitted

When `studio:` is omitted entirely, tools MUST use the defaults `name: SpecScore.Studio` and `url: https://specscore.studio/`. The studio toolbar MUST be rendered in artifact documents with these defaults.

#### REQ: studio-explicit-values

When `studio:` is present as a mapping, both `name` (string) and `url` (valid URL) MUST be present. Partial mappings (e.g., `studio: { name: X }` without a URL) are a hard error. Defaults do not apply to a partially-specified `studio:` block — defaults apply only when the whole block is absent.

#### REQ: studio-null-opts-out

When `studio:` is explicitly set to `null` (YAML `null`, `~`, or an empty value), tools MUST NOT render any studio toolbar in artifact documents. This is the only way to suppress the toolbar; omitting the block falls back to defaults.

#### REQ: studio-toolbar-mandatory-unless-opted-out

Tools rendering artifact documents that support a studio toolbar (e.g., the lint `studio-toolbar` rule operating on feature READMEs) MUST emit such a toolbar unless `studio: null` is explicitly set. Implementations MUST NOT silently omit the toolbar for any other reason. The toolbar's rendering convention — its line position, byte form, verb set, separator, and URL grammar — is owned by the [studio-toolbar](../studio-toolbar/README.md) Feature, not by this Feature; `repo-config` owns only the `studio:` schema (defaults, opt-out, explicit-values rules, trailing-slash validation).

#### REQ: studio-url-trailing-slash

`studio.url` MUST end with exactly one trailing `/` character. A `studio.url` value with no trailing slash (e.g., `https://specscore.studio`) or with multiple trailing slashes (e.g., `https://specscore.studio//`) is a hard error. This schema-level validation is owned by `repo-config`; the [studio-toolbar](../studio-toolbar/README.md) Feature consumes the validated value and strips the single trailing `/` before joining it with the toolbar URL path grammar.

### Grade values

The optional `grade:` block declares the repository's allowed `**Grade:**` body-metadata values:

```yaml
grade:
  values: [A, B, C, D, F]
```

`grade.values` is an OPTIONAL list of allowed grade tokens. When the block (or its `values:` key) is omitted, the built-in default `A, B, C, D, F` applies. This Feature owns only the presence of the `grade:` key in the `specscore.yaml` schema; the value-set shape rules, the default, and how `**Grade:**` is parsed, placed, and validated are owned by the [Grade body-metadata field](../canonical-grade-metadata-field/README.md) Feature.

### Publication policy

The optional `publication:` block declares project-level publication policy for SpecScore producers. It is consumed by SpecStudio skills and by `specscore` CLI helpers that need durable defaults for whether artifact changes are left unstaged, staged, committed, or pushed.

```yaml
publication:
  default:
    actions: [stage]
  events:
    idea.approved:
      actions: [stage, commit]
  commands:
    implement:
      events:
        feature.approved:
          actions: [stage, commit]
  push:
    allow_branches: ["feature/*"]
    deny_branches: ["main", "release/*"]
```

#### REQ: publication-block-optional

The `publication:` block MAY be omitted entirely. Absence means no project-level publication policy is configured; consumers fall back to user, session, task, or built-in defaults owned by their respective features.

#### REQ: publication-schema-delegated

When `publication:` is present, its durable shape and validation rules are owned by the [Publication Policy Config](../publication-policy-config/README.md) Feature. Repo Config owns the top-level key and the requirement that unknown fields under `publication:` round-trip unchanged with the rest of `specscore.yaml`.

### Recaps and journal storage targets

The optional `recaps:` and `journal:` blocks declare where machine-generated recap/verify artifacts and the activity journal are stored. Per [Decision D-0002](../../decisions/0002-recap-storage-surfacing-and-producer-gate-boundary.md), storage **always targets the hub repo** — a configured location, default a dedicated repo, that **may be configured to equal the current code repo** (covering the case before the Hub backend exists and repos that are not SpecScore-managed). The drift recap, verify report, and session recap all resolve their storage target through these blocks.

```yaml
recaps:
  enabled: true
  repo: ~/work-hub     # hub-repo path (may point at the current repo); layered-config-gated
  user: alice          # explicit username override; layered-config-gated
  users:               # email→username map (never raw emails in paths)
    alice@example.com: alice

journal:
  enabled: true
  repo: ~/work-hub     # layered-config-gated (Phase 2)
```

#### REQ: recaps-journal-storage-target

`specscore.yaml` MAY carry `recaps:` and `journal:` blocks. Storage always targets the hub repo — `recaps.repo` / `journal.repo` is the configured hub-repo location (default a dedicated repo) and may be configured to equal the current code repo, in which case artifacts are written into it directly. There is no separate `local` mode: writing into the current repo is the degenerate case where the hub repo points at it, so resolution stays uniform (one storage target, one read path). The durable shape of these blocks is owned by the [session-recap](../session-recap/README.md), [drift-recap](../drift-recap/README.md), and [journal-and-summary](../journal-and-summary/README.md) artifacts; Repo Config owns only the presence of the top-level keys and that unknown fields under them round-trip unchanged.

#### REQ: recaps-journal-layered-config-gated

Because `recaps.repo`, `recaps.user`, and `journal.repo` are per-user/per-machine values that MUST NOT be committed into a shared project `specscore.yaml`, the loader MUST reject them with a clear error pointing at the [**layered-config**](../layered-config/README.md) Feature (**Under Review; not yet implemented**) (`specscore.local.yaml` → `specscore.yaml` → `~/.specscore.yaml`) until that Feature ships. `recaps.enabled` / `journal.enabled` alone are accepted in the project file. This is the same posture the activity-journal cross-repo key carries.

### Unknown fields

Orchestration tools (e.g., Synchestra) MAY extend `specscore.yaml` with additional fields at any level. SpecScore tooling preserves them.

```yaml
# Example: orchestrator extension
state_repo: https://github.com/org/project-state
planning:
  auto_create: false
```

#### REQ: unknown-fields-preserved

SpecScore MUST ignore unknown fields at the root, inside `project:`, inside any `project.repositories` entry, and inside any module entry. Unknown fields MUST round-trip unchanged on read/write. They MUST NOT cause a validation warning or error. Unknown fields under `publication:` MUST also round-trip unchanged per [Publication Policy Config](../publication-policy-config/README.md).

### Example

A maximal `specscore.yaml` exercising every documented field:

```yaml
# SpecScore Repo Config Schema: https://specscore.md/repo-config

project:
  title: Acme Service
  host: github.com
  org: acme
  repo: service
  repositories:
    - url: https://github.com/acme/service-api
      title: Service API
      roles: [code]
    - url: https://github.com/acme/service-web
      title: Service Web
      roles: [code]
    - url: https://github.com/acme/service-spec
      title: Service Spec
      comment: SpecScore-managed spec repo for the service
      roles: [specification, code]
    - url: https://github.com/acme/service-state
      title: Service State
      roles: [state]

projects:
  - https://github.com/acme/platform
  - ./packages/shared

specs_dir_name: specs
docs_dir_name: docs

studio:
  name: SpecScore.Studio
  url: https://specscore.studio/

grade:
  values: [A, B, C, D, F]

publication:
  default:
    actions: [stage]
  events:
    idea.approved:
      actions: [stage, commit]
  push:
    allow_branches: ["feature/*"]
    deny_branches: ["main", "release/*"]

modules:
  - name: Highlevel
  - name: Backend
    path: backend
    code:
      - pkg/
      - go.mod
      - go.sum
  - path: frontend
```

## Dependencies

- [feature](../feature/README.md)
- [source-references](../source-references/README.md)

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Source References](../source-references/README.md) | `project.host`, `project.org`, `project.repo` provide overrides for git-remote inference when expanding same-repo references. |
| [Feature](../feature/README.md) | Feature READMEs carry the studio toolbar (see [studio-toolbar](../studio-toolbar/README.md)) rendered using `studio.name` and `studio.url` from this block. |
| [Studio Toolbar](../studio-toolbar/README.md) | The `studio:` block defined here provides `studio.name` and `studio.url` consumed by the studio-toolbar Feature's renderer and lint rule. |
| [Document Types Registry](../document-types-registry/README.md) | This feature is registered in the canonical document-types table; its consumer path is `specscore.yaml`. |
| [Grade body-metadata field](../canonical-grade-metadata-field/README.md) | The optional `grade:` block defined here carries `grade.values`; that Feature owns the field's parsing, placement, default, and lint semantics. |
| [Publication Policy Config](../publication-policy-config/README.md) | The optional `publication:` block defined here carries project-level publication policy; that Feature owns policy shape, action validation, command/event scopes, branch rules, and user-config parity. |
| [Layered Config](../layered-config/README.md) | Supplies local, repository, organization, and user origins and enforces that `repo_checkouts` never enters committed project config. |
| [Plan](../plan/README.md) | Consumes the resolved destination and source identity to select the authoritative Plan namespace while leaving source Features read-only. |

## Acceptance Criteria

### AC: file-recognized

**Requirements:** repo-config#req:file-name, repo-config#req:schema-header-comment, repo-config#req:empty-config-valid

A file named `specscore.yaml` at the repository root with the exact schema-header comment on line 1 is recognized as a SpecScore repo config. A file with any other name, or with a missing or misplaced header comment, is not.

### AC: identity-inferred

**Requirements:** repo-config#req:project-block-optional, repo-config#req:project-title-default, repo-config#req:source-reference-overrides, repo-config#req:project-repositories

Project identity (title, host, org, repo) is inferred from the working directory and git remote when fields are omitted. Explicit values override inference. `project.repositories` round-trips without being interpreted by SpecScore tooling.

### AC: repositories-role-tagged

**Requirements:** repo-config#req:repositories-entry-shape, repo-config#req:repositories-roles-list, repo-config#req:repositories-roles-enum

**Given** a `specscore.yaml` whose `project.repositories` is a list of role-tagged object entries (including at least one entry with multi-valued `roles`, e.g. `[specification, code]`, and at least one optional `title` / `comment` field)
**When** tools load the config and round-trip it back to disk
**Then** the on-disk bytes after the round-trip are identical to the input for every well-formed entry, every entry's `url`, `title`, `comment`, and `roles` list (including multi-valued combinations) is preserved, and any unknown fields inside an entry survive unchanged.

### AC: repositories-shape-errors

**Requirements:** repo-config#req:repositories-entry-shape, repo-config#req:repositories-roles-list, repo-config#req:repositories-roles-enum

**Given** a `specscore.yaml` whose `project.repositories` contains any of the following malformed entries — a flat-string entry (e.g. `- https://example.com/repo`), an entry missing `url`, an entry with `roles:` omitted, an entry with `roles: []`, an entry with a scalar `roles: code` instead of a list, or an entry containing an unknown role value (e.g. `roles: [helm-chart]`)
**When** tools load the config
**Then** the load fails with a hard error that names the offending entry by its index and the violated REQ (one of `repositories-entry-shape`, `repositories-roles-list`, `repositories-roles-enum`); no implicit defaults are applied; the file is not silently rewritten.

### AC: related-projects-validated

**Requirements:** repo-config#req:projects-list, repo-config#req:projects-local-path-must-resolve

Each `projects:` entry is recognized as a URL or a local directory path. Local paths are validated to exist and to contain a nested `specscore.yaml`; entries that fail validation are reported as errors.

### AC: plan-route-resolved

**Requirements:** repo-config#req:plans-repo-project-selection, repo-config#req:plan-repos-aggregate-routing, repo-config#req:plan-route-required, repo-config#req:plan-config-precedence

A local or committed scalar `plans_repo` selects one Plan destination and overrides organization and user mappings. Without either form of explicit routing, every Plan operation fails before Plan access, while unrelated artifact operations continue normally. Shorthand and full repository identities normalize to the same host/owner/repo identity, and conflicting mappings fail.

### AC: plan-checkout-separated

**Requirements:** repo-config#req:repo-checkouts-machine-local

An external Plan destination resolves only through an absolute machine-local `repo_checkouts` entry whose checkout root and origin match the destination. Committed checkout paths, relative or nested paths, wrong origins, and unresolved symlinks fail without touching either repository.

### AC: modules-resolved

**Requirements:** repo-config#req:modules-default, repo-config#req:module-name-deduction, repo-config#req:module-path-resolution, repo-config#req:module-code-list, repo-config#req:module-paths-unique, repo-config#req:module-paths-non-nested

Modules resolve their specs, docs, and code paths relative to `module.path`. Module names are deduced from `path` when omitted. Duplicate or nested explicit module paths produce hard errors. The implicit-root module (no `path:`) coexists with explicit-path modules.

### AC: studio-rules

**Requirements:** repo-config#req:studio-default-when-omitted, repo-config#req:studio-explicit-values, repo-config#req:studio-null-opts-out, repo-config#req:studio-toolbar-mandatory-unless-opted-out, repo-config#req:studio-url-trailing-slash

**Given** a `specscore.yaml` with any of the following `studio:` configurations — the block omitted entirely, the block present as a mapping with both `name` and `url`, the block present as a partial mapping (only `name` or only `url`), the block set to `null` / `~` / empty, or the block present with a `studio.url` lacking or doubling its trailing slash
**When** tools load the config and render artifact documents
**Then** an omitted `studio:` block applies defaults (`name: SpecScore.Studio`, `url: https://specscore.studio/`) and the studio toolbar is rendered; a full mapping uses its explicit values and the toolbar is rendered; a partial mapping is a hard error and no defaults are applied; `studio: null` suppresses the studio toolbar entirely; and a `studio.url` without a single trailing `/` is a hard error. The studio toolbar is mandatory in artifact documents unless explicitly opted out via `studio: null`.

### AC: extensions-preserved

**Requirements:** repo-config#req:unknown-fields-preserved

Unknown fields at any level survive read/write without warnings. Orchestrators can extend the config without breaking SpecScore validation.

### AC: publication-block-delegated

**Requirements:** repo-config#req:publication-block-optional, repo-config#req:publication-schema-delegated

**Given** a `specscore.yaml` with a top-level `publication:` block containing event policy, command policy, push branch rules, and an unknown future field
**When** SpecScore tooling reads and writes an unrelated config key
**Then** the `publication:` block remains a top-level project config block, its known fields are available to consumers under the Publication Policy Config schema, and the unknown future field round-trips unchanged.

### AC: recaps-journal-storage-target-gated

**Requirements:** repo-config#req:recaps-journal-storage-target, repo-config#req:recaps-journal-layered-config-gated

**Given** a project `specscore.yaml` with a `recaps:` block setting `enabled: true`, and a second one additionally setting `recaps.repo` while the layered-config Feature has not shipped
**When** the config loader runs
**Then** `recaps.enabled` is accepted (storage always targets the hub repo, which may be configured to equal the current code repo), and `recaps.repo` (likewise `recaps.user` / `journal.repo`) is rejected with a clear error pointing at the layered-config dependency; unknown fields under the blocks round-trip unchanged.

### AC: path-overrides-loaded

**Requirements:** repo-config#req:path-overrides-optional, repo-config#req:ideas-path-field

**Given** a `specscore.yaml` with one module declaring `path_overrides.ideas_path: ideas` plus an unknown override key, one module with no `path_overrides`, and one module declaring `path_overrides.ideas_path: /ideas` (absolute)
**When** the config loader runs
**Then** the first module's `ideas_path` is read as `ideas` and its unknown key round-trips unchanged, the second module defaults to `spec/ideas`, and the third fails with a hard error naming that module and the absolute-path violation.

## Open Questions

- The schema-pointer URL `https://specscore.md/repo-config` does not follow the registry convention `https://specscore.md/{type}-specification` used by other Document-Kind features. Should it be renamed to `https://specscore.md/repo-config-specification` for consistency, or is the shorter form preferred for human typing in YAML headers?
- Should `code:` entries support glob patterns (e.g., `pkg/**/*.go`), or must they be literal file/directory paths? Defer until lint actually consumes them.
- Should an explicit `studio.name` be allowed to differ from the host implied by `studio.url` (e.g., `name: MyCompanyDocs` with a SpecScore.Studio-hosted URL)? Spec currently allows it freely.
- How should tooling handle a repository that has no `specscore.yaml` at all — refuse to operate, or assume defaults?
- Should organization configuration eventually have an explicit discovery registry instead of being anchored inside the organization directory containing the canonical source clone? Current behavior uses the canonical-clone anchor so worktrees agree.

---
*This document follows the https://specscore.md/feature-specification*
