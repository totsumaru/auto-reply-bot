package message_update

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/handler/shared"
)

// URL制限について確認します
func URL(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// 処理内容はmessage_createと同じのため、sharedにて共通化しています。
	if m.Message != nil {
		shared.CheckAndHandleURLContainMessage(s, m.Message)
	}
}
