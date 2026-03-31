# Scenario: Cross-repo reference requires @host/org/repo suffix

**Validates:** [source-references#req:cross-repo-suffix](../README.md#req-cross-repo-suffix)

## Steps

GIVEN a source file contains `// specscore:feature/agent-skills@github.com/acme/orchestrator`
WHEN the resolver processes the reference
THEN the reference targets the repository `github.com/acme/orchestrator`
AND the resolved path is `spec/features/agent-skills`
AND the expanded URL is `https://specscore.org/github.com/acme/orchestrator/spec/features/agent-skills`
