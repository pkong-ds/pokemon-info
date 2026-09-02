class PokemonInfo < Formula
  desc "Terminal Pokedex with TUI browsing and CLI lookups"
  homepage "https://github.com/pkong-ds/pokemon-info"
  version "0.3.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/pkong-ds/pokemon-info/releases/download/v0.3.0/pokemon-info-0.3.0-darwin-arm64.tar.gz"
      sha256 "bebd81aca711b1a550c6f05e51343fb5f2310da9a74d4b3eb08b2eacce2201f5"
    end
    on_intel do
      url "https://github.com/pkong-ds/pokemon-info/releases/download/v0.3.0/pokemon-info-0.3.0-darwin-amd64.tar.gz"
      sha256 "2253b3fe619571252895b764e4456d2f33d4c889b1128bdcf3c1f4546329a906"
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
