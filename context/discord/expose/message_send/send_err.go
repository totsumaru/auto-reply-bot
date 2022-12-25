package message_send

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/shared/now"
	"log"
)

// エラーメッセージを送信します
func SendErrMsg(s *discordgo.Session, e error) {
	// エラーメッセージを送信します
	embedInfo := &discordgo.MessageEmbed{
		Title:       "エラーが発生しました",
		Description: e.Error(),
		Color:       conf.ColorRed,
		Timestamp:   now.GetNowTimeStamp(),
	}

	_, err := s.ChannelMessageSendEmbed(conf.ErrMsgChannelID, embedInfo)
	if err != nil {
		log.Println(err)
	}
}
