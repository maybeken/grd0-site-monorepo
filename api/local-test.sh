#!/bin/bash
# Run the API development loop.
#
#   ./local-test.sh              run Air in Docker, or with the local Go
#                                toolchain when Docker is not available
#   ./local-test.sh --no-docker  run Air directly with the local Go toolchain
#
# Both modes use .air.toml and .env.dev (APP_ENV=LOCAL, PORT=:8080), rebuilding
# and restarting the API when Go sources change. Extra arguments are forwarded
# to Docker Compose or Air.
set -euo pipefail

cd "$(dirname "$0")"

run_local=false
if [ "${1:-}" = "--no-docker" ]; then
	run_local=true
	shift
elif ! command -v docker >/dev/null 2>&1; then
	echo "docker not found; running with the local Go toolchain" >&2
	run_local=true
fi

if [ "$run_local" = true ]; then
	if ! command -v go >/dev/null 2>&1; then
		echo "go not found in PATH" >&2
		exit 1
	fi

	exec go tool air "$@"
fi

exec docker compose -f compose.dev.yaml up --build "$@"
