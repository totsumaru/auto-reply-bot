package interaction_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/cmd"
)

// help コマンドが実行された時のハンドラーです
func HelpHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	cmd.CmdHelp.Handler(s, m)
}
