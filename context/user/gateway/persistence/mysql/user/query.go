package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/user/gateway/persistence/mysql"
	"log"
)

// クエリです
type Query struct {
	mysql.Infra
	ctx context.Context
}

// クエリを生成します
func NewQuery(ctx context.Context) (*Query, error) {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return nil, errors.NewError("型アサーションに失敗しました")
	}

	if tx == nil {
		return nil, errors.NewError("txがnilです")
	}

	q := &Query{}
	q.Tx = tx
	q.ctx = ctx

	return q, nil
}

// IDからユーザーを取得します
func (q *Query) FindByID(id string) (map[string]interface{}, error) {
	qs := fmt.Sprintf(
		`SELECT %s, %s, %s, %s FROM %s WHERE %s=?`,
		Conf.ColumnName.ID,
		Conf.ColumnName.Content,
		Conf.ColumnName.Created,
		Conf.ColumnName.Updated,
		Conf.TableName,
		Conf.ColumnName.ID,
	)

	s, err := q.Tx.Prepare(qs)
	if err != nil {
		return nil, errors.NewError("ステートメントの作成に失敗しました", err)
	}
	defer q.CloseStmt(s)

	rows, err := s.Query(id)
	if err != nil {
		return nil, errors.NewError("クエリが失敗しました", err)
	}

	res := make([]*Row, 0)
	for rows.Next() {
		ro := &Row{}
		if err = rows.Scan(&ro.ID, &ro.Content, &ro.Created, &ro.Updated); err != nil {
			return nil, errors.NewError("ユーザーが見つかりません", err)
		}

		res = append(res, ro)
	}

	if len(res) == 0 {
		return nil, errors.NotFoundErr
	}

	if len(res) > 1 {
		return nil, errors.NewError("ユーザーが複数見つかりました", err)
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(res[0].Content), &m); err != nil {
		return nil, errors.NewError("mapに変換できません", err)
	}

	return m, nil
}

// 全てのユーザーを取得します
func (q *Query) FindAll() ([]map[string]interface{}, error) {
	qs := fmt.Sprintf(
		`SELECT %s, %s, %s, %s FROM %s`,
		Conf.ColumnName.ID,
		Conf.ColumnName.Content,
		Conf.ColumnName.Created,
		Conf.ColumnName.Updated,
		Conf.TableName,
	)

	s, err := q.Tx.Prepare(qs)
	if err != nil {
		return nil, errors.NewError("ステートメントの作成に失敗しました", err)
	}
	defer q.CloseStmt(s)

	rows, err := s.Query()
	if err != nil {
		return nil, errors.NewError("クエリが失敗しました", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Fatal("rowsを閉じれません", err)
		}
	}()

	m := make([]map[string]interface{}, 0)
	for rows.Next() {
		ro := &Row{}
		if err = rows.Scan(&ro.ID, &ro.Content, &ro.Created, &ro.Updated); err != nil {
			return nil, errors.NewError("行のデータを取得できません", err)
		}

		var data map[string]interface{}
		if err = json.Unmarshal([]byte(ro.Content), &data); err != nil {
			return nil, errors.NewError("mapに変換できません", err)
		}

		m = append(m, data)
	}

	return m, nil
}
