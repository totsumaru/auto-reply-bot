package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/cmd"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/critical"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーを削除するコマンドです
var CmdDeleteServer = cmd.CMD{
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
			if m.Member.User.ID != conf.TotsumaruDiscordID {
				if err := message_send.SendEphemeralReply(s, m, "権限がありません"); err != nil {
					message_send.SendInteractionErrMsg(s, m, err)
					return
				}
				return
			}
		}

		ctx, tx, err := shared.CreateDBTx()
		if err != nil {
			message_send.SendInteractionErrMsg(s, m, err)
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
			if err := critical.LeaveFromServer(s, id); err != nil {
				return errors.NewError("サーバーからbotを削除できません", err)
			}

			return nil
		})()

		if bffErr != nil {
			// ロールバックを実行します
			txErr := tx.Rollback()
			if txErr != nil {
				msg := errors.NewError("ロールバックに失敗しました。データに不整合が発生している可能性があります。")
				message_send.SendInteractionErrMsg(s, m, msg)
				return
			}

			message_send.SendInteractionErrMsg(s, m, err)
			return
		}

		if txErr := tx.Commit(); txErr != nil {
			message_send.SendInteractionErrMsg(s, m, err)
			return
		}

		gl, err := s.Guild(id)
		if err != nil {
			message_send.SendInteractionErrMsg(s, m, err)
			return
		}

		msg := fmt.Sprintf("ID: %s, Name: %s を削除しました", gl.ID, gl.Name)
		if err := message_send.SendReplyInteraction(s, m, msg); err != nil {
			message_send.SendInteractionErrMsg(s, m, err)
			return
		}
	},
}
