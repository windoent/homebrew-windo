class Windo < Formula
  desc "A CLI tool for generating kratos project templates"
  homepage "https://github.com/windoent/homebrew-windo"

  version "0.0.1"

  on_macos do
    on_intel do
      url "https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-darwin-amd64"
      sha256 "ac9986f44fd699e095652ecd6674d2db81038e97d08d6d023c9eb4d3f046283a"
    end

    on_arm do
      url "https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-darwin-arm64"
      sha256 "000727afd77132af9e120e5a8eb84151a15db945140ff593d7dc0d102ccbed71"
    end
  end

  on_linux do
    url "https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-linux-amd64"
    sha256 "30271cce73115e51c2101971ad4b404c699b3fd9d277dfc11f186b8d85d3e8d2"
  end

  def install
    if OS.mac? && Hardware::CPU.intel?
      bin.install "windo-darwin-amd64" => "windo"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "windo-darwin-arm64" => "windo"
    elsif OS.linux?
      bin.install "windo-linux-amd64" => "windo"
    end
  end

  test do
    system "#{bin}/windo", "version"
  end
end