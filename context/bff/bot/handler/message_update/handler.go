package message_update

import (
	"github.com/bwmarrin/discordgo"
)

// メッセージが更新された時のハンドラーです
func Handler(s *discordgo.Session, m *discordgo.MessageUpdate) {
	URL(s, m)
}
