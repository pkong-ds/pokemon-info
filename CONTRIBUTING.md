# Contributing to pokemon-info

Thanks for helping! A few pointers before you start:

- The shared vocabulary lives in [CONTEXT.md](CONTEXT.md) — when a term
  gets resolved in discussion, it gets defined there.
- Architectural decisions live in [docs/adr/](docs/adr/).
- `make test` runs go vet and go test; CI runs the same.
- `make build` produces `dist/pokemon-info`.

## Regenerating the catalog

The name catalog (`cmd/pokemon-info/pokemons.yaml`, `moves.yaml`) is
embedded in the binary at build time and goes stale between releases.
Regenerate it with the Prepare tool (`cmd/prepare`), a maintainer-only
command that never ships in the snap:

```
make catalog
```

## Regenerating the demo GIFs

The README GIFs are recorded with
[vhs](https://github.com/charmbracelet/vhs) from the tapes in
`docs/tapes/`:

```
make demos
```

Requires `vhs` installed (Homebrew: `brew install vhs`). Re-record
whenever the TUI or CLI output changes shape.
