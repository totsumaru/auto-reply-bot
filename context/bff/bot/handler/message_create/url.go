package message_create

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/check"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/manage_message"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/rule"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"strings"
)

// 不正URLが送信された時の削除後のメッセージです
const InvalidURLReplyTmpl = `
許可されていないURLが投稿されたので、メッセージを削除しました。
許可されているURL: %s
`

// アラートチャンネルに送信するメッセージです
const AlertChannelMessageTmpl = `
以下の内容で不正なURLが送信されたので、botが削除しました。
---
[チャンネル]
<#%s>
---
[送信された内容(URL注意)]
%s
`

// URL制限について確認します
func URL(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TEST SERVERはカウントしません
	if m.GuildID == conf.TestServerID {
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

	content := m.Content

	storeRes, err := v1.GetStoreRes(m.GuildID)
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("IDでサーバーを取得できません", err), guildName)
		return
	}

	ok, err := isAllowedURLMessage(s, storeRes, m.Author.ID, m.ChannelID, m.Content)
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("IDでサーバーを取得できません", err), guildName)
		return
	}
	if !ok {
		// 不正URLの含まれたメッセージを削除します
		{
			if err := manage_message.DeleteMessage(s, m.ChannelID, m.Message.ID); err != nil {
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

			req := message_send.SendMessageEmbedReq{
				ChannelID: m.ChannelID,
				Title:     "許可されていないURLです",
				Content:   fmt.Sprintf(InvalidURLReplyTmpl, allowURLs),
				Color:     conf.ColorBlack,
			}
			if err = message_send.SendMessageEmbed(s, req); err != nil {
				message_send.SendErrMsg(s, errors.NewError("埋め込みメッセージを送信できません", err), guildName)
			}
		}

		// アラートチャンネルに詳細を送信します
		{
			alertCh := storeRes.Rule.URL.AlertChannelID
			if alertCh != "none" {
				req := message_send.SendMessageEmbedReq{
					ChannelID: alertCh,
					Title:     "[ALERT]許可されていないURLが送信されました",
					Content:   fmt.Sprintf(AlertChannelMessageTmpl, m.ChannelID, content),
					Color:     conf.ColorRed,
				}
				if err = message_send.SendMessageEmbed(s, req); err != nil {
					message_send.SendErrMsg(s, errors.NewError("埋め込みメッセージを送信できません", err), guildName)
				}
			}
		}
	}
}

// メッセージが許可されているか検証します
func isAllowedURLMessage(
	s *discordgo.Session,
	storeRes v1.Res,
	authorID string,
	channelID string,
	msg string,
) (bool, error) {
	urlRule := storeRes.Rule.URL

	// URL制限していない場合はここで終了
	if !urlRule.IsRestrict {
		return true, nil
	}

	// httpが何個含まれているか確認(含まれていなければここで終了)
	httpCount := strings.Count(msg, rule.URLPrefixHTTP)
	httpsCount := strings.Count(msg, rule.URLPrefixHTTPS)
	urlCount := httpCount + httpsCount

	if urlCount == 0 {
		return true, nil
	}

	// 許可されているチャンネルの場合はここで終了
	for _, chID := range urlRule.AllowChannelID {
		if chID == channelID {
			return true, nil
		}
	}

	// 許可されているロールの場合はここで終了
	for _, roleID := range urlRule.AllowRoleID {
		ok, err := check.HasRole(s, storeRes.ID, authorID, roleID)
		if err != nil {
			return false, errors.NewError("ロール所持確認ができません", err)
		}
		if ok {
			return true, nil
		}
	}

	allowURLCount := 0

	// YouTubeのURLの個数をカウントに追加
	if urlRule.IsYoutubeAllow {
		allowURLCount += strings.Count(msg, rule.YoutubeURL)
	}
	// TwitterのURLの個数をカウントに追加
	if urlRule.IsTwitterAllow {
		allowURLCount += strings.Count(msg, rule.TwitterURL)
	}
	// GIFのURLの個数をカウントに追加
	if urlRule.IsGIFAllow {
		allowURLCount += strings.Count(msg, rule.GIFURL)
	}

	// httpの個数と、許可されたURLの個数が一致した場合はOK
	if urlCount == allowURLCount {
		return true, nil
	}

	return false, nil
}
