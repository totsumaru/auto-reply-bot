package message_create

import (
	"github.com/bwmarrin/discordgo"
)

// メッセージが作成された時のハンドラーです
func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	Reply(s, m)
	Weather(s, m)
}
