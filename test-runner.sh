#!/bin/bash

set -eu

readonly SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

if [[ ! -f "${SCRIPT_DIR}"/.env ]]; then
  cp "${SCRIPT_DIR}"/.env.dist "${SCRIPT_DIR}"/.env
fi

. "${SCRIPT_DIR}"/.env

# shellcheck disable=SC2046
export $(xargs <"${SCRIPT_DIR}"/.env)

# =============================================================
# userコンテキストのテストを実行します
# =============================================================
cd "${SCRIPT_DIR}"/context/user/.docker/mysql
docker-compose up -d --build
docker-compose exec mysql /working/.docker/mysql/setup.sh
cd "${SCRIPT_DIR}"
go test -count 1 -p 1 "${SCRIPT_DIR}"/context/user/...
# =============================================================

# 他のコンテキストがある場合は、この下に追加します
