#!/bin/sh

set -eu

echo 'Waiting for MySQL to be available'

maxTries=20
while [ ${maxTries} -gt 0 ] && ! mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "SHOW DATABASES;" >/dev/null 2>&1; do
  sleep 5
  maxTries=$(expr "$maxTries" - 1)
done

if [ "${maxTries}" -le 0 ]; then
  echo >&2 'error: unable to contact MySQL after 20 tries'
  exit 1
fi

echo 'Executing sql ...'

mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "DROP SCHEMA IF EXISTS ${DB_NAME};"
mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "CREATE SCHEMA ${DB_NAME} DEFAULT CHARACTER SET utf8mb4;"
mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "CREATE USER IF NOT EXISTS '${DB_USER_NAME}'@'%' IDENTIFIED BY '${DB_USER_PASSWORD}';"
mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "GRANT ALL PRIVILEGES ON *.* TO '${DB_USER_NAME}'@'%' WITH GRANT OPTION;"
mysql -u root -p${MYSQL_ROOT_PASSWORD} ${DB_NAME} </working/.docker/mysql/bff.sql

# memo:権限確認のコマンド
# select user, host from mysql.user;
