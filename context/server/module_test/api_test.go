package module_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	mysqlServer "github.com/techstart35/auto-reply-bot/context/server/gateway/persistence/mysql/server"
	"github.com/techstart35/auto-reply-bot/context/shared/map/gen"
	"github.com/techstart35/auto-reply-bot/context/shared/map/seeker"
	"github.com/techstart35/auto-reply-bot/context/shared/now"
	"github.com/techstart35/auto-reply-bot/context/shared/test"
	"os"
	"testing"
	"time"
)

// テスト用の値です
var (
	// 2022年5月5日 AM5:00
	TestNow = time.Date(2022, time.Month(5), 5, 5, 0, 0, 0, time.Local)
	TestID  = "984614055681613864" // `TestServer`のDiscordのサーバーID
)

// セットアップ用関数
func setup(t *testing.T) (context.Context, func()) {
	// serve.shでgoをコンテナで起動した場合のために設定しています。
	// goをコンテナで起動した場合、`DB_HOST`は`host.docker.internal`で設定されます。
	// そうするとテストでlocalhostが接続できなくなってしまうので、ここで強制的に環境変数を書き換えます。
	if err := os.Setenv("DB_HOST", "localhost"); err != nil {
		t.Fatal(err)
	}

	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	// Setup
	{
		test.OpenDB(t)

		// DBのテーブルを初期化する
		if _, err = test.DB.Exec("TRUNCATE servers"); err != nil {
			t.Fatal(err)
		}

		tx, err = test.DB.Begin()
		if err != nil {
			t.Fatal("エラーが返された")
		}
		ctx = context.WithValue(context.Background(), "tx", tx)
	}

	// Teardown
	return ctx, func() {
		if err = tx.Commit(); err != nil {
			t.Fatal("エラーが返された")
		}

		if _, err = test.DB.Exec("TRUNCATE servers"); err != nil {
			t.Fatal(err)
		}

		test.CloseDB(t)
	}
}

// テスト用にサーバーを登録します
func RegisterServer(t *testing.T, m map[string]interface{}) {
	j, err := json.Marshal(m)
	if err != nil {
		t.Fatal("構造体をJSONに変換できなかった", err)
	}

	qs := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s) VALUES (?, ?, ?, ?)`,
		mysqlServer.Conf.TableName,
		mysqlServer.Conf.ColumnName.ID,
		mysqlServer.Conf.ColumnName.Content,
		mysqlServer.Conf.ColumnName.Created,
		mysqlServer.Conf.ColumnName.Updated,
	)

	s, err := test.DB.Prepare(qs)
	if err != nil {
		t.Fatal("エラーが返された", err)
	}

	_, err = s.Exec(
		seeker.Str(m, []string{"id", "value"}),
		string(j),
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		t.Fatal("エラーが返された", err)
	}
}

// 初期状態のサーバーのモック用データを作成します
func ServerInitialMock(id string) map[string]interface{} {
	mock := map[string]interface{}{}

	gen.Gen(mock, []string{"id", "value"}, id)

	return mock
}

// 全ての外部APIコールをモックします
func MockExternal(
	mockResDateTime time.Time,
) {
	// 現在日時を取得する関数をモックします
	now.Now = func() time.Time {
		return mockResDateTime
	}
}
