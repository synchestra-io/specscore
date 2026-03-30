# SpecScore

**SpecScore** — Specification Patterns for AI-Driven Development

An open-source framework and mental model for writing specifications that AI agents can understand and execute.

## What Is SpecScore?

SpecScore teaches you to think about **specifications as machine-readable blueprints**. It defines:

- **Format:** How to structure and write specifications (directory layout, files, conventions)
- **Schema:** What makes a valid specification (required metadata, validation rules)
- **Mental Model:** Best practices for specifications in AI-driven development
- **Standards:** How specifications link to code and to each other

## Getting Started

**New to SpecScore?**
- Start with [Features](spec/features/README.md) to understand the mental model
- Read [feature/README.md](spec/features/feature/README.md) to write your first feature spec
- Use the SpecScore CLI to validate: `specscore lint`

**Already familiar with Synchestra specs?**
- SpecScore is the open-source foundation that Synchestra is built on
- Format is identical. You're already using SpecScore patterns.

## The SpecScore Approach

Specifications define *what* needs doing:
- **Features:** User-facing capabilities and requirements
- **Architecture:** System design and technical decisions
- **Plans:** Executable steps that bridge specs to implementation
- **Traceability:** Links from specs to code and back

This clarity creates three layers of value:

1. **For Planning:** Clear specs prevent misunderstanding and rework
2. **For Execution:** Machines (agents, tests) understand specifications automatically
3. **For Validation:** Rehearse tests can verify specs work end-to-end

## The Orchestra

SpecScore is part of an ecosystem:

| Tool | What It Does |
|------|-------------|
| **SpecScore** | Defines specification format and mental model |
| **Rehearse** | Tests and validates SpecScore specs automatically |
| **Synchestra** | Orchestrates multi-agent execution of SpecScore specs |

**Use standalone:** SpecScore works with any orchestration tool.
**Use together:** SpecScore + Rehearse + Synchestra = full spec-driven development lifecycle.

## Repository Structure

```
specscore/
├── README.md                    # This file
├── LICENSE                      # Apache 2.0
├── spec/                        # SpecScore specification
│   ├── features/               # Core features
│   │   ├── feature/           # Feature specification format
│   │   ├── architecture/       # Architecture specification format
│   │   ├── development-plan/   # Planning document format
│   │   └── project-definition/ # Project structure and config
│   └── README.md
├── tools/                       # SpecScore CLI and validators
│   └── (coming soon)
├── examples/                    # Example SpecScore projects
│   └── (coming soon)
├── docs/                        # User guides and tutorials
│   └── (coming soon)
└── .github/                     # GitHub templates, workflows
```

## Fair Questions

**"Isn't this the same as Synchestra's spec format?"**

SpecScore is what Synchestra's specification layer is *built on*. Synchestra adds orchestration, tasks, and coordination on top. You can use SpecScore specs independently with any tool — but they're optimized for Synchestra.

**"Can I use SpecScore without Synchestra?"**

Yes. Use SpecScore specs with:
- Your own orchestration (Linear, Jira, homegrown)
- Rehearse for automated testing
- Any CI/CD system that can read YAML/Markdown

**"What about vendor lock-in?"**

SpecScore specs are:
- Stored as plain text (YAML/Markdown)
- Schema-published and open
- Tool-agnostic format
- No proprietary extensions required

You own your specs. Use them anywhere.

## Links

- [SpecScore Features](spec/features/README.md)
- [Rehearse](https://github.com/synchestra-io/rehearse) — Test your specs
- [Synchestra](https://synchestra.io) — Orchestrate your specs
- [Synchestra Marketing](https://github.com/synchestra-io/synchestra-marketing) — Ecosystem positioning

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
