#!/bin/bash
#
# release.sh - Create and push a new release tag to GitHub
#
# Usage:
#   ./release.sh <version>      # Create and push tag (e.g., ./release.sh v1.2.3)
#   ./release.sh --list         # List existing tags
#   ./release.sh --latest       # Show latest tag
#   ./release.sh --help         # Show help
#
# The version should follow semantic versioning: vMAJOR.MINOR.PATCH
# Pre-release versions can use: v1.2.3-beta.1, v1.2.3-rc.1, etc.
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_usage() {
    echo -e "${BLUE}Usage:${NC}"
    echo "  $0 <version>      Create and push a new release tag"
    echo "  $0 --list         List existing tags"
    echo "  $0 --latest       Show the latest tag"
    echo "  $0 --help         Show this help message"
    echo ""
    echo -e "${BLUE}Examples:${NC}"
    echo "  $0 v1.0.0         Create release v1.0.0"
    echo "  $0 v1.2.3-beta.1  Create pre-release v1.2.3-beta.1"
    echo ""
    echo -e "${BLUE}Version format:${NC}"
    echo "  vMAJOR.MINOR.PATCH[-PRERELEASE]"
    echo ""
}

list_tags() {
    echo -e "${BLUE}Existing release tags:${NC}"
    git tag -l 'v*' --sort=-v:refname | head -20
    echo ""
    total=$(git tag -l 'v*' | wc -l | tr -d ' ')
    echo -e "${YELLOW}Total: $total tags${NC}"
}

show_latest() {
    latest=$(git tag -l 'v*' --sort=-v:refname | head -1)
    if [ -z "$latest" ]; then
        echo -e "${YELLOW}No release tags found.${NC}"
        echo "Start with: $0 v0.1.0"
    else
        echo -e "${GREEN}Latest release: $latest${NC}"
        # Show commit info for the tag
        echo ""
        git log -1 --format="  Commit: %h%n  Date: %ci%n  Author: %an%n  Message: %s" "$latest" 2>/dev/null || true
    fi
}

validate_version() {
    local version=$1
    # Check if version matches semantic versioning pattern
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
        echo -e "${RED}Error: Invalid version format '$version'${NC}"
        echo "Version must match: vMAJOR.MINOR.PATCH[-PRERELEASE]"
        echo "Examples: v1.0.0, v2.1.3, v1.0.0-beta.1, v1.0.0-rc.1"
        exit 1
    fi
}

check_clean_working_tree() {
    if ! git diff-index --quiet HEAD -- 2>/dev/null; then
        echo -e "${YELLOW}Warning: You have uncommitted changes.${NC}"
        echo ""
        git status --short
        echo ""
        read -p "Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Aborted."
            exit 1
        fi
    fi
}

check_tag_exists() {
    local version=$1
    if git rev-parse "$version" >/dev/null 2>&1; then
        echo -e "${RED}Error: Tag '$version' already exists.${NC}"
        echo ""
        echo "Existing tags:"
        git tag -l 'v*' --sort=-v:refname | head -5
        exit 1
    fi
}

create_and_push_tag() {
    local version=$1
    
    echo -e "${BLUE}Creating release tag: $version${NC}"
    echo ""
    
    # Show what will be tagged
    echo "This will tag the current commit:"
    git log -1 --format="  Commit: %H%n  Date: %ci%n  Author: %an%n  Message: %s"
    echo ""
    
    # Confirm
    read -p "Create and push tag '$version'? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
    
    # Create annotated tag
    echo ""
    echo -e "${BLUE}Creating annotated tag...${NC}"
    git tag -a "$version" -m "Release $version"
    
    # Push tag to origin
    echo -e "${BLUE}Pushing tag to origin...${NC}"
    git push origin "$version"
    
    echo ""
    echo -e "${GREEN}✓ Successfully created and pushed tag: $version${NC}"
    echo ""
    echo "GitHub Actions will now build and create the release."
    echo "Check the progress at: https://github.com/$(git remote get-url origin | sed -E 's|.*github.com[:/](.*)\.git|\1|')/actions"
}

# Main script
case "${1:-}" in
    --help|-h)
        print_usage
        exit 0
        ;;
    --list|-l)
        list_tags
        exit 0
        ;;
    --latest)
        show_latest
        exit 0
        ;;
    "")
        echo -e "${RED}Error: No version specified.${NC}"
        echo ""
        print_usage
        exit 1
        ;;
    *)
        VERSION="$1"
        
        # Validate version format
        validate_version "$VERSION"
        
        # Check for uncommitted changes
        check_clean_working_tree
        
        # Check if tag already exists
        check_tag_exists "$VERSION"
        
        # Create and push the tag
        create_and_push_tag "$VERSION"
        ;;
esac
