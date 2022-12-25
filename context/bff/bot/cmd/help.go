package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/cmd"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"net/url"
)

// DiscordログインのURLテンプレートです
//
// fmt.Sprintf(DiscordLoginURLTmpl,{エンコードしたリダイレクトURL})
const DiscordLoginURLTmpl = "https://discord.com/api/oauth2/authorize?client_id=1055348253614419989&redirect_uri=%s&response_type=code&scope=identify"

const msg = `
決められた条件に一致するコメントが送信された場合、
自動で返信をするbotです。

■ 各種URL
・[管理者に問い合わせ(Twitter)](https://twitter.com/totsumaru_dot)

■ お知らせ
・設定画面は以下のボタンからアクセスしてください
・導入のご依頼はTwitterのDMからお願いします
`

// 設定に関する情報を取得します
var CmdHelp = cmd.CMD{
	Name:        CMDNameHelp,
	Description: "設定に関する情報を表示します",
	Handler: func(s *discordgo.Session, m *discordgo.InteractionCreate) {
		// 検証します
		{
			// コマンドが正しいかを検証します
			if m.Interaction.ApplicationCommandData().Name != CMDNameHelp {
				return
			}
		}

		redirectURL := shared.CreateDiscordLoginRedirectURL(m.GuildID)

		discordLoginURL := fmt.Sprintf(
			DiscordLoginURLTmpl,
			url.QueryEscape(redirectURL),
		)

		fmt.Println("URL: ", discordLoginURL)

		if err := message_send.SendEmbedEphemeralReplyWithURLBtn(
			s, m, "botについて", msg, discordLoginURL, conf.ColorYellow,
		); err != nil {
			message_send.SendErrMsg(s, errors.NewError("helpコマンドに対する返信を遅れません", err))
			return
		}
	},
}
