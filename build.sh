#!/bin/bash
# Build script for AlternateDNS
# Usage: ./build.sh [version] [commit] [date]
# Example: ./build.sh 1.0.0 abc123 2025-01-01

# Get version info from arguments or use defaults
VERSION=${1:-dev}
GIT_COMMIT=${2:-unknown}
BUILD_DATE=${3:-$(date +%Y-%m-%d)}

echo "Building AlternateDNS for Linux/macOS (portable)..."
echo "Version: $VERSION"
echo "Commit: $GIT_COMMIT"
echo "Build Date: $BUILD_DATE"
echo ""
echo "Checking for C compiler (required for Fyne GUI)..."
echo ""

# Check for gcc
if ! command -v gcc &> /dev/null; then
    echo "ERROR: C compiler (gcc) not found!"
    echo ""
    echo "Fyne requires CGO which needs a C compiler."
    echo ""
    echo "To fix this, install GCC:"
    echo "  Linux (Ubuntu/Debian): sudo apt-get install build-essential"
    echo "  Linux (Fedora/RHEL):   sudo dnf install gcc"
    echo "  Linux (Arch):          sudo pacman -S base-devel"
    echo "  macOS:                 xcode-select --install"
    echo ""
    echo "After installing, try again."
    exit 1
fi

echo "C compiler found!"
echo ""
echo "Starting build (this may take a minute)..."
echo ""

# Create dist folder if it doesn't exist
mkdir -p dist

# Enable CGO for Fyne (required on Linux/macOS too)
export CGO_ENABLED=1

# Build with same flags as Windows (minus -H windowsgui which is Windows-only)
go build -ldflags="-s -w -X main.Version=$VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o dist/AlternateDNS

if [ $? -eq 0 ]; then
    # Copy default_config.yaml to dist folder for reference
    if [ -f "default_config.yaml" ]; then
        cp "default_config.yaml" "dist/default_config.yaml"
        echo "Copied default_config.yaml to dist folder"
    fi
    
    echo ""
    echo "Build successful! dist/AlternateDNS created."
    echo "File size:"
    ls -lh dist/AlternateDNS
    echo ""
else
    echo "Build failed!"
    echo ""
    echo "Check the error messages above."
    exit 1
fi
