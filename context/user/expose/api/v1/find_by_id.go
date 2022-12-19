package v1

import (
	"context"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/user/gateway/di"
)

// IDでユーザーを取得します
func FindByID(ctx context.Context, id string) (Res, error) {
	res := Res{}

	q, err := di.InitQuery(ctx)
	if err != nil {
		return res, errors.NewError("クエリーを初期化できません", err)
	}

	m, err := q.FindByID(id)
	if err != nil {
		return res, errors.NewError("IDでユーザーを取得できません", err)
	}

	res, err = CreateRes(m)
	if err != nil {
		return res, errors.NewError("レスポンスを作成できません", err)
	}

	return res, nil
}
