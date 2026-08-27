---
format: https://specscore.md/plan-specification
status: Approved
---

# Plan: Grade body-metadata field

**Status:** Approved
**Source Feature:** canonical-grade-metadata-field
**Date:** 2026-05-29
**Owner:** alex
**Supersedes:** —

## Summary

Implements the `canonical-grade-metadata-field` Feature: an optional `**Grade:**` body-metadata line on any header-block artifact, validated by `specscore spec lint` against a value set that defaults to `A, B, C, D, F` and is configurable via `grade.values` in `specscore.yaml`. Work lands primarily in `specscore-cli` (config parsing, header-block parsing, lint rule) plus a canonical-doc update in this spec repo.

## Approach

Grouped by implementation layer in dependency order: the config layer first (resolves the effective value set everything downstream consumes), then the canonical schema definition and header-block parsing (recognizes and structurally validates the field), then value validation against the effective set, then the cross-kind/decoupling guarantees that constrain when the rule fires. The canonical-doc update (publishing `**Grade:**` and `grade.values` into `feature-specification.md` and related docs) has no dedicated AC — it is the means by which the field becomes canonical — so it is folded into Task 2 as supporting work rather than a standalone AC-less task. No ACs are deferred.

## Tasks

### Task 1: Config layer — `grade.values` parsing and effective-set resolution

**Verifies:** canonical-grade-metadata-field#ac:grade-values-shape-errors

In the `specscore-cli` config loader, parse the optional `grade:` block with a `values:` key. Resolve the effective value set: use `grade.values` when present, otherwise the built-in default `A, B, C, D, F`. Validate shape — `grade.values` must be a non-empty YAML list of non-empty scalar tokens; an empty list, a scalar, or a list with an empty/non-scalar entry is a hard error that names the violated `grade-values-shape` rule and substitutes no implicit set. Duplicate tokens emit an advisory and are de-duplicated.

### Task 2: Canonical schema definition and header-block parsing of `**Grade:**`

**Verifies:** canonical-grade-metadata-field#ac:grade-absent-is-valid, canonical-grade-metadata-field#ac:grade-single-token-enforced, canonical-grade-metadata-field#ac:grade-placement-enforced

Define `**Grade:**` in the canonical Feature specification (`feature-specification.md` and the related published specification docs) as an optional body-metadata line. In the `specscore-cli` header-block parser: recognize an optional `**Grade:**` line (absence produces no warning or error); enforce exactly one non-empty token (empty or multi-token value is a hard error); and enforce placement as the last line of the contiguous header block, after `**Status:**` and any `**Source Ideas:**` / `**Supersedes:**` lines. A misplaced or out-of-block `**Grade:**` is a hard error, and `specscore spec lint --fix` normalizes a misplaced line to the canonical last position.

### Task 3: Lint value validation against the effective set

**Verifies:** canonical-grade-metadata-field#ac:grade-default-scale-validated, canonical-grade-metadata-field#ac:grade-custom-scale-validated

Add the lint check that validates a parsed `**Grade:**` value against the effective value set resolved in Task 1. A value in the set passes; a value outside it is a hard error naming the artifact, the offending value, and the effective set. Cover both paths: the built-in default `A, B, C, D, F` when no `grade:` block is configured, and a repo-configured set (e.g. `grade.values: [1, 2, 3, 4, 5]`).

### Task 4: Apply the rule kind-agnostically and decouple it from Status and workflow

**Verifies:** canonical-grade-metadata-field#ac:grade-on-any-header-block-kind, canonical-grade-metadata-field#ac:grade-decoupled-from-workflow-and-status

Wire the Grade lint rule so it applies uniformly to every artifact kind that has a header block (Feature, Plan, Idea, Decision, Task) with no per-kind allow-list or exemption, reusing the shared header-block parser rather than re-deriving placement logic. Ensure the rule never gates on `Status` or on any reviewer-gate/workflow artifact: a valid grade on a `Draft` artifact in a repository that declares no reviewer gates passes validation.

## Open Questions

None at this time. (Resolved: the built-in default value set is `A, B, C, D, F` — no `E` — following the academic letter-grade convention; repos needing a different scale declare `grade.values`.)

---
*This document follows the https://specscore.md/plan-specification*
