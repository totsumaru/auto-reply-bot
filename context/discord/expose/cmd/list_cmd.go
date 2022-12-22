package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"os"
)

// コマンドの一覧を取得します
//
// コマンド名:id のmapを返します。
func ListCmd(s *discordgo.Session, guildID string) (map[string]string, error) {
	appID := os.Getenv("DISCORD_APPLICATION_ID")

	res := map[string]string{}

	apps, err := s.ApplicationCommands(appID, guildID)
	if err != nil {
		return res, errors.NewError("コマンドの一覧を取得できません", err)
	}

	for _, v := range apps {
		res[v.Name] = v.ID
	}

	return res, nil
}
