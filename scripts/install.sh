#!/usr/bin/env sh

# references https://github.com/anomalyco/opentui/blob/59c8a83680a4357570b108ccfa6b472353968a15/packages/core/src/examples/install.sh

set -e

echo "Installing cURLy! The snake game with cURL"
echo ""

REPO="gjtiquia/cURLy"
GITHUB_API="https://api.github.com/repos/$REPO"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
x86_64 | amd64) ARCH="amd64" ;;
arm64 | aarch64) ARCH="arm64" ;;
esac

case "$OS" in
darwin) OS="darwin" ;;
linux) OS="linux" ;;
mingw* | cygwin* | msys*) OS="windows" ;;
esac

PLATFORM="${OS}/${ARCH}"
echo "Detected platform: $PLATFORM"

# Get the latest stable release
RELEASE_DATA=$(curl -s "$GITHUB_API/releases/latest")
VERSION=$(echo "$RELEASE_DATA" | grep '"tag_name"' | cut -d '"' -f 4)
if [ -z "$VERSION" ]; then
    printf "${RED}Error: Failed to fetch latest release information${NC}\n"
    exit 1
fi

# Remove 'v' prefix if present
VERSION_NO_V="${VERSION#v}"

printf "${BLUE}Version: $VERSION_NO_V${NC}\n"

# Construct download URL
ASSET_NAME="cURLy_${OS}_${ARCH}"
[ "$OS" = "windows" ] && ASSET_NAME="${ASSET_NAME}.exe"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/${VERSION}/${ASSET_NAME}"

# Download the file
echo ""
echo "Downloading $ASSET_NAME from $DOWNLOAD_URL..."
if ! curl -fsSLO "$DOWNLOAD_URL"; then
    printf "${RED}Error: Failed to download! Platform $PLATFORM unsupported!${NC}\n"
    exit 1
fi
echo "Download successful!"

# Rename downloaded file
EXEC_NAME="cURLy"
[ "$OS" = "windows" ] && EXEC_NAME="${EXEC_NAME}.exe"
mv "$ASSET_NAME" "$EXEC_NAME"

# Make executable (if not on Windows)
if [ "$OS" != "windows" ]; then
    chmod +x "$EXEC_NAME"
fi

# cannot run directly cuz when piping to bash it does not have any terminal info, ie. size=0x0
# echo "running $EXEC_NAME"
# ./$EXEC_NAME

echo ""
echo "Run the following command to play cURLy!"
echo ""
printf "${GREEN}./$EXEC_NAME${NC}\n"
echo ""
