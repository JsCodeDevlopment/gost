class Gost < Formula
  desc "A powerful, opinionated Go boilerplate engineered with a NestJS-like architecture."
  homepage "https://github.com/JsCodeDevlopment/gost"
  url "https://github.com/JsCodeDevlopment/gost/releases/download/v1.1.0/gost_darwin_amd64.tar.gz"
  version "1.1.0"
  sha256 "376075701e8421340189559048239f1518400199931839763139047d9128a48d"
  author "JsCodeDevlopment"
  license "MIT"
  

  def install
    bin.install "gost"
  end

  test do
    system "#{bin}/gost", "--version"
  end
end
