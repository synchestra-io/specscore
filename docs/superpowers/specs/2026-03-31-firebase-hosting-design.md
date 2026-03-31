# Firebase Hosting for specscore.org

## Goal

Deploy the SpecScore static site (built by `tools/site-generator/` into `public/`) to Firebase Hosting under the `synchestra-io` GCP project, served at the `specscore-org` hosting site. Deployment is automated via GitHub Actions on every push to `main`.

## Firebase Configuration

Two files added to the repository root:

**`firebase.json`**

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

- `cleanUrls: true` — serves `slug.html` for requests to `/slug`, matching the URL scheme defined in the static site generator design.
- `trailingSlash: false` — canonical URLs have no trailing slash.
- `public: "public"` — matches the build output directory from `npm run build`.

**`.firebaserc`**

```json
{
  "projects": {
    "default": "synchestra-io"
  }
}
```

Binds the repo to the `synchestra-io` Firebase project.

## GitHub Actions

The existing `site-ci.yml` workflow gains a `deploy` job that runs after `test-and-build` succeeds, but only on push to `main` (not on PRs or other branches).

### Job structure

```
test-and-build  →  deploy (main only)
```

The `deploy` job re-runs the full build rather than passing an artifact. This keeps the job self-contained and avoids artifact upload/download complexity for a build that takes under a minute.

### Authentication

Workload Identity Federation (WIF) — keyless, no long-lived credentials stored as secrets. Same pattern used by the synchestra repo.

Two GitHub Actions secrets required:
- `WIF_PROVIDER` — Workload Identity Provider resource name
- `WIF_SERVICE_ACCOUNT` — service account email with Firebase Hosting deploy permissions (`roles/firebasehosting.admin` on the `synchestra-io` project)

The `id-token: write` permission is set on the `deploy` job to enable OIDC token issuance.

### Deploy job (to be added to `site-ci.yml`)

```yaml
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
      uses: FirebaseExtended/action-hosting-deploy@v0
      with:
        repoToken: ${{ secrets.GITHUB_TOKEN }}
        firebaseServiceAccount: ''
        projectId: synchestra-io
        target: specscore-org
        channelId: live
```

## What This Design Does Not Cover

- **Custom domain DNS setup** — pointing `specscore.org` to the Firebase Hosting site is a manual step in the Firebase console and DNS registrar.
- **WIF infrastructure setup** — the Workload Identity Pool, Provider, and service account bindings must exist in GCP before the workflow runs. This mirrors the existing synchestra repo setup.
- **Firebase Hosting site creation** — the `specscore-org` hosting site must be created in the Firebase console under the `synchestra-io` project before the first deploy.
