package v1

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/user/gateway/di"
)

// ユーザーを登録します
func CreateUser(s *discordgo.Session, ctx context.Context, id string) (Res, error) {
	res := Res{}

	if err := validator.New().Var(id, "required"); err != nil {
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

	appResID, err := a.CreateUser(id)
	if err != nil {
		return res, errors.NewError("ユーザーを作成できません", err)
	}

	m, err := q.FindByID(appResID)
	if err != nil {
		return res, errors.NewError("IDでユーザーを取得できません", err)
	}

	res, err = CreateRes(m)
	if err != nil {
		return res, errors.NewError("レスポンスを作成できません", err)
	}

	return res, nil
}
