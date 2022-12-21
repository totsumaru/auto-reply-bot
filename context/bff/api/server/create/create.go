package create

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/discord"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/db"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/shared/map/gen"
	"net/http"
	"os"
)

// サーバーを作成します
func ServerCreate(e *gin.Engine) {
	e.POST("/server/create", postServerCreate)
}

// サーバーを作成します
func postServerCreate(c *gin.Context) {
	var Token = "Bot " + os.Getenv("APP_BOT_TOKEN")

	session, err := discordgo.New(Token)
	session.Token = Token
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

	var (
		apiRes = v1.Res{}
	)

	bffErr := (func() error {
		id := c.Query("id")

		apiRes, err = v1.CreateServer(session, ctx, id)
		if err != nil {
			return errors.NewError("サーバーを作成できません", err)
		}

		return nil
	})()

	if bffErr != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			discord.SendErrMsg(
				session,
				errors.NewError("ロールバックに失敗しました。データに不整合が発生している可能性があります。"),
			)
			return
		}

		discord.SendErrMsg(session, bffErr)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	if txErr := tx.Commit(); txErr != nil {
		discord.SendErrMsg(session, err)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	res := map[string]interface{}{}
	gen.Gen(res, []string{"id"}, apiRes.ID)
	gen.Gen(res, []string{"admin_role_id"}, apiRes.AdminRoleID)

	c.JSON(http.StatusOK, res)
}
