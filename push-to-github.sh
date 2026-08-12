#!/bin/bash
# Run this from inside the extracted "repo" folder
set -e

echo "🔧 Initializing WarnetechVPN repository..."

git init
git branch -M main

git remote remove origin 2>/dev/null || true
git remote add origin https://github.com/tewartech-node/warnetech-vpn.git

git add .
git commit -m "Initial commit: WarnetechVPN scaffold (Go backend + Android app)"

echo ""
echo "✅ Repo initialized and committed locally."
echo ""
echo "Now push it up:"
echo "  git push -u origin main"
echo ""
echo "If the GitHub repo already has commits (e.g. a README you created on GitHub),"
echo "pull first to avoid conflicts:"
echo "  git pull origin main --allow-unrelated-histories"
echo "  git push -u origin main"
