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

// TODO: エラーの見直し。このエラーは、コマンド実行者に表示されるため、使用箇所は「それでいいのか？」を全て確認。

// インタラクションの失敗メッセージを送信します
func SendInteractionErrMsg(s *discordgo.Session, i *discordgo.InteractionCreate, e error) {
	embed := &discordgo.MessageEmbed{
		Title:       "エラーが発生しました",
		Description: e.Error(),
		Color:       conf.ColorRed,
		Timestamp:   now.GetNowTimeStamp(),
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsCrossPosted,
		},
	}
	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		log.Println(err)
	}
}
