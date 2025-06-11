#!/bin/sh
# Cross-compile script for building binaries for Linux/amd64, Linux/arm64, macOS/amd64 and macOS/arm64

set -eu

# Your CLI application name (without extension)
APP_NAME="sitedog"

# Output directory for builds
OUTPUT_DIR="dist"
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Array of target GOOS/GOARCH platforms
PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"

for PLATFORM in $PLATFORMS; do
  GOOS=$(echo "$PLATFORM" | cut -d/ -f1)
  GOARCH=$(echo "$PLATFORM" | cut -d/ -f2)
  BINARY_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
  # Add .exe for Windows (optional, not used in this case)
  # if [ "$GOOS" = "windows" ]; then
  #   BINARY_NAME+=".exe"
  # fi

  echo "Building for $GOOS/$GOARCH → $OUTPUT_DIR/$BINARY_NAME"
  env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTPUT_DIR/$BINARY_NAME" .
done

echo "Done! Binaries are in the $OUTPUT_DIR folder:"
ls -1 "$OUTPUT_DIR"
