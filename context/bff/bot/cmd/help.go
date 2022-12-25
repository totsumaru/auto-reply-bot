package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/cmd"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"os"
)

const msg = `
あらかじめ決められた条件に一致するコメントが送信された場合、
自動で返信をするbotです。

■ 各種URL
・[設定はこちら](%s)
・[管理者に問い合わせ(Twitter)](https://twitter.com/totsumaru_dot)
・[botの導入はこちらから](https://localhost:8080)
`

// TODO: 返信はEphemeralの埋め込みにしたい

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

		req := message_send.SendMessageEmbedReq{
			ChannelID: m.ChannelID,
			Title:     "botについて",
			Content:   fmt.Sprintf(msg, os.Getenv("FE_ROOT_URL")),
			Color:     conf.ColorYellow,
		}

		if err := message_send.SendMessageEmbed(s, req); err != nil {
			message_send.SendErrMsg(s, errors.NewError("helpコマンドに対する返信を遅れません", err))
			return
		}
	},
}
