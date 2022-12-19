#!/bin/bash

git pull origin main

# Goのコンテナのみ再起動します。
# 環境変数はgodotenvを使用してGoの内部で取得しているため、
# 再設定する必要はありません。
docker-compose up -d --build