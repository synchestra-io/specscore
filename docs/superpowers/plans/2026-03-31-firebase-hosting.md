# Firebase Hosting Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Configure Firebase Hosting for the `specscore-org` site under the `synchestra-io` project, and add a GitHub Actions deploy job to `site-ci.yml` that deploys on every push to `main` using Workload Identity Federation.

**Architecture:** Two config files (`firebase.json`, `.firebaserc`) are added to the repo root. The existing `site-ci.yml` workflow gains a `deploy` job that depends on `test-and-build`, rebuilds the site, authenticates via WIF using `google-github-actions/auth@v2`, and deploys using `firebase-tools` CLI directly.

**Tech Stack:** Firebase Hosting, Firebase CLI (`firebase-tools`), GitHub Actions, `google-github-actions/auth@v2` (Workload Identity Federation)

**Spec:** `docs/superpowers/specs/2026-03-31-firebase-hosting-design.md`

---

## File Structure

```
firebase.json                        # Firebase Hosting config — create
.firebaserc                          # Firebase project binding — create
.github/workflows/site-ci.yml        # Add deploy job — modify
```

---

### Task 1: Add Firebase config files

**Files:**
- Create: `firebase.json`
- Create: `.firebaserc`

- [ ] **Step 1: Create `firebase.json`**

Write to `firebase.json` in the repo root:

```json
{
  "hosting": {
    "site": "specscore-org",
    "public": "public",
    "ignore": ["firebase.json", "**/.*"],
    "cleanUrls": true,
    "trailingSlash": false
  }
}
```

- [ ] **Step 2: Create `.firebaserc`**

Write to `.firebaserc` in the repo root:

```json
{
  "projects": {
    "default": "synchestra-io"
  }
}
```

- [ ] **Step 3: Validate both files are valid JSON**

Run from the repo root:

```bash
node -e "JSON.parse(require('fs').readFileSync('firebase.json','utf8')); console.log('firebase.json OK')"
node -e "JSON.parse(require('fs').readFileSync('.firebaserc','utf8')); console.log('.firebaserc OK')"
```

Expected output:
```
firebase.json OK
.firebaserc OK
```

- [ ] **Step 4: Commit**

```bash
git add firebase.json .firebaserc
git commit -m "feat: add Firebase Hosting config for specscore-org"
```

---

### Task 2: Add deploy job to site-ci.yml

**Files:**
- Modify: `.github/workflows/site-ci.yml`

The current `site-ci.yml` has a single `test-and-build` job. Add a `deploy` job after it. The full updated file is shown below — replace the entire file contents.

- [ ] **Step 1: Replace `.github/workflows/site-ci.yml` with the updated workflow**

```yaml
name: Site Generator CI

on:
  push:
    paths:
      - '.github/workflows/site-ci.yml'
      - 'tools/site-generator/**'
      - 'spec/**'
  pull_request:
    paths:
      - '.github/workflows/site-ci.yml'
      - 'tools/site-generator/**'
      - 'spec/**'
  workflow_dispatch:

permissions:
  contents: read

concurrency:
  group: site-ci-${{ github.ref }}
  cancel-in-progress: true

jobs:
  test-and-build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
          cache-dependency-path: tools/site-generator/package-lock.json

      - name: Install dependencies
        working-directory: tools/site-generator
        run: npm ci

      - name: Install Chrome with system deps
        run: npx puppeteer browsers install chrome --install-deps

      - name: Run tests
        working-directory: tools/site-generator
        run: npm test

      - name: Build site
        working-directory: tools/site-generator
        run: npm run build

      - name: Verify no mermaid code blocks in HTML
        run: |
          if grep -r '```mermaid' public/*.html; then
            echo "ERROR: Found unrendered mermaid blocks in HTML output"
            exit 1
          fi

      - name: Verify mermaid SVGs rendered
        run: |
          for f in public/feature-specification.html public/development-plan-specification.html; do
            if ! grep -q '<svg' "$f"; then
              echo "ERROR: No SVG found in $f"
              exit 1
            fi
          done

  deploy:
    needs: test-and-build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main' && github.event_name == 'push'

    permissions:
      contents: read
      id-token: write

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
          cache-dependency-path: tools/site-generator/package-lock.json

      - name: Install dependencies
        working-directory: tools/site-generator
        run: npm ci

      - name: Install Chrome with system deps
        run: npx puppeteer browsers install chrome --install-deps

      - name: Build site
        working-directory: tools/site-generator
        run: npm run build

      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
          service_account: ${{ secrets.WIF_SERVICE_ACCOUNT }}

      - name: Deploy to Firebase Hosting
        run: npx firebase-tools@latest deploy --only hosting:specscore-org --non-interactive
```

- [ ] **Step 2: Validate the YAML is well-formed**

```bash
node -e "
const yaml = require('js-yaml');
const fs = require('fs');
try {
  yaml.load(fs.readFileSync('.github/workflows/site-ci.yml', 'utf8'));
  console.log('YAML OK');
} catch(e) {
  console.error(e.message); process.exit(1);
}
"
```

If `js-yaml` is not available locally, use this alternative:

```bash
python3 -c "import yaml, sys; yaml.safe_load(open('.github/workflows/site-ci.yml')); print('YAML OK')"
```

Expected output:
```
YAML OK
```

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/site-ci.yml
git commit -m "feat: add Firebase Hosting deploy job to site-ci workflow"
```

---

## Prerequisites (manual steps before first deploy)

These must be done in GCP/Firebase console before the workflow can succeed — they cannot be automated by this plan:

1. **Firebase Hosting site created** — In the Firebase console under `synchestra-io`, create a Hosting site with the ID `specscore-org`.

2. **Workload Identity Pool and Provider configured** — A WIF pool and OIDC provider must exist in the `synchestra-io` GCP project, bound to `https://token.actions.githubusercontent.com` for the `synchestra-io/specscore` repo. (Mirror the setup from the synchestra repo.)

3. **Service account with Firebase Hosting permissions** — A service account with `roles/firebasehosting.admin` on the `synchestra-io` project, with an IAM binding allowing the WIF provider to impersonate it.

4. **GitHub Actions secrets set** — In the `synchestra-io/specscore` repo settings, add:
   - `WIF_PROVIDER`: the Workload Identity Provider resource name (format: `projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_ID/providers/PROVIDER_ID`)
   - `WIF_SERVICE_ACCOUNT`: the service account email (format: `NAME@synchestra-io.iam.gserviceaccount.com`)
