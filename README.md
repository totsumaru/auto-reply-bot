# auto-reply-bot

指定した文字列に反応して、自動でメッセージを返信するbotです。

## 用語について

| 用語    | 説明                | 補足  | 
|-------|-------------------|-----|
| Dev   | botの開発者           |     | 
| Admin | 各サーバーのbotの管理者ユーザー |     | 

## botの導入URL

[botの導入はこちら](https://discord.com/api/oauth2/authorize?client_id=1055348253614419989&permissions=412317207552&scope=bot)

#### 権限

- Read Messages/View Channels
- Send Messages
- Send Messages in Threads
- Embed Links
- Read Message History
- Use External Emojis
- Use External Stickers

## API

### 1. サーバーの情報を取得(Adminのみ)

```
GET /server
```

##### クエリパラメーター

| パラメーター | 内容                | 必須   | 例                              |
|--------|-------------------|------|--------------------------------|
| id     | DiscordのサーバーID    | true | 984614055681613864             |
| code   | Discordログイン後のcode | true | N5GeO3MTBvAyIPMvhjUNItkqrLg2aA |

##### レスポンス

```json
{
  "id": "1055315616002740294",
  "admin_role_id": "1055362036495826964",
  "block": [
    {
      "keyword": [
        "hello",
        "world"
      ],
      "reply": [
        "good",
        "very good"
      ],
      "is_all_match": true,
      "is_random": true,
      "is_embed": false
    }
  ]
}
```

### 2. 設定を更新(Adminのみ)

```
POST /server/config
```

##### クエリパラメーター

| パラメーター | 内容                | 必須   | 例                  |
|--------|-------------------|------|--------------------|
| id     | DiscordのサーバーID    | true | 984614055681613864 |

##### ヘッダー

| key   | value | 必須   | 例                              |
|-------|-------|------|--------------------------------|
| token | トークン  | true | HVBJiJU3JtJxAeXg0mTOavM5R0lty3 |

##### Body

```json
{
  "admin_role_id": "1055362036495826964",
  "block": [
    {
      "keyword": [
        "hello",
        "world"
      ],
      "reply": [
        "good",
        "very good"
      ],
      "is_all_match": true,
      "is_random": true,
      "is_embed": false
    }
  ]
}
```

##### レスポンス

```json
{
  "id": "1055315616002740294",
  "admin_role_id": "1055362036495826964",
  "block": [
    {
      "keyword": [
        "hello",
        "world"
      ],
      "reply": [
        "good",
        "very good"
      ],
      "is_all_match": true,
      "is_random": true,
      "is_embed": false
    }
  ]
}
```

## Discordコマンド

### 1. 新規サーバー作成(Devのみ)

```
/create <サーバーID>
```

- 使用条件
    - Devのみ実行可能
    - `TEST SERVER`でのみ実行可能

### 2. サーバー削除(Devのみ)

```
/delete <サーバーID>
```

- 使用条件
    - Devのみ実行可能
    - `TEST SERVER`でのみ実行可能

### 3. ヘルプ(Adminのみ)

```
/help
```

- 使用条件
    - Adminのみ実行可能
    - 全てのサーバーで実行可能
