package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/now"
	"log"
)

// エラーメッセージを送信します
var SendErrMsg = func(s *discordgo.Session, e error) {
	// エラーメッセージを送信します
	embedInfo := &discordgo.MessageEmbed{
		Title:       "エラーが発生しました",
		Description: e.Error(),
		Color:       ColorRed,
		Timestamp:   now.GetNowTimeStamp(),
	}

	_, err := s.ChannelMessageSendEmbed(ErrMsgChannelID, embedInfo)
	if err != nil {
		log.Println(err)
	}
}
