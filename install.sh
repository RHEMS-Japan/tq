#!/bin/bash

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}Installing tq - TOON Query Processor${NC}\n"

# Check prerequisites
check_command() {
    if ! command -v "$1" &> /dev/null; then
        echo -e "${RED}Error: $1 is not installed${NC}"
        echo "$2"
        exit 1
    fi
}

echo "Checking prerequisites..."
check_command "node" "Install Node.js from https://nodejs.org/"
check_command "jq" "Install jq: brew install jq (macOS) or apt-get install jq (Ubuntu)"
check_command "go" "Install Go from https://golang.org/dl/"
echo -e "${GREEN}✓ All prerequisites found${NC}\n"

# Install dependencies
echo "Installing Node.js dependencies..."
npm install --silent
echo -e "${GREEN}✓ Dependencies installed${NC}\n"

# Build
echo "Building tq..."
go build -o tq ./cmd/tq
echo -e "${GREEN}✓ Build complete${NC}\n"

# Install location
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
TQ_DIR="$HOME/.tq"

mkdir -p "$INSTALL_DIR"
mkdir -p "$TQ_DIR/scripts"

# Copy files
echo "Installing files..."
cp tq "$INSTALL_DIR/"
cp scripts/*.js "$TQ_DIR/scripts/"
cp -r node_modules "$TQ_DIR/"

echo -e "${GREEN}✓ Installation complete!${NC}\n"
echo "  Binary: $INSTALL_DIR/tq"
echo "  Scripts: $TQ_DIR/scripts"
echo ""

# Check PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}⚠ Add $INSTALL_DIR to your PATH:${NC}"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
fi

echo "Test installation:"
echo "  tq --version"
