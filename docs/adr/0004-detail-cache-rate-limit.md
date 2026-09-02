# ADR 0004: Disk detail cache, bounded in-memory cache, and a client-side rate limit

## Status

Accepted

## Context

Every detail view is fetched live from PokeAPI (ADR 0002), and a
Pokémon detail is a chain of three requests: `/pokemon/{id}`, then
`/pokemon-species/{id}`, then `/evolution-chain/{id}`. Two problems:

1. **Repeat fetches.** The same Pokémon looked up tomorrow costs the
   server (and the user) the same three requests again. There was no
   persistence across sessions — only the TUI's in-memory session map.
2. **Server load.** The TUI fetches on selection with a 250ms debounce;
   rapid scrolling can fire 3-request chains at a sustained 12 req/s.
   PokeAPI is a community-run free service. Being a good citizen is a
   requirement, not a nicety — especially before publicizing the tool
   (e.g. a Show-and-tell post on the PokeAPI repo).

PokeAPI is the only runtime upstream: sprite URLs appear in responses
but are printed, never fetched; the Colorscripts and the Catalog are
embedded at build time.

Alternatives considered:

1. **HTTP conditional requests** (ETag / If-Modified-Since). Correct,
   but adds conditional logic and depends on server header behavior we
   do not control. Overkill for data that changes at generational speed.
2. **No caching, just rate limiting.** Protects the server, but wastes
   requests on data that is effectively immutable between generations.
3. **Aggressive per-request spacing** (fixed 200ms between all
   requests). Protects the server, but taxes every lookup with ~400ms
   of artificial delay — even a single legitimate 3-request chain.

## Decision

Three layers, all best-effort and invisible when broken:

1. **Detail Cache (disk).** Raw JSON responses cached at
   `os.UserCacheDir()/pokemon-info/`, keyed by URL path
   (`pokemon-25.json`). TTL by file mod time, **30 days default**.
   Configurable from the TUI via `/config` (persisted at
   `os.UserConfigDir()/pokemon-info/config.json`, accepts `Nd` day
   shorthand or Go durations); the `POKEMON_INFO_CACHE_TTL` env var
   overrides for automation. Any cache error (unreadable dir, bad
   write) is silently
   skipped — the fetch falls through to the network. PokeAPI data
   changes at generational speed; 30 days is safely fresh.
2. **In-memory session LRU.** The Browser's session map becomes a
   fixed-capacity LRU (100 entries) so long browsing sessions cannot
   grow memory without bound. Eviction only kicks in past the cap;
   disk cache makes evicted re-fetches cheap.
3. **Token-bucket rate limit** on every API request: **burst 3**
   (= one full detail chain, so a single lookup pays no artificial
   delay) with a **5 req/s sustained refill**. Rapid scrolling drains
   the burst and then throttles to the refill rate.

A **`/clearcache`** slash command in the TUI clears the disk cache and
both Browsers' in-memory caches, with a brief confirmation flash. It
exists because TTL-based invalidation alone leaves no user escape
hatch for "I know this changed." A **`/config`** slash command edits
the cache TTL in the TUI and persists it; the config file lives
outside the cache directory so clearing the cache never resets
settings. The rate limit is deliberately **not** user-configurable:
it is a server-citizenship contract, not a preference.

## Consequences

- First lookup of any entry: network speed, no rate-limit tax (the
  burst covers the chain). Every later lookup within the TTL: a disk
  read.
- A user can see up to 30-day-old details; `/clearcache` (or TTL
  expiry) forces fresh data. No ETag revalidation — accepted staleness
  in exchange for zero server-side conditional traffic.
- Sustained browsing pressure on PokeAPI is capped at ~5 req/s per
  client, regardless of scroll speed.
- Cache files are plain JSON and safe to delete by hand; the command
  and the manual path are equivalent.
- The Prepare tool does not use the cache — it must always see fresh
  list data.
