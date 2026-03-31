# Scenario: Host/org/repo inferred from git remote origin

**Validates:** [source-references#req:git-remote-default](../README.md#req-git-remote-default)

## Steps

GIVEN the git remote `origin` URL is `git@github.com:acme/myproject.git`
AND no project config override exists
AND a source file contains `// specscore:feature/cli/task/claim`
WHEN the resolver expands the reference
THEN the host is resolved as `github.com`
AND the org is resolved as `acme`
AND the repo is resolved as `myproject`
AND the expanded URL is `https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim`
