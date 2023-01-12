package nickname

import (
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/check"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/convert"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/initiate"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	discordCtxAPINickname "github.com/techstart35/auto-reply-bot/context/discord/expose/nickname"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"net/http"
)

// ニックネームを変更します
func Nickname(e *gin.Engine) {
	e.POST("/server/nickname", postNickname)
}

// レスポンスです
type Res struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

// ニックネームを変更します
func postNickname(c *gin.Context) {
	session, err := initiate.CreateSession()
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("セッションを作成できません", err), "none")
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	ctx, _, err := shared.CreateDBTx()
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("DBのTxを作成できません", err), "none")
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	id := c.Query("id")
	token := c.GetHeader("Token")
	nickname := c.Query("name")

	// クエリパラメータに指定されたサーバーです
	guildName, err := guild.GetGuildName(session, id)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("ギルド名を取得できません", err), "")
		return
	}

	// 認証されているユーザーかを検証します
	{
		tmpRes, err := v1.FindByID(ctx, id)
		if err != nil {
			message_send.SendErrMsg(session, errors.NewError("IDでサーバーを取得できません", err), guildName)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		userID, err := convert.TokenToDiscordID(token)
		if err != nil {
			message_send.SendErrMsg(session, errors.NewError("トークンをDiscordIDに変換できません", err), guildName)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		ok, err := check.HasRole(session, id, userID, tmpRes.AdminRoleID)
		if err != nil {
			message_send.SendErrMsg(session, errors.NewError("ロールの所有確認に失敗しました", err), guildName)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		guildOwnerID, err := guild.GetGuildOwnerID(session, id)
		if err != nil {
			message_send.SendErrMsg(session, errors.NewError("オーナーIDを取得できません", err), guildName)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		if !(ok || userID == guildOwnerID || userID == conf.TotsumaruDiscordID) {
			message_send.SendErrMsg(session, errors.NewError("管理者ロールを持っていません"), guildName)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}
	}

	if err = discordCtxAPINickname.UpdateNickname(session, id, nickname); err != nil {
		message_send.SendErrMsg(session, errors.NewError("botのニックネームを変更できません", err), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	res := Res{}
	res.ID = id
	res.Nickname = nickname

	c.JSON(http.StatusOK, res)
}
