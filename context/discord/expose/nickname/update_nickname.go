package nickname

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// ニックネームを変更します
func UpdateNickname(s *discordgo.Session, guildID, nickname string) error {
	if err := s.GuildMemberNickname(guildID, "@me", nickname); err != nil {
		return errors.NewError("ニックネームを変更できません", err)
	}

	return nil
}
