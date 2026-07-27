#!/usr/bin/env bash
set -e

echo "Installing local consumer dependencies with scripts enabled..."
npm install --foreground-scripts

echo "Granting execution permissions to local binaries..."
chmod -R +x ./node_modules/.bin/ || true

echo "Running Pact contract tests..."
npx jest --config pact.config.js
