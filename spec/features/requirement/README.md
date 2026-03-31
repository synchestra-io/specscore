# Feature: Requirement

**Status:** Conceptual

## Summary

A requirement is a discrete, testable rule or condition that a system must satisfy. Requirements live as named subsections within a feature's Behavior section — they are a naming convention, not a separate file artifact. Each requirement is addressable by ID, enabling traceability from acceptance criteria and scenarios back to the specific obligation they verify.

## Problem

SpecScore features describe behavior in prose Behavior sections. When a feature has many behavioral rules, individual obligations are not addressable — acceptance criteria cannot trace back to the specific rule they verify, and tooling cannot enumerate or lint individual requirements. This makes it hard to answer: "Which specific rule does this AC verify?" or "Are all behavioral rules covered by ACs?"

## Design Philosophy

Requirements are the **precision layer** within features. A feature's Behavior section explains *how something works* narratively; requirements mark the specific *rules* within that narrative that the system must satisfy.

Requirements are lightweight by design — a heading convention, not a new file type. This keeps them close to the behavioral context they formalize and avoids duplicating content across artifacts.

```mermaid
graph LR
    F["Feature<br/>(what)"]
    R["Requirement<br/>(rule)"]
    AC["Acceptance Criteria<br/>(condition)"]
    S["Scenario<br/>(example)"]

    F -->|contains| R
    R -->|verified by| AC
    AC -->|validated by| S
```

## Behavior

### Requirement format

Requirements live within a feature's `## Behavior` section, scoped under topic headings. Topic headings (`###`) provide narrative context; requirements (`#### REQ:`) state the enforceable rules within that topic:

```markdown
## Behavior

### Item management

Todos are created, edited, and deleted through the CLI.

#### REQ: title-required

A todo item MUST have a non-empty title. Creating a todo without a title MUST be rejected with an error message.

#### REQ: slug-format

Feature slugs MUST be lowercase, hyphen-separated, and URL-safe. Underscores, spaces, and special characters are not permitted.
```

The `REQ:` prefix distinguishes requirements from other subsections. Topic headings without the `REQ:` prefix are organizational — they group related requirements and provide context, but are not themselves requirements.

Requirements are always **one heading level below their containing topic**. In the typical case (`### topic` under `## Behavior`), requirements use `####`. The heading level may vary if the topic is nested deeper, but requirements are always direct children of their topic.

### Requirement identification

Requirements are identified by their feature path and slug:

```
{feature-path}#req:{slug}
```

| Feature path | Requirement slug | Full ID |
|---|---|---|
| `todo-item/manage` | `title-required` | `todo-item/manage#req:title-required` |
| `todo-item/completion` | `timestamp-on-complete` | `todo-item/completion#req:timestamp-on-complete` |
| `todo-list` | `default-filter-active` | `todo-list#req:default-filter-active` |

### Requirement slugs

Slugs follow the same rules as feature slugs:
- Lowercase letters, numbers, and hyphens only
- No underscores, spaces, or special characters
- URL-safe and path-safe

### RFC 2119 language

Requirements SHOULD use RFC 2119 keywords (MUST, MUST NOT, SHOULD, SHOULD NOT, MAY) to express obligation levels. This is a convention, not a strict rule — natural language is acceptable when the intent is clear.

| Keyword | Meaning |
|---|---|
| MUST / MUST NOT | Absolute requirement or prohibition |
| SHOULD / SHOULD NOT | Recommended but exceptions exist |
| MAY | Truly optional |

### Requirement granularity

Each requirement expresses a **single testable obligation**. If a requirement contains multiple independent conditions, split it into separate requirements.

**Good** — single obligation:
```markdown
#### REQ: title-required
A todo item MUST have a non-empty title.
```

**Bad** — multiple obligations bundled:
```markdown
#### REQ: title-rules
A todo item MUST have a non-empty title, MUST NOT exceed 200 characters,
and SHOULD be unique within the list.
```

Split the bad example into `title-required`, `title-max-length`, and `title-unique`.

### Referencing requirements from ACs and scenarios

Acceptance criteria and scenarios reference requirements using the `{feature-path}#req:{slug}` identifier.

**From an inline AC** (in the Acceptance Criteria section of a feature README):

```markdown
### AC: item-validation

**Requirements:** todo-item/manage#req:title-required, todo-item/manage#req:title-max-length

Creating a todo without a title or with a title exceeding the limit is rejected.
```

ACs are optional. They bundle related requirements into composite verification conditions. When an AC does not add value beyond the REQ itself, scenarios MAY reference the requirement directly.

**From a scenario** (in `_tests/`):

```markdown
**Validates:** [todo-item/manage#req:title-required](../README.md#req-title-required)
```

Or referencing a bundled AC:

```markdown
**Validates:** [todo-item/manage#ac:item-validation](../README.md#ac-item-validation)
```

### Parent features and requirements

A parent feature MAY define requirements that apply broadly to its sub-features. Sub-features define their own requirements for behavior specific to them. There is no inheritance — a sub-feature's ACs must explicitly reference the parent requirement ID if they verify a parent-level rule.

## Structural Rules

1. **Requirements are scoped under topic headings.** A topic heading provides narrative context; requirements are one heading level below (typically `#### REQ: {slug}` under a `###` topic within `## Behavior`).
2. **The `REQ:` prefix is case-sensitive** and followed by a space and the slug.
3. **Requirement slugs are unique within a feature.** Two requirements in the same feature cannot share a slug.
4. **Requirements live in Behavior sections only.** The `REQ:` convention is not valid outside `## Behavior`.
5. **Each requirement is a single testable obligation.** Multi-condition requirements should be split.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Feature](../feature/README.md) | Requirements live within a feature's Behavior section as `#### REQ:` subsections under topic headings. |
| [Acceptance Criteria](../acceptance-criteria/README.md) | ACs optionally bundle requirements via the `**Requirements:**` metadata field. |
| [Scenario](../scenario/README.md) | Scenarios validate requirements directly or through bundled ACs — completing the traceability chain. |

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Acceptance criteria not yet defined for this feature.
- Should tooling enforce that every requirement has at least one AC, or is it acceptable to have requirements without ACs (e.g., during early specification)?
- Should requirements support a status independent of their parent feature (e.g., a requirement could be marked `deprecated` while the feature remains `stable`)?
