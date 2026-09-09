---
format: https://specscore.md/scenarios-index-specification
---

# Repo Config — Test Scenarios

Scenarios that verify the requirements defined in [Repo Config](../README.md).

| Scenario | Validates |
|----------|-----------|
| [minimal-config-valid](minimal-config-valid.md) | [REQ: file-name](../README.md#req-file-name), [REQ: empty-config-valid](../README.md#req-empty-config-valid) |
| [missing-header-comment-fails](missing-header-comment-fails.md) | [REQ: schema-header-comment](../README.md#req-schema-header-comment) |
| [header-comment-not-on-line-1-fails](header-comment-not-on-line-1-fails.md) | [REQ: schema-header-comment](../README.md#req-schema-header-comment) |
| [project-defaults-from-git](project-defaults-from-git.md) | [REQ: project-block-optional](../README.md#req-project-block-optional), [REQ: project-title-default](../README.md#req-project-title-default) |
| [project-overrides-respected](project-overrides-respected.md) | [REQ: source-reference-overrides](../README.md#req-source-reference-overrides) |
| [project-repositories-preserved](project-repositories-preserved.md) | [REQ: project-repositories](../README.md#req-project-repositories) |
| [projects-url-entries-preserved](projects-url-entries-preserved.md) | [REQ: projects-list](../README.md#req-projects-list) |
| [projects-local-path-resolves](projects-local-path-resolves.md) | [REQ: projects-list](../README.md#req-projects-list), [REQ: projects-local-path-must-resolve](../README.md#req-projects-local-path-must-resolve) |
| [projects-local-path-without-specscore-yaml-fails](projects-local-path-without-specscore-yaml-fails.md) | [REQ: projects-local-path-must-resolve](../README.md#req-projects-local-path-must-resolve) |
| [plan-route-precedence](plan-route-precedence.md) | [REQ: plans-repo-project-selection](../README.md#req-plans-repo-project-selection), [REQ: plan-repos-aggregate-routing](../README.md#req-plan-repos-aggregate-routing), [REQ: plan-route-required](../README.md#req-plan-route-required), [REQ: plan-config-precedence](../README.md#req-plan-config-precedence) |
| [plan-checkout-machine-local](plan-checkout-machine-local.md) | [REQ: repo-checkouts-machine-local](../README.md#req-repo-checkouts-machine-local) |
| [specs-dir-name-default](specs-dir-name-default.md) | [REQ: specs-dir-name-default](../README.md#req-specs-dir-name-default) |
| [docs-dir-name-default](docs-dir-name-default.md) | [REQ: docs-dir-name-default](../README.md#req-docs-dir-name-default) |
| [modules-default-to-root](modules-default-to-root.md) | [REQ: modules-default](../README.md#req-modules-default) |
| [module-name-deduced-from-path](module-name-deduced-from-path.md) | [REQ: module-name-deduction](../README.md#req-module-name-deduction) |
| [module-paths-resolved](module-paths-resolved.md) | [REQ: module-path-resolution](../README.md#req-module-path-resolution) |
| [module-paths-duplicate-fails](module-paths-duplicate-fails.md) | [REQ: module-paths-unique](../README.md#req-module-paths-unique) |
| [module-paths-nested-fails](module-paths-nested-fails.md) | [REQ: module-paths-non-nested](../README.md#req-module-paths-non-nested) |
| [module-root-coexists-with-explicit-paths](module-root-coexists-with-explicit-paths.md) | [REQ: module-paths-non-nested](../README.md#req-module-paths-non-nested) |
| [code-paths-are-scan-roots](code-paths-are-scan-roots.md) | [REQ: module-code-list](../README.md#req-module-code-list) |
| [unknown-fields-preserved](unknown-fields-preserved.md) | [REQ: unknown-fields-preserved](../README.md#req-unknown-fields-preserved) |
| [studio-defaults-when-omitted](studio-defaults-when-omitted.md) | [REQ: studio-default-when-omitted](../README.md#req-studio-default-when-omitted), [REQ: studio-toolbar-mandatory-unless-opted-out](../README.md#req-studio-toolbar-mandatory-unless-opted-out) |
| [studio-explicit-overrides](studio-explicit-overrides.md) | [REQ: studio-explicit-values](../README.md#req-studio-explicit-values) |
| [studio-partial-block-fails](studio-partial-block-fails.md) | [REQ: studio-explicit-values](../README.md#req-studio-explicit-values) |
| [studio-null-suppresses-toolbar](studio-null-suppresses-toolbar.md) | [REQ: studio-null-opts-out](../README.md#req-studio-null-opts-out) |
| [publication-block-delegated](publication-block-delegated.md) | [REQ: publication-block-optional](../README.md#req-publication-block-optional), [REQ: publication-schema-delegated](../README.md#req-publication-schema-delegated) |

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/scenarios-index-specification*
