# Glossary

The shared language of this repo. When a term is resolved in discussion, it
gets defined here. Nothing else lives in this file — no specs, no plans.

## Canonical name

**pokemon-info** — the single name used by the snap, the GitHub repository,
the binary, and shell completion. See ADR 0001 for why it is not `poke`.

## Catalog

**Catalog** — the complete list of Pokémon and Move entries (id, slug, name,
url) embedded in the binary at build time. Powers name completion and
name-to-URL resolution. Deliberately contains no stats or details. It goes
stale between releases; a snap refresh fixes it.

## Prepare tool

**Prepare tool** — the maintainer-only command that regenerates the Catalog
from PokeAPI. Lives in the repository, never ships in the snap. Not part of
the user-facing product.
