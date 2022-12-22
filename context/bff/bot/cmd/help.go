package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord/cmd"
)

// 設定に関する情報を取得します
var CmdHelp = cmd.CMD{
	Name:        CMDNameHelp,
	Description: "設定に関する情報を表示します",
	Handler: func(s *discordgo.Session, m *discordgo.InteractionCreate) {

	},
}
