---
name: release
description: Use when cutting a pokemon-info release — walks the real version-bump, tag, build, publish, and tap/bucket-mirror chain end to end.
---

# Release

Releases are tag-driven (ADR 0003). Everything below is what actually
happens, verified against the Makefile, the workflows in
`.github/workflows/`, and prior release commits.

## Pre-flight

- `make test` passes (`go vet ./... && go test ./...`) — the release
  workflow runs exactly this before building.
- Decide the version `X.Y.Z`. The current one lives in `version:` in
  `snapcraft.yaml` (and equals the latest `v*` git tag).

## Steps

1. **Bump the snap version.** Set `version:` in `snapcraft.yaml` to
   `X.Y.Z` and commit as `Bump snap version to X.Y.Z` (matches the
   `0.2.1` and `0.4.0` bumps in git history).
2. **Tag and push.** `git tag vX.Y.Z && git push origin vX.Y.Z`. The tag
   (pattern `v*.*.*`) triggers both `release.yml` and `snap.yml`.
3. **`release.yml`** (automatic): runs vet + test, cross-compiles
   linux/darwin × amd64/arm64 and windows/amd64 (`CGO_ENABLED=0`,
   version injected via `-X main.version=${GITHUB_REF_NAME}`), zips/tars
   the archives with LICENSE and README, writes `checksums.txt`, and
   creates the GitHub release with `gh release create --generate-notes`.
4. **`snap.yml`** (automatic): builds and publishes the **amd64** snap to
   the `stable` channel via the `SNAP_STORE_TOKEN` secret. No workflow
   publishes an arm64 snap — arm64 snaps are not covered by CI.
5. **Fill the tap formula and scoop manifest.** Once the release exists,
   download `checksums.txt` from the release and fill:
   - `packaging/homebrew/Formula/pokemon-info.rb` — `version` plus the
     four `sha256` fields (darwin/amd64, darwin/arm64).
   - `packaging/scoop/bucket/pokemon-info.json` — `version` plus the
     windows-amd64 `hash`.
   Commit (prior style: `Fill v0.3.0 checksums in the tap formula and scoop
   manifest`).
6. **Mirror to the published repos.** Apply the same edits to
   `pkong-ds/homebrew-tap` (`Formula/pokemon-info.rb`) and
   `pkong-ds/scoop-bucket` (`bucket/pokemon-info.json`) and push — those
   repos are what `brew` and `scoop` actually consume. The `packaging/`
   copies in this repo are the maintained source.
7. **Verify.** `pokemon-info version` on each channel reports the tag.

## Notes

- Do not create the GitHub release by hand — `release.yml` owns it.
- Keep `snapcraft.yaml` `version:` and the git tag in sync: the snap
  binary reports `v${CRAFT_PROJECT_VERSION}` from that field.
- Build artifacts (`dist/`, `*.snap`) are gitignored; never commit them.
