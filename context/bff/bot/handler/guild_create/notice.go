package guild_create

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const ContentTmpl = `
サーバー名: **%s**
ID: %s
`

// 新規サーバーに導入されたことを通知します
func NoticeHandler(s *discordgo.Session, m *discordgo.GuildCreate) {
	req := message_send.SendMessageEmbedReq{
		ChannelID: conf.ErrMsgChannelID,
		Title:     "新規サーバーに追加されました",
		Content:   fmt.Sprintf(ContentTmpl, m.Guild.Name, m.Guild.ID),
		Color:     conf.ColorGreen,
	}

	if err := message_send.SendMessageEmbed(s, req); err != nil {
		message_send.SendErrMsg(
			s, errors.NewError("新規サーバー導入時の通知を送信できません", err), m.Guild.Name,
		)
	}
}
