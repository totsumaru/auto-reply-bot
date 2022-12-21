# auto-reply-bot

指定した文字列に反応して、自動でメッセージを返信するbotです。

## 用語について

| 用語    | 説明                | 補足  | 
|-------|-------------------|-----|
| Dev   | botの開発者           |     | 
| Admin | 各サーバーのbotの管理者ユーザー |     | 

## botの導入URL

[botの導入はこちら](https://discord.com/api/oauth2/authorize?client_id=1054635548142223400&permissions=412317207552&scope=bot)

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

| パラメーター | 内容                | 必須   | 例                  |
|--------|-------------------|------|--------------------|
| id     | DiscordのサーバーID    | true | 984614055681613864 |
| code   | Discordログイン後のcode | true | 123                |

##### レスポンス

```json
{
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

| key   | value             | 必須   | 例          |
|-------|-------------------|------|------------|
| token | トークン              | true | abcdxyz... |

##### Body

```json
{
}
```

##### レスポンス

```json
{
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
