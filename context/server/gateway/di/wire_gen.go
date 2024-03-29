// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/server/app"
	"github.com/techstart35/auto-reply-bot/context/server/gateway/persistence/mysql/server"
)

// Injectors from wire.go:

// アプリケーションサービスの作成
func InitApp(ctx context.Context, session *discordgo.Session) (*app.App, error) {
	repository, err := server.NewRepository(ctx)
	if err != nil {
		return nil, err
	}
	appApp := app.NewApp(repository, session)
	return appApp, nil
}

// クエリの作成
func InitQuery(ctx context.Context) (*server.Query, error) {
	query, err := server.NewQuery(ctx)
	if err != nil {
		return nil, err
	}
	return query, nil
}
