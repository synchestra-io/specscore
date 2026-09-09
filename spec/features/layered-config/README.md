---
format: https://specscore.md/feature-specification
status: In Review
---

# Feature: Layered Config

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/specscore/specscore/spec/features/layered-config?op=explore) | [Edit](https://specscore.studio/app/github.com/specscore/specscore/spec/features/layered-config?op=edit) | [Ask question](https://specscore.studio/app/github.com/specscore/specscore/spec/features/layered-config?op=ask) | [Request change](https://specscore.studio/app/github.com/specscore/specscore/spec/features/layered-config?op=request-change) |

**Status:** In Review
**Source Ideas:** —

## Summary

**Layered config** resolves SpecScore configuration from four origins — `specscore.local.yaml` (repo, uncommitted) → `specscore.yaml` (repo, committed) → organization `.specscore.yaml` → `~/.specscore.yaml` (user/machine) — with the most specific origin winning per key and maps deep-merged. It also introduces **machine-scoped keys** whose filesystem paths MUST NOT be committed into shared project config. Plan routing uses the same origins with field-specific shapes: local/project `plans_repo`, organization/user `plan_repos`, and machine-local `repo_checkouts`.

## Problem

Today SpecScore reads exactly one config file: the committed `specscore.yaml`. That is correct for project-shared settings but wrong for anything per-user or per-machine. A developer who wants their drift/verify/session recaps and activity journal aggregated into a personal hub repo (`~/work-log`, `~/specscore-hub`) has nowhere to put that path: committing it into the shared `specscore.yaml` is hostile to teammates, and there is no per-user override layer. Multiple Features ([journal-and-summary](../journal-and-summary/README.md) Phase 2, [session-recap](../session-recap/README.md), [repo-config](../repo-config/README.md)) already declare keys they must reject until a layered resolver exists — `recaps.repo`, `recaps.user`, `journal.repo`, `journal.stream`. This Feature provides the one resolver they all consume, plus the rule that keeps machine-scoped keys out of the shared file.

## Behavior

### Layer set and discovery

#### REQ: layer-set

Configuration MUST be resolved from four origins, in decreasing specificity:

1. **`specscore.local.yaml`** — at the repo root; per-developer-per-checkout; **uncommitted** (gitignored).
2. **`specscore.yaml`** — at the repo root; project-shared; **committed**.
3. **organization `.specscore.yaml`** — inside the organization directory containing the canonical source clone; shared machine policy for that organization's checkouts.
4. **`~/.specscore.yaml`** — in the user's home directory; per-user/per-machine defaults across all repos.

The repo root is the git top-level (`git rev-parse --show-toplevel`); when the working directory is not a git repo, the current directory is used for the repo-root layers. Linked worktrees discover the organization origin from the canonical clone rather than from the worktree's parent directory. Any layer file that is absent is treated as an empty layer (not an error).

#### REQ: precedence

Resolution MUST apply **most-specific-wins** per key: `specscore.local.yaml` overrides `specscore.yaml`, which overrides organization config, which overrides `~/.specscore.yaml`. A key set in a more specific origin fully determines that key's value (subject to merge semantics below). A Feature MAY restrict a key to a subset of origins; [Repo Config](../repo-config/README.md) does so for Plan routing and checkout paths.

#### REQ: merge-semantics

Merging across layers MUST be a **deep merge of mapping nodes**: nested maps combine key-by-key across layers. Scalars and sequences are **replaced wholesale** by the more specific layer (sequences are not concatenated). An explicit `null` in a more specific layer MUST override (clear) a value set by a less specific layer, distinct from the key being absent.

#### REQ: degrade-single-file

When only `specscore.yaml` is present (today's state), resolution MUST return exactly its contents — layered config is backward compatible and introduces no behavior change for repos that adopt neither the local nor the home layer.

### User-scoped keys

#### REQ: gated-key-registry

The resolver MUST support keys declared **machine-scoped**. A machine-scoped key carries a per-user/per-machine value (typically a filesystem path or identity) and MUST NOT appear in the committed `specscore.yaml`. Owning Features register which keys are machine-scoped; this Feature owns the mechanism, not the list. The initial registered set is `recaps.repo`, `recaps.user`, `journal.repo`, `journal.stream`, and `repo_checkouts`.

#### REQ: gate-committed-file

If a machine-scoped key is present in the committed `specscore.yaml`, the loader MUST reject it with a hard error that names the key, states it was found in the committed project file, and points the user to an allowed machine-local origin. The same key set in `specscore.local.yaml`, organization config, or `~/.specscore.yaml` MUST be accepted when its owning Feature allows that origin.

#### REQ: unblocks-gated-features

Once this Feature lands, the previously documented-but-rejected keys (`journal.repo`, `journal.stream`, `recaps.repo`, `recaps.user`) MUST resolve through the layered stack when set in an allowed layer — superseding the interim "rejected with a pointer to layered-config" posture those Features carry. The committed-file gate (above) remains in force.

### Local layer hygiene

#### REQ: local-gitignored

`specscore.local.yaml` MUST be gitignored. SpecScore MUST ensure this: when it scaffolds or first resolves config, it MUST add `specscore.local.yaml` to the repo's `.gitignore` if absent, and if `specscore.local.yaml` is already tracked by git it MUST warn that the file is committed (defeating its per-user purpose) with instructions to untrack it.

### Introspection

#### REQ: config-resolve-command

`specscore config show [--origin]` MUST print the fully resolved configuration. With `--origin`, it MUST annotate each key with the layer it was resolved from (`local` / `project` / `organization` / `home`), so precedence and gating outcomes are debuggable. `specscore config get <dotted.key>` MUST print a single resolved value (and its origin with `--origin`).

## Acceptance Criteria

### AC: four-origins-discovered

**Requirements:** layered-config#req:layer-set

**Given** a git repo containing `specscore.local.yaml` and `specscore.yaml` at its root, an organization `.specscore.yaml` inside the organization directory containing the canonical source clone, and a `~/.specscore.yaml` in the user's home
**When** SpecScore resolves config from anywhere inside the repo
**Then** all four files are discovered, linked worktrees use the canonical-clone organization origin, and a missing file is treated as an empty layer rather than erroring.

### AC: most-specific-wins

**Requirements:** layered-config#req:precedence

**Given** a key `studio.theme` set in user, organization, committed project, and local config
**When** config is resolved
**Then** resolution selects local, then project, then organization, then user as each more-specific value is removed.

### AC: deep-merge-and-null-override

**Requirements:** layered-config#req:merge-semantics

**Given** `~/.specscore.yaml` setting `journal: {enabled: true, aggregation_periods: [day, week]}` and `specscore.yaml` setting `journal: {aggregation_periods: [day, week, month]}`, plus a separate case where `specscore.local.yaml` sets `journal.enabled: null`
**When** config is resolved
**Then** the maps deep-merge to `{enabled: true, aggregation_periods: [day, week, month]}` (the sequence is replaced wholesale, not concatenated), and in the separate case the explicit `null` clears `enabled` rather than being ignored.

### AC: single-file-unchanged

**Requirements:** layered-config#req:degrade-single-file

**Given** a repo with only `specscore.yaml` (no local or home layer)
**When** config is resolved
**Then** the result equals the contents of `specscore.yaml` exactly — no behavior change for repos that adopt neither extra layer.

### AC: machine-scoped-key-rejected-in-committed-file

**Requirements:** layered-config#req:gated-key-registry, layered-config#req:gate-committed-file

**Given** `recaps.repo` or `repo_checkouts` (machine-scoped keys) set in the committed `specscore.yaml`
**When** the loader runs
**Then** it rejects the key with an error naming it, stating it was found in the committed project file, and pointing to an allowed machine-local origin; the same key set in local, organization, or user config is accepted when its owning Feature allows that origin.

### AC: gated-features-unblocked

**Requirements:** layered-config#req:unblocks-gated-features

**Given** `journal.repo` set in `~/.specscore.yaml` and `recaps.repo` set in `specscore.local.yaml`, with layered-config present
**When** the journal and recaps loaders resolve their config
**Then** both keys resolve through the layered stack (no longer rejected as "pending layered-config"), while a `journal.repo` placed in the committed `specscore.yaml` is still rejected by the committed-file gate.

### AC: local-file-gitignored

**Requirements:** layered-config#req:local-gitignored

**Given** a repo whose `.gitignore` does not list `specscore.local.yaml`
**When** SpecScore scaffolds or first resolves config
**Then** `specscore.local.yaml` is added to `.gitignore`; and given a repo where `specscore.local.yaml` is already git-tracked, SpecScore warns that the file is committed with instructions to untrack it.

### AC: config-show-reports-origins

**Requirements:** layered-config#req:config-resolve-command

**Given** a resolved config drawn from all four origins
**When** the user runs `specscore config show --origin` (and `specscore config get <dotted.key> --origin`)
**Then** each resolved key is annotated with its source origin (`local` / `project` / `organization` / `home`), and the single-key form prints that key's resolved value and origin.

## Rehearse Integration

Every AC above is testable through fixture config trees (combinations of `specscore.local.yaml` / `specscore.yaml` / `~/.specscore.yaml`) plus CLI invocation and resolved-value/origin inspection. Rehearse scenario stubs are deferred to the Plan phase to match the current repo convention.

## Not Doing

- **Per-block layering rules** (layering scoped to only the `journal:` or `recaps:` block) — rejected in the source Ideas; layering applies uniformly across all of `specscore.yaml` via one code path.
- **A system-wide layer** — local, project, organization, and user origins cover the required scopes; a host-global origin needs a separate use case.
- **Network-fetched or remote config layers** — all layers are local files.
- **Encryption / secret management for config values** — machine-scoped keys hold paths and identities, not secrets; secret handling is out of scope.
- **Auto-creating `~/.specscore.yaml` or `specscore.local.yaml`** — SpecScore reads them when present and ensures `.gitignore` hygiene; it does not fabricate them (the `specscore config` surface / init flows may scaffold templates separately).
- **Migrating existing `specscore.yaml` keys** — adopting layers is opt-in; nothing is moved automatically.

## Open Questions

- The exact registration mechanism by which an owning Feature declares a key machine-scoped (a static list in the loader vs. a per-key annotation in each Feature's config schema) — decide at plan time; the initial set is fixed (`recaps.repo`, `recaps.user`, `journal.repo`, `journal.stream`, `repo_checkouts`).
- Whether `~/.specscore.yaml` should itself be allowed to set machine-scoped keys for *all* repos at once (current design: yes — it is the per-user defaults layer) vs. requiring per-repo `specscore.local.yaml` for path values.
- Whether to also honor an `XDG_CONFIG_HOME`-based location (`~/.config/specscore/config.yaml`) as an alias for the home layer; defer unless requested.

---
*This document follows the https://specscore.md/feature-specification*
