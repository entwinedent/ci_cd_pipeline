#!/usr/bin/env bash
set -e

echo "Navigating to E2E directory and installing correct dependencies..."
cd "$(dirname "$0")"

# Ensure core packages are explicitly installed if missing
npm install @playwright/test playwright --save-dev

echo "Ensuring local binaries have correct execution permissions..."
chmod -R +x ./node_modules/.bin/ || true

echo "Installing Playwright browsers and dependencies..."
npx playwright install --with-deps

echo "Running Playwright tests..."
npx playwright test
