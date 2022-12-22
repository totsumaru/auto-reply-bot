package guild_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/bot/cmd"
	discordCmd "github.com/techstart35/auto-reply-bot/context/discord/expose/cmd"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーにコマンドを追加します
func AddCommandHandler(s *discordgo.Session, m *discordgo.GuildCreate) {
	// =========================================
	// TESTサーバーの場合は以下のコマンドを追加します
	// =========================================
	if m.Guild.ID == conf.TestServerID {
		// create-server: サーバーを作成します
		_, err := discordCmd.AddCommand(s, m.Guild.ID, cmd.CmdCreateServer)
		if err != nil {
			message_send.SendErrMsg(s, errors.NewError("create-serverコマンドを追加できません", err))
			return
		}

		// delete-server: 登録を削除 & サーバーからbotを削除します
		_, err = discordCmd.AddCommand(s, m.Guild.ID, cmd.CmdDeleteServer)
		if err != nil {
			message_send.SendErrMsg(s, errors.NewError("delete-serverコマンドを追加できません", err))
			return
		}
	}

	// =========================================
	// 全てのサーバーに以下のコマンドを追加します
	// =========================================

	// help: 設定に関する情報を返します
	_, err := discordCmd.AddCommand(s, m.Guild.ID, cmd.CmdDeleteServer)
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("delete-serverコマンドを追加できません", err))
		return
	}
}
