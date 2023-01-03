# auto-reply-bot

指定した文字列に反応して、自動でメッセージを返信するbotです。

## 用語について

| 用語    | 説明                   | 補足  | 
|-------|----------------------|-----|
| Dev   | botの開発者              |     | 
| Owner | 各サーバーの所有者(主にサーバー作成者) |     | 
| Admin | 各サーバーのbotの管理者ユーザー    |     | 

## botの導入URL

[本番用のbotの導入はこちら](https://discord.com/api/oauth2/authorize?client_id=1056843645967413309&permissions=412317207552&scope=bot)

[テスト用のbotの導入はこちら](https://discord.com/api/oauth2/authorize?client_id=1055348253614419989&permissions=412317207552&scope=bot)

## 新規サーバーへの導入手順

1. TwitterのDMに依頼をもらう
2. 無料期間を案内した後、[Google Form](https://forms.gle/6pmaX1bX7bdzvvGi9)を送る → 回答後に導入URLあり
3. ※ユーザー対応: URLからbotの導入
4. `/create-server`コマンドでDBにレコードを作成
5. [botのコンパネ](https://discord.com/developers/applications/1056843645967413309/oauth2/general)からリダイレクトURLを設定
6. [Notion](https://www.notion.so/bot-92d5460b67d3407893343008d1821a49)に情報を追加

## 有料移行時の手順

1. 終了5日前に有料に移行するか聞く（しない場合はbotを削除してもらって終了）
2. stripeのURL（）を送付

## インフラ構成

- FE: Cloudflare Pages
    - ドメインはCloudflare Pagesのデフォルトを使用（[auto-reply-bot.pages.dev](https://auto-reply-bot.pages.dev)）
- BE: さくらのVPS
    - ドメインは独自ドメイン（さくらのドメイン → Cloudflare DNS → バックエンドIPアドレス）
- ドメイン
    - さくらのドメイン（[auto-reply-bot](https://auto-reply-bot)）

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

##### アクセスできるユーザー

- Dev
- サーバーオーナー
- 管理者ロールを持つユーザー

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
      "match_condition": "all-match",
      "is_random": true,
      "is_embed": false
    }
  ],
  "token": "abcd0123",
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
    }
  ]
}
```

### 2. 設定を更新(Adminのみ)

```
POST /server/config
```

##### アクセスできるユーザー

- Dev
- サーバーオーナー
- 管理者ロールを持つユーザー

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
      "match_condition": "all-match",
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
      "match_condition": "all-match",
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

##### アクセスできるユーザー

- Devのみ

##### アクセスできるサーバー

- `TEST SERVER`でのみ実行可能

### 2. サーバー削除(Devのみ)

```
/delete <サーバーID>
```

##### アクセスできるユーザー

- Devのみ

##### アクセスできるサーバー

- `TEST SERVER`でのみ実行可能

### 3. ヘルプ(Adminのみ)

```
/help
```

##### アクセスできるユーザー

- Dev
- サーバーオーナー
- 管理者ロールを持つユーザー

##### アクセスできるサーバー

- 全てのサーバーで実行可能
