# auto-reply-bot

指定した文字列に反応して、自動でメッセージを返信するbotです。

## ドキュメントURL(Notion)

[https://daffy-hamburger-7f6.notion.site/bot-d7292773a40e405eab3fc8310a83ddd7](https://daffy-hamburger-7f6.notion.site/bot-d7292773a40e405eab3fc8310a83ddd7)

## 用語について

| 用語    | 説明                   | 補足  | 
|-------|----------------------|-----|
| Dev   | botの開発者              |     | 
| Owner | 各サーバーの所有者(主にサーバー作成者) |     | 
| Admin | 各サーバーのbotの管理者ユーザー    |     | 

## botの導入URL

[本番用のbotの導入はこちら](https://discord.com/api/oauth2/authorize?client_id=1056843645967413309&permissions=8&scope=bot)

[テスト用のbotの導入はこちら](https://discord.com/api/oauth2/authorize?client_id=1055348253614419989&permissions=8&scope=bot)

## 新規サーバーへの導入手順

1. TwitterのDMに依頼をもらう
2. [Google Form](https://forms.gle/6pmaX1bX7bdzvvGi9)を送る → 回答後に導入URLあり
3. ※ユーザー対応: URLからbotの導入
4. `/create-server`コマンドでDBにレコードを作成
5. ユーザーのコンパネから管理者ロールを設定

## 有料移行時の手順

未定

## インフラ構成

- FE: Cloudflare Pages
    - ドメインはCloudflare Pagesのデフォルトを使用（[auto-reply-bot.pages.dev](https://auto-reply-bot.pages.dev)）
- BE: さくらのVPS
    - ドメインは独自ドメイン（さくらのドメイン → Cloudflare DNS → バックエンドIPアドレス）
- ドメイン
    - さくらのドメイン（[auto-reply-bot](https://auto-reply-bot)）

#### 権限

[ デフォルト ]

- Admin

Permission

```
8
```

[ 必要最低限 ]

- Read Messages/View Channels
- Send Messages
- Send Messages in Threads
- Manage Messages
- Embed Links
- Read Message History
- Use External Emojis
- Use External Stickers

Permission

```
412317215744
```

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
  "admin_role_id": "1056464506554957824",
  "block": [
    {
      "name": "hello",
      "keyword": [
        "1",
        "2",
        "3",
        "4",
        "5",
        "6",
        "7",
        "8",
        "9",
        "10"
      ],
      "reply": [
        "1",
        "2",
        "3",
        "4",
        "5",
        "6",
        "7",
        "8",
        "9",
        "10"
      ],
      "match_condition": "one-contain",
      "is_random": true,
      "is_embed": true
    },
    {
      "name": "おはよう",
      "keyword": [
        "おはよう"
      ],
      "reply": [
        "おはようございます！"
      ],
      "match_condition": "all-contain",
      "is_random": false,
      "is_embed": false
    }
  ],
  "token": "hAYNUK7nWPkNYd20Fmo9cJ7cvDsk8M",
  "server_name": "TEST 2",
  "avatar_url": "",
  "role": [
    {
      "id": "1055362036495826964",
      "name": "自動返信botの管理者"
    },
    {
      "id": "1056464506554957824",
      "name": "テストロールです"
    },
    {
      "id": "1056544962973532196",
      "name": "[test]自動返信bot"
    },
    {
      "id": "1056894585156145207",
      "name": "Comment-bot"
    }
  ],
  "channel": [
    {
      "id": "1055315616002740297",
      "name": "一般"
    },
    {
      "id": "1055359277683986433",
      "name": "掲示板"
    }
  ],
  "nickname": "",
  "rule": {
    "url": {
      "is_restrict": true,
      "is_youtube_allow": false,
      "is_twitter_allow": true,
      "is_gif_allow": false,
      "is_opensea_allow": false,
      "is_discord_allow": false,
      "allow_role_id": [
        "1056464506554957824"
      ],
      "allow_channel_id": []
    }
  }
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
  "admin_role_id": "1056464506554957824",
  "block": [
    {
      "name": "hello",
      "keyword": [
        "1",
        "2",
        "3",
        "4",
        "5",
        "6",
        "7",
        "8",
        "9",
        "10"
      ],
      "reply": [
        "1",
        "2",
        "3",
        "4",
        "5",
        "6",
        "7",
        "8",
        "9",
        "10"
      ],
      "match_condition": "one-contain",
      "is_random": true,
      "is_embed": true
    },
    {
      "name": "おはよう",
      "keyword": [
        "おはよう"
      ],
      "reply": [
        "おはようございます！"
      ],
      "match_condition": "all-contain",
      "is_random": false,
      "is_embed": false
    }
  ],
  "rule": {
    "url": {
      "is_restrict": true,
      "is_youtube_allow": false,
      "is_twitter_allow": true,
      "is_gif_allow": false,
      "is_opensea_allow": false,
      "is_discord_allow": false,
      "allow_role_id": [
        "1056464506554957824"
      ],
      "allow_channel_id": []
    }
  }
}
```

##### レスポンス

```json
{
  "id": "1055315616002740294",
  "admin_role_id": "1056464506554957824",
  "block": [
    {
      "name": "hello",
      "keyword": [
        "1",
        "2",
        "3",
        "4",
        "5",
        "6",
        "7",
        "8",
        "9",
        "10"
      ],
      "reply": [
        "1",
        "2",
        "3",
        "4",
        "5",
        "6",
        "7",
        "8",
        "9",
        "10"
      ],
      "match_condition": "one-contain",
      "is_random": true,
      "is_embed": true
    },
    {
      "name": "おはよう",
      "keyword": [
        "おはよう"
      ],
      "reply": [
        "おはようございます！"
      ],
      "match_condition": "all-contain",
      "is_random": false,
      "is_embed": false
    }
  ],
  "server_name": "TEST 2",
  "avatar_url": "",
  "role": [
    {
      "id": "1055362036495826964",
      "name": "自動返信botの管理者"
    },
    {
      "id": "1056464506554957824",
      "name": "テストロールです"
    },
    {
      "id": "1056544962973532196",
      "name": "[test]自動返信bot"
    },
    {
      "id": "1056894585156145207",
      "name": "Comment-bot"
    }
  ],
  "channel": [
    {
      "id": "1055315616002740297",
      "name": "一般"
    },
    {
      "id": "1055359277683986433",
      "name": "掲示板"
    }
  ],
  "rule": {
    "url": {
      "is_restrict": true,
      "is_youtube_allow": false,
      "is_twitter_allow": true,
      "is_gif_allow": false,
      "is_opensea_allow": false,
      "is_discord_allow": false,
      "allow_role_id": [
        "1056464506554957824"
      ],
      "allow_channel_id": []
    }
  }
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

### 3. ヘルプ(Admin,Owner,Dev)

```
/help
```

##### アクセスできるユーザー

- Dev
- サーバーオーナー
- 管理者ロールを持つユーザー

##### アクセスできるサーバー

- 全てのサーバーで実行可能

## セキュリティについて

### 1.VPSがハックされた場合

[ ハッカーができること ]

- botを使用してスキャムURLを流す(everyoneも可能)

[ 盗まれる情報 ]

- botのアクセストークン
- プログラム（書き換え可能）

[ ハックされた時の対応 ]

- アクセストークンを再発行（https://discord.com/developers/applications/1056843645967413309/bot）

### 2.Discordがハックされた場合

[ ハッカーができること ]

- Totsumaru#7777 のアカウントを使用
- Totsumaru#7777 のパスワードを変更
- botの情報にアクセス & アクセストークンを再設定

[ ハックされた時の対応 ]

- パスワードの再設定
- アクセストークンを再発行（https://discord.com/developers/applications/1056843645967413309/bot）
