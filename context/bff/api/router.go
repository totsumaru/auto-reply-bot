package api

import (
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/api/user"
	"github.com/techstart35/auto-reply-bot/context/bff/api/user/create"
)

// ルートを設定します
func RegisterRouter(e *gin.Engine) {
	Route(e)
	user.User(e)
	create.UserCreate(e)
}

// ルートです
//
// Note: この関数は削除しても問題ありません
func Route(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello",
		})
	})
}
