package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"os"
)

// サーバーにコマンドを追加します
func AddCommand(s *discordgo.Session, guildID string, cmd CMD) (string, error) {
	appID := os.Getenv("DISCORD_APPLICATION_ID")

	c := &discordgo.ApplicationCommand{
		Name:        cmd.Name,
		Description: cmd.Description,
		Options:     cmd.Options,
	}

	apCmd, err := s.ApplicationCommandCreate(appID, guildID, c)
	if err != nil {
		return "", errors.NewError("サーバーにコマンドを追加できません", err)
	}

	return apCmd.ID, nil
}
