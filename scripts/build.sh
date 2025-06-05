#!/bin/sh
# Скрипт cross-compile для сборки бинарников под Linux/amd64, Linux/arm64, macOS/amd64 и macOS/arm64

set -eu

# Название вашего CLI-приложения (без расширения)
APP_NAME="sitedog"

# Директория для выходных сборок
OUTPUT_DIR="dist"
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Массив целевых платформ GOOS/GOARCH
PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"

for PLATFORM in $PLATFORMS; do
  GOOS=$(echo "$PLATFORM" | cut -d/ -f1)
  GOARCH=$(echo "$PLATFORM" | cut -d/ -f2)
  BINARY_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
  # Для Windows добавлять .exe (необязательно, в данном случае не используется)
  # if [ "$GOOS" = "windows" ]; then
  #   BINARY_NAME+=".exe"
  # fi

  echo "Собираем для $GOOS/$GOARCH → $OUTPUT_DIR/$BINARY_NAME"
  env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTPUT_DIR/$BINARY_NAME" .
done

echo "Готово! Бинарники лежат в папке $OUTPUT_DIR:"
ls -1 "$OUTPUT_DIR"
