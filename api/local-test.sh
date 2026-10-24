#!/bin/bash
# Run the API locally.
#
#   ./local-test.sh              build and run in Docker, or with the local Go
#                                toolchain when Docker is not available
#   ./local-test.sh --no-docker  run directly with the local Go toolchain
#
# Both modes read ./api.env when it exists. Port 8080 is used outside Docker
# so the server can run as a non-root user; set PORT in api.env to change it.
set -euo pipefail

cd "$(dirname "$0")"

run_local=false
if [ "${1:-}" = "--no-docker" ]; then
	run_local=true
elif ! command -v docker >/dev/null 2>&1; then
	echo "docker not found; running with the local Go toolchain" >&2
	run_local=true
fi

if [ "$run_local" = true ]; then
	if ! command -v go >/dev/null 2>&1; then
		echo "go not found in PATH" >&2
		exit 1
	fi

	if [ -f ./api.env ]; then
		set -a
		. ./api.env
		set +a
	fi

	export APP_ENV=LOCAL
	export PORT="${PORT:-:80}"

	exec go run .
fi

docker build -t grd0/api:latest . && \
docker run --rm -p 80:80 \
	-e APP_ENV=LOCAL \
	-v ./data/api.db:/opt/app/data/api.db \
	-v ./data/files/files.json:/opt/app/data/files/files.json \
	--env-file ./api.env \
	grd0/api
