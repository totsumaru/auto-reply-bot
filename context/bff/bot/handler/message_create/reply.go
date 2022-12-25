package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"math/rand"
	"strings"
	"time"
)

// 送信されたメッセージが条件と一致する場合は返信を送信します
func Reply(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TEST SERVERはカウントしません
	if m.GuildID == conf.TestServerID {
		return
	}

	// Botユーザーはカウントしません
	if m.Author.Bot {
		return
	}

	content := m.Content

	ctx, _, err := shared.CreateDBTx()
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("Txを作成できません", err))
		return
	}

	apiRes, err := v1.FindByID(ctx, m.GuildID)
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("IDでサーバーを取得できません", err))
		return
	}

	for _, block := range apiRes.Block {
		mustReply := true

		if block.IsAllMatch {
			// [全て含む場合]1つでも含んでいないキーワードがあれば終了
			for _, keyword := range block.Keyword {
				if !strings.Contains(content, keyword) {
					mustReply = false
					break
				}
			}
		} else {
			isContain := false
			// [1つでも含む場合]1つでも含んでいるキーワードがあれば終了
			for _, keyword := range block.Keyword {
				if strings.Contains(content, keyword) {
					isContain = true
					break
				}
			}
			// 1つも含んでいない場合
			mustReply = isContain
		}

		if mustReply {
			// 返信を送信します。
			//
			// ランダムに返信を返すかを確認します。
			msg := block.Reply[0]
			if block.IsRandom {
				msg = getRandomMessage(block.Reply)
			}

			if block.IsEmbed {
				// 埋め込みメッセージを送信します
				req := message_send.SendReplyEmbedReq{
					ChannelID: m.ChannelID,
					Content:   msg,
					Color:     conf.ColorCyan,
					Reference: m.Reference(),
				}
				if err = message_send.SendReplyEmbed(s, req); err != nil {
					message_send.SendErrMsg(s, errors.NewError("埋め込みの返信を送信できません", err))
					return
				}
			} else {
				// 通常のテキストメッセージを送信します
				if err = message_send.SendReply(s, m.GuildID, m.ChannelID, m.ID, msg); err != nil {
					message_send.SendErrMsg(s, errors.NewError("返信を送信できません", err))
					return
				}
			}
		}
	}
}

// メッセージのスライスからランダムに1つ取得します
func getRandomMessage(messages []string) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(messages))

	return messages[index]
}
