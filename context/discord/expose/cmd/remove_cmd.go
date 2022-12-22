package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"os"
)

// サーバーからコマンドを削除します
func RemoveCommand(s *discordgo.Session, guildID string, cmdID string) error {
	appID := os.Getenv("DISCORD_APPLICATION_ID")

	if err := s.ApplicationCommandDelete(appID, guildID, cmdID); err != nil {
		return errors.NewError("サーバーからコマンドを削除できません", err)
	}

	return nil
}
