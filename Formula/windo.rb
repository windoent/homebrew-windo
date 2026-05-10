class Windo < Formula
  desc "A CLI tool for generating kratos project templates"
  homepage "https://github.com/windoent/homebrew-windo"

  version "0.0.1"

  on_macos do
    on_intel do
      url "https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-darwin-amd64"
      sha256 "640e1db029a9899ae449cc2ced29e6c090e4a10da8e9a3154ad6cb5ca68f3349"
    end

    on_arm do
      url "https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-darwin-arm64"
      sha256 "5474936dd95f71006fa18bc0aca42b5184a781cc488a83be7dfbd28c4f379640"
    end
  end

  on_linux do
    url "https://github.com/windoent/homebrew-windo/releases/download/v0.0.1/windo-linux-amd64"
    sha256 "e692cf3533b7220789befc2498129dba3cbd9bf18b1a1e8b2cf421eabb67ff2a"
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