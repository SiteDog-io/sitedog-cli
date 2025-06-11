#!/bin/sh

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Remove binaries from user directories
if [ -f "$HOME/bin/sitedog" ]; then
    rm -f "$HOME/bin/sitedog"
    echo "Removed $HOME/bin/sitedog"
fi
if [ -f "$HOME/.local/bin/sitedog" ]; then
    rm -f "$HOME/.local/bin/sitedog"
    echo "Removed $HOME/.local/bin/sitedog"
fi

# Remove entire ~/.sitedog directory
if [ -d "$HOME/.sitedog" ]; then
    rm -rf "$HOME/.sitedog"
    echo "Removed $HOME/.sitedog directory"
fi

# Remove ~/.sitedog/bin from PATH in all common rc files
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