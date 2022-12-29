package shared

import (
	"context"
	"database/sql"
	"github.com/techstart35/auto-reply-bot/context/shared/db"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// DBのトランザクションを作成します
func CreateDBTx() (context.Context, *sql.Tx, error) {
	ctx := context.Background()
	tx := &sql.Tx{}

	conf, err := db.NewConf()
	if err != nil {
		return ctx, tx, errors.NewError("DBのConfを作成できません", err)
	}

	database, err := db.NewDB(conf)
	if err != nil {
		return ctx, tx, errors.NewError("DBを作成できません", err)
	}

	tx, err = database.Begin()
	if err != nil {
		return ctx, tx, errors.NewError("DB開始のTxを作成できません", err)
	}

	ctx = context.WithValue(context.Background(), "tx", tx)

	return ctx, tx, nil
}
