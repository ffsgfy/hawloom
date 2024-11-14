#!/bin/bash
#
# Runs sqlc with the given arguments

docker compose run -it --rm sqlc $@
