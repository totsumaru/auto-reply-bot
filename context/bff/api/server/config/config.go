package config

import (
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/initiate"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
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
	session, err := initiate.CreateSession()
	if err != nil {
		message_send.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	ctx, tx, err := shared.CreateDBTx()
	if err != nil {
		message_send.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	id := c.Query("id")
	// TODO: コメントアウト解除（FEの実装のため一時的にコメントアウト）
	//token := c.GetHeader("token")

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
