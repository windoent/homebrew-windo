class Windo < Formula
  desc "A CLI tool for generating kratos project templates"
  homepage "https://git.code.tencent.com/windo-/homebrew-windo"

  version "0.0.1"

  on_macos do
    on_intel do
      url "https://git.code.tencent.com/windo-/homebrew-windo/uploads/2/tags/untagged-FFB5ADF1FDB14C09B7BE4B513B9646BC/windo-darwin-amd64"
      sha256 "14381d1aaa3fbc46ad624a255a021f0c1610c9de0f8c26c7e4202abe1dd6cad1"
    end

    on_arm do
      url "https://git.code.tencent.com/windo-/homebrew-windo/uploads/2/tags/untagged-DEFD94CBAF4541EBAFAD1F62A2AEBC5A/windo-darwin-arm64"
      sha256 "691d0ce95c4ccd5a5fa6760fc5bab39b517886515a19a51636f53aaac2cd9cbf"
    end
  end

  on_linux do
    url "https://git.code.tencent.com/windo-/homebrew-windo/uploads/2/tags/untagged-E9F92923BFAB47528CEE884A795CEF4A/windo-linux-amd64"
    sha256 "555cfd52a27eba327fa6223653295de4b0c6288248cbabf58f93bc201d6ce5ba"
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