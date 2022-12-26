#!/bin/bash

# ================================================
# このファイルはVPSにアップロード後に修正しないでください。
# ================================================
# NOTE:
#  - VPSで最初の1回だけ実行します。
#  - 実行後、このファイルを削除します。
# ================================================

# 環境変数を設定します
set -eu

readonly SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

if [[ ! -f "${SCRIPT_DIR}"/.env ]]; then
  cp "${SCRIPT_DIR}"/.env.dist "${SCRIPT_DIR}"/.env
fi

. "${SCRIPT_DIR}"/.env

# shellcheck disable=SC2046
export $(xargs <"${SCRIPT_DIR}"/.env)

# 環境変数がprd以外は実行できないようにします
if ! [ "${ENV}" = "prd" ]; then
  echo "prd環境以外では実行できません"
  exit 1
fi

# bffのMySQLのコンテナを起動
cd "${SCRIPT_DIR}"/context/bff/.docker/mysql
docker-compose up -d --build
docker-compose exec mysql /working/.docker/mysql/setup.sh

# Goのコンテナを起動
cd "${SCRIPT_DIR}"
docker-compose up -d --build

# `serve.sh`を削除
cd "${SCRIPT_DIR}"
rm serve.sh
echo "serve.sh を削除しました"

# `start.sh`を削除
cd "${SCRIPT_DIR}"
rm start.sh
echo "start.sh を削除しました"

# `test-runner.sh`を削除
cd "${SCRIPT_DIR}"
rm test-runner.sh
echo "test-runner.sh を削除しました"

# このファイルを削除
cd "${SCRIPT_DIR}"
rm setup-vps.sh
echo "setup-vps.sh を削除しました"
