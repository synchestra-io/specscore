# Scenarios: Idea

Test scenarios for the [Idea](../README.md) specification.

| Scenario | Validates |
|---|---|
| [idea-in-wrong-location-rejected](idea-in-wrong-location-rejected.md) | [idea#req:idea-location](../README.md#req-idea-location) |
| [invalid-slug-rejected](invalid-slug-rejected.md) | [idea#req:slug-format](../README.md#req-slug-format) |
| [idea-directory-rejected](idea-directory-rejected.md) | [idea#req:single-file](../README.md#req-single-file) |
| [missing-idea-prefix](missing-idea-prefix.md) | [idea#ac:idea-header](../README.md#ac-idea-header), [idea#req:title-format](../README.md#req-title-format) |
| [missing-required-header-field](missing-required-header-field.md) | [idea#ac:idea-header](../README.md#ac-idea-header), [idea#req:header-fields](../README.md#req-header-fields) |
| [slug-derived-from-filename](slug-derived-from-filename.md) | [idea#ac:idea-header](../README.md#ac-idea-header), [idea#req:id-is-slug](../README.md#req-id-is-slug) |
| [author-edits-promotes-to-rejected](author-edits-promotes-to-rejected.md) | [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:promotes-to-managed](../README.md#req-promotes-to-managed) |
| [missing-required-section](missing-required-section.md) | [idea#ac:idea-structure](../README.md#ac-idea-structure), [idea#req:required-sections](../README.md#req-required-sections) |
| [empty-not-doing-rejected](empty-not-doing-rejected.md) | [idea#ac:idea-structure](../README.md#ac-idea-structure), [idea#req:not-doing-non-empty](../README.md#req-not-doing-non-empty) |
| [missing-must-be-true-assumption](missing-must-be-true-assumption.md) | [idea#ac:idea-structure](../README.md#ac-idea-structure), [idea#req:must-be-true-present](../README.md#req-must-be-true-present) |
| [missing-hmw-framing-warning](missing-hmw-framing-warning.md) | [idea#req:hmw-framing](../README.md#req-hmw-framing) |
| [invalid-status-rejected](invalid-status-rejected.md) | [idea#ac:idea-header](../README.md#ac-idea-header), [idea#req:status-values](../README.md#req-status-values) |
| [specified-without-promotion-rejected](specified-without-promotion-rejected.md) | [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:specified-requires-promotion](../README.md#req-specified-requires-promotion) |
| [archived-idea-outside-archived-dir-rejected](archived-idea-outside-archived-dir-rejected.md) | [idea#ac:archival](../README.md#ac-archival), [idea#req:archived-location](../README.md#req-archived-location) |
| [archived-without-reason-rejected](archived-without-reason-rejected.md) | [idea#ac:archival](../README.md#ac-archival), [idea#req:archive-reason](../README.md#req-archive-reason) |
| [supersedes-non-archived-rejected](supersedes-non-archived-rejected.md) | [idea#ac:archival](../README.md#ac-archival), [idea#req:supersedes-target-archived](../README.md#req-supersedes-target-archived) |
| [invalid-related-ideas-syntax-rejected](invalid-related-ideas-syntax-rejected.md) | [idea#ac:related-ideas](../README.md#ac-related-ideas), [idea#req:related-ideas-format](../README.md#req-related-ideas-format) |
| [related-ideas-broken-slug-rejected](related-ideas-broken-slug-rejected.md) | [idea#ac:related-ideas](../README.md#ac-related-ideas), [idea#req:related-ideas-target-exists](../README.md#req-related-ideas-target-exists) |
| [dependency-cycle-accepted](dependency-cycle-accepted.md) | [idea#ac:related-ideas](../README.md#ac-related-ideas), [idea#req:cycles-allowed](../README.md#req-cycles-allowed) |
| [specified-derivation-from-feature-reference](specified-derivation-from-feature-reference.md) | [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:specified-derivation](../README.md#req-specified-derivation) |
| [author-writes-specified-rejected](author-writes-specified-rejected.md) | [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:specified-not-author-set](../README.md#req-specified-not-author-set) |
| [sync-drift-rejected-fix-repairs](sync-drift-rejected-fix-repairs.md) | [idea#ac:sync-strictness](../README.md#ac-sync-strictness), [idea#req:sync-lint-strict](../README.md#req-sync-lint-strict) |
| [feature-references-draft-idea-rejected](feature-references-draft-idea-rejected.md) | [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:feature-cross-reference](../README.md#req-feature-cross-reference) |
| [hand-written-idea-lints-clean](hand-written-idea-lints-clean.md) | [idea#ac:authoring-independence](../README.md#ac-authoring-independence), [idea#req:authoring-agnostic](../README.md#req-authoring-agnostic) |
| [scaffold-command-produces-lint-clean-file](scaffold-command-produces-lint-clean-file.md) | [idea#ac:scaffold-behavior](../README.md#ac-scaffold-behavior), [idea#req:scaffold-command](../README.md#req-scaffold-command) |
| [unlisted-idea-in-index](unlisted-idea-in-index.md) | [idea#req:index-completeness](../README.md#req-index-completeness) |
| [archived-index-chronological-ordering](archived-index-chronological-ordering.md) | [idea#ac:archival](../README.md#ac-archival), [idea#req:archived-index-chronological](../README.md#req-archived-index-chronological) |

## Outstanding Questions

None at this time.
