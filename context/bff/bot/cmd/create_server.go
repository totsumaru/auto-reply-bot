package cmd

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/db"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"log"
)

// サーバーを作成するコマンドです
var CmdCreateServer = CMD{
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
		// Devであるかを検証します
		{
			if m.Member.User.ID != discord.TotsumaruDiscordID {
				if err := discord.SendEphemeralReply(s, m, "権限がありません"); err != nil {
					log.Fatalln(err)
				}
			}
		}

		conf, err := db.NewConf()
		if err != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}

		database, err := db.NewDB(conf)
		if err != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}

		tx, err := database.Begin()
		if err != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}

		ctx := context.WithValue(context.Background(), "tx", tx)

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

		gl, err := s.Guild(apiRes.ID)
		if err != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}

		msg := fmt.Sprintf("ID: %s, Name: %s を追加しました", gl.ID, gl.Name)
		if err := discord.SendReplyInteraction(s, m, msg); err != nil {
			discord.SendInteractionErrMsg(s, m, err)
			return
		}
	},
}
