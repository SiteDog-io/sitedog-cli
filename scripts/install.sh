#!/bin/sh

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# Directories
INSTALL_DIR="$HOME/.sitedog/bin"
TEMPLATES_DIR="$HOME/.sitedog"
mkdir -p "$INSTALL_DIR"
mkdir -p "$TEMPLATES_DIR"

# Download latest release from GitHub (без jq)
REPO="SiteDog-io/sitedog-cli"
API_URL="https://api.github.com/repos/$REPO/releases/latest"

# Получаем ссылку на нужный бинарник из релиза (без jq)
ASSET_URL=$(curl -s "$API_URL" | grep 'browser_download_url' | grep "$BIN_NAME" | head -n1 | cut -d '"' -f 4)

if [ -z "$ASSET_URL" ]; then
    echo -e "${RED}Error: Could not find asset $BIN_NAME in the latest release${NC}"
    exit 1
fi

echo "Downloading $BIN_NAME from $ASSET_URL..."
curl -sL "$ASSET_URL" -o "$INSTALL_DIR/sitedog"

# Check if file was downloaded
if [ ! -f "$INSTALL_DIR/sitedog" ]; then
    echo -e "${RED}Error: Failed to download sitedog${NC}"
    exit 1
fi

# Make file executable
chmod +x "$INSTALL_DIR/sitedog"
echo "Installed sitedog to $INSTALL_DIR/sitedog"

# Download demo.html.tpl template
TPL_NAME="demo.html.tpl"
TPL_URL=$(curl -s "$API_URL" | grep 'browser_download_url' | grep "$TPL_NAME" | head -n1 | cut -d '"' -f 4)

if [ -z "$TPL_URL" ]; then
    echo -e "${RED}Error: Could not find asset $TPL_NAME in the latest release${NC}"
    exit 1
fi

echo "Downloading $TPL_NAME from $TPL_URL..."
curl -sL "$TPL_URL" -o "$TEMPLATES_DIR/demo.html.tpl"

# Check if template was downloaded
if [ ! -f "$TEMPLATES_DIR/demo.html.tpl" ]; then
    echo -e "${RED}Error: Failed to download demo.html.tpl${NC}"
    exit 1
fi

echo "Installed demo.html.tpl to $TEMPLATES_DIR/demo.html.tpl"

# Try to create symlink in /usr/local/bin
SYMLINK_OK=0
if [ -w /usr/local/bin ]; then
    ln -sf "$INSTALL_DIR/sitedog" /usr/local/bin/sitedog && SYMLINK_OK=1
else
    if sudo ln -sf "$INSTALL_DIR/sitedog" /usr/local/bin/sitedog; then
        SYMLINK_OK=1
    fi
fi

if [ $SYMLINK_OK -eq 1 ]; then
    echo "Symlink created: /usr/local/bin/sitedog -> $INSTALL_DIR/sitedog"
else
    echo "\033[0;33mNo permissions to create symlink in /usr/local/bin.\033[0m"
    echo "Please add $INSTALL_DIR to your PATH. For example, add this line to your shell rc file (e.g., ~/.bashrc or ~/.zshrc):"
    echo 'export PATH="$HOME/.sitedog/bin:$PATH"'
fi

echo "${GREEN}SiteDog has been installed successfully!${NC}"
echo "Try: sitedog help" 