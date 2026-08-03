#!/bin/bash
#
# build.sh - Build the pc binary with version information
#
# Usage:
#   ./build.sh              # Build with latest tag as version
#   ./build.sh v1.2.3       # Build with specific version
#   ./build.sh --dev        # Build with dev version (includes commit hash)
#

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Determine version
if [ "$1" == "--dev" ]; then
    # Dev build: use tag + commit hash + dirty flag
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
elif [ -n "$1" ]; then
    # Specific version provided
    VERSION="$1"
else
    # Use latest tag
    VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
fi

# Build flags
LDFLAGS="-s -w -X main.Version=${VERSION}"

echo -e "${BLUE}Building pc version: ${VERSION}${NC}"

# Build
go build -ldflags="${LDFLAGS}" -o pc .

echo -e "${GREEN}✓ Built: ./pc (${VERSION})${NC}"

# Show version
./pc -version
