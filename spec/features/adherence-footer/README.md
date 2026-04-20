# Feature: Adherence Footer

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fadherence-footer) — graph, discussions, approvals

**Status:** Draft
**Source Ideas:** adherence-footer-and-doc-type-registry

## Summary

The adherence footer is a single italic line at the bottom of every SpecScore document that declares — in one machine-verifiable, human-readable sentence — which specification format the document conforms to. It is the mechanism every Document-Kind feature uses so authors do not re-specify the rule per document type, and so lint can validate conformance from the footer alone.

This feature defines the mechanism once: what the footer looks like, where it lives, how it is matched by tooling, and the policy that footer URLs are unversioned. Each Document-Kind feature declares its own URL in a short local `REQ: adherence-footer` that delegates here.

## Problem

Without a shared mechanism, every Document-Kind feature re-specifies the footer rule in full — syntax, placement, bare-URL requirement, lint behavior, unversioned policy. The rule duplicates across `feature`, `plan`, `idea`, `task`, `scenario`, `project-definition`, and every Index-Kind feature. A policy change (e.g., "authors MAY reword the prose") requires editing N features. Worse, inconsistencies creep in — one feature insists the URL be bare while another quietly accepts Markdown links — and lint has no single definition to enforce against.

The adherence footer is one of the most common structural rules in SpecScore. It deserves its own Meta-Kind feature.

## Behavior

### Footer shape

Every SpecScore document of a Document Kind or Index Kind carries a single-line adherence footer at the very bottom of the file.

#### REQ: footer-shape

The footer MUST be a single italic line preceded by a horizontal rule (`---`) on the line above. The italic line MUST contain a bare URL of the form `https://specscore.md/{type}-specification`, where `{type}` is the document-type slug defined by the Document-Kind or Index-Kind feature.

The recommended form is:

```markdown
---
*This document follows the https://specscore.md/{type}-specification*
```

Authors MAY reword the surrounding prose — for example, "This index follows the …" on an Index-Kind document — but the URL MUST be present verbatim.

#### REQ: bare-url

The URL MUST appear in bare form. It MUST NOT be wrapped in Markdown link syntax (`[text](url)`, `<url>`, or reference-style links). The bare form keeps the URL clickable in terminals, diff viewers, and other renderers that detect URLs but do not parse Markdown.

#### REQ: trailing-slash

The footer URL MAY include or omit a trailing slash. `https://specscore.md/feature-specification` and `https://specscore.md/feature-specification/` are both valid. Lint treats them as equivalent.

### Placement

The footer MUST be the last meaningful content in the file. Nothing follows it except an optional trailing newline.

#### REQ: last-line

The adherence footer MUST be the final non-empty line of the document. Content appearing after the footer — additional paragraphs, a second footer, a signature — is a validation error.

### Per-document-type URL

Each Document-Kind or Index-Kind feature declares its consumer URL in a local `REQ: adherence-footer`. The body of that REQ MUST be a two-line form that delegates to this feature and declares the URL.

#### REQ: delegation-form

The `REQ: adherence-footer` body in a Document-Kind or Index-Kind feature MUST take the form:

```markdown
#### REQ: adherence-footer

Every {TYPE} document MUST end with an adherence footer per the [Adherence Footer feature](../adherence-footer/README.md). The footer URL MUST be `https://specscore.md/{type}-specification`.
```

The exact prose MAY vary, but the REQ MUST (a) reference this feature by relative link and (b) declare the exact URL for its document type. The URL MUST match the URL declared for that document type in the [Document Types Registry](../document-types-registry/README.md).

### Unversioned URL policy

Footer URLs are stable forever. There is no versioning scheme.

#### REQ: unversioned

The footer URL MUST NOT carry a version suffix (`/v1/`, `-v2`, `?version=`). SpecScore commits to additive-only evolution of document formats; breaking changes are not planned. If a breaking change ever becomes unavoidable, the canonical URL will be reissued rather than branched into parallel versions.

### Lint matching

Tooling that validates adherence footers operates on the URL, not the surrounding prose.

#### REQ: lint-matches-url

`specscore lint` MUST match footers by the presence of the exact URL string `https://specscore.md/{type}-specification` (trailing slash optional) anywhere in the document. The prose around the URL is not part of the match.

#### REQ: lint-severity

A missing adherence footer in a document of a Document or Index Kind MUST be a lint error. A footer whose URL does not match the document's expected URL (determined by the document's Kind and Consumer Path via the [Document Types Registry](../document-types-registry/README.md)) MUST be a lint error.

#### REQ: fix-inserts-only

`specscore lint --fix` MUST insert a missing adherence footer when the document has none. `--fix` MUST NOT rewrite a footer whose URL is wrong — a wrong URL may indicate the document has been mis-classified, and auto-rewriting would mask the underlying bug. Wrong-URL violations remain hard errors and require author action.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Feature](../feature/README.md) | Feature READMEs carry an adherence footer per this mechanism. The Feature feature's local `REQ: adherence-footer` declares the URL `https://specscore.md/feature-specification`. |
| [Document Types Registry](../document-types-registry/README.md) | The registry is the canonical source of each document type's footer URL. Lint cross-checks that every Document-Kind feature's local `REQ: adherence-footer` URL matches its registry row. |
| [Plan](../plan/README.md), [Idea](../idea/README.md), [Task](../task/README.md), [Scenario](../scenario/README.md), [Project Definition](../project-definition/README.md), [Plans Index](../plans-index/README.md) | Each declares its consumer URL via a local two-line `REQ: adherence-footer` that delegates here. |

## Acceptance Criteria

### AC: footer-shape-valid

**Requirements:** adherence-footer#req:footer-shape, adherence-footer#req:bare-url, adherence-footer#req:last-line

A valid adherence footer is a single italic line preceded by `---`, containing the bare specification URL, positioned as the final non-empty line of the document. A footer violating any of these is rejected by lint.

### AC: delegation-form-enforced

**Requirements:** adherence-footer#req:delegation-form

Every Document-Kind and Index-Kind feature's local `REQ: adherence-footer` uses the two-line delegation form that references this feature and declares the exact URL matching the Document Types Registry. Drift between the declared URL and the registry row is a lint error.

### AC: lint-behavior

**Requirements:** adherence-footer#req:lint-matches-url, adherence-footer#req:lint-severity, adherence-footer#req:fix-inserts-only

Lint matches footers by URL string (not prose), errors on missing or mismatched footers, and only inserts missing footers on `--fix` — never rewriting wrong URLs.

## Outstanding Questions

- Should the two-line delegation form be template-enforced by lint (exact regex match on the REQ body), or left as a convention that reviewers police?
- Should the `trailing-slash` allowance be narrowed over time in favor of a single canonical form, or is the tolerant match good forever?
- When `specscore lint --fix` inserts a missing footer, how does it decide which URL to use? Presumably from the document's Kind and Consumer Path in the registry — but this needs to be stated explicitly and tested.

---
*This document follows the https://specscore.md/feature-specification*
