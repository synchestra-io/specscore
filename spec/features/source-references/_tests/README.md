# Scenarios: Source References

Test scenarios for the [Source References](../README.md) specification.

| Scenario | Validates |
|---|---|
| [type-prefix-expands](type-prefix-expands.md) | [source-references#req:type-prefix-expansion](../README.md#req-type-prefix-expansion) |
| [unknown-prefix-falls-back-to-path](unknown-prefix-falls-back-to-path.md) | [source-references#req:resolution-order](../README.md#req-resolution-order) |
| [type-prefix-fallback-on-missing-path](type-prefix-fallback-on-missing-path.md) | [source-references#req:resolution-order](../README.md#req-resolution-order) |
| [short-notation-expanded-to-url](short-notation-expanded-to-url.md) | [source-references#ac:canonical-expansion](../README.md#ac-canonical-expansion), [source-references#req:canonical-url-form](../README.md#req-canonical-url-form), [source-references#req:auto-expansion](../README.md#req-auto-expansion) |
| [url-omits-type-prefix](url-omits-type-prefix.md) | [source-references#req:url-structure](../README.md#req-url-structure) |
| [comment-prefix-detected](comment-prefix-detected.md) | [source-references#req:comment-prefix-required](../README.md#req-comment-prefix-required) |
| [bare-reference-ignored](bare-reference-ignored.md) | [source-references#req:comment-prefix-required](../README.md#req-comment-prefix-required) |
| [string-literal-ignored](string-literal-ignored.md) | [source-references#req:comment-prefix-required](../README.md#req-comment-prefix-required) |
| [nonexistent-reference-errors](nonexistent-reference-errors.md) | [source-references#req:nonexistent-is-error](../README.md#req-nonexistent-is-error) |
| [missing-context-errors](missing-context-errors.md) | [source-references#req:unresolvable-context-error](../README.md#req-unresolvable-context-error) |
| [git-remote-inferred](git-remote-inferred.md) | [source-references#req:git-remote-default](../README.md#req-git-remote-default) |
| [config-overrides-remote](config-overrides-remote.md) | [source-references#req:config-override](../README.md#req-config-override) |
| [cross-repo-suffix-required](cross-repo-suffix-required.md) | [source-references#req:cross-repo-suffix](../README.md#req-cross-repo-suffix) |
| [only-known-type-prefixes](only-known-type-prefixes.md) | [source-references#req:type-prefix-set](../README.md#req-type-prefix-set) |

## Outstanding Questions

None at this time.

---
*This document follows the https://specscore.md/scenarios-index-specification*
