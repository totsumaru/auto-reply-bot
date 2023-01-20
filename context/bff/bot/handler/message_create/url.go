package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/handler/shared"
)

// URL制限について確認します
func URL(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 処理内容はmessage_updateと同じのため、sharedにて共通化しています。
	if m.Message != nil {
		shared.CheckAndHandleURLContainMessage(s, m.Message)
	}
}
