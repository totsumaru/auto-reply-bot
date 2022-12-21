#!/bin/bash

set -eu

readonly SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

if [[ ! -f "${SCRIPT_DIR}"/.env ]]; then
  cp "${SCRIPT_DIR}"/.env.dist "${SCRIPT_DIR}"/.env
fi

. "${SCRIPT_DIR}"/.env

# shellcheck disable=SC2046
export $(xargs <"${SCRIPT_DIR}"/.env)

# Goをローカルで起動
go run "${SCRIPT_DIR}"/main.go
