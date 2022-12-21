package cmd

import "github.com/bwmarrin/discordgo"

// 設定に関する情報を取得します
var CmdHelp = CMD{
	Name:        CMDNameHelp,
	Description: "設定に関する情報を表示します",
	Handler: func(s *discordgo.Session, m *discordgo.InteractionCreate) {

	},
}
