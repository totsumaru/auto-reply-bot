package message_send

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

// URLボタン付きの他人には見えない返信を送信します
func SendEmbedEphemeralReplyWithURLBtn(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	title string,
	content string,
	url string,
	color int,
) error {
	btn := discordgo.Button{
		Label: "設定画面はこちら",
		Style: discordgo.LinkButton,
		URL:   url,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn},
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: content,
		Color:       color,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{actions},
			Embeds:     []*discordgo.MessageEmbed{embed},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	}

	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("非公開のレスポンスを送信できません", err)
	}

	return nil
}
