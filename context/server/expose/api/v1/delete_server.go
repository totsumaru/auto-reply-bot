package v1

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/server/gateway/di"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーを削除します
//
// Devのみ実行可能です。
func DeleteServer(s *discordgo.Session, ctx context.Context, serverID string) error {
	if err := validator.New().Var(serverID, "required"); err != nil {
		return errors.NewError("リクエストが不正です", err)
	}

	a, err := di.InitApp(ctx, s)
	if err != nil {
		return errors.NewError("アプリケーションを初期化できません", err)
	}

	if err = a.DeleteServer(serverID); err != nil {
		return errors.NewError("サーバーを削除できません", err)
	}

	// storeから値を削除します
	if err = removeStore(serverID); err != nil {
		return errors.NewError("storeから値を削除できません", err)
	}

	return nil
}
