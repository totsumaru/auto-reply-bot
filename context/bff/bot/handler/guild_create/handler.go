package guild_create

import (
	"github.com/bwmarrin/discordgo"
)

// botが追加された時のハンドラーです
func Handler(s *discordgo.Session, m *discordgo.GuildCreate) {
	AddCommandHandler(s, m)
	NoticeHandler(s, m)
}
