package cmd

import "github.com/bwmarrin/discordgo"

// コマンド名の一覧です
const (
	CMDNameCreateServer = "create-server"
	CMDNameDeleteServer = "delete-server"
	CMDNameHelp         = "help"
)

// コマンドの定義です
type CMD struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(s *discordgo.Session, m *discordgo.InteractionCreate)
}
