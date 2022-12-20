//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/google/wire"
	"github.com/techstart35/auto-reply-bot/context/server/app"
	"github.com/techstart35/auto-reply-bot/context/server/gateway/persistence/mysql/server"
)

// アプリケーションサービスの作成
func InitApp(ctx context.Context, session *discordgo.Session) (*app.App, error) {
	wire.Build(
		app.NewApp,
		server.NewRepository,
		wire.Bind(new(app.Repository), new(*server.Repository)),
	)
	return nil, nil
}

// クエリの作成
func InitQuery(ctx context.Context) (*server.Query, error) {
	wire.Build(
		server.NewQuery,
	)

	return nil, nil
}
