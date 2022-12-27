package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/check"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/cmd"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"net/url"
	"os"
)

// DiscordログインのURLテンプレートです
//
// fmt.Sprintf(DiscordLoginURLTmpl,{ENVのDISCORD_CLIENT_ID},{エンコードしたリダイレクトURL})
const DiscordLoginURLTmpl = "https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify"

const msg = `
決められた条件に一致するコメントが送信された場合、
自動で返信をするbotです。

■ 各種URL
・[管理者に問い合わせ(Twitter)](https://twitter.com/totsumaru_dot)

■ お知らせ
・設定画面は以下のボタンからアクセスしてください
・導入のご依頼はTwitterのDMからお願いします
`

// 設定に関する情報を取得します
var CmdHelp = cmd.CMD{
	Name:        CMDNameHelp,
	Description: "設定に関する情報を表示します",
	Handler: func(s *discordgo.Session, m *discordgo.InteractionCreate) {
		// コマンドが実行されたサーバーのIDです
		guildName, err := guild.GetGuildName(s, m.GuildID)
		if err != nil {
			message_send.SendErrMsg(s, errors.NewError("ギルド名を取得できません", err), "")
			return
		}

		ctx, _, err := shared.CreateDBTx()
		if err != nil {
			message_send.SendEphemeralInteractionErrMsg(s, m, fmt.Errorf("エラーが発生しました"))
			message_send.SendErrMsg(s, errors.NewError("DBトランザクションを作成できません", err), guildName)
			return
		}

		// 検証します
		{
			// コマンドが正しいかを検証します
			if m.Interaction.ApplicationCommandData().Name != CMDNameHelp {
				return
			}

			// Dev,Owner,Adminのどれかであることを確認します
			{
				tmpRes, err := v1.FindByID(ctx, m.GuildID)
				if err != nil {
					message_send.SendEphemeralInteractionErrMsg(s, m, fmt.Errorf("エラーが発生しました"))
					message_send.SendErrMsg(s, errors.NewError("IDでサーバーを取得できません", err), guildName)
					return
				}

				ok, err := check.HasRole(s, m.GuildID, m.Member.User.ID, tmpRes.AdminRoleID)
				if err != nil {
					message_send.SendEphemeralInteractionErrMsg(s, m, fmt.Errorf("エラーが発生しました"))
					message_send.SendErrMsg(s, errors.NewError("ロールの所有を確認できません", err), guildName)
					return
				}

				guildOwnerID, err := guild.GetGuildOwnerID(s, m.GuildID)
				if err != nil {
					message_send.SendEphemeralInteractionErrMsg(s, m, fmt.Errorf("エラーが発生しました"))
					message_send.SendErrMsg(s, errors.NewError("ギルドのオーナーIDを取得できません", err), guildName)
					return
				}

				userID := m.Member.User.ID
				if !(ok || userID == guildOwnerID || userID == conf.TotsumaruDiscordID) {
					message_send.SendEphemeralInteractionErrMsg(s, m, fmt.Errorf("権限がありません"))
					// Devにエラーメッセージは送信しません
					return
				}
			}
		}

		redirectURL := shared.CreateDiscordLoginRedirectURL(m.GuildID)

		discordLoginURL := fmt.Sprintf(
			DiscordLoginURLTmpl,
			os.Getenv("DISCORD_CLIENT_ID"),
			url.QueryEscape(redirectURL),
		)

		if err := message_send.SendEmbedEphemeralReplyWithURLBtn(
			s, m, "botについて", msg, discordLoginURL, conf.ColorYellow,
		); err != nil {
			message_send.SendErrMsg(s, errors.NewError("helpコマンドに対する返信を遅れません", err), guildName)
			return
		}
	},
}
