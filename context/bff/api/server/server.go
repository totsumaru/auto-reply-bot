package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/db"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"net/http"
)

// サーバーを取得します
func Server(e *gin.Engine) {
	e.GET("/server", getServer)
}

// レスポンスです
type GetServerRes struct {
	ID          string              `json:"id"`
	AdminRoleID string              `json:"admin_role_id"`
	Block       []GetServerBlockRes `json:"block"`
}

// ブロックのレスポンスです
type GetServerBlockRes struct {
	Keyword    []string `json:"keyword"`
	Reply      []string `json:"reply"`
	IsAllMatch bool     `json:"is_all_match"`
	IsRandom   bool     `json:"is_random"`
	IsEmbed    bool     `json:"is_embed"`
}

// サーバーを取得します
func getServer(c *gin.Context) {
	session, err := discord.CreateSession()
	if err != nil {
		discord.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	conf, err := db.NewConf()
	if err != nil {
		discord.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	database, err := db.NewDB(conf)
	if err != nil {
		discord.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	tx, err := database.Begin()
	if err != nil {
		discord.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	ctx := context.WithValue(context.Background(), "tx", tx)

	id := c.Query("id")
	code := c.Query("code")

	var (
		apiRes = v1.Res{}
	)

	bffErr := (func() error {
		apiRes, err = v1.FindByID(ctx, id)
		if err != nil {
			return errors.NewError("IDでサーバーを取得できません", err)
		}

		// 認証されているユーザーかを検証します
		{
			token, err := discord.CodeToToken(code, id)
			if err != nil {
				return errors.NewError("codeからtokenに変換できません", err)
			}

			userID, err := discord.TokenToDiscordID(token)
			if err != nil {
				return errors.NewError("tokenからDiscordIDに変換できません", err)
			}

			ok, err := discord.HasRole(session, id, userID, apiRes.AdminRoleID)
			if err != nil {
				return errors.NewError("ロールの確認ができません", err)
			}

			if !ok {
				return errors.AuthErr
			}
		}

		return nil
	})()

	if bffErr != nil {
		// 認証されていない場合は、ここで終了します
		if bffErr == errors.AuthErr {
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		// ロールバックを実行します
		txErr := tx.Rollback()
		if txErr != nil {
			discord.SendErrMsg(
				session,
				errors.NewError("ロールバックに失敗しました。データに不整合が発生している可能性があります。"),
			)
			return
		}

		// その他のエラーはここでエラーを返します
		discord.SendErrMsg(session, bffErr)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	if txErr := tx.Commit(); txErr != nil {
		discord.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	res := GetServerRes{}
	res.ID = apiRes.ID
	res.AdminRoleID = apiRes.AdminRoleID

	for _, v := range apiRes.Block {
		blockRes := GetServerBlockRes{}
		blockRes.Keyword = v.Keyword
		blockRes.Reply = v.Reply
		blockRes.IsAllMatch = v.IsAllMatch
		blockRes.IsRandom = v.IsRandom
		blockRes.IsEmbed = v.IsEmbed

		res.Block = append(res.Block, blockRes)
	}

	c.JSON(http.StatusOK, res)
}
