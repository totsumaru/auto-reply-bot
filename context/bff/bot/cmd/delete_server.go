package cmd

import "github.com/bwmarrin/discordgo"

// サーバーを削除するコマンドです
var CmdDeleteServer = CMD{
	Name:        CMDNameDeleteServer,
	Description: "登録を削除 & サーバーからbotを削除します",
	Options: []*discordgo.ApplicationCommandOption{
		// 引数1: サーバーID
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "server-id",
			Description:  "削除するサーバーのIDを入力します",
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			Required:     true,
		},
	},
	Handler: func(s *discordgo.Session, m *discordgo.InteractionCreate) {

	},
}
