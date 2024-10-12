#!/bin/bash
#
# This script can be used to pass custom arguments to migrate

docker compose run -it --rm db-migrate -path=/migrations/ -database=${POSTGRES_URI} $@
