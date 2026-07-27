#!/usr/bin/env bash
set -e

# Ensure we are inside the consumer directory
cd "$(dirname "$0")"

echo "Installing consumer dependencies..."
npm install --force

echo "Granting execution permissions to local node binaries..."
chmod +x ./node_modules/.bin/* || true

echo "Running Pact contract tests..."
npx jest --config pact.config.js
