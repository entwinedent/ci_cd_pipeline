#!/usr/bin/env bash
set -e

echo "Installing E2E dependencies..."
npm install --ignore-scripts

echo "Ensuring local binaries have correct execution permissions..."
chmod +x ./node_modules/.bin/*

echo "Installing Playwright browsers..."
npx playwright install --with-deps

echo "Running Playwright tests..."
npx playwright test
