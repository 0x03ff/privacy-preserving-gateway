#!/bin/bash
set -e

# Automatically move into the directory where this script lives
cd "$(dirname "$0")"

# Define output path relative to scripts/fame/
OUTPUT_DIR="../../static/wasm"

echo "Compiling FAME to WebAssembly..."
GOOS=js GOARCH=wasm go build -o "$OUTPUT_DIR/fame.wasm" main.go

echo "Copying wasm_exec.js from Go installation..."
# Updated path for Go 1.24+ / Go 1.27
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$OUTPUT_DIR/"

echo "WASM build successfully generated in $OUTPUT_DIR!"