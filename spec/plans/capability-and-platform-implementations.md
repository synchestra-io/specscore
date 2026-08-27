---
format: https://specscore.md/plan-specification
status: Approved
---

# Plan: Capability and Platform Implementations

**Status:** Approved
**Source Feature:** capability-and-platform-implementations
**Date:** 2026-06-03
**Owner:** alexander.trakhimenok
**Supersedes:** —

## Summary

Implements the `capability-and-platform-implementations` Feature: Capability/Implementation role classification, the `**Implements:**` reference (reusing source-references notation, including the cross-repo `@host/org/repo` suffix), and `specscore spec lint` validation of the author-declared Implementation Matrix. Work lands entirely in `specscore-cli` (header-block parsing, reference resolution, and new lint rules); the canonical convention itself is already published by the source Feature.

## Approach

Grouped by implementation layer in dependency order. Role classification lands first because every downstream rule needs to know whether a Feature is a Capability (has an `## Implementation Matrix`) or an Implementation (has `**Implements:**`). Reference handling follows in two tasks — first accepting well-formed same-repo and cross-repo references, then the error paths (unresolved, and resolves-to-non-Capability) which depend on both classification and parsing. Matrix validation is its own group: structural shape/vocabulary/brief checks, then the deliberate no-rollup guarantee under `lint --fix`. The `## Other Platforms` no-duplication rule is independent and small. A final integration task exercises the single-repo co-located layout end to end. No ACs are deferred.

## Tasks

### Task 1: Role classification — Capability vs Implementation

**Verifies:** capability-and-platform-implementations#ac:capability-marker, capability-and-platform-implementations#ac:implementation-marker

In `specscore-cli`, classify a Feature by structural markers: a Feature containing an `## Implementation Matrix` section is a Capability; a Feature whose header block contains an `**Implements:**` line is an Implementation. Expose this classification to the lint layer so downstream rules can branch on role. The two roles are independent flags (a Feature is normally one or the other, but classification reports each marker on its own merits).

### Task 2: Parse and resolve well-formed `**Implements:**` references

**Verifies:** capability-and-platform-implementations#ac:implements-single-same-repo, capability-and-platform-implementations#ac:implements-cross-repo-suffix

Parse the `**Implements:**` header-block line using the existing source-references `specscore:` grammar. Accept a same-repo reference (e.g. `specscore:feature/dashboards`) that resolves to an existing Capability with no violation, and accept a well-formed cross-repo reference carrying the `@{host}/{org}/{repo}` suffix as valid without attempting to fetch the remote repo. Reuse the source-references parser rather than duplicating notation logic.

### Task 3: `**Implements:**` resolution errors

**Verifies:** capability-and-platform-implementations#ac:implements-unresolved-same-repo, capability-and-platform-implementations#ac:implements-non-capability

Add the lint errors for invalid same-repo `**Implements:**` targets: a reference that resolves to nothing is an error ("does not resolve"), and a reference that resolves to an existing Feature lacking an `## Implementation Matrix` section is an error ("target is not a Capability"). Cross-repo references are exempt from these resolution checks (liveness is out of scope); only their notation well-formedness, validated in Task 2, applies.

### Task 4: Implementation Matrix structural validation

**Verifies:** capability-and-platform-implementations#ac:matrix-shape, capability-and-platform-implementations#ac:matrix-bad-status, capability-and-platform-implementations#ac:matrix-brief-single-line

Add lint validation of a Capability's `## Implementation Matrix` table: require the columns `Platform`, `Status`, `Brief`, `Link` (a missing required column is an error); constrain each `Status` cell to the vocabulary `Full | Partial | Planned | Absent` (any other value is an error naming the allowed set); and enforce that each `Brief` cell is a single line (multi-line behavior in a Brief is an error).

### Task 5: Author-declared matrix — no rollup under `lint --fix`

**Verifies:** capability-and-platform-implementations#ac:matrix-no-rollup

Guarantee that `specscore spec lint --fix` validates the matrix shape (Task 4) but never mutates `Status` cells — no rollup from same-repo Implementation lifecycle, no auto-population, no auto-correction. Add a regression test where a same-repo Implementation's lifecycle differs from the declared parity status and assert `--fix` leaves the matrix unchanged.

### Task 6: `## Other Platforms` links-only enforcement

**Verifies:** capability-and-platform-implementations#ac:other-platforms-links-only

Add the lint rule that an Implementation's optional `## Other Platforms` section may carry links but MUST NOT restate a parity status; a status token appearing in that section is an error directing the author to the Capability's Implementation Matrix as the single source of parity status.

### Task 7: Single-repo co-located layout — integration

**Verifies:** capability-and-platform-implementations#ac:colocated-layout-valid

Add an integration test for the single-repository special case: a Capability and its CLI and Web Implementations co-located under one feature folder. Assert `specscore spec lint` reports no violations and that each Implementation's same-repo `**Implements:**` reference resolves to the co-located Capability, confirming Tasks 1–4 compose correctly in the co-located layout.

## Open Questions

- The three cross-repo open questions carried from the Feature (reconciliation mechanism, Capability repo ownership, version-pinning) are out of this Plan's scope; they remain Feature-level and would seed a follow-on Plan if pursued.

---
*This document follows the https://specscore.md/plan-specification*
