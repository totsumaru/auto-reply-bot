package message_send

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// 埋め込みメッセージを送信するリクエストです
type SendMessageEmbedReq struct {
	ChannelID string
	Title     string
	Content   string
	Color     int
}

// 埋め込みメッセージを送信します
func SendMessageEmbed(s *discordgo.Session, req SendMessageEmbedReq) error {
	embed := &discordgo.MessageEmbed{
		Title:       req.Title,
		Description: req.Content,
		Color:       req.Color,
	}

	_, err := s.ChannelMessageSendEmbed(req.ChannelID, embed)
	if err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}
