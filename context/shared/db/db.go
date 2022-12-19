package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	// MySQLを利用するときに必要（削除するとコンパイルできない）
	_ "github.com/go-sql-driver/mysql"
)

var (
	mu   sync.Mutex
	db   *sql.DB
	open = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.Open(driverName, dataSourceName)
	}
)

// データベースへ接続してデータベースのハンドラを取得します
//
// 変数dbによってコネクションを再利用するするため、この関数はsync.Mutexによって排他制御されます。
//
// 返されたハンドラは変数dbに保存され、既に存在する場合は変数dbからハンドラを返します。
// これは、sql.Open関数が呼び出されるとConnection Pool（持続的接続）が有効なハンドラを生成し、
// コネクションを再利用するデザインになっているからです。
//
// この関数のコール時には、データベースの接続を作成せずに、引数を検証するだけの場合があります。
//
// ref: https://golang.org/pkg/database/sql/#Open
func NewDB(c *Config) (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	// すでにハンドラが存在している場合はそのまま返す
	if db != nil {
		return db, nil
	}

	var err error

	db, err = open(
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
		return nil, errors.NewError("Open関数の呼び出しに失敗しました", err)
	}

	return db, nil
}
