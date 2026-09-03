---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: Rule

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/specscore/specscore/spec/features/rule?op=explore) | [Edit](https://specscore.studio/app/github.com/specscore/specscore/spec/features/rule?op=edit) | [Ask question](https://specscore.studio/app/github.com/specscore/specscore/spec/features/rule?op=ask) | [Request change](https://specscore.studio/app/github.com/specscore/specscore/spec/features/rule?op=request-change) |

**Status:** Approved
**Source Ideas:** —

## Summary

A Rule is one normative sentence — MUST or MUST NEVER, in plain words — carrying the scope it binds, the reason it exists, the sources that produced it, and the control that enforces it. It is the smallest kind in the spec tree and the only prescriptive one.

A Rule has **two forms and one identity**. An *inline* Rule is exactly one row in the [Rules Index](../rules-index/README.md) and nothing else. A *detailed* Rule keeps that identical row, linked to `spec/rules/<slug>/README.md`, which adds the reason, worked compliant and violating examples, agent instructions, exceptions, and supersession. The index row is the source of truth for every field it carries.

## Problem

Durable operating knowledge — *"never mock an extension backend"*, *"always format before building"* — accumulates in per-agent memory files and personal notes. Those do not transfer: a new session, machine, or teammate starts blind and re-learns the rule the expensive way, by repeating the mistake it encodes.

The other kinds each hold part of the answer and none of it whole. A [Lesson](../lesson/README.md) explains the *process gap* a defect exposed; a [Decision](../decision/README.md) explains *why a choice was made*; an [Idea](../idea/README.md) explains *a direction*. None of them is the single transferable sentence you can hand to a newcomer with no other context. A Lesson is a narrative about the past; a Rule is an instruction about the next commit.

The two forms exist because the two failure modes are opposite. Give every Rule a directory and recording one becomes ceremony, so nobody records the ninety per cent that really are one sentence. Give no Rule a document and the ten per cent that need a worked example — where the boundary falls, what a compliant change looks like — have nowhere to put it, and the Rule gets read three different ways. One row is the floor; a document is the opt-in ceiling.

## Behavior

### Location and identity

#### REQ: rule-location

An inline Rule exists ONLY as a row in `spec/rules/README.md`. A detailed Rule adds `spec/rules/<slug>/README.md` and its index row's identity cell becomes a link to it. No other location is a Rule.

#### REQ: rule-id

A Rule's identity is its slug: lowercase, hyphen-separated, URL-safe, matching `^[a-z0-9]+(-[a-z0-9]+)*$`. For a detailed Rule the slug is its directory name, and the document's first H1 MUST read `# Rule: <title>` with a non-empty title.

### The index row is the source of truth

#### REQ: row-is-authoritative

Every field the index row carries — `Status`, `Statement`, `Scope`, `Enforcement`, `Control`, `Sources` — is authoritative. A detail document repeats those six in its header for readability, and they MUST match. Repair flows one way only: from the row into the document, never the reverse.

Two representations of one Rule, each authoritative to whoever happened to open it, is precisely how one Rule quietly becomes two.

### Detail document shape

A detail document declares twelve bold metadata fields in this exact order, then three sections:

```markdown
---
format: https://specscore.md/rule-specification
status: Active
---

# Rule: <title>

**Status:** Active
**Date:** 2026-09-03
**Owner:** <person>
**Statement:** <one normative sentence, MUST or NEVER>
**Scope:** fleet
**Enforcement:** Enforced
**Control:** <the mechanism that refuses>
**Sources:** lesson:<slug>, decision:0012
**Why:** <the reason, plus one concrete use case>
**Exceptions:** <the escape hatch and who may use it, or "none">
**Supersedes:** —
**Superseded By:** —

## Instructions

<what an agent should actually do to comply>

## Examples

### Compliant

<a short worked example that satisfies the rule>

### Violation

<the closest thing that looks fine and still breaks it>

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/rule-specification*
```

#### REQ: detail-required-fields

All twelve fields MUST be present, in the order above, and none but `Scope` and `Sources` may repeat. `Status`, `Date`, `Owner`, `Statement`, `Enforcement`, `Why` and `Exceptions` MUST carry content. `Control`, `Sources`, `Supersedes` and `Superseded By` MUST carry the em-dash sentinel `—` when empty, so "absent" and "the author forgot" never look alike. `Date` MUST be `YYYY-MM-DD`.

#### REQ: detail-required-sections

`## Instructions`, `## Examples` and `## Open Questions` MUST be present, and `## Examples` MUST carry both `### Compliant` and `### Violation`. An example set showing only the happy path leaves the reader guessing at the boundary the Rule actually draws.

### Statement

#### REQ: statement-is-one-normative-sentence

`Statement` MUST be one sentence in the imperative, using MUST or NEVER, in plain words. A Rule that needs paragraphs to state is not yet a Rule — it is a Decision or an Idea that has not been distilled.

### Status ladder

#### REQ: rule-statuses

`Status` MUST be one of `Draft` (written down, not yet binding), `Active` (binding within its declared Scope), or `Superseded` (replaced by another Rule, named in `**Superseded By:**`).

`Superseded` is a detail-document status: supersession pointers live there, so an inline row at `Superseded` has nowhere to name its successor and MUST be reported.

### Scope grammar and matching

#### REQ: scope-grammar

`Scope` is a comma-separated list of one or more entries, each of which MUST be `fleet`, `product:<name>`, `repo:<owner>/<repository>`, or `path:<glob>`. A bare unprefixed token other than the `fleet` keyword MUST be rejected rather than guessed at: a mis-scoped Rule binds work it was never meant to bind, which is worse than one that fails to match, because nobody goes looking for it.

#### REQ: scope-matching

Resolving a Scope against a path follows this table:

| Scope | Matches |
|---|---|
| `fleet` | everything |
| `path:<glob>` | doublestar match against the path AND against every trailing suffix of it, so a repo-relative pattern still matches an absolute path a caller holds. Deliberately generous: `path:cli/**` also matches `vendor/x/cli/y.go`. |
| `product:<name>` | `<name>` appears as a whole path segment |
| `repo:<owner>/<name>` | `<owner>` and `<name>` appear as **consecutive** whole path segments. A bare `<name>` MUST NOT match. |

The repository rule is the strict one on purpose. Matching a bare repository name would make every Rule scoped to a repo called `docs`, `api`, `web` or `cli` bind every path in the fleet containing that directory — and would let a Rule scoped to one owner bind another's.

### Enforcement ladder and its control

#### REQ: enforcement-ladder

`Enforcement` MUST be one of:

| Tier | Meaning |
|---|---|
| `Stated` | An agent or human is told. Nothing refuses. |
| `Enforced` | A named control refuses. |
| `Automated` | A named control refuses and repairs, with no human in the loop. |

#### REQ: enforced-requires-a-control

`Enforced` and `Automated` MUST name a non-empty `Control` — a concrete mechanism such as a CLI verb, a hook profile, a CI check, a lint rule, or a review probe. `Stated` MUST NOT be required to name one, because the absence of a control is exactly what that tier means.

A Rule claiming a tier no mechanism backs is the worst failure in this kind: it reads as binding, so a reviewer stops looking for the gate, and no gate exists. **`Control` MUST describe what exists today, not what a specification promises.** A control that is designed but unshipped belongs in `Sources`, so the record shows where enforcement is coming from without claiming it has arrived.

### Sources

#### REQ: source-references

Each `Sources` entry MUST be `lesson:<slug>`, `decision:<NNNN>` or `decision:<NNNN-slug>`, `idea:<slug>`, or an `http(s)` URL. Typed references MUST resolve to an existing artifact in the same spec tree; a free URL is validated syntactically only, because resolving one would make validation depend on the network.

A Source records what *produced* the Rule, so a reader can always get back to *why* without the Rule restating it.

### Reciprocity with Lessons

#### REQ: lesson-rule-pair

A Lesson MAY carry an optional `**Promotes To:** rule:<slug>` field, placed after its canonical relation block so an existing Lesson stays valid without it. That pointer and the Rule's `lesson:<slug>` Source are one relation and MUST agree in both directions.

A Lesson carries exactly **one** promotion pointer, so it can be the Source of at most one Rule, while a Rule may cite several Lessons. That constraint is what keeps *"which Rule did this Lesson become?"* answerable rather than ambiguous. A Lesson that merely informed a Rule without producing it belongs in that Rule's `**Why:**`, not its `Sources`.

### Reciprocity with agent skills

#### REQ: rule-skill-pair

A detail document's `## Instructions` MAY reference `skill:<name>`, and an agent skill MAY list `rule:<slug>` under its own `## Rules` heading. Both directions MUST resolve and MUST reciprocate: a skill that silently outlives the Rule constraining it is the failure this pair catches.

Only references under a skill's `## Rules` heading count as declarations — a `rule:` token in prose is discussion. A skill may bind only a *detailed* Rule, because an inline Rule has nowhere to name the skill back.

### Supersession

#### REQ: supersession-is-asymmetric

`**Superseded By:**` points FORWARD to the Rule that replaced this one and MUST resolve to a Rule listed in the index; a retirement whose destination does not exist leaves a reader nowhere. `**Supersedes:**` points BACKWARD at what this Rule replaced, and that Rule may legitimately have been deleted — the backward pointer is often the only surviving record of a retired Rule — so an absent target is history, not a defect. A malformed slug is still a defect.

This asymmetry is deliberate. Requiring both to resolve would make keeping the record of a retired Rule illegal.

### Validation

#### REQ: rule-lint-family

A conforming tool MUST validate Rules with the `R-` family:

| Rule | Checks |
|---|---|
| `R-001` | detail-document field set, order, duplication, non-empty required values, `YYYY-MM-DD` Date, required sections and both example subsections, non-empty `# Rule:` title |
| `R-002` | `Status` is in the vocabulary |
| `R-003` | index table shape: canonical header, seven cells per row, unique and slug-sorted rows, a Statement on every row |
| `R-004` | a linked row has its detail document, and a detail document has its row |
| `R-005` | `Enforcement` is in the vocabulary, and `Enforced`/`Automated` name a `Control` |
| `R-006` | every `Scope` entry is unique and well-formed |
| `R-007` | every `Sources` entry is unique, well-formed, and — when typed — resolves |
| `R-008` | the lesson↔rule pair reciprocates in both directions |
| `R-009` | supersession resolves per REQ:supersession-is-asymmetric, is inverse-consistent and acyclic; a `Superseded` inline row is reported |
| `R-010` | the rule↔skill pair reciprocates in both directions |
| `R-011` | a detail document's mirrored header equals its index row |

### CLI behavior

#### REQ: rule-cli-commands

A conforming CLI SHOULD expose `rule new` (inline by default; `--detailed` opts in), `expand`, `list`, `show`, `update`, `delete`, `promote` and `lint`. Every verb SHOULD be non-interactive and accept `--format text|yaml|json`.

`rule list` SHOULD read only the index and never open a detail document, so it stays cheap enough to run at the start of every working session and print the Rules that apply to the files at hand.

#### REQ: rule-exit-codes

| Code | Condition |
|---|---|
| `0` | Success. A listing may be empty, and a plain `list`/`show` exits `0` even when a row carries a scope error. |
| `1` | Conflict: lint reported an error-severity violation; a create target already exists (including as an unparseable row) without `--force`; a mutating verb was asked to rewrite an index holding a row it cannot read; or a scope query could not evaluate a scope it needed. |
| `2` | Invalid arguments, including a control-requiring tier with no control and a typed Source that does not resolve. |
| `3` | The named Rule has no row at all. |
| `4` | Invalid state: live links block a delete; a Lesson already promotes elsewhere; a document-only edit or `Superseded` was requested on an inline Rule; or the addressed Rule's row exists but does not parse. |
| `10` | Unexpected I/O or parse failure. |

### Adherence footer

#### REQ: adherence-footer

A Rule detail document MUST end with the adherence footer naming `https://specscore.md/rule-specification`, matching its frontmatter `format:` field per the [artifact-frontmatter-convention](../artifact-frontmatter-convention/README.md).

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Rules Index](../rules-index/README.md) | Holds the authoritative row for every Rule, inline or detailed. |
| [Lesson](../lesson/README.md) | Promotes into a Rule; the pair is checked in both directions by `R-008`. |
| [Decision](../decision/README.md) | A common Source: a Decision explains a choice, a Rule states what to do about it. |
| [Source References](../source-references/README.md) | Shares the typed-reference grammar `Sources` uses. |
| [Status Vocabulary](../status-vocabulary/README.md) | The `Draft` / `Active` / `Superseded` ladder. |
| [artifact-frontmatter-convention](../artifact-frontmatter-convention/README.md) | Governs `format:`/`status:` frontmatter and the footer mirror. |

