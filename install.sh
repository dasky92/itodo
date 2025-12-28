#!/bin/bash

set -e

OWNER="dasky92"
REPO="itodo"
BINARY_NAME="itodo"
GITHUB_REPO="$OWNER/$REPO"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Detect OS and Arch
OS=$(uname -s)
ARCH=$(uname -m)

case $OS in
    Linux)
        OS_TYPE="Linux"
        ;;
    Darwin)
        OS_TYPE="Darwin"
        ;;
    *)
        error "Unsupported operating system: $OS"
        ;;
esac

case $ARCH in
    x86_64)
        ARCH_TYPE="x86_64"
        ;;
    aarch64|arm64)
        ARCH_TYPE="aarch64"
        ;;
    *)
        error "Unsupported architecture: $ARCH"
        ;;
esac

log "Detected OS: $OS_TYPE, Arch: $ARCH_TYPE"

# Get latest release tag
log "Fetching latest release version..."
LATEST_RELEASE_URL="https://api.github.com/repos/$GITHUB_REPO/releases/latest"

# Use curl to get the latest release data
RELEASE_DATA=$(curl -sL "$LATEST_RELEASE_URL")

# Check if curl failed or if we got a "Not Found" response
if echo "$RELEASE_DATA" | grep -q "Not Found"; then
     error "Repository not found or no releases available."
fi

# Extract tag name using basic text processing (avoiding jq dependency)
TAG=$(echo "$RELEASE_DATA" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
    error "Could not find latest release tag. Please check if there are any releases in the repository."
fi

# Remove 'v' prefix for version number if present (assuming standard GoReleaser naming convention)
VERSION=${TAG#v}

log "Latest version: $TAG"

# Construct asset name
# Assuming GoReleaser default format: ProjectName_Version_Os_Arch.tar.gz
ASSET_NAME="${BINARY_NAME}_${VERSION}_${OS_TYPE}_${ARCH_TYPE}.tar.gz"
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$TAG/$ASSET_NAME"

TEMP_DIR=$(mktemp -d)
# Ensure cleanup happens on exit
trap 'rm -rf "$TEMP_DIR"' EXIT

log "Downloading $ASSET_NAME..."
HTTP_STATUS=$(curl -sL -w "%{http_code}" -o "$TEMP_DIR/$ASSET_NAME" "$DOWNLOAD_URL")

if [ "$HTTP_STATUS" -ne 200 ]; then
    error "Download failed with status $HTTP_STATUS. URL: $DOWNLOAD_URL"
fi

log "Extracting..."
tar -xzf "$TEMP_DIR/$ASSET_NAME" -C "$TEMP_DIR"

# Determine install directory
INSTALL_DIR="/usr/local/bin"
USE_SUDO=""

if [ ! -w "$INSTALL_DIR" ]; then
    if command -v sudo >/dev/null 2>&1; then
        log "$INSTALL_DIR is not writable. Will use sudo."
        USE_SUDO="sudo"
    else
        error "$INSTALL_DIR is not writable and sudo is not available. Please run this script as root or check permissions."
    fi
fi

# Find the binary in the extracted files (handle potential subdirectories)
EXTRACTED_BIN=$(find "$TEMP_DIR" -type f -name "$BINARY_NAME" | head -n 1)

if [ -z "$EXTRACTED_BIN" ]; then
    error "Binary '$BINARY_NAME' not found in the downloaded archive."
fi

log "Installing $BINARY_NAME to $INSTALL_DIR..."
$USE_SUDO mv "$EXTRACTED_BIN" "$INSTALL_DIR/$BINARY_NAME"
$USE_SUDO chmod +x "$INSTALL_DIR/$BINARY_NAME"

success "$BINARY_NAME installed successfully to $INSTALL_DIR/$BINARY_NAME"
success "Run '$BINARY_NAME' to get started!"
