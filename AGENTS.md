# AGENTS.md

pokemon-info — a terminal Pokédex: a single Go binary with an embedded name
catalog (instant offline completion), a Bubble Tea TUI for browsing, and live
PokeAPI details. Distributed via GitHub Releases, snap, Homebrew tap, Scoop
bucket, and a curl installer.

## Layout

- `cmd/pokemon-info/` — the user-facing binary (cobra CLI + bubbletea TUI)
- `cmd/pokemon-info/pokemons.yaml`, `moves.yaml` — the generated Catalog,
  embedded at build time
- `cmd/prepare/` — the Prepare tool: regenerates the Catalog from PokeAPI
  (maintainer-only, never shipped in the snap)
- `scripts/install.sh` — the curl one-liner installer
- `packaging/homebrew/Formula/`, `packaging/scoop/bucket/` — the Homebrew
  Formula and Scoop manifest, maintained here and mirrored to the
  `pkong-ds/homebrew-tap` and `pkong-ds/scoop-bucket` repos
- `snapcraft.yaml` — the snap definition; its `version:` field is the
  release version
- `docs/adr/` — accepted decisions; `docs/tapes/` — vhs tapes behind the
  README demo GIFs
- `CONTEXT.md` — the glossary (glossary only)
- `.github/workflows/` — `ci.yml`, `release.yml`, `snap.yml`

## Daily commands (verified)

```
make build      # dist/pokemon-info, version from git describe --tags
make test       # go vet ./... && go test ./...
go build ./...  # compile check; CI runs the same plus vet and test
make catalog    # regenerate the Catalog from live PokeAPI (needs network)
make demos      # re-record docs/demo.gif and docs/moves.gif via vhs
make clean      # rm -rf dist
```

`make demos` requires [vhs](https://github.com/charmbracelet/vhs)
(`brew install vhs`).

## Release flow (what actually happens)

Releases are tag-driven (ADR 0003). See the `release` skill for the full
checklist. In brief:

1. Bump `version:` in `snapcraft.yaml`; commit as `Bump snap version to X.Y.Z`
   (matches prior release commits).
2. Tag `vX.Y.Z` and push the tag.
3. `release.yml` runs vet + test, cross-compiles linux/darwin × amd64/arm64
   plus windows/amd64, then creates the GitHub release with the archives
   and `checksums.txt` (`gh release create --generate-notes`).
4. `snap.yml` (same tag trigger) publishes the **amd64** snap to the
   `stable` channel via `SNAP_STORE_TOKEN`. No workflow publishes the
   arm64 snap.
5. Fill the version and sha256 fields in
   `packaging/homebrew/Formula/pokemon-info.rb` and
   `packaging/scoop/bucket/pokemon-info.json` from `checksums.txt`, commit,
   and mirror to `pkong-ds/homebrew-tap` and `pkong-ds/scoop-bucket`.
6. `pokemon-info version` reports the tag everywhere: it is ldflags-injected
   (`main.version`); snap builds inject `CRAFT_PROJECT_VERSION`.

## Git conventions

Commit subjects are imperative and capitalized, no prefixes:
`Add disk detail cache, token-bucket rate limit, and in-memory LRU`,
`Fix TUI filter: forward list filter messages and track selection by slug`,
`Bump snap version to 0.4.0`. A colon separates topic from detail; subjects
run long rather than splitting.

## Hard rules

- Never hand-edit `pokemons.yaml` / `moves.yaml` — regenerate via the
  Prepare tool (`make catalog`). See ADR 0002.
- `CONTEXT.md` is glossary-only. New terms get defined there; nothing else
  goes in the file.
- New durable decisions go in `docs/adr/` as the next numbered ADR,
  following the existing format (Status / Context / Decision /
  Consequences).
- Build artifacts (`dist/`, `*.snap`) are gitignored — never commit them.
- The Catalog badge counts at the top of README.md move only when the
  Catalog is regenerated.
