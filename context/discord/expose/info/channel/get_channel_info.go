package channel

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// チャンネルIDからチャンネルを取得します
func GetChannel(s *discordgo.Session, channelID string) (*discordgo.Channel, error) {
	c, err := s.Channel(channelID)
	if err != nil {
		return nil, errors.NewError("チャンネルIDからチャンネルを取得できません", err)
	}

	return c, nil
}

// サーバーのチャンネルを全て取得します
func GetAllChannels(s *discordgo.Session, guildID string) ([]*discordgo.Channel, error) {
	ch, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, errors.NewError("サーバーの全てのチャンネルを取得できません", err)
	}

	return ch, nil
}
