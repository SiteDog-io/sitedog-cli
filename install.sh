#!/bin/sh

set -e

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Определяем тип системы
if [ "$(uname)" = "Darwin" ]; then
    IS_MACOS=true
else
    IS_MACOS=false
fi

# Проверяем наличие Ruby
if ! ruby -v >/dev/null 2>&1; then
    echo -e "${RED}Error: Ruby is not installed${NC}"
    if [ "$IS_MACOS" = true ]; then
        echo "On macOS, Ruby comes pre-installed. If you're seeing this error, please install Ruby:"
        echo "brew install ruby"
    else
        echo "Please install Ruby first: https://www.ruby-lang.org/en/downloads/"
    fi
    exit 1
fi

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

# Определяем директорию для установки
if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
    cp sitedog "$INSTALL_DIR/sitedog"
else
    # Устанавливаем в домашнюю директорию
    if [ "$IS_MACOS" = true ]; then
        INSTALL_DIR="$HOME/bin"
    else
        INSTALL_DIR="$HOME/.local/bin"
    fi
    mkdir -p "$INSTALL_DIR"
    cp sitedog "$INSTALL_DIR/sitedog"
    
    # Определяем конфигурационный файл оболочки
    if [ -f "$HOME/.zshrc" ]; then
        SHELL_RC="$HOME/.zshrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        SHELL_RC="$HOME/.bash_profile"
    elif [ -f "$HOME/.bashrc" ]; then
        SHELL_RC="$HOME/.bashrc"
    else
        SHELL_RC="$HOME/.bash_profile"
    fi

    # Проверяем, что директория в PATH
    if echo ":$PATH:" | grep -q ":$INSTALL_DIR:"; then
        :
    else
        echo -e "${YELLOW}Warning: $INSTALL_DIR is not in your PATH${NC}"
        echo "Add this line to your $SHELL_RC:"
        echo -e "${GREEN}export PATH=\"\$PATH:$INSTALL_DIR\"${NC}"
        echo "Then run: source $SHELL_RC"
    fi
fi

# Создаем директорию для шаблонов и копируем demo.html.erb
TEMPLATES_DIR="$HOME/.sitedog"
mkdir -p "$TEMPLATES_DIR"
cp demo.html.erb "$TEMPLATES_DIR/"

# Очищаем временную директорию
cd - > /dev/null
rm -rf "$TMP_DIR"

echo -e "${GREEN}SiteDog has been installed successfully!${NC}"
echo "You can now use 'sitedog' command from anywhere."
echo "Try: sitedog help" 