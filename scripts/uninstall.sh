#!/bin/sh

set -e

# Colors for output
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Remove binary or symlink from /usr/local/bin
if [ -L "/usr/local/bin/sitedog" ] || [ -f "/usr/local/bin/sitedog" ]; then
    echo "Removing /usr/local/bin/sitedog (may require sudo)..."
    sudo rm -f /usr/local/bin/sitedog
    echo "Removed /usr/local/bin/sitedog"
fi

# Remove demo.html.tpl
if [ -f "$HOME/.sitedog/demo.html.tpl" ]; then
    rm -f "$HOME/.sitedog/demo.html.tpl"
    echo "Removed $HOME/.sitedog/demo.html.tpl"
fi

# Remove ~/.sitedog/bin directory if empty
if [ -d "$HOME/.sitedog/bin" ] && [ ! "$(ls -A $HOME/.sitedog/bin)" ]; then
    rmdir "$HOME/.sitedog/bin"
    echo "Removed empty $HOME/.sitedog/bin directory"
fi

# Remove ~/.sitedog directory if empty
if [ -d "$HOME/.sitedog" ] && [ ! "$(ls -A $HOME/.sitedog)" ]; then
    rmdir "$HOME/.sitedog"
    echo "Removed empty $HOME/.sitedog directory"
fi

# Remove ~/.sitedog/bin from PATH in all common rc files (backward compatibility)
for RC in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile" "$HOME/.bash_profile" "$HOME/.config/fish/config.fish"; do
    if [ -f "$RC" ]; then
        # Determine OS for correct sed usage
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