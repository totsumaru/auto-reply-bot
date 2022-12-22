package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーを削除するコマンドです
var CmdDeleteServer = CMD{
	Name:        CMDNameDeleteServer,
	Description: "登録を削除 & サーバーからbotを削除します",
	Options: []*discordgo.ApplicationCommandOption{
		// 引数1: サーバーID
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "server-id",
			Description:  "削除するサーバーのIDを入力します",
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			Required:     true,
		},
	},
	Handler: func(s *discordgo.Session, m *discordgo.InteractionCreate) {
		// 検証します
		{
			// コマンドが正しいかを検証します
			if m.Interaction.ApplicationCommandData().Name != CMDNameDeleteServer {
				return
			}

			// Devであるかを検証します
			if m.Member.User.ID != discord.TotsumaruDiscordID {
				if err := discord.SendEphemeralReply(s, m, "権限がありません"); err != nil {
					discord.SendInteractionErrMsg(s, m, err)
					return
				}
				return
			}
		}

		ctx, tx, err := shared.CreateDBTx()
		if err != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}

		var id string
		for _, v := range m.Interaction.ApplicationCommandData().Options {
			if v.Name == "server-id" {
				id = v.Value.(string)
			}
		}

		bffErr := (func() error {
			// DBから削除します
			if err = v1.DeleteServer(s, ctx, id); err != nil {
				return errors.NewError("サーバーを作成できません", err)
			}

			// サーバーからbotを削除します
			if err := discord.LeaveFromServer(s, id); err != nil {
				return errors.NewError("サーバーからbotを削除できません", err)
			}

			return nil
		})()

		if bffErr != nil {
			// ロールバックを実行します
			txErr := tx.Rollback()
			if txErr != nil {
				msg := errors.NewError("ロールバックに失敗しました。データに不整合が発生している可能性があります。")
				discord.SendInteractionErrMsg(s, m, msg)
				return
			}

			discord.SendInteractionErrMsg(s, m, err)
			return
		}

		if txErr := tx.Commit(); txErr != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}

		gl, err := s.Guild(id)
		if err != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}

		msg := fmt.Sprintf("ID: %s, Name: %s を削除しました", gl.ID, gl.Name)
		if err := discord.SendReplyInteraction(s, m, msg); err != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}
	},
}
