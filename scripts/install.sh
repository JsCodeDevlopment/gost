#!/bin/sh
set -e

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
esac

BINARY_NAME="gost"
if [ "$OS" = "mingw64_nt-10.0" ] || [ "$OS" = "msys_nt-10.0" ]; then
    OS="windows"
    BINARY_NAME="gost.exe"
fi

REPO="JsCodeDevlopment/gost"
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "❌ Failed to fetch latest release from $REPO"
    exit 1
fi

echo "🚀 Installing Gost CLI $LATEST_TAG for $OS-$ARCH..."

URL="https://github.com/$REPO/releases/download/$LATEST_TAG/gost_${OS}_${ARCH}"
if [ "$OS" = "windows" ]; then
    URL="${URL}.exe"
fi

TEMP_BIN="/tmp/$BINARY_NAME"
curl -L -o "$TEMP_BIN" "$URL"
chmod +x "$TEMP_BIN"

INSTALL_DIR="/usr/local/bin"
if [ "$OS" = "windows" ]; then
    echo "✅ Downloaded to $TEMP_BIN. Please move it to your PATH."
else
    sudo mv "$TEMP_BIN" "$INSTALL_DIR/$BINARY_NAME"
    echo "✅ Gost CLI installed successfully to $INSTALL_DIR/$BINARY_NAME"
fi

echo "✨ Try running: gost init my-project"
