package server

import (
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/initiate"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"net/http"
)

// サーバーを取得します
func Server(e *gin.Engine) {
	e.GET("/server", getServer)
}

// レスポンスです
type Res struct {
	ID          string     `json:"id"`
	AdminRoleID string     `json:"admin_role_id"`
	Block       []resBlock `json:"block"`
	// 以下はComputedです
	ServerName string    `json:"server_name"`
	AvatarURL  string    `json:"avatar_url"`
	Role       []resRole `json:"role"`
}

// ブロックのレスポンスです
type resBlock struct {
	Name       string   `json:"name"`
	Keyword    []string `json:"keyword"`
	Reply      []string `json:"reply"`
	IsAllMatch bool     `json:"is_all_match"`
	IsRandom   bool     `json:"is_random"`
	IsEmbed    bool     `json:"is_embed"`
}

// ロールのレスポンスです
type resRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// サーバーを取得します
func getServer(c *gin.Context) {
	session, err := initiate.CreateSession()
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("セッションを作成できません", err))
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	ctx, tx, err := shared.CreateDBTx()
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("DBのTxを作成できません", err))
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	id := c.Query("id")
	// TODO: コメントアウト解除（FEの実装のため一時的にコメントアウト）
	//code := c.Query("code")

	var (
		token string
	)

	// TODO: コメントアウト解除（FEの実装のため一時的にコメントアウト）
	// 認証されているユーザーかを検証します
	//{
	//	tmpRes, err := v1.FindByID(ctx, id)
	//	if err != nil {
	//		message_send.SendErrMsg(session, err)
	//		c.JSON(http.StatusUnauthorized, "認証されていません")
	//		return
	//	}
	//
	//	token, err = convert.CodeToToken(code, id)
	//	if err != nil {
	//		// codeの不正に関してはエラー通知しません
	//		c.JSON(http.StatusUnauthorized, "認証されていません")
	//		return
	//	}
	//
	//	fmt.Println("token: ", token)
	//
	//	userID, err := convert.TokenToDiscordID(token)
	//	if err != nil {
	//		message_send.SendErrMsg(session, err)
	//		c.JSON(http.StatusUnauthorized, "認証されていません")
	//		return
	//	}
	//
	//	ok, err := check.HasRole(session, id, userID, tmpRes.AdminRoleID)
	//	if err != nil {
	//		message_send.SendErrMsg(session, err)
	//		c.JSON(http.StatusUnauthorized, "認証されていません")
	//		return
	//	}
	//
	//	guild, err := session.Guild(id)
	//	if err != nil {
	//		message_send.SendErrMsg(session, err)
	//		c.JSON(http.StatusUnauthorized, "認証されていません")
	//		return
	//	}
	//
	//	if !(ok || userID == guild.OwnerID || userID == conf.TotsumaruDiscordID) {
	//		message_send.SendErrMsg(session, errors.NewError("管理者ロールを持っていません"))
	//		c.JSON(http.StatusUnauthorized, "認証されていません")
	//		return
	//	}
	//}

	var (
		apiRes = v1.Res{}
	)

	bffErr := (func() error {
		apiRes, err = v1.FindByID(ctx, id)
		if err != nil {
			return errors.NewError("IDでサーバーを取得できません", err)
		}

		return nil
	})()

	if bffErr != nil {
		// ロールバックを実行します
		txErr := tx.Rollback()
		if txErr != nil {
			message_send.SendErrMsg(
				session,
				errors.NewError("ロールバックに失敗しました。データに不整合が発生している可能性があります。", txErr),
			)
			return
		}

		message_send.SendErrMsg(session, errors.NewError("バックエンドの処理でエラーが発生しました", bffErr))
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	if txErr := tx.Commit(); txErr != nil {
		message_send.SendErrMsg(session, errors.NewError("Commitに失敗しました", txErr))
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	guildName, err := guild.GetGuildName(session, apiRes.ID)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("ギルド名を取得できません", err))
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	avatarURL, err := guild.GetAvatarURL(session, apiRes.ID)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("アバターURLを取得できません", err))
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	allRoles, err := guild.GetAllRoles(session, apiRes.ID)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("全てのロールを取得できません", err))
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	res := Res{}
	res.ID = apiRes.ID
	res.AdminRoleID = apiRes.AdminRoleID
	res.Block = []resBlock{}
	res.ServerName = guildName
	res.AvatarURL = avatarURL
	res.Role = []resRole{}

	// レスポンスにブロックを追加します
	for _, v := range apiRes.Block {
		blockRes := resBlock{}
		blockRes.Name = v.Name
		blockRes.Keyword = v.Keyword
		blockRes.Reply = v.Reply
		blockRes.IsAllMatch = v.IsAllMatch
		blockRes.IsRandom = v.IsRandom
		blockRes.IsEmbed = v.IsEmbed

		res.Block = append(res.Block, blockRes)
	}

	// レスポンスにロールを追加します
	for roleID, roleName := range allRoles {
		tmpRole := resRole{
			ID:   roleID,
			Name: roleName,
		}

		res.Role = append(res.Role, tmpRole)
	}

	c.Header("token", token)
	c.JSON(http.StatusOK, res)
}
