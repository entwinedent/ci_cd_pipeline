#!/usr/bin/env bash
set -e

# Ensure we are inside the consumer directory
cd "$(dirname "$0")"

echo "Installing consumer dependencies with scripts enabled..."
npm install --allow-scripts

echo "Running Pact contract tests..."
npx jest --config pact.config.js
