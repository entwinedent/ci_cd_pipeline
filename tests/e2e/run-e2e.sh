#!/usr/bin/env bash
set -e

echo "Navigating to E2E directory..."
cd "$(dirname "$0")"

echo "Installing Playwright dependencies..."
rm -rf node_modules package-lock.json
npm install

echo "Ensuring local binaries have correct execution permissions..."
chmod +x ./node_modules/.bin/* || true

echo "Installing Playwright browsers..."
npx playwright install --with-deps

echo "Running Playwright tests..."
npx playwright test
