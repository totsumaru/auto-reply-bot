package guild_create

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/cmd"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"log"
	"os"
)

// サーバーにコマンドを追加します
func AddCommandHandler(s *discordgo.Session, m *discordgo.GuildCreate) {
	appID := os.Getenv("DISCORD_APPLICATION_ID")

	// =========================================
	// TESTサーバーの場合は以下のコマンドを追加します
	// =========================================

	if m.Guild.ID == conf.TestServerID {
		fmt.Println("専用コマンドを追加します id: ", m.Guild.ID)
		// create-server: サーバーを作成します
		{
			c := &discordgo.ApplicationCommand{
				Name:        cmd.CmdCreateServer.Name,
				Description: cmd.CmdCreateServer.Description,
				Options:     cmd.CmdCreateServer.Options,
			}

			_, err := s.ApplicationCommandCreate(appID, m.Guild.ApplicationID, c)
			if err != nil {
				log.Fatalln(err)
			}
		}

		// delete-server: 登録を削除 & サーバーからbotを削除します
		{
			c := &discordgo.ApplicationCommand{
				Name:        cmd.CmdDeleteServer.Name,
				Description: cmd.CmdDeleteServer.Description,
				Options:     cmd.CmdDeleteServer.Options,
			}

			_, err := s.ApplicationCommandCreate(appID, m.Guild.ID, c)
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
		c := &discordgo.ApplicationCommand{
			Name:        cmd.CmdHelp.Name,
			Description: cmd.CmdHelp.Description,
		}

		_, err := s.ApplicationCommandCreate(appID, m.Guild.ID, c)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
