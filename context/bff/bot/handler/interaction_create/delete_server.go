package interaction_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/cmd"
)

// delete-server コマンドが実行された時のハンドラーです
func DeleteServerHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	cmd.CmdDeleteServer.Handler(s, m)
}
