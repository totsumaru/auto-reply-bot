# auto-reply-bot

指定した文字列に反応して、自動でメッセージを返信するbotです。

## 用語について

| 用語    | 説明                | 補足  | 
|-------|-------------------|-----|
| Dev   | botの開発者           |     | 
| Admin | 各サーバーのbotの管理者ユーザー |     | 

## botの導入URL

[botの導入はこちら](https://discord.com/api/oauth2/authorize?client_id=1055348253614419989&permissions=412317207552&scope=bot)

TODO: 一時的なURLのため削除  
[管理者のログインはこちら](https://discord.com/api/oauth2/authorize?client_id=1055348253614419989&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fserver%3Fid%3D1055315616002740294&response_type=code&scope=identify)

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
      "name": "あいさつ",
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
  ],
  "server_name": "TEST 2",
  "avatar_url": "https://cdn.discordapp.com/icons/1055315616002740294/c17fe110e848098db614687645f17586.png",
  "role": [
    {
      "id": "1055315616002740294",
      "name": "@everyone"
    },
    {
      "id": "1055350145975328863",
      "name": "[test]自動返信bot"
    },
    {
      "id": "1055362036495826964",
      "name": "自動返信botの管理者"
    },
    {
      "id": "1056464506554957824",
      "name": "テストロールです"
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
      "name": "あいさつ",
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
      "name": "あいさつ",
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
  ],
  "server_name": "TEST 2",
  "avatar_url": "https://cdn.discordapp.com/icons/1055315616002740294/c17fe110e848098db614687645f17586.png",
  "role": [
    {
      "id": "1055315616002740294",
      "name": "@everyone"
    },
    {
      "id": "1055350145975328863",
      "name": "[test]自動返信bot"
    },
    {
      "id": "1055362036495826964",
      "name": "自動返信botの管理者"
    },
    {
      "id": "1056464506554957824",
      "name": "テストロールです"
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
