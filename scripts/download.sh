#!/bin/bash

REPO="studio-b12/goat"
NAME="goat"

VERSION="$1"

if [ -d "$HOME/.local/bin" ]; then
    OUTDIR="$HOME/.local/bin/$NAME"
else
    if [ "$(id -u)" != "0" ]; then
        echo "This script needs to be ran as root."
        exit 1
    fi
    OUTDIR="/usr/local/bin/$NAME"
fi

check_installed() {
    which "$1" > /dev/null 2>&1 || {
        echo "'$1' needs to be installed to run this script."
        exit 1
    }
}

check_installed "curl"
check_installed "jq"

if [ -z "$VERSION" ]; then
    VERSION=$(curl -Ls "https://api.github.com/repos/$REPO/releases" | jq -r ".[0].tag_name")
fi

os=$(uname -s)
arch=$(uname -m)

case "$os $arch" in

    "Linux x86_64") 
        assetv="linux-amd64" ;;
    "Linux aarch64") 
        assetv="linux-arm64" ;;
    "Darwin x86_64") 
        assetv="darwin-amd64" ;;
    "Darwin aarch64") 
        assetv="darwin-arm64" ;;
    "Windows x86_64") 
        assetv="windows-amd64.exe" ;;
    "Windows aarch64") 
        assetv="windows-arm64.exe" ;;
    *)
        echo "Unsupported OS/Arch kombination: $os/$arch"
        exit 1 
        ;;

esac

curl -Lso "$OUTDIR" "https://github.com/$REPO/releases/download/$VERSION/$NAME-$assetv"
chmod +x "$OUTDIR"
