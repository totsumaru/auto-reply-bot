package message_create

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/handler/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/manage_message"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"strings"
)

// URL制限について確認します
func URL(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TEST SERVERはカウントしません
	if m.GuildID == conf.TestServerID {
		return
	}

	// Webhookはカウントしません
	if m.Author == nil {
		return
	}

	// Botユーザーはカウントしません
	if m.Author.Bot {
		return
	}

	guildName, err := guild.GetGuildName(s, m.GuildID)
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("ギルド名を取得できません", err), "")
		return
	}

	storeRes, err := v1.GetStoreRes(m.GuildID)
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("IDでサーバーを取得できません", err), guildName)
		return
	}

	ok, err := shared.IsAllowedURLMessage(s, storeRes, m.Author.ID, m.ChannelID, m.Content)
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("IDでサーバーを取得できません", err), guildName)
		return
	}
	if !ok {
		// 不正URLの含まれたメッセージを削除します
		{
			if err := manage_message.DeleteMessage(s, m.ChannelID, m.ID); err != nil {
				message_send.SendErrMsg(s, errors.NewError("メッセージを削除できません", err), guildName)
				return
			}
		}

		// 投稿されたチャンネルにメッセージを返します
		{
			allowURLs := make([]string, 0)
			if storeRes.Rule.URL.IsYoutubeAllow {
				allowURLs = append(allowURLs, "YouTube")
			}
			if storeRes.Rule.URL.IsTwitterAllow {
				allowURLs = append(allowURLs, "Twitter")
			}
			if storeRes.Rule.URL.IsGIFAllow {
				allowURLs = append(allowURLs, "GIF")
			}
			if storeRes.Rule.URL.IsOpenseaAllow {
				allowURLs = append(allowURLs, "Opensea")
			}
			if storeRes.Rule.URL.IsDiscordAllow {
				allowURLs = append(allowURLs, "Discord")
			}

			fixedContent := strings.Replace(
				m.Content,
				"http",
				"\n**[URLが含まれています: 信頼できる場合のみ、アクセスしてください👇]**\n ⚠️ http",
				-1,
			)

			req := message_send.SendMessageEmbedWithIconReq{
				ChannelID: m.ChannelID,
				Content: fmt.Sprintf(
					shared.InvalidURLReplyTmpl,
					fixedContent,
				),
				Color:      conf.ColorGray,
				Name:       m.Author.Username,
				IconURL:    m.Author.AvatarURL(""),
				FooterText: fmt.Sprintf("スキャム対策として、このサーバーでは%s以外のURLはbotが監視しています。", allowURLs),
			}
			if err = message_send.SendMessageEmbedWithIcon(s, req); err != nil {
				message_send.SendErrMsg(s, errors.NewError("埋め込みメッセージを送信できません", err), guildName)
			}
		}
	}
}
