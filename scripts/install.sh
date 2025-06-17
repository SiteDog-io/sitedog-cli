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
curl -sL "$ASSET_URL" -o sitedog

# Check if file was downloaded
if [ ! -f sitedog ]; then
    echo -e "${RED}Error: Failed to download sitedog${NC}"
    exit 1
fi

# Download demo.html.tpl template из того же релиза
TPL_NAME="demo.html.tpl"
TPL_URL=$(curl -s "$API_URL" | grep 'browser_download_url' | grep "$TPL_NAME" | head -n1 | cut -d '"' -f 4)

if [ -z "$TPL_URL" ]; then
    echo -e "${RED}Error: Could not find asset $TPL_NAME in the latest release${NC}"
    exit 1
fi

echo "Downloading $TPL_NAME from $TPL_URL..."
curl -sL "$TPL_URL" -o demo.html.tpl

# Check if template was downloaded
if [ ! -f demo.html.tpl ]; then
    echo -e "${RED}Error: Failed to download demo.html.tpl${NC}"
    exit 1
fi

# Make file executable
chmod +x sitedog

# Install binary to /usr/local/bin
echo "Installing sitedog to /usr/local/bin (may require sudo)..."
sudo cp sitedog /usr/local/bin/sitedog
sudo chmod +x /usr/local/bin/sitedog

echo "Installed sitedog to /usr/local/bin/sitedog"

# Create templates directory and copy demo.html.tpl
TEMPLATES_DIR="$HOME/.sitedog"
mkdir -p "$TEMPLATES_DIR"
cp demo.html.tpl "$TEMPLATES_DIR/"

echo "Installed demo.html.tpl to $TEMPLATES_DIR/demo.html.tpl"

# Clean up temporary directory
cd - > /dev/null
rm -rf "$TMP_DIR"

echo "${GREEN}SiteDog has been installed successfully!${NC}"
echo "Try: sitedog help" 