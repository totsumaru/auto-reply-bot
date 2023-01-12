package api

import (
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/api/server"
	"github.com/techstart35/auto-reply-bot/context/bff/api/server/config"
	nickname "github.com/techstart35/auto-reply-bot/context/bff/api/server/nockname"
)

// ルートを設定します
func RegisterRouter(e *gin.Engine) {
	Route(e)
	server.Server(e)
	config.ServerConfig(e)
	nickname.Nickname(e)
}

// ルートです
//
// Note: この関数は削除しても問題ありません
func Route(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.Header("hello", "world")
		c.JSON(200, gin.H{
			"message": "hello",
		})
	})
}
