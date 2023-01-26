package nickname

import (
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/initiate"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	discordCtxAPINickname "github.com/techstart35/auto-reply-bot/context/discord/expose/nickname"
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

	// リクエストを検証します
	if id == "" || token == "" {
		c.JSON(http.StatusBadRequest, "リクエストが不正です")
		return
	}

	// クエリパラメータに指定されたサーバーです
	guildName, err := guild.GetGuildName(session, id)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("ギルド名を取得できません", err), "")
		return
	}

	// 認証されているユーザーかを検証します
	ok, err := shared.IsAuthorizedUser(ctx, session, id, token)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("認証されているかの確認に失敗しました", err), guildName)
		c.JSON(http.StatusUnauthorized, "認証されていません")
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, "認証されていません")
		return
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
