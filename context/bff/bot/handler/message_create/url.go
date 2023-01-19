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
許可の無いURLが送信されました。
元メッセージを削除し、全てのURLを無効化しました。
▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
%s
▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
許可されているURL: %s
送信者: <@%s>
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
			if storeRes.Rule.URL.IsOpenseaAllow {
				allowURLs = append(allowURLs, "Opensea")
			}
			if storeRes.Rule.URL.IsDiscordAllow {
				allowURLs = append(allowURLs, "Discord")
			}

			fixedContent := strings.Replace(m.Content, "http", "h ttp", -1)

			req := message_send.SendMessageEmbedReq{
				ChannelID: m.ChannelID,
				Content: fmt.Sprintf(
					InvalidURLReplyTmpl,
					fixedContent,
					allowURLs,
					m.Author.ID,
				),
				Color: conf.ColorBlack,
			}
			if err = message_send.SendMessageEmbed(s, req); err != nil {
				message_send.SendErrMsg(s, errors.NewError("埋め込みメッセージを送信できません", err), guildName)
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
		allowURLCount += strings.Count(msg, rule.YoutubeShareURL)
		allowURLCount += strings.Count(msg, rule.YoutubeWWWURL)
	}
	// TwitterのURLの個数をカウントに追加
	if urlRule.IsTwitterAllow {
		allowURLCount += strings.Count(msg, rule.TwitterURL)
		allowURLCount += strings.Count(msg, rule.TwitterWWWURL)
	}
	// GIFのURLの個数をカウントに追加
	if urlRule.IsGIFAllow {
		allowURLCount += strings.Count(msg, rule.GIFURL)
		allowURLCount += strings.Count(msg, rule.GIFWWWURL)
	}
	// OpenseaのURLの個数をカウントに追加
	if urlRule.IsOpenseaAllow {
		allowURLCount += strings.Count(msg, rule.OpenseaURL)
		allowURLCount += strings.Count(msg, rule.OpenseaWWWURL)
		allowURLCount += strings.Count(msg, rule.OpenseaTestnetURL)
	}
	// DiscordのURLの個数をカウントに追加
	if urlRule.IsDiscordAllow {
		allowURLCount += strings.Count(msg, rule.DiscordURL)
		allowURLCount += strings.Count(msg, rule.DiscordWWWURL)
	}

	// httpの個数と、許可されたURLの個数が一致した場合はOK
	if urlCount <= allowURLCount {
		return true, nil
	}

	return false, nil
}
