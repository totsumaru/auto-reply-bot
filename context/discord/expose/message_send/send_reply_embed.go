package message_send

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// 埋め込みの返信を送信するリクエストです
type SendReplyEmbedReq struct {
	ChannelID string
	Content   string
	Color     int
	Reference *discordgo.MessageReference
	Thumbnail *discordgo.MessageEmbedThumbnail // サムネを入れる場合のみ
}

// 埋め込みの返信を送信します
func SendReplyEmbed(s *discordgo.Session, req SendReplyEmbedReq) error {
	embed := &discordgo.MessageEmbed{
		Description: req.Content,
		Color:       req.Color,
		Thumbnail:   req.Thumbnail,
	}

	_, err := s.ChannelMessageSendEmbedReply(req.ChannelID, embed, req.Reference)
	if err != nil {
		return errors.NewError("埋め込みの返信を送信できません", err)
	}

	return nil
}
