#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
BASE_URL="${1:-http://localhost:3000}"
DURATION="${2:-10s}"
THREADS="${3:-4}"
CONNECTIONS="${4:-50}"

echo "=========================================================="
echo "URL Shortener PostgreSQL Benchmark Suite"
echo "Target: $BASE_URL | Duration: $DURATION | Threads: $THREADS | Conns: $CONNECTIONS"
echo "=========================================================="

echo ""
echo "[1/4] Generating sample short URL for GET benchmark..."
SHORT_RESP=$(curl -s -X POST "$BASE_URL/" -H "Content-Type: application/json" -d '{"url":"https://example.com/postgres-bench-target"}')
SHORT_CODE=$(echo "$SHORT_RESP" | grep -oP 'http://localhost:3000/\K[a-zA-Z0-9]+' || echo "9XkbVH9nh9l")
echo "Generated Short Code: $SHORT_CODE"

echo ""
echo "----------------------------------------------------------"
echo "[2/4] Benchmarking GET /:shortUrl (URL Redirection)"
echo "----------------------------------------------------------"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency "$BASE_URL/$SHORT_CODE"

echo ""
echo "----------------------------------------------------------"
echo "[3/4] Benchmarking POST / (Dynamic Unique URLs)"
echo "----------------------------------------------------------"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" -s "$DIR/post_shorten_dynamic.lua" --latency "$BASE_URL/"

echo ""
echo "----------------------------------------------------------"
echo "[4/4] Benchmarking POST / (Static URL with Upsert)"
echo "----------------------------------------------------------"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" -s "$DIR/post_shorten.lua" --latency "$BASE_URL/"

echo ""
echo "=========================================================="
echo "PostgreSQL Benchmarks Completed!"
echo "=========================================================="
