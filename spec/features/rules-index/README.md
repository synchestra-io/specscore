---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: Rules Index

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/specscore/specscore/spec/features/rules-index?op=explore) | [Edit](https://specscore.studio/app/github.com/specscore/specscore/spec/features/rules-index?op=edit) | [Ask question](https://specscore.studio/app/github.com/specscore/specscore/spec/features/rules-index?op=ask) | [Request change](https://specscore.studio/app/github.com/specscore/specscore/spec/features/rules-index?op=request-change) |

**Status:** Approved
**Source Ideas:** —

## Summary

The Rules Index at `spec/rules/README.md` is unusual among SpecScore indexes: it is both the inventory *and* the primary artifact. An inline Rule lives in exactly one row here and nowhere else, and for a detailed Rule this row remains authoritative for every field it carries.

## Problem

Most indexes are a derived view; losing one costs a rebuild. This one holds content that exists nowhere else, which changes what correctness means for it in two ways.

First, a row that cannot be parsed cannot be discarded. The most ordinary authoring slip there is — an unescaped `|` pasted into a Statement — would otherwise make a Rule disappear from a kind whose entire purpose is that operating knowledge stops evaporating.

Second, because a detailed Rule repeats the row's fields in its document, the two can disagree. Without one authoritative side, one Rule becomes two, each convincing to whoever opened it first.

## Behavior

This feature inherits location, section, and footer rules from [Index](../index/README.md). It defines the Rules-specific columns, the identity-cell grammar, and the preservation guarantee.

### Index columns

#### REQ: rule-index-columns

The `## Rules` table MUST carry exactly this header:

```
| Rule | Status | Scope | Enforcement | Control | Sources | Statement |
```

Every data row MUST have exactly seven cells. `Status`, `Scope`, `Enforcement`, `Control` and `Sources` follow the grammars in [Rule](../rule/README.md); `Statement` is the Rule's one normative sentence. An empty `Control` or `Sources` cell carries the em-dash sentinel `—`.

#### REQ: identity-cell-grammar

The first cell MUST be either a bare canonical slug — an *inline* Rule — or `[<slug>](<slug>/README.md)` — a *detailed* Rule whose document exists at that path. No other shape is valid.

The identity cell is the one place a guess would silently rename a Rule, so it is never inferred. The presence of a link is also the index's own statement about which Rules have more to read, which is what lets a reader skim the table and know where to open.

#### REQ: escaped-free-text

A pipe inside a free-text cell (`Control`, `Sources`, `Statement`) MUST be written escaped as `\|`, and MUST survive a read/write round trip as one cell. A hand-wrapped value MUST reach its row joined into one line: reading only a first physical line truncates a Statement mid-sentence, which is how an index row becomes a confidently wrong instruction.

### Ordering and uniqueness

#### REQ: sorted-and-unique

Rows MUST be unique by slug and sorted by slug, so a reviewer sees a stable diff and a reader can find a Rule by scanning.

### The index row is the source of truth

#### REQ: index-row-is-authoritative

For a detailed Rule, this row is authoritative for `Status`, `Statement`, `Scope`, `Enforcement`, `Control` and `Sources`. A tool that repairs a disagreement MUST rewrite the *document* from the row and MUST NEVER rewrite the row from the document.

The single exception is a detail document with no row at all: there, and only there, a document may seed a row — because the alternative is an artifact invisible to every reader of the index.

### An unparseable row is never dropped

#### REQ: preserve-unparseable-rows

A row-like line the seven-column contract cannot represent MUST be preserved verbatim by every write path and reported. No operation may reduce the Rule set by writing this file.

A tool MAY repair such a line only when the repair is unambiguous — surplus cells caused by an unescaped `|` in the `Statement`, which is provable only because the four columns between the identity cell and the Statement have closed grammars that must all validate first. The same line shape caused by a pipe inside `Control` MUST be refused rather than silently reinterpreted. Everything else is carried through and left to a human.

#### REQ: differing-duplicates-are-kept

A duplicated slug MAY be deduplicated only when the rows are byte-identical. Two rows that disagree are ambiguous, so both MUST be kept and reported: choosing one would discard a Rule's content on a coin flip.

### Worked example

```markdown
---
format: https://specscore.md/rules-index-specification
---

# Rules

Normative one-sentence rules with their scope, enforcing control, and sources.

## Rules

| Rule | Status | Scope | Enforcement | Control | Sources | Statement |
|---|---|---|---|---|---|---|
| deps-propagate-at-stream-end | Draft | fleet | Stated | — | https://example.com/wb/dependency-streams | Propagate a dependency bump in one wave at the end of a stream, and NEVER per tag mid-stream. |
| [firestore-no-tx-read-after-write](firestore-no-tx-read-after-write/README.md) | Active | fleet | Enforced | go-core >= v0.66.7 in-memory DB rejects a transaction read that follows a write | lesson:adapter-tests-disagreed-on-one-contract | NEVER read inside a Firestore transaction after writing in it; do every read first. |
| never-gc-apply-on-a-shared-vm | Draft | fleet | Stated | — | https://example.com/wb/decisions/0001 | NEVER run a garbage-collecting --apply on a shared VM until the liveness of every candidate has been measured. |

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/rules-index-specification*
```

The first and third rows are inline Rules — those two sentences exist nowhere else in the repository. The second is detailed: its name links to a document that carries the reason, worked examples and agent instructions, and whose header repeats these six cells.

### Adherence footer

#### REQ: adherence-footer

The index MUST end with the adherence footer naming `https://specscore.md/rules-index-specification`, matching its frontmatter `format:` field. As a status-less Index-Kind document it MUST NOT carry a frontmatter `status:` field.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Rule](../rule/README.md) | Defines the field grammars this index's columns carry, and the detail-document form a linked row points at. |
| [Index](../index/README.md) | Shared index shape: location, sections, completeness, footer. |
| [artifact-frontmatter-convention](../artifact-frontmatter-convention/README.md) | Governs the `format:` field and the status-less classification. |

## Acceptance Criteria

### AC: both-forms-share-one-table

**Requirements:** rules-index#req:rule-index-columns, rules-index#req:identity-cell-grammar

Given one inline Rule and one detailed Rule, when the index is read, then both appear as seven-cell rows in the same table, the inline one with a bare slug and the detailed one with a link to its document.

### AC: an-unparseable-row-survives-a-rewrite

**Requirements:** rules-index#req:preserve-unparseable-rows

Given an index carrying a hand-written row whose Statement holds an unescaped `|`, when any tool rewrites the index for an unrelated Rule, then that line is still present, byte-for-byte, and is reported.

### AC: only-the-unambiguous-repair-is-applied

**Requirements:** rules-index#req:preserve-unparseable-rows

Given that same row, when an automatic repair runs, then the surplus pipe is escaped and the Statement round-trips to its original text; given a row whose parse failure is not attributable to the Statement, the row is unchanged and still reported.

### AC: differing-duplicates-are-both-kept

**Requirements:** rules-index#req:differing-duplicates-are-kept

Given two rows for one slug that disagree, when an automatic repair runs, then both rows remain and the duplicate is still reported; two byte-identical rows are collapsed to one.

## Open Questions

- The preservation guarantee means an index can hold a line no tool understands, indefinitely, as long as it is reported. Is a machine-readable quarantine marker worth adding so a reader can tell a known-broken line from one nobody has looked at yet?

---
*This document follows the https://specscore.md/feature-specification*
