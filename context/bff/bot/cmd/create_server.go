package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/cmd"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーを作成するコマンドです
var CmdCreateServer = cmd.CMD{
	Name:        CMDNameCreateServer,
	Description: "サーバーを作成します",
	Options: []*discordgo.ApplicationCommandOption{
		// 引数1: サーバーID
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "server-id",
			Description:  "新規作成するサーバーのIDを入力します",
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			Required:     true,
		},
	},
	// コマンドが実行された時の処理です
	Handler: func(s *discordgo.Session, m *discordgo.InteractionCreate) {
		// 検証します
		{
			// コマンドが正しいかを検証します
			if m.Interaction.ApplicationCommandData().Name != CMDNameCreateServer {
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

		var (
			apiRes = v1.Res{}
		)

		var id string
		for _, v := range m.Interaction.ApplicationCommandData().Options {
			if v.Name == "server-id" {
				id = v.Value.(string)
			}
		}

		bffErr := (func() error {
			apiRes, err = v1.CreateServer(s, ctx, id)
			if err != nil {
				return errors.NewError("サーバーを作成できません", err)
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

		gl, err := s.Guild(apiRes.ID)
		if err != nil {
			message_send.SendInteractionErrMsg(s, m, err)
			return
		}

		msg := fmt.Sprintf("ID: %s, Name: %s を追加しました", gl.ID, gl.Name)
		if err := message_send.SendReplyInteraction(s, m, msg); err != nil {
			message_send.SendInteractionErrMsg(s, m, err)
			return
		}
	},
}
