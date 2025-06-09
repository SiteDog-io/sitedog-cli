#!/bin/sh

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Удаляем бинарники из пользовательских директорий
if [ -f "$HOME/bin/sitedog" ]; then
    rm -f "$HOME/bin/sitedog"
    echo "Removed $HOME/bin/sitedog"
fi
if [ -f "$HOME/.local/bin/sitedog" ]; then
    rm -f "$HOME/.local/bin/sitedog"
    echo "Removed $HOME/.local/bin/sitedog"
fi

# Удаляем всю директорию ~/.sitedog
if [ -d "$HOME/.sitedog" ]; then
    rm -rf "$HOME/.sitedog"
    echo "Removed $HOME/.sitedog directory"
fi

# Удаляем ~/.sitedog/bin из PATH во всех популярных rc-файлах
for RC in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile" "$HOME/.bash_profile" "$HOME/.config/fish/config.fish"; do
    if [ -f "$RC" ]; then
        # Определяем ОС для корректного использования sed
        case "$(uname)" in
            Darwin*)
                sed -i '' '/sitedog\/bin.*PATH/d' "$RC"
                sed -i '' '/# Added by sitedog installer/d' "$RC"
                ;;
            *)
                sed -i.bak '/sitedog\/bin.*PATH/d' "$RC"
                sed -i '/# Added by sitedog installer/d' "$RC"
                ;;
        esac
        # shellcheck disable=SC1090
        . "$RC"
    fi
done


echo "${GREEN}SiteDog has been fully uninstalled!${NC}" 