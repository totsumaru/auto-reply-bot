package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// 他人には見えない返信を送信します
func SendEphemeralReply(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	content string,
) error {
	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("非公開のレスポンスを送信できません", err)
	}

	return nil
}
