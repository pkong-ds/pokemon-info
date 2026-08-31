# ADR 0002: Embedded catalog for completion, live PokeAPI for details

## Status

Accepted

## Context

The tool needs two kinds of data: a list of all Pokémon/Move names (for
shell completion and name lookup), and per-entry details (stats, types,
evolution chains, move effects). Alternatives considered:

1. **All-online**: fetch the name list from PokeAPI and cache it. Keeps names
   always fresh, but shell completion invokes the binary on every TAB press —
   a network round-trip there makes completion slow or broken, and adds a
   cache directory with new failure modes.
2. **All-offline**: embed every detail. The binary balloons by megabytes and
   a data refresh means a new release — which defeats the point of being a
   PokeAPI client.

Today the name list (id, slug, name, url) is generated into two YAML files
and embedded at build time; all details are fetched live.

## Decision

Keep the hybrid: the Catalog (names/slugs/URLs) is embedded at build time and
used for instant, offline shell completion and name resolution. Details are
always fetched live from PokeAPI. Catalog staleness is accepted and resolved
by shipping a new release (snap refresh).

## Consequences

- TAB completion is instant and works with no network.
- `pokemon-info <name>` requires network for details. The asymmetry
  (offline completion, online details) is intentional, not a bug.
- When a new Pokémon generation releases, the Catalog is stale until a
  maintainer runs the Prepare tool and cuts a new release.
- The Prepare tool is a maintainer-only repository tool and is never packaged
  inside the snap.
