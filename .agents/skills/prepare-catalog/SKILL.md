---
name: prepare-catalog
description: Use when the embedded Pokémon/move name catalog needs regenerating — a new generation released, counts drifted, or before a release that should ship fresh names.
---

# Prepare the Catalog

The Catalog (`cmd/pokemon-info/pokemons.yaml` and `moves.yaml`) is the
complete list of Pokémon and Move entries (id, slug, name, url) embedded
in the binary at build time. It powers offline shell completion and
name-to-URL resolution, and deliberately contains no details (ADR 0002).
It goes stale between releases — this skill refreshes it.

## When to run

- A new Pokémon generation (or move set) appeared on PokeAPI.
- Before cutting a release, if names may have drifted.
- Never as a fix for anything else — the Catalog is only names.

## Steps

1. Run the Prepare tool via the Makefile (needs network; hits live
   PokeAPI):

   ```
   make catalog
   ```

   which runs:

   ```
   go run ./cmd/prepare --output-format yaml --resource pokemon --output-file cmd/pokemon-info/pokemons.yaml
   go run ./cmd/prepare --output-format yaml --resource move --output-file cmd/pokemon-info/moves.yaml
   ```

2. Inspect the diff of the two YAML files. Only entry counts/ids/names
   should change. `go test ./...` still passes.
3. If the counts changed, update the Catalog badge at the top of
   `README.md` (the `1302 Pokémon · 937 moves` badge) to match.
4. Commit the regenerated files (prior subjects are plain and
   descriptive, e.g. `Regenerate the catalog for gen 10`). The new names
   reach users only when the next release ships (ADR 0002) — see the
   `release` skill.

## Hard rules

- Never hand-edit `pokemons.yaml` / `moves.yaml` — always regenerate via
  the Prepare tool (`cmd/prepare`). See ADR 0002.
- The Prepare tool is maintainer-only and never ships inside the snap.
