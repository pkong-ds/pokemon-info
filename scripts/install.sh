#!/bin/sh
# pokemon-info installer — downloads the latest release from GitHub.
# Usage: curl -fsSL https://raw.githubusercontent.com/pkong-ds/pokemon-info/main/scripts/install.sh | sh
# Override the install directory with: BIN_DIR=... curl ... | sh
set -eu

repo="pkong-ds/pokemon-info"

os=$(uname -s)
case "$os" in
    Darwin) os="darwin" ;;
    Linux) os="linux" ;;
    *)
        echo "error: unsupported OS ($os). On Windows use Scoop:" >&2
        echo "  scoop bucket add pkong-ds https://github.com/pkong-ds/scoop-bucket" >&2
        echo "  scoop install pokemon-info" >&2
        exit 1
        ;;
esac

arch=$(uname -m)
case "$arch" in
    x86_64) arch="amd64" ;;
    aarch64 | arm64) arch="arm64" ;;
    *) echo "error: unsupported architecture ($arch)" >&2; exit 1 ;;
esac

resolved=$(curl -fsSLI -o /dev/null -w '%{url_effective}' "https://github.com/${repo}/releases/latest")
case "$resolved" in
    *"/releases/tag/"*) tag=${resolved##*/} ;;
    *)
        echo "error: no releases published yet — build from source: https://github.com/${repo}#install" >&2
        exit 1
        ;;
esac
version=${tag#v}

archive="pokemon-info-${version}-${os}-${arch}.tar.gz"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading ${archive} (${tag})..."
curl -fsSL -o "${tmp}/${archive}" "https://github.com/${repo}/releases/download/${tag}/${archive}"
curl -fsSL -o "${tmp}/checksums.txt" "https://github.com/${repo}/releases/download/${tag}/checksums.txt"

expected=$(grep " ${archive}\$" "${tmp}/checksums.txt" | cut -d' ' -f1)
if [ -z "$expected" ]; then
    echo "error: ${archive} not found in checksums.txt" >&2
    exit 1
fi
printf '%s  %s\n' "$expected" "$archive" | {
    if command -v sha256sum >/dev/null 2>&1; then (cd "$tmp" && sha256sum -c - >/dev/null)
    else (cd "$tmp" && shasum -a 256 -c - >/dev/null); fi
}

bin_dir=${BIN_DIR:-"${HOME}/.local/bin"}
mkdir -p "$bin_dir"
tar -xzf "${tmp}/${archive}" -C "$tmp"
mv "${tmp}/pokemon-info-${version}-${os}-${arch}/pokemon-info" "${bin_dir}/pokemon-info"
chmod +x "${bin_dir}/pokemon-info"

echo "Installed: ${bin_dir}/pokemon-info ($("${bin_dir}/pokemon-info" version))"
case ":${PATH}:" in
    *":${bin_dir}:"*) ;;
    *) echo "warning: ${bin_dir} is not on your PATH — add it to your shell profile" >&2 ;;
esac

echo "Shell completion (optional):"
echo "  zsh:  pokemon-info completion zsh > \"\${fpath[1]}/_pokemon-info\""
echo "  bash: pokemon-info completion bash > ~/.local/share/bash-completion/completions/pokemon-info"
echo "  fish: pokemon-info completion fish > ~/.config/fish/completions/pokemon-info.fish"
