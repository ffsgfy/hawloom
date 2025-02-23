#!/bin/bash
#
# Builds and runs a local instance of the app
# Used for development

set -e
cd "$(dirname "$0")/.."

scripts/sqlc.sh
scripts/templ.sh
scripts/tailwind.sh
go build ./cmd/hawloom

POSTGRES_HOST=127.0.0.1 ./hawloom -c ./config/dev.json
