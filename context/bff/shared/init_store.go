package shared

import (
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// storeの値を初期化します
//
// DBの最新情報をstoreに保存します。
func InitStore() error {
	// 全ての値を取得します
	ctx, _, err := CreateDBTx()
	if err != nil {
		return errors.NewError("DBのTxを作成できません", err)
	}

	res, err := v1.FindAll(ctx)
	if err != nil {
		return errors.NewError("全ての値を取得できません", err)
	}

	// storeに保存します
	v1.InitStore(res)

	return nil
}
