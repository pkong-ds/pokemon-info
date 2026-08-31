# pokemon-info

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/pokemon-info)

A terminal Pokédex: an interactive TUI for browsing, and a one-shot CLI
for lookups — every Pokémon and every move, with data from
[PokeAPI](https://pokeapi.co).

![TUI demo](docs/demo.gif)

```
$ pokemon-info                 # interactive TUI (default in a terminal)
$ pokemon-info pikachu          # one-shot CLI lookup
$ pokemon-info moves thunderbolt
```

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

## Interactive TUI

Bare `pokemon-info` opens the TUI straight into the Pokémon browser;
`esc` reaches the slash-command menu (`/pokemons`, `/moves`, `/help`,
`/quit`). Browse with `↑↓`/`jk`, filter with `/`, cycle Pokémon art
with `f`. The name catalog and art are embedded, so browsing and
filtering work offline; details are fetched live from PokeAPI. Press
`/help` in-app for every keybinding.

![Moves browser](docs/moves.gif)

The moves browser is one `esc` away: down to `/moves`, enter — then the
same live `/`-filtering (typing narrows as you go, details follow the
selection).

## Command line

```
$ pokemon-info pi<TAB>
$ pokemon-info Pikachu

$ pokemon-info moves thunder<TAB>
$ pokemon-info moves thunderbolt
```

English display names ("Thunderbolt"), PokeAPI slugs ("mr-mime"), and
unique prefixes ("zek" → Zekrom) all resolve; ambiguous prefixes list
their candidates. Completion candidates are embedded in the binary, so
TAB is instant and offline. Details are fetched live from PokeAPI.

## Shell completion

Completion candidates come from the binary itself, so every shell gets
the same instant, offline suggestions.

**zsh** (drops into a directory already on your `fpath` — no `.zshrc`
edits needed):

```
pokemon-info completion zsh > "${fpath[1]}/_pokemon-info"
```

**bash** (non-snap installs; the snap wires bash completion
automatically):

```
pokemon-info completion bash > ~/.local/share/bash-completion/completions/pokemon-info
```

**fish**:

```
pokemon-info completion fish > ~/.config/fish/completions/pokemon-info.fish
```

If zsh completions don't appear after (re)installing, clear the
completion cache and restart your shell: `rm -f ~/.zcompdump*` then
`exec zsh`.

## Data freshness

The name catalog is baked into the binary at build time. When a new
Pokémon generation releases, names appear in completions after a `snap
refresh` (or a rebuilt binary). Fetched details are always live.

## Credits

- [PokeAPI](https://pokeapi.co) — all live data, and the embedded name
  catalog.
- [poketex](https://github.com/ckaznable/poketex) by ckaznable — the TUI
  browsing experience this interface is inspired by, and the tree we
  bundled the gen 9 art from.
- [pokemon-colorscripts](https://gitlab.com/phoneybadger/pokemon-colorscripts)
  by phoneybadber (MIT) — the embedded ANSI Pokémon art (gen 1–8).
- Caruban's generation 9 resource pack (bundled via
  [poketex #58](https://github.com/ckaznable/poketex/pull/58)) — gen 9
  art.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) — the glossary lives in
[CONTEXT.md](CONTEXT.md), decisions in [docs/adr/](docs/adr/), and the
demo GIFs are re-recorded with `make demos`.

## License

[MIT](LICENSE)
