#!/bin/bash

set -e

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Определяем тип системы
if [[ "$OSTYPE" == "darwin"* ]]; then
    IS_MACOS=true
else
    IS_MACOS=false
fi

# Определяем директорию установки
if [ -f "/usr/local/bin/sitedog" ]; then
    INSTALL_DIR="/usr/local/bin"
    echo -e "${YELLOW}Found sitedog in $INSTALL_DIR${NC}"
    sudo rm -f "$INSTALL_DIR/sitedog"
    echo -e "${GREEN}Removed sitedog from $INSTALL_DIR${NC}"
else
    # Проверяем в домашней директории
    if [ "$IS_MACOS" = true ]; then
        INSTALL_DIR="$HOME/bin"
    else
        INSTALL_DIR="$HOME/.local/bin"
    fi

    if [ -f "$INSTALL_DIR/sitedog" ]; then
        echo -e "${YELLOW}Found sitedog in $INSTALL_DIR${NC}"
        rm -f "$INSTALL_DIR/sitedog"
        echo -e "${GREEN}Removed sitedog from $INSTALL_DIR${NC}"
    else
        echo -e "${YELLOW}sitedog not found in $INSTALL_DIR${NC}"
    fi
fi

# Удаляем директорию с шаблонами
TEMPLATES_DIR="$HOME/.sitedog"
if [ -d "$TEMPLATES_DIR" ]; then
    echo -e "${YELLOW}Found templates in $TEMPLATES_DIR${NC}"
    rm -rf "$TEMPLATES_DIR"
    echo -e "${GREEN}Removed templates directory${NC}"
else
    echo -e "${YELLOW}No templates directory found${NC}"
fi

# Удаляем Go-бинарник
if [ -f "/usr/local/bin/sitedog" ]; then
    sudo rm -f /usr/local/bin/sitedog
    echo "Removed Go binary from /usr/local/bin/sitedog"
fi
if [ -f "$HOME/.local/bin/sitedog" ]; then
    rm -f "$HOME/.local/bin/sitedog"
    echo "Removed Go binary from $HOME/.local/bin/sitedog"
fi

echo -e "${GREEN}SiteDog has been uninstalled successfully!${NC}" 