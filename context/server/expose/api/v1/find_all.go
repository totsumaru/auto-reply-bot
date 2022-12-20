package v1

import (
	"context"
	"github.com/techstart35/auto-reply-bot/context/server/gateway/di"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// 全てのサーバーを取得します
func FindAll(ctx context.Context) ([]Res, error) {
	res := make([]Res, 0)

	q, err := di.InitQuery(ctx)
	if err != nil {
		return res, errors.NewError("クエリーを初期化できません", err)
	}

	ms, err := q.FindAll()
	if err != nil {
		return res, errors.NewError("IDでサーバーを取得できません", err)
	}

	for _, m := range ms {
		r, err := CreateRes(m)
		if err != nil {
			return res, errors.NewError("レスポンスを作成できません", err)
		}

		res = append(res, r)
	}

	return res, nil
}
