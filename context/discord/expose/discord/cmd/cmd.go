package cmd

import "github.com/bwmarrin/discordgo"

// コマンドの定義です
type CMD struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(s *discordgo.Session, m *discordgo.InteractionCreate)
}
