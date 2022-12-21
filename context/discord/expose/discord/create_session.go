package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"os"
)

// セッションを作成します
func CreateSession() (*discordgo.Session, error) {
	var Token = "Bot " + os.Getenv("APP_BOT_TOKEN")

	session, err := discordgo.New(Token)
	session.Token = Token
	if err != nil {
		return nil, errors.NewError("Discordを初期化できません")
	}

	return session, nil
}
