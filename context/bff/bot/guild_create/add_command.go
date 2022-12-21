package guild_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord"
	"log"
	"os"
)

// サーバーにコマンドを追加します
func AddCommandHandler(s *discordgo.Session, m *discordgo.GuildCreate) {
	appID := os.Getenv("DISCORD_APPLICATION_ID")

	// =========================================
	// TESTサーバーの場合は以下のコマンドを追加します
	// =========================================

	if m.Guild.ID == discord.TestServerID {
		// create-server: サーバーを作成します
		{
			opt1 := &discordgo.ApplicationCommandOption{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "server-id",
				Description:  "新規作成するサーバーのIDを入力します",
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				Required:     true,
			}

			cmd := &discordgo.ApplicationCommand{
				Name:        "create-server",
				Description: "サーバーを作成します",
				Options:     []*discordgo.ApplicationCommandOption{opt1},
			}

			_, err := s.ApplicationCommandCreate(appID, m.Guild.ApplicationID, cmd)
			if err != nil {
				log.Fatalln(err)
			}
		}

		// delete-server: 登録を削除 & サーバーからbotを削除します
		{
			opt1 := &discordgo.ApplicationCommandOption{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "server-id",
				Description:  "削除するサーバーのIDを入力します",
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				Required:     true,
			}

			cmd := &discordgo.ApplicationCommand{
				Name:        "delete-server",
				Description: "登録を削除 & サーバーからbotを削除します",
				Options:     []*discordgo.ApplicationCommandOption{opt1},
			}

			_, err := s.ApplicationCommandCreate(appID, m.Guild.ID, cmd)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	// =========================================
	// 全てのサーバーに以下のコマンドを追加します
	// =========================================

	// help: 設定に関する情報を返します
	{
		cmd := &discordgo.ApplicationCommand{
			Name:        "help",
			Description: "設定に関する情報を表示します",
		}

		_, err := s.ApplicationCommandCreate(appID, m.Guild.ID, cmd)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
