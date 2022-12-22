package critical

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーからbotを削除します
func LeaveFromServer(s *discordgo.Session, guildID string) error {
	if err := s.GuildLeave(guildID); err != nil {
		return errors.NewError("サーバーからbotを削除できません", err)
	}

	return nil
}
