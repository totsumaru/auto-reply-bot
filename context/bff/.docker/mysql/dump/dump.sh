#!/bin/bash

# ============================================
# このスクリプトは、ホスト側で実行します
# ============================================

# このスクリプトと同じ階層に移動します
readonly SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR"

# 変数を設定します
CONTAINER_NAME='auto-reply-bot-mysql'
OUTPUT_DIR="logs"
OUTPUT_FILE="$OUTPUT_DIR/$(date "+%Y_%m_%d_%H_%M_%S").sql" # 例: 2022_10_18_00_47_28.sql

# logsディレクトリが存在しない場合は作成します
if [ ! -d "$OUTPUT_DIR" ]; then
  mkdir "$OUTPUT_DIR"
fi

if [ -z "$CONTAINER_NAME" ]; then
  echo 'コンテナ名が設定されていません.'
  exit 1
fi

if [ -z "$OUTPUT_FILE" ]; then
  echo 'アウトプットファイル名が設定されていません'
  exit 1
fi

# mysqlコンテナの中に入って、マウントされたディレクトリにdumpファイルを出力します
docker exec ${CONTAINER_NAME} sh -c \
  'exec mysqldump --all-databases --lock-all-tables -uroot -p"$MYSQL_ROOT_PASSWORD"' >"$OUTPUT_FILE"

# 30件を超えたファイルを削除します
cd $OUTPUT_DIR
rm -f "$(ls -t . | tail -n+31)"

# Discordに通知します
#curl -X POST -H "Content-Type: application/json" -d '{"content":"DBのダンプが完了しました"}' \
#  https://discord.com/api/webhooks/xxx
