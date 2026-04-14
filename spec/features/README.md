# SpecScore Features

Core features of the SpecScore specification framework:

| Feature | Status | Description |
|---------|--------|-------------|
| [feature](feature/README.md) | Stable | Feature specification structure, metadata, lifecycle, and conventions |
| [acceptance-criteria](acceptance-criteria/README.md) | Stable | Acceptance criteria format and conventions |
| [requirement](requirement/README.md) | Stable | Discrete testable rules within feature Behavior sections |
| [scenario](scenario/README.md) | Stable | Concrete Given/When/Then behavior examples in `_tests/` directories |
| [source-references](source-references/README.md) | Stable | Code-to-spec traceability via inline annotations |
| [plan](plan/README.md) | Stable | Planning documents that bridge specs to execution |
| [task](task/README.md) | Stable | Discrete units of work within a plan |
| [project-definition](project-definition/README.md) | Stable | SpecScore project configuration and root structure |

## Feature Hierarchy

```
spec/features/
├── feature/               # How to structure and write features
├── requirement/           # How to define addressable rules in Behavior sections
├── acceptance-criteria/   # How to define abstract verification conditions
├── scenario/              # How to write concrete behavior examples
├── source-references/     # How to link code to specifications
├── plan/                  # How to structure planning documents
├── task/                  # How to define discrete units of work within a plan
└── project-definition/    # Project config and top-level structure
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
*This document follows the https://specscore.md/feature-specification*
