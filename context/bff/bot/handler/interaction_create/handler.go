package interaction_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/cmd"
)

// コマンドが実行された時のハンドラーです
func Handler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	cmd.CmdCreateServer.Handler(s, m)
	cmd.CmdDeleteServer.Handler(s, m)
	cmd.CmdHelp.Handler(s, m)
}
