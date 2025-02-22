#!/bin/bash

echo "🚀 Installing Navi File Manager..."

# Check if Go is installed
if ! command -v go &>/dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19 or higher."
    exit 1
fi

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Clone the repository
echo "📥 Cloning repository..."
git clone https://github.com/Tony-ArtZ/Navi.git
cd Navi

# Build the project
echo "🛠️ Building Navi..."
go build -o navi

# Install the binary
echo "📦 Installing Navi..."
sudo mv navi /usr/local/bin/

# Cleanup
cd ..
rm -rf "$TMP_DIR"

echo "✅ Navi has been installed successfully!"
echo "🎮 Run 'navi' to start the file manager"
