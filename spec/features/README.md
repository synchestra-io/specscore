# SpecScore Features

Core features of the SpecScore specification framework:

| Feature | Status | Description |
|---------|--------|-------------|
| [feature](feature/README.md) | Defined | Feature specification structure and conventions |
| [architecture](architecture/README.md) | Defined | Architecture specification structure and conventions |
| [development-plan](development-plan/README.md) | Defined | Planning documents that bridge specs to task execution |
| [project-definition](project-definition/README.md) | Defined | SpecScore project configuration and root structure |

## Feature Hierarchy

```
spec/features/
├── feature/               # How to structure and write features
├── architecture/          # How to structure architecture specs
├── development-plan/      # How to structure planning documents
└── project-definition/    # Project config and top-level structure
```

## Integration with Orchestration Tools

SpecScore specs are format-agnostic. These features define the mental model and conventions:

- **Standalone:** Use SpecScore specs with any orchestration tool (Linear, Jira, custom)
- **With Rehearse:** Add automated testing and validation to SpecScore specs
- **With Synchestra:** Add multi-agent orchestration and coordination

See [synchestra.io](https://synchestra.io) for orchestration on top of SpecScore.
