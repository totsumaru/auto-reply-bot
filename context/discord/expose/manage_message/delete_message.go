package manage_message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// メッセージを削除します
func DeleteMessage(s *discordgo.Session, channelID, messageID string) error {
	if err := s.ChannelMessageDelete(channelID, messageID); err != nil {
		return errors.NewError("メッセージを削除できません", err)
	}

	return nil
}
