#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

pnpm --dir excalidraw-web build
pnpm --dir vue-web build
