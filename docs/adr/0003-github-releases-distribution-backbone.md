# ADR 0003: GitHub Releases as the distribution backbone

## Status

Accepted

## Context

Through v0.2.x the snap was the only install path, and only on
`latest/edge`: Linux-only, unstable-channel flag, and the README showed
`sudo snap install pokemon-info` without `--edge` — a command that fails.
The snap upload flow has now been experienced end-to-end (the original
primary goal, ADR 0001); the next goal is users, which means macOS and
Windows installs, auto-updates, and a working one-liner.

The binary is a single static Go executable (`CGO_ENABLED=0`), so packaging
is cross-compiling a handful of archives and handing them to each
platform's package manager.

Alternatives considered:

1. **Share first, package later.** Visitors who cannot install in seconds
   bounce; a Show HN moment is one-shot.
2. **Homebrew core / Scoop extras directly.** Both have acceptance bars
   (notability: stars, forks, adoption); a day-old repository does not
   qualify.
3. **goreleaser.** Does the same job in one config, but hand-rolling the
   release workflow is closer to this project's learn-the-flow goal and has
   fewer moving parts.

## Decision

GitHub Releases are the backbone: a `v*.*.*` tag triggers a workflow that
cross-compiles (linux/darwin × amd64/arm64, windows/amd64), attaches the
archives plus `checksums.txt`, and creates the release. Every other
channel consumes those assets:

- `scripts/install.sh` — the curl one-liner. Resolves the latest tag via
  the `releases/latest` redirect (no API calls), verifies sha256 against
  `checksums.txt`, installs to `~/.local/bin` (overridable with `BIN_DIR`).
- The Homebrew tap `pkong-ds/homebrew-tap` — self-published formula; wires
  zsh/bash/fish completions at install time via
  `generate_completions_from_executable`.
- The Scoop bucket `pkong-ds/scoop-bucket` — self-published manifest with
  checkver against this repository's releases.
- The snap graduates to the `stable` channel and remains the
  Linux-package-manager path.

Submissions to homebrew-core / scoop-extras are deferred until the project
has the notability they require.

## Consequences

- `pokemon-info version` is the single source of the running version in
  every channel; it is ldflags-injected (`main.version`), and snap builds
  inject `CRAFT_PROJECT_VERSION`.
- Catalog staleness (ADR 0002) is fixed by any channel's upgrade: `snap
  refresh`, `brew upgrade`, `scoop update`, or re-running the installer.
- Releasing is: tag `vX.Y.Z`, push (the workflow publishes the release),
  fill the tap formula and scoop manifest hashes from `checksums.txt`,
  then publish the snap to stable.
- The tap and bucket exist only because homebrew-core and scoop-extras are
  out of reach; they retire when the project is accepted upstream.
