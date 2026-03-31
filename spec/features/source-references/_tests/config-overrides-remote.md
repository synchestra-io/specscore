# Scenario: Project config overrides git remote inference

**Validates:** [source-references#req:config-override](../README.md#req-config-override)

## Steps

GIVEN the git remote `origin` URL is `git@github.com:dev-fork/myproject.git`
AND `specscore-project.yaml` contains:
```yaml
project:
  host: github.com
  org: acme
  repo: myproject
```
AND a source file contains `// specscore:feature/cli/task/claim`
WHEN the resolver expands the reference
THEN the expanded URL uses the config values, not the git remote
AND the expanded URL is `https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim`
