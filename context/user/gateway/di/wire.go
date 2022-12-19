//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/google/wire"
	"github.com/techstart35/auto-reply-bot/context/user/app"
	"github.com/techstart35/auto-reply-bot/context/user/gateway/persistence/mysql/user"
)

// アプリケーションサービスの作成
func InitApp(ctx context.Context, session *discordgo.Session) (*app.App, error) {
	wire.Build(
		app.NewApp,
		user.NewRepository,
		wire.Bind(new(app.Repository), new(*user.Repository)),
	)
	return nil, nil
}

// クエリの作成
func InitQuery(ctx context.Context) (*user.Query, error) {
	wire.Build(
		user.NewQuery,
	)

	return nil, nil
}
