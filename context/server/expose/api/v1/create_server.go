package v1

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/server/gateway/di"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーを登録します
//
// Devのみ実行可能です。
func CreateServer(s *discordgo.Session, ctx context.Context, serverID string) (Res, error) {
	res := Res{}

	if err := validator.New().Var(serverID, "required"); err != nil {
		return res, errors.NewError("リクエストが不正です", err)
	}

	a, err := di.InitApp(ctx, s)
	if err != nil {
		return res, errors.NewError("アプリケーションを初期化できません", err)
	}

	q, err := di.InitQuery(ctx)
	if err != nil {
		return res, errors.NewError("クエリーを初期化できません", err)
	}

	appResID, err := a.CreateServer(serverID)
	if err != nil {
		return res, errors.NewError("サーバーを作成できません", err)
	}

	m, err := q.FindByID(appResID)
	if err != nil {
		return res, errors.NewError("IDでサーバーを取得できません", err)
	}

	res, err = CreateRes(m)
	if err != nil {
		return res, errors.NewError("レスポンスを作成できません", err)
	}

	// storeに値を保存します
	if err = updateStore(res); err != nil {
		return res, errors.NewError("storeを更新できません", err)
	}

	return res, nil
}
