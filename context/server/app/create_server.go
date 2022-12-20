package app

import (
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーを新規作成します
//
// IDを返します。
func (a *App) CreateServer(serverID string) (string, error) {
	i, err := server.NewID(serverID)
	if err != nil {
		return "", errors.NewError("idを作成できません", err)
	}

	u, err := server.NewServer(i)
	if err != nil {
		return "", errors.NewError("サーバーを作成できません", err)
	}

	if err = a.Repo.Create(u); err != nil {
		return "", errors.NewError("サーバーを登録できません", err)
	}

	return u.ID().String(), nil
}
