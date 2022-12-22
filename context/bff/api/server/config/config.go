package config

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord/convert"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/db"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"net/http"
)

// サーバーの設定を更新します
func ServerConfig(e *gin.Engine) {
	e.POST("/server/config", postServerConfig)
}

// リクエストBodyです
type ReqConfig struct {
	AdminRoleID string `json:"admin_role_id"`
	Block       []struct {
		Keyword    []string `json:"keyword"`
		Reply      []string `json:"reply"`
		IsAllMatch bool     `json:"is_all_match"`
		IsRandom   bool     `json:"is_random"`
		IsEmbed    bool     `json:"is_embed"`
	} `json:"block"`
}

// レスポンスです
type ResGetServer struct {
	ID          string              `json:"id"`
	AdminRoleID string              `json:"admin_role_id"`
	Block       []ResGetServerBlock `json:"block"`
}

// ブロックのレスポンスです
type ResGetServerBlock struct {
	Keyword    []string `json:"keyword"`
	Reply      []string `json:"reply"`
	IsAllMatch bool     `json:"is_all_match"`
	IsRandom   bool     `json:"is_random"`
	IsEmbed    bool     `json:"is_embed"`
}

// サーバーの設定を更新します
func postServerConfig(c *gin.Context) {
	session, err := discord.CreateSession()
	if err != nil {
		message_send.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	conf, err := db.NewConf()
	if err != nil {
		message_send.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	database, err := db.NewDB(conf)
	if err != nil {
		message_send.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	tx, err := database.Begin()
	if err != nil {
		message_send.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	ctx := context.WithValue(context.Background(), "tx", tx)

	id := c.Query("id")
	token := c.GetHeader("token")

	// 認証されているユーザーかを検証します
	{
		tmpRes, err := v1.FindByID(ctx, id)
		if err != nil {
			message_send.SendErrMsg(session, err)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		userID, err := convert.TokenToDiscordID(token)
		if err != nil {
			message_send.SendErrMsg(session, err)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		ok, err := discord.HasRole(session, id, userID, tmpRes.AdminRoleID)
		if err != nil {
			message_send.SendErrMsg(session, err)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		if !ok {
			message_send.SendErrMsg(session, errors.NewError("管理者ロールを持っていません"))
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}
	}

	var (
		apiRes = v1.Res{}
	)

	bffErr := (func() error {
		req := &ReqConfig{}
		if err = c.BindJSON(req); err != nil {
			return errors.NewError("リクエストをJSONにバインドできません", err)
		}

		apiReqBlocks := make([]v1.BlockReq, 0)

		for _, rb := range req.Block {
			apiBlockReq := v1.BlockReq{}
			apiBlockReq.Keyword = rb.Keyword
			apiBlockReq.Reply = rb.Reply
			apiBlockReq.IsAllMatch = rb.IsAllMatch
			apiBlockReq.IsRandom = rb.IsRandom
			apiBlockReq.IsEmbed = rb.IsEmbed

			apiReqBlocks = append(apiReqBlocks, apiBlockReq)
		}

		apiRes, err = v1.UpdateConfig(session, ctx, id, req.AdminRoleID, apiReqBlocks)
		if err != nil {
			return errors.NewError("設定を更新できません", err)
		}

		return nil
	})()

	if bffErr != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			message_send.SendErrMsg(
				session,
				errors.NewError("ロールバックに失敗しました。データに不整合が発生している可能性があります。"),
			)
			return
		}

		message_send.SendErrMsg(session, bffErr)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	if txErr := tx.Commit(); txErr != nil {
		message_send.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	res := ResGetServer{}
	res.ID = apiRes.ID
	res.AdminRoleID = apiRes.AdminRoleID
	res.Block = []ResGetServerBlock{}

	for _, v := range apiRes.Block {
		blockRes := ResGetServerBlock{}
		blockRes.Keyword = v.Keyword
		blockRes.Reply = v.Reply
		blockRes.IsAllMatch = v.IsAllMatch
		blockRes.IsRandom = v.IsRandom
		blockRes.IsEmbed = v.IsEmbed

		res.Block = append(res.Block, blockRes)
	}

	c.JSON(http.StatusOK, res)
}
