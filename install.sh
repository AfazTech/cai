#!/usr/bin/env bash

set -e

REPO="AfazTech/cai"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cai"

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)
        OS_NAME="linux"
        ;;

    Darwin)
        OS_NAME="darwin"
        ;;

    *)
        echo "Error: Unsupported operating system: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64|amd64)
        ARCH_NAME="amd64"
        ;;

    aarch64|arm64)
        ARCH_NAME="arm64"
        ;;

    *)
        echo "Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

URL="https://github.com/${REPO}/releases/latest/download/cai-${OS_NAME}-${ARCH_NAME}"

echo "Downloading CAI..."
echo "OS: $OS_NAME"
echo "Architecture: $ARCH_NAME"

TMP_FILE="$(mktemp)"

curl -fL "$URL" -o "$TMP_FILE"

chmod +x "$TMP_FILE"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
fi

echo ""
echo "CAI installed successfully!"
echo ""
echo "Run:"
echo "  cai"
