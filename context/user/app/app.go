package app

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/user/domain/model/user"
)

// リポジトリのインターフェイスです
type Repository interface {
	Create(u *user.User) error
	Update(u *user.User) error
	Delete(discordID user.ID) error
	FindByID(id user.ID) (*user.User, error)
	FindAll() (map[string]*user.User, error)
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
