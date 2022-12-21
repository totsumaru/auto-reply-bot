package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// ユーザーが指定のロールを保持しているかを確認します
func HasRole(
	s *discordgo.Session,
	guildID string,
	userID string,
	roleID string,
) (bool, error) {
	m, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false, errors.NewError("Discordのメンバーを取得できません", err)
	}

	for _, rID := range m.Roles {
		if rID == roleID {
			return true, nil
		}
	}

	return false, nil
}
