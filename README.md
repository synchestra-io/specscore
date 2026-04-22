# SpecScore

**The open specification standard for AI-driven development**

SpecScore is an open specification format that makes requirements machine-readable without making them human-unreadable. It's Markdown and YAML — version-controlled, portable, no vendor lock-in. Use it standalone or with any orchestration tool.

[Read the Spec](spec/README.md) · [specscore.md](https://specscore.md) · [Ecosystem](docs/ecosystem.md)

---

## Reference CLI

The reference implementation for working with SpecScore repositories is [`specscore-cli`](https://github.com/synchestra-io/specscore-cli):

```bash
curl -fsSL https://specscore.md/get-cli | sh
```

```bash
specscore spec lint              # validate the current spec tree
specscore feature list           # list features
specscore feature show <slug>    # inspect a feature
specscore task list              # show the task board
specscore version                # full build identity
```

Full installation guide: [docs/installation.md](docs/installation.md). Source: <https://github.com/synchestra-io/specscore-cli>.

---

## What SpecScore Defines

SpecScore provides a structured format for:

- **Features** — user-facing capabilities with requirements and acceptance criteria
- **Requirements** — scoped, testable conditions that define done
- **Acceptance Criteria** — machine-readable conditions tied to features and tasks
- **Plans** — ordered sequences of tasks that bridge specs to implementation
- **Tasks** — atomic units of work assigned to agents or people
- **Source References** — traceable links from specs to code and back
- **Project Definition** — root configuration that ties a project together

---

## Role-Based Guides

SpecScore is designed for every role on a product team:

| Role | Guide |
|------|-------|
| Developers | [docs/for/developers.md](docs/for/developers.md) |
| Product Owners | [docs/for/product-owners.md](docs/for/product-owners.md) |
| QA Engineers | [docs/for/qas.md](docs/for/qas.md) |
| Business Analysts | [docs/for/business-analysts.md](docs/for/business-analysts.md) |
| Project Managers | [docs/for/project-managers.md](docs/for/project-managers.md) |
| Architects | [docs/for/architects.md](docs/for/architects.md) |

---

## Ecosystem

SpecScore is the foundation layer of a three-tool stack:

| Tool | Role |
|------|------|
| **SpecScore** | Open specification format — the standard itself |
| **Rehearse** | Validates and tests SpecScore specs automatically |
| **Synchestra** | Orchestrates multi-agent execution of SpecScore specs |

Use SpecScore standalone with any tool, or pair it with [Rehearse](https://github.com/synchestra-io/rehearse) and [Synchestra](https://synchestra.io) for a full spec-driven development lifecycle.

See [docs/ecosystem.md](docs/ecosystem.md) for details.

---

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

---

## License

This repository's specification text and documentation are licensed under [CC BY 4.0](LICENSE). The reference CLI implementation [`specscore-cli`](https://github.com/synchestra-io/specscore-cli) is licensed separately under Apache-2.0.

---

## History

The reference CLI `specscore` was extracted from this repository on 2026-04-22. Its source code, releases, and the `v0.x` release tags that were originally created here now live at [`synchestra-io/specscore-cli`](https://github.com/synchestra-io/specscore-cli). This repository's `v*` tags were removed at the same time, since they tagged CLI releases that no longer reside here. Engineering history (commits) for the extracted code is preserved in `specscore-cli` via `git filter-repo`.
