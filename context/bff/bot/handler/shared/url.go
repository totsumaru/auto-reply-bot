package shared

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/check"
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

// メッセージが許可されているか検証します
func IsAllowedURLMessage(
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
