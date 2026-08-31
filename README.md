# pokemon-info

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/pokemon-info)

A terminal Pokédex. Run it bare for an interactive TUI, or stay on the
command line: type a few letters, hit TAB, and pick from every Pokémon
or move name — then get stats, types, abilities, evolution chains, move
effects, and flavor text straight from [PokeAPI](https://pokeapi.co).

```
$ pokemon-info                 # interactive TUI (default in a terminal)
$ pokemon-info pikachu          # classic one-shot CLI lookup
$ pokemon-info moves thunderbolt
```

## Interactive TUI

Running `pokemon-info` with no arguments in a terminal opens the TUI
(`pokemon-info ui` does the same explicitly; piped invocations print help
instead) straight into the Pokémon browser. Press `esc` to reach the
slash-command menu.

Slash commands:

- `/pokemons` — browse all Pokémon: fuzzy filter, embedded truecolor ASCII
  art, stat bars, type badges, abilities, evolution chain, and flavor text
  in the detail pane.
- `/moves` — browse all moves: type/damage class, power/accuracy/PP,
  effect, meta, stat changes, flavor text, learned-by count.
- `/help` — keybindings; `/quit` — exit.

Keybindings inside a browser: `↑↓`/`jk` move, `/` filter, `enter` reload
details, `f` cycle Pokémon art (small → large → off), `esc` back to the
menu, `q` quit. Details are fetched live from PokeAPI as you browse (the
name catalog and art are embedded, so they work offline). Art wider than
the detail pane is hidden automatically.

## Credits

- [PokeAPI](https://pokeapi.co) — all live data, and the embedded name
  catalog.
- [poketex](https://github.com/ckaznable/poketex) by ckaznable — the TUI
  browsing experience this interface is inspired by, and the tree we
  bundled the gen 9 art from.
- [pokemon-colorscripts](https://gitlab.com/phoneybadber/pokemon-colorscripts)
  by phoneybadber (MIT) — the embedded ANSI Pokémon art (gen 1–8).
- Caruban's generation 9 resource pack (bundled via
  [poketex #58](https://github.com/ckaznable/poketex/pull/58)) — gen 9 art.

## Command line

```
$ pokemon-info pi<TAB>
$ pokemon-info Pikachu

$ pokemon-info moves thunder<TAB>
$ pokemon-info moves thunderbolt
```

Name completion is embedded in the binary, so TAB is instant and works
offline. Details are fetched live from PokeAPI when you look something up.

## What it shows

- **Pokémon**: base stats table (with total), types, abilities, height,
  weight, base experience, official artwork link, and the full evolution
  chain with evolution requirements.
- **Moves**: type, damage class, power, accuracy, PP, priority, generation,
  effect text, meta (ailment, drain, flinch chance...), stat changes, flavor
  text, and how many Pokémon learn it.

Both English display names ("Thunderbolt") and PokeAPI slugs
("thunderbolt", "mr-mime") are accepted as input. A unique name prefix
resolves too ("zek" → Zekrom); ambiguous prefixes list their candidates.

## Install

From the Snap Store:

```
sudo snap install pokemon-info
```

From source:

```
git clone https://github.com/pkong-ds/pokemon-info
cd pokemon-info
make build
# put dist/pokemon-info somewhere on your PATH
```

## Shell completion

Completion candidates come from the binary itself, so every shell gets the
same instant, offline suggestions.

**zsh** (drops the completion into a directory already on your `fpath` — no
`.zshrc` edits needed):

```
pokemon-info completion zsh > "${fpath[1]}/_pokemon-info"
```

**bash** (non-snap installs; the snap wires bash completion automatically):

```
pokemon-info completion bash > ~/.local/share/bash-completion/completions/pokemon-info
```

**fish**:

```
pokemon-info completion fish > ~/.config/fish/completions/pokemon-info.fish
```

If zsh completions don't appear after (re)installing, clear the completion
cache and restart your shell: `rm -f ~/.zcompdump*` then `exec zsh`.

## Data freshness

The name catalog is baked into the binary at build time. When a new Pokémon
generation releases, names appear in completions after a `snap refresh` (or a
rebuilt binary). Fetched details are always live.

## Maintainers

To regenerate the embedded catalog from PokeAPI:

```
make catalog
```

This runs the Prepare tool (`cmd/prepare`), a maintainer-only command that
never ships in the snap. See [CONTEXT.md](CONTEXT.md) for the project's
shared vocabulary and [docs/adr/](docs/adr/) for the decisions behind it.

## License

[MIT](LICENSE)
