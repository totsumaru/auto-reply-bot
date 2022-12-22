package message_send

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// 返信を送信します
func SendReply(
	s *discordgo.Session,
	guildID string,
	channelID string,
	messageID string,
	content string,
) error {
	_, err := s.ChannelMessageSendReply(channelID, content, &discordgo.MessageReference{
		MessageID: messageID,
		ChannelID: channelID,
		GuildID:   guildID,
	})
	if err != nil {
		return errors.NewError("返信を送信できません", err)
	}

	return nil
}
