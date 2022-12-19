package module_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/techstart35/auto-reply-bot/context/shared/map/gen"
	"github.com/techstart35/auto-reply-bot/context/shared/map/seeker"
	"github.com/techstart35/auto-reply-bot/context/shared/now"
	"github.com/techstart35/auto-reply-bot/context/shared/test"
	mysqlUser "github.com/techstart35/auto-reply-bot/context/user/gateway/persistence/mysql/user"
	"os"
	"testing"
	"time"
)

// テスト用の値です
var (
	// 2022年5月5日 AM5:00
	TestNow = time.Date(2022, time.Month(5), 5, 5, 0, 0, 0, time.Local)
	TestID  = "1122"
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
		if _, err = test.DB.Exec("TRUNCATE users"); err != nil {
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

		if _, err = test.DB.Exec("TRUNCATE users"); err != nil {
			t.Fatal(err)
		}

		test.CloseDB(t)
	}
}

// テスト用にユーザーを登録します
func RegisterUser(t *testing.T, m map[string]interface{}) {
	j, err := json.Marshal(m)
	if err != nil {
		t.Fatal("構造体をJSONに変換できなかった", err)
	}

	qs := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s) VALUES (?, ?, ?, ?)`,
		mysqlUser.Conf.TableName,
		mysqlUser.Conf.ColumnName.ID,
		mysqlUser.Conf.ColumnName.Content,
		mysqlUser.Conf.ColumnName.Created,
		mysqlUser.Conf.ColumnName.Updated,
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

// 初期状態のユーザーのモック用データを作成します
func UserInitialMock() map[string]interface{} {
	mock := map[string]interface{}{}

	gen.Gen(mock, []string{"id", "value"}, TestID)

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
