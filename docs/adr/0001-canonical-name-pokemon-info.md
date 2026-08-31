# ADR 0001: The tool is named `pokemon-info`, not `poke`

## Status

Accepted

## Context

The core idea was pitched as `poke <pokemon-name-fragment>` with autocomplete.
Three constraints collided with that name:

1. Snap names only allow lowercase letters, digits, and hyphens. The existing
   binary name `pokemon_info` (underscore) is illegal as a snap name, so a
   rename was forced regardless.
2. `poke` collides with GNU poke, a well-known binary editor. A namespace
   fight risks confusing users searching for either tool.
3. The primary goal of this project is to experience the snap store upload
   flow end-to-end. A name dispute would derail that goal for zero gain.

The snap name, GitHub repo name, binary name, and zsh compdef name should all
be one name to avoid a spread of aliases nobody can remember.

## Decision

The canonical name is `pokemon-info`, used by the snap, the GitHub repository,
the binary, and generated shell completions.

## Consequences

- No collision risk; name is available for store registration.
- The product feel comes from the autocomplete experience, not from a short
  command name; typing `pokemon-info pi<TAB>` is acceptable.
- If a shorter name is ever wanted, a shell alias is cheap to document and
  does not require a rename anywhere.
- The old name `pokemon_info` appears in git history, the local directory name,
  and the generated YAML/Makefiles and must be updated before release.
