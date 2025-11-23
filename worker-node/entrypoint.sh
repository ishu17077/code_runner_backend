#!/bin/sh
set -e

echo "Restricting internet access for executor"

iptables -A OUTPUT -m owner --uid-owner 6969 --gid-owner 7070 -j DROP

echo "Launching worker node..."
# ? Replaces this process then with docker exect that will be executed later $@
exec "$@"