#!/usr/bin/env bash
set -e

echo "Starting K6 Load Tests..."
mkdir -p tests/load/results

# Run k6 and export summary json for artifacts
k6 run --out json=tests/load/results/k6-results.json tests/load/load-test.js
