package interaction_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/cmd"
)

// create-server コマンドが実行された時のハンドラーです
func CreateServerHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	cmd.CmdCreateServer.Handler(s, m)
}
