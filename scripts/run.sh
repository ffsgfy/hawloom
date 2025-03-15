#!/bin/bash
#
# Builds and runs a local instance of the app
# Used for development; pressing enter rebuilds and restarts the server

set -e
cd "$(dirname "$0")/.."

while true; do
    scripts/sqlc.sh
    scripts/templ.sh
    scripts/tailwind.sh
    go build ./cmd/hawloom

    POSTGRES_HOST=127.0.0.1 ./hawloom -c ./config/dev.json &
    pid=$!

    restart=0
    if read _; then
        restart=1
    fi

    kill $pid
    wait $pid || true

    if [ $restart -eq 0 ]; then
        exit 0
    fi
done
