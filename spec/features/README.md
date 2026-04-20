# SpecScore Features

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures) — graph, discussions, approvals

Core features of the SpecScore specification framework. This table is the canonical [Document Types Registry](document-types-registry/README.md) — every SpecScore document type is listed here with its Kind, specification URL, and Consumer Path.

| Feature | Status | Kind | URL | Consumer Path | Index | Description |
|---------|--------|------|-----|---------------|-------|-------------|
| [idea](idea/README.md) | Conceptual | Document | `https://specscore.md/idea-specification` | `spec/ideas/*.md` | [ideas-index](ideas-index/README.md) | Pre-spec one-pager that refines a concept into a promotable artifact |
| [feature](feature/README.md) | Stable | Document | `https://specscore.md/feature-specification` | `spec/features/**/README.md` | [features-index](features-index/README.md) | Feature specification structure, metadata, lifecycle, and conventions |
| [acceptance-criteria](acceptance-criteria/README.md) | Stable | Structure | `https://specscore.md/acceptance-criteria-specification` | — | — | Acceptance criteria format and conventions |
| [requirement](requirement/README.md) | Stable | Structure | `https://specscore.md/requirement-specification` | — | — | Discrete testable rules within feature Behavior sections |
| [scenario](scenario/README.md) | Stable | Document | `https://specscore.md/scenario-specification` | `spec/features/**/_tests/*.md` | [scenarios-index](scenarios-index/README.md) | Concrete Given/When/Then behavior examples in `_tests/` directories |
| [source-references](source-references/README.md) | Stable | Structure | `https://specscore.md/source-references-specification` | — | — | Code-to-spec traceability via inline annotations |
| [plan](plan/README.md) | Stable | Document | `https://specscore.md/plan-specification` | `spec/plans/**/README.md` | [plans-index](plans-index/README.md) | Planning documents that bridge specs to execution |
| [plans-index](plans-index/README.md) | Draft | Index | `https://specscore.md/plans-index-specification` | `spec/plans/README.md` | — | Canonical index of all Plan documents in a repo |
| [ideas-index](ideas-index/README.md) | Draft | Index | `https://specscore.md/ideas-index-specification` | `spec/ideas/README.md` | — | Canonical index of all active Idea documents in a repo |
| [features-index](features-index/README.md) | Draft | Index | `https://specscore.md/features-index-specification` | `spec/features/README.md` | — | Canonical index of all top-level features in a repo |
| [scenarios-index](scenarios-index/README.md) | Draft | Index | `https://specscore.md/scenarios-index-specification` | `spec/features/**/_tests/README.md` | — | Per-feature index of scenarios inside `_tests/` directories |
| [task](task/README.md) | Stable | Document | `https://specscore.md/task-specification` | `spec/plans/**/tasks/*.md` | — | Discrete units of work within a plan |
| [project-definition](project-definition/README.md) | Stable | Document | `https://specscore.md/project-definition-specification` | `specscore-spec-repo.yaml` | — | SpecScore project configuration and root structure |
| [adherence-footer](adherence-footer/README.md) | Draft | Meta | — | — | — | The shared footer mechanism every Document-Kind feature delegates to |
| [document-types-registry](document-types-registry/README.md) | Draft | Meta | — | — | — | This registry — canonical list of SpecScore document types |
| [index](index/README.md) | Draft | Meta | — | — | — | Shared shape of every Index-Kind feature — required sections, completeness, footer delegation |

## Feature Hierarchy

```
spec/features/
├── idea/                      # How to structure pre-spec ideation artifacts
├── feature/                   # How to structure and write features
├── requirement/               # How to define addressable rules in Behavior sections
├── acceptance-criteria/       # How to define abstract verification conditions
├── scenario/                  # How to write concrete behavior examples
├── source-references/         # How to link code to specifications
├── plan/                      # How to structure planning documents
├── plans-index/               # How to structure the plans index
├── ideas-index/               # How to structure the ideas index
├── features-index/            # How to structure this features index
├── scenarios-index/           # How to structure per-feature scenarios indexes
├── task/                      # How to define discrete units of work within a plan
├── project-definition/        # Project config and top-level structure
├── adherence-footer/          # Shared footer mechanism
├── document-types-registry/   # Canonical list of document types (this registry)
└── index/                     # Shared shape of every Index-Kind feature
```

## Integration with Orchestration Tools

SpecScore specs are format-agnostic. These features define the mental model and conventions:

- **Standalone:** Use SpecScore specs with any orchestration tool (Linear, Jira, custom)
- **With Rehearse:** Add automated testing and validation to SpecScore specs
- **With Synchestra:** Add multi-agent orchestration and coordination

See [synchestra.io](https://synchestra.io) for orchestration on top of SpecScore.

## Outstanding Questions

None at this time.

---
*This document follows the https://specscore.md/features-index-specification*
