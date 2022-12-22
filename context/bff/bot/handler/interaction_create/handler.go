package interaction_create

import (
	"github.com/bwmarrin/discordgo"
)

// コマンドが実行された時のハンドラーです
func Handler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	CreateServerHandler(s, m)
	DeleteServerHandler(s, m)
	HelpHandler(s, m)
}
