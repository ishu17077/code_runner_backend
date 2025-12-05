#!/usr/bin/env bash
set -euo pipefail

WORKERS=${WORKERS:-6}
URL="http://127.0.0.1:30080/submission/test/private"
DATA_FILE="/tmp/stress_payload.json"
UA_PREFIX="stress-worker"
LOGDIR="/tmp/stress.json"

cat > "$DATA_FILE" <<'JSON'
{
  "problem_id": "69",
    "language": "Python",
    "code": "bWVtb3J5X2NvbnN1bWVyID0gW10Kc3RyaW5nX3NpemVfYnl0ZXMgPSA1MTIKbnVtX3N0cmluZ3MgPSBpbnQoKDEwMDAwMDAwMDAwMDAwMDAwMDAwMDAgKiAxMDI0ICogMTAyNCkgLyBzdHJpbmdfc2l6ZV9ieXRlcykKbWVtb3J5X2NvbnN1bWVyID0gW10KZm9yIGkgaW4gcmFuZ2UobnVtX3N0cmluZ3MpOgogICAgbWVtb3J5X2NvbnN1bWVyLmFwcGVuZCgiQSIgKiBzdHJpbmdfc2l6ZV9ieXRlcyk=",
  "tests":[
    {"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},
    {"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"}
  ]
}
JSON

echo "Starting $WORKERS workers (press Ctrl-C to stop)"

# Ensure child curl loops are killed on INT/TERM by matching the data file on their cmdline
cleanup() {
  echo "Stopping workers..."
  pkill -f -- "$DATA_FILE" || true
  sleep 1
}
trap cleanup SIGINT SIGTERM

# Use xargs -P to run N parallel workers without using '&'
seq 1 "$WORKERS" | xargs -n1 -P "$WORKERS" -I{} bash -c '
  while true; do
    curl -s -X POST "'"$URL"'" \
      -H "Content-Type: application/json" \
      -H "User-Agent: '"$UA_PREFIX"'-{}" \
      -H "Connection: close" \
      --data @'"$DATA_FILE"' >> /tmp/stress.json
    # small sleep to avoid tight busy-loop: conditional
    # sleep 0.01
  done
' {}