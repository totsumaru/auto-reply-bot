package test

import (
	"database/sql"
	"fmt"
	"github.com/techstart35/auto-reply-bot/context/shared/db"
	"testing"
)

var DB *sql.DB

// データベースへ接続します
func OpenDB(t *testing.T) *sql.DB {
	var err error

	c, err := db.NewConf()
	if err != nil {
		t.Fatal(err)
		return nil
	}

	DB, err = sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			c.EnvKeyDBUserName,
			c.EnvKeyDBUserPassword,
			c.EnvKeyDBHost,
			c.EnvKeyDBPort,
			c.EnvKeyDBName,
		),
	)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	return DB
}

// データベースから切断します
func CloseDB(t *testing.T) {
	if DB != nil {
		if err := DB.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
