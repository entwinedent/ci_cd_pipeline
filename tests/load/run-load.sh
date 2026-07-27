#!/usr/bin/env bash
set -e

echo "Starting K6 Load Tests..."
mkdir -p tests/load/results

# Run k6 and export summary json for artifacts
k6 run --summary-export=tests/load/results/k6-summary.json tests/load/load-test.js
