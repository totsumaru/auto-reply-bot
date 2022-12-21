package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// メッセージを送信します
func SendMessage(s *discordgo.Session, channelID, text string) error {
	_, err := s.ChannelMessageSend(channelID, text)
	if err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}
