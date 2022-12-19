package app

import (
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/user/domain/model/user"
)

// ユーザーを新規作成します
//
// IDを返します。
func (a *App) CreateUser(id string) (string, error) {
	i, err := user.NewID(id)
	if err != nil {
		return "", errors.NewError("idを作成できません", err)
	}

	u, err := user.NewUser(i)
	if err != nil {
		return "", errors.NewError("ユーザーを作成できません", err)
	}

	if err = a.Repo.Create(u); err != nil {
		return "", errors.NewError("ユーザーを登録できません", err)
	}

	return u.ID().String(), nil
}
