package app

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server"
)

// リポジトリのインターフェイスです
type Repository interface {
	Create(u *server.Server) error
	Update(u *server.Server) error
	Delete(discordID model.ID) error
	FindByID(id model.ID) (*server.Server, error)
	FindAll() (map[string]*server.Server, error)
}

// アプリケーションです
type App struct {
	Repo    Repository
	Session *discordgo.Session
}

// アプリケーションを作成します
func NewApp(repo Repository, s *discordgo.Session) *App {
	app := &App{}
	app.Repo = repo
	app.Session = s

	return app
}
