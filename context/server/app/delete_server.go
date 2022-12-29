package app

import (
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// サーバーを削除します
//
// IDを返します。
func (a *App) DeleteServer(serverID string) error {
	i, err := server.NewID(serverID)
	if err != nil {
		return errors.NewError("idを作成できません", err)
	}

	u, err := server.NewServer(i)
	if err != nil {
		return errors.NewError("サーバーを作成できません", err)
	}

	// 存在確認をします
	_, err = a.Repo.FindByID(i)
	if err != nil {
		return errors.NewError("削除予定のサーバーを取得できません", err)
	}

	if err = a.Repo.Delete(u.ID()); err != nil {
		return errors.NewError("サーバーを削除できません", err)
	}

	return nil
}
