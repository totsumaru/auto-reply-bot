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
		return errors.NewError("埋め込みメッセージを送信できません", err)
	}

	return nil
}

// アイコン付きの埋め込みメッセージを送信するリクエストです
type SendMessageEmbedWithIconReq struct {
	ChannelID  string
	Title      string
	Content    string
	Color      int
	IconURL    string
	Name       string
	FooterText string
}

// アイコン付きの埋め込みメッセージを送信します
func SendMessageEmbedWithIcon(s *discordgo.Session, req SendMessageEmbedWithIconReq) error {
	author := &discordgo.MessageEmbedAuthor{}
	if req.Name != "" {
		author.Name = req.Name
		author.IconURL = req.IconURL
	}

	embed := &discordgo.MessageEmbed{
		Title:       req.Title,
		Description: req.Content,
		Color:       req.Color,
		Author:      author,
		Footer: &discordgo.MessageEmbedFooter{
			Text: req.FooterText,
		},
	}

	_, err := s.ChannelMessageSendEmbed(req.ChannelID, embed)
	if err != nil {
		return errors.NewError("アイコン付きの埋め込みメッセージを送信できません", err)
	}

	return nil
}