## Acceptance Criteria

### AC: inline-rule-needs-only-a-row

**Requirements:** rule#req:rule-location, rule#req:row-is-authoritative

Given a repository with a rules index, when an author records a Rule with only a slug and a statement, then one row appears in `spec/rules/README.md`, no directory is created, and validation passes.

### AC: enforced-without-a-control-is-refused

**Requirements:** rule#req:enforced-requires-a-control

Given a Rule whose `Enforcement` is `Enforced` and whose `Control` is the em-dash sentinel, when validation runs, then `R-005` reports it, stating that an enforced Rule with no control is a stated Rule wearing a stronger label.

### AC: repo-scope-does-not-bind-a-bare-name

**Requirements:** rule#req:scope-matching

Given a Rule scoped `repo:acme/docs`, when a scope query resolves it against `docs/x.md` or `otherorg/docs/x.md`, then it does not match; against `projects/acme/docs/x.md` it does.

### AC: the-document-mirrors-the-row

**Requirements:** rule#req:row-is-authoritative

Given a detail document whose `**Status:**` was hand-edited to disagree with its index row, when validation runs, then `R-011` names both values; an automatic repair rewrites the *document* from the row, leaves the row unchanged, and is idempotent.

### AC: the-lesson-pair-holds-both-ways

**Requirements:** rule#req:lesson-rule-pair

Given a Rule listing `lesson:l` and a Lesson `l` with no `**Promotes To:**`, when validation runs, then `R-008` reports it; the mirror case — a Lesson promoting to a Rule that does not cite it — is reported against the Lesson.

### AC: a-retired-rule-may-still-be-named

**Requirements:** rule#req:supersession-is-asymmetric

Given a detail document whose `**Supersedes:**` names a Rule that has been deleted, when validation runs, then no violation is reported for it, while an unresolvable `**Superseded By:**` still is.

## Open Questions

- A Lesson carries one promotion pointer, so it can Source at most one Rule. Should `**Promotes To:**` become a list once real usage shows a Lesson that genuinely produced two independent Rules?
- Supersession lives only in the detail document, so an inline Rule must be expanded before it can be retired. Should the row carry supersession too?
- `path:` globs and `product:`/`repo:` segment matching are deliberately generous and deliberately strict respectively. Is one explicit "anchored" glob form worth adding for a Rule that must bind exactly one directory?

---
*This document follows the https://specscore.md/feature-specification*
