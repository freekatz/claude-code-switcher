#!/bin/bash
set -e

# Build script for ccs
# Usage: ./scripts/build.sh [version] [output]
# Example: ./scripts/build.sh v1.0.0 ccs-darwin-arm64

VERSION=${1:-dev}
OUTPUT=${2:-ccs}

# Build info
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Package path
PKG="github.com/katz/ccs/cmd"

# Build with ldflags
LDFLAGS="-s -w \
  -X '${PKG}.version=${VERSION}' \
  -X '${PKG}.buildTime=${BUILD_TIME}' \
  -X '${PKG}.gitBranch=${GIT_BRANCH}' \
  -X '${PKG}.gitCommit=${GIT_COMMIT}'"

echo "Building ${OUTPUT}..."
echo "  Version:    ${VERSION}"
echo "  Build time: ${BUILD_TIME}"
echo "  Git branch: ${GIT_BRANCH}"
echo "  Git commit: ${GIT_COMMIT}"
echo "  GOOS:       ${GOOS:-$(go env GOOS)}"
echo "  GOARCH:     ${GOARCH:-$(go env GOARCH)}"

go build -v -ldflags="${LDFLAGS}" -o "${OUTPUT}" .

echo "Done: ${OUTPUT}"
