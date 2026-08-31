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

## TUI

**TUI** — the interactive terminal interface. Launched by default when the
binary runs with no arguments in a terminal (explicit form: `pokemon-info
ui`). Piped invocations never launch it. Browsing inspired by poketex.

## Browser

**Browser** — a TUI page listing all Catalog entries of one kind (Pokémon or
Moves) with a live-fetched detail pane. `/pokemons` and `/moves` open
Browsers from the landing page.

## Colorscripts

**Colorscripts** — the embedded ANSI half-block Pokémon art (small and
large sets, regular colors), compressed from pokemon-colorscripts
(phoneybadber) plus poketex's gen 9 pack (Caruban). Shown above the stats
in the Pokémon Browser; the `f` key cycles small → large → off. Art wider
than the detail pane is skipped.
