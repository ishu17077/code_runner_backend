#!/usr/bin/env bash
set -euo pipefail

WORKERS=${WORKERS:-6}
URL="http://127.0.0.1:31897/submission/test"
DATA_FILE="/tmp/stress_payload.json"
UA_PREFIX="stress-worker"
LOGDIR="/tmp/stress-logs"
mkdir -p "$LOGDIR"

cat > "$DATA_FILE" <<'JSON'
{
  "problem_id":"69",
  "language":"Java",
  "code":"aW1wb3J0IGphdmEudXRpbC5TY2FubmVyOwoKcHVibGljIGNsYXNzIFNvbHV0aW9uIHsKICAgIHB1YmxpYyBzdGF0aWMgdm9pZCBtYWluKFN0cmluZ1tdIGFyZ3MpIHsKICAgICAgICBTY2FubmVyIHNjID0gbmV3IFNjYW5uZXIoU3lzdGVtLmluKTsKICAgICAgICBpbnQgeCA9IHNjLm5leHRJbnQoKTsKICAgICAgICBpZiAoeCAlIDIgPT0gMCkgewogICAgICAgICAgICBTeXN0ZW0ub3V0LnByaW50bG4oIlllcyIpOwogICAgICAgIH0gZWxzZSB7CiAgICAgICAgICAgIFN5c3RlbS5vdXQucHJpbnRsbigiTm8iKTsKICAgICAgICB9CiAgICAgICAgc2MuY2xvc2UoKTsKICAgIH0KfQo=",
  "tests":[
    {"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes\nYes","test_id":"1"},
    {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No\nNo","test_id":"2"}
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
	  --data @'"$DATA_FILE"' >/dev/null
	  
    # small sleep to avoid tight busy-loop if you need it
    # sleep 0.01
  done
' {}
