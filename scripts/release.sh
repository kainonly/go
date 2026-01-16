#!/bin/bash
# Release script for creating new versions
# Usage: ./scripts/release.sh <version>
# Example: ./scripts/release.sh v1.0.0

set -e

VERSION=${1}
if [ -z "$VERSION" ]; then
    echo "Error: Version number is required"
    echo "Usage: ./scripts/release.sh <version>"
    echo "Example: ./scripts/release.sh v1.0.0"
    exit 1
fi

echo "🚀 Starting release process for $VERSION"

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "⚠️  Warning: You are not on main branch (current: $CURRENT_BRANCH)"
    read -p "Do you want to continue? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "❌ Error: Working directory is not clean"
    echo "Please commit or stash your changes first"
    exit 1
fi

# Pull latest changes
echo "📥 Pulling latest changes..."
git pull origin main

# Run tests
echo "🧪 Running tests..."
go test ./... -v

if [ $? -ne 0 ]; then
    echo "❌ Tests failed! Please fix the issues before releasing"
    exit 1
fi

# Create git tag
echo "🏷️  Creating git tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION

See CHANGELOG.md for details."

# Push tag
echo "📤 Pushing tag to remote..."
git push origin "$VERSION"

# Create GitHub release using gh CLI
if command -v gh &> /dev/null; then
    echo "📦 Creating GitHub release..."

    # Read release notes if exists
    RELEASE_NOTES_FILE="RELEASE_NOTES_${VERSION}.md"
    if [ -f "$RELEASE_NOTES_FILE" ]; then
        NOTES=$(cat "$RELEASE_NOTES_FILE")
    else
        NOTES="Release $VERSION

See [CHANGELOG.md](https://github.com/kainonly/go/blob/main/CHANGELOG.md) for details."
    fi

    gh release create "$VERSION" \
        --title "$VERSION" \
        --notes "$NOTES"

    echo "✅ GitHub release created successfully!"
else
    echo "⚠️  gh CLI not found. Please create the release manually at:"
    echo "   https://github.com/kainonly/go/releases/new?tag=$VERSION"
fi

echo ""
echo "✨ Release $VERSION completed successfully!"
echo ""
echo "📝 Next steps:"
echo "1. Verify the release at https://github.com/kainonly/go/releases"
echo "2. Update documentation if needed"
echo "3. Announce the release to users"
