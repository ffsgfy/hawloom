#!/bin/bash
#
# Generates Tailwind CSS styles using the standalone CLI tool
#
# The required binary may be downloaded from the official repository:
# https://github.com/tailwindlabs/tailwindcss/releases

set -e
cd "$(dirname "$0")/.."

tailwindcss -i ./internal/ui/css/style.css -o ./static/style.css
