package message_send

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/shared/now"
	"log"
)

// エラーメッセージを送信します
func SendErrMsg(s *discordgo.Session, e error, serverName string) {
	// エラーメッセージを送信します
	embedInfo := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("エラー発生: @%s", serverName),
		Description: e.Error(),
		Color:       conf.ColorRed,
		Timestamp:   now.GetNowTimeStamp(),
	}

	_, err := s.ChannelMessageSendEmbed(conf.ErrMsgChannelID, embedInfo)
	if err != nil {
		log.Println(err)
	}
}

// インタラクションの失敗メッセージを送信します
func SendEphemeralInteractionErrMsg(s *discordgo.Session, i *discordgo.InteractionCreate, e error) {
	embed := &discordgo.MessageEmbed{
		Title:       "ERROR",
		Description: e.Error(),
		Color:       conf.ColorRed,
		Timestamp:   now.GetNowTimeStamp(),
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}
	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		log.Println(err)
	}
}
