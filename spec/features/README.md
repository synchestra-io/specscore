# SpecScore Features

Core features of the SpecScore specification framework:

| Feature | Status | Description |
|---------|--------|-------------|
| [feature](feature/README.md) | Conceptual | Feature specification structure, metadata, lifecycle, and conventions |
| [acceptance-criteria](acceptance-criteria/README.md) | Conceptual | Acceptance criteria format and conventions |
| [requirement](requirement/README.md) | Conceptual | Discrete testable rules within feature Behavior sections |
| [scenario](scenario/README.md) | Conceptual | Concrete Given/When/Then behavior examples in `_tests/` directories |
| [source-references](source-references/README.md) | Conceptual | Code-to-spec traceability via inline annotations |
| [development-plan](development-plan/README.md) | Conceptual | Planning documents that bridge specs to execution |
| [project-definition](project-definition/README.md) | Conceptual | SpecScore project configuration and root structure |

## Feature Hierarchy

```
spec/features/
├── feature/               # How to structure and write features
├── requirement/           # How to define addressable rules in Behavior sections
├── acceptance-criteria/   # How to define abstract verification conditions
├── scenario/              # How to write concrete behavior examples
├── source-references/     # How to link code to specifications
├── development-plan/      # How to structure planning documents
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
