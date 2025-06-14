#!/bin/sh

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux) BIN_NAME="sitedog-linux-$ARCH" ;;
    darwin) BIN_NAME="sitedog-darwin-$ARCH" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

echo "Detected platform: $OS/$ARCH"
echo "Downloading $BIN_NAME..."

curl -sL "https://gist.github.com/qelphybox/fe278d331980a1ce09c3d946bbf0b83b/raw/$BIN_NAME" -o sitedog

# Check if file was downloaded
if [ ! -f sitedog ]; then
    echo -e "${RED}Error: Failed to download sitedog${NC}"
    exit 1
fi

# Download demo.html.tpl template
echo "Downloading demo template..."
curl -sL https://gist.github.com/qelphybox/fe278d331980a1ce09c3d946bbf0b83b/raw/demo.html.tpl -o demo.html.tpl

# Check if template was downloaded
if [ ! -f demo.html.tpl ]; then
    echo -e "${RED}Error: Failed to download demo.html.tpl${NC}"
    exit 1
fi

# Make file executable
chmod +x sitedog

# Install binary to ~/.sitedog/bin
INSTALL_DIR="$HOME/.sitedog/bin"
mkdir -p "$INSTALL_DIR"
cp sitedog "$INSTALL_DIR/sitedog"
echo "Installed sitedog to $INSTALL_DIR/sitedog"

# Create templates directory and copy demo.html.tpl
TEMPLATES_DIR="$HOME/.sitedog"
mkdir -p "$TEMPLATES_DIR"
cp demo.html.tpl "$TEMPLATES_DIR/"

# Add ~/.sitedog/bin to PATH if not already there
SHELL_NAME=$(basename "$SHELL")
case "$SHELL_NAME" in
    zsh)
        RC_FILE="$HOME/.zshrc"
        ;;
    bash)
        RC_FILE="$HOME/.bashrc"
        ;;
    fish)
        RC_FILE="$HOME/.config/fish/config.fish"
        ;;
    *)
        RC_FILE="$HOME/.profile"
        ;;
esac

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "\n# Added by sitedog installer" >> "$RC_FILE"
    echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$RC_FILE"
    echo "${GREEN}Added $INSTALL_DIR to PATH in $RC_FILE${NC}"
    # shellcheck disable=SC1090
    . "$RC_FILE"
else
    echo "${YELLOW}$INSTALL_DIR already in PATH${NC}"
fi

# Clean up temporary directory
cd - > /dev/null
rm -rf "$TMP_DIR"

echo "${GREEN}SiteDog has been installed successfully!${NC}"
echo "Try: sitedog help" 