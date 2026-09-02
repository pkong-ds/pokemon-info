class PokemonInfo < Formula
  desc "Terminal Pokedex with TUI browsing and CLI lookups"
  homepage "https://github.com/pkong-ds/pokemon-info"
  version "0.3.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/pkong-ds/pokemon-info/releases/download/v0.3.0/pokemon-info-0.3.0-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_OF_pokemon-info-0.3.0-darwin-arm64.tar.gz"
    end
    on_intel do
      url "https://github.com/pkong-ds/pokemon-info/releases/download/v0.3.0/pokemon-info-0.3.0-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_OF_pokemon-info-0.3.0-darwin-amd64.tar.gz"
    end
  end

  def install
    bin.install "pokemon-info"
    generate_completions_from_executable(bin/"pokemon-info", "completion", shells: [:bash, :zsh, :fish])
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/pokemon-info version")
  end
end
