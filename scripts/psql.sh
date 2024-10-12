#!/bin/bash
#
# Drops into a psql shell connected to the main database

docker compose exec -it db psql --user ${POSTGRES_USER}
