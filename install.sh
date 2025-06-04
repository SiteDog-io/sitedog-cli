#!/bin/sh

set -e

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Создаем временную директорию
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Скачиваем sitedog
echo "Downloading sitedog..."
curl -sL https://gist.github.com/qelphybox/fe278d331980a1ce09c3d946bbf0b83b/raw/sitedog -o sitedog

# Проверяем, что файл скачался
if [ ! -f sitedog ]; then
    echo -e "${RED}Error: Failed to download sitedog${NC}"
    exit 1
fi

# Скачиваем шаблон demo.html.erb
echo "Downloading demo template..."
curl -sL https://gist.github.com/qelphybox/fe278d331980a1ce09c3d946bbf0b83b/raw/demo.html.erb -o demo.html.erb

# Проверяем, что шаблон скачался
if [ ! -f demo.html.erb ]; then
    echo -e "${RED}Error: Failed to download demo.html.erb${NC}"
    exit 1
fi

# Делаем файл исполняемым
chmod +x sitedog

# Устанавливаем бинарник в ~/.sitedog/bin
INSTALL_DIR="$HOME/.sitedog/bin"
mkdir -p "$INSTALL_DIR"
cp sitedog "$INSTALL_DIR/sitedog"
echo "Installed sitedog to $INSTALL_DIR/sitedog"

# Создаем директорию для шаблонов и копируем demo.html.erb
TEMPLATES_DIR="$HOME/.sitedog"
mkdir -p "$TEMPLATES_DIR"
cp demo.html.erb "$TEMPLATES_DIR/"

# Добавляем ~/.sitedog/bin в PATH, если его там нет
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

# Очищаем временную директорию
cd - > /dev/null
rm -rf "$TMP_DIR"

echo "${GREEN}SiteDog has been installed successfully!${NC}"
echo "Try: sitedog help" 