#!/usr/bin/env bash
set -e

echo "Installing local consumer dependencies..."
npm install

echo "Granting execution permissions to local binaries..."
chmod +x ./node_modules/.bin/jest

echo "Running Pact contract tests..."
./node_modules/.bin/jest
