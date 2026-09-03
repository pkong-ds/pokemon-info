---
name: record-demos
description: Use when the README demo GIFs are stale — after any TUI or CLI output change that alters what the demos show.
---

# Record the demo GIFs

The README embeds two GIFs — `docs/demo.gif` (Pokémon browser) and
`docs/moves.gif` (moves browser) — recorded with
[vhs](https://github.com/charmbracelet/vhs) from the tapes in
`docs/tapes/`.

## When to re-record

Whenever the TUI layout, colors, or CLI output changes shape. The GIFs
are the README's main pitch; stale demos misrepresent the tool.

## Steps

1. Install vhs if needed: `brew install vhs`.
2. Build a current binary: `make build` (the tapes invoke the built
   binary).
3. Re-record both tapes:

   ```
   make demos
   ```

   which runs `vhs docs/tapes/demo.tape` and `vhs docs/tapes/moves.tape`,
   overwriting `docs/demo.gif` and `docs/moves.gif` in place.
4. Check the new GIFs actually show the intended flow (the tapes drive
   the TUI; timing can break if output speed changed).
5. Commit the updated GIFs together with the change that altered the
   output.
