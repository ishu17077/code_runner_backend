# #!/usr/bin/env bash
# set -euo pipefail

# WORKERS=${WORKERS:-8}
# URL="http://127.0.0.1:30080/submission/test/private"
# DATA_FILE="/tmp/stress_payload.json"
# UA_PREFIX="stress-worker"
# LOGDIR="/tmp/stress.json"

# cat > "$DATA_FILE" <<'JSON'
# {
#   "problem_id": "69",
#     "language": "C",
#     "code": "I2luY2x1ZGUgPHN0ZGlvLmg+CmludCBtYWluKCkgewogICAgaW50IHJlczsKICAgIHNjYW5mKCIlZCIsICZyZXMpOwogICAgaWYgKHJlcyAlIDIgPT0gMCkgewogICAgICAgIHByaW50ZigiWWVzXG5ZZXMiKTsKICAgIH0gZWxzZSB7CiAgICAgICAgcHJpbnRmKCJOb1xuTm8iKTsKICAgIH0KfQ==",
#   "tests":[
#     {"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},
#     {"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"},{"problem_id":"69","is_public":true,"stdin":"12\n","expected_output":"Yes","test_id":"1"},
#     {"problem_id":"69","is_public":true,"stdin":"11\n","expected_output":"No","test_id":"2"} 
#   ]
# }
# JSON

# REQS=0
# START_TIME=$SECONDS

# echo "Starting $WORKERS workers (press Ctrl-C to stop)"

# # Ensure child curl loops are killed on INT/TERM by matching the data file on their cmdline
# cleanup() {
#   echo "Stopping workers..."
#   echo "Time taken: $((SECONDS-START_TIME)) seconds"
#   echo "Requests Sent: $((REQS))"
#   pkill -f -- "$DATA_FILE" || true
#   sleep 1
# }
# trap cleanup SIGINT SIGTERM




# seq 1 "$WORKERS" | xargs -n1 -P "$WORKERS" -I{} bash -c '
#   while true;
#    do
#      curl -w "\n" -X POST "'"$URL"'" \
#       -H "Content-Type: application/json" \
#       -H "User-Agent: '"$UA_PREFIX"'-{}" \
#       -H "Connection: close" \
#       --data @'"$DATA_FILE"' # >> /tmp/stress.json;
   
#     # sleep 0.01
#   done
# ' {}

hey -n 500 -c 50 -m POST -H "Content-Type: application/json" -H "Connection: close" -D ./examples/stress_payload.json http://127.0.0.1:30080/submission/test/private