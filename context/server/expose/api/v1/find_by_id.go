package v1

import (
	"context"
	"github.com/techstart35/auto-reply-bot/context/server/gateway/di"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// IDでサーバーを取得します
func FindByID(ctx context.Context, serverID string) (Res, error) {
	res := Res{}

	q, err := di.InitQuery(ctx)
	if err != nil {
		return res, errors.NewError("クエリーを初期化できません", err)
	}

	m, err := q.FindByID(serverID)
	if err != nil {
		return res, errors.NewError("IDでサーバーを取得できません", err)
	}

	res, err = CreateRes(m)
	if err != nil {
		return res, errors.NewError("レスポンスを作成できません", err)
	}

	return res, nil
}
