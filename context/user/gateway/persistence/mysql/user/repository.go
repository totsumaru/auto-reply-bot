package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/user/domain/model/user"
	"github.com/techstart35/auto-reply-bot/context/user/gateway/persistence/mysql"
	"log"
	"time"
)

// リポジトリです
type Repository struct {
	mysql.Infra
}

// リポジトリを生成します
func NewRepository(ctx context.Context) (*Repository, error) {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return nil, errors.NewError("型アサーションに失敗しました")
	}

	if tx == nil {
		return nil, errors.NewError("txがnilです")
	}

	r := new(Repository)
	r.Tx = tx

	return r, nil
}

// 新規作成します
func (r *Repository) Create(u *user.User) error {
	qs := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s) VALUES (?, ?, ?, ?)`,
		Conf.TableName,
		Conf.ColumnName.ID,
		Conf.ColumnName.Content,
		Conf.ColumnName.Created,
		Conf.ColumnName.Updated,
	)

	d, err := r.Tx.Prepare(qs)
	if err != nil {
		return errors.NewError("ステートメントの作成に失敗しました", err)
	}
	defer func() { r.CloseStmt(d) }()

	b, err := json.Marshal(u)
	if err != nil {
		return errors.NewError("請求構造体をJSONに変換できません", err)
	}

	_, err = d.Exec(
		u.ID().String(),
		string(b),
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		return errors.NewError("クエリーが失敗しました", err)
	}

	return nil
}

// 更新します
func (r *Repository) Update(u *user.User) error {
	qs := fmt.Sprintf(
		`UPDATE %s SET %s = ?, %s = ?, %s = ? WHERE %s =?`,
		Conf.TableName,
		Conf.ColumnName.Content,
		Conf.ColumnName.Created,
		Conf.ColumnName.Updated,
		Conf.ColumnName.ID,
	)

	d, err := r.Tx.Prepare(qs)
	if err != nil {
		return errors.NewError("ステートメントの作成に失敗しました", err)
	}
	defer func() { r.CloseStmt(d) }()

	b, err := json.Marshal(u)
	if err != nil {
		return errors.NewError("請求構造体をJSONに変換できません", err)
	}

	_, err = d.Exec(
		string(b),
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"),
		u.ID().String(),
	)
	if err != nil {
		return errors.NewError("クエリーが失敗しました", err)
	}

	return nil
}

// 削除します
func (r *Repository) Delete(id user.ID) error {
	qs := fmt.Sprintf(
		`DELETE FROM %s WHERE %s =?`,
		Conf.TableName,
		Conf.ColumnName.ID,
	)

	d, err := r.Tx.Prepare(qs)
	if err != nil {
		return errors.NewError("ステートメントの作成に失敗しました")
	}
	defer func() { r.CloseStmt(d) }()

	_, err = d.Exec(id.String())
	if err != nil {
		return errors.NewError("クエリーが失敗しました")
	}

	return nil
}

// DiscordIDでユーザーを取得します
func (r *Repository) FindByID(id user.ID) (*user.User, error) {
	qs := fmt.Sprintf(
		`SELECT %s, %s, %s, %s FROM %s WHERE %s = ? FOR UPDATE`,
		Conf.ColumnName.ID,
		Conf.ColumnName.Content,
		Conf.ColumnName.Created,
		Conf.ColumnName.Updated,
		Conf.TableName,
		Conf.ColumnName.ID,
	)

	s, err := r.Tx.Prepare(qs)
	if err != nil {
		return nil, errors.NewError("ステートメントの作成に失敗しました", err)
	}
	defer func() { r.CloseStmt(s) }()

	rows, err := s.Query(id.String())
	if err != nil {
		return nil, errors.NewError("クエリーが失敗しました", err)
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
		return nil, errors.NewError("ユーザーが複数見つかりました")
	}

	uu := &user.User{}
	if err = json.Unmarshal([]byte(res[0].Content), uu); err != nil {
		return nil, errors.NewError("マッピングできません", err)
	}

	return uu, nil
}

// 全てのコレクションを取得します
//
// ID:User のmapを返します
func (r *Repository) FindAll() (map[string]*user.User, error) {
	res := map[string]*user.User{}

	qs := fmt.Sprintf(
		`SELECT %s, %s, %s, %s FROM %s FOR UPDATE`,
		Conf.ColumnName.ID,
		Conf.ColumnName.Content,
		Conf.ColumnName.Created,
		Conf.ColumnName.Updated,
		Conf.TableName,
	)

	s, err := r.Tx.Prepare(qs)
	if err != nil {
		return res, errors.NewError("ステートメントの作成に失敗しました", err)
	}
	defer r.CloseStmt(s)

	rows, err := s.Query()
	if err != nil {
		return res, errors.NewError("クエリが失敗しました", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Fatal("rowsを閉じれません", err)
		}
	}()

	for rows.Next() {
		ro := &Row{}
		if err = rows.Scan(&ro.ID, &ro.Content, &ro.Created, &ro.Updated); err != nil {
			return res, errors.NewError("行のデータを取得できません", err)
		}

		data := &user.User{}
		if err = json.Unmarshal([]byte(ro.Content), &data); err != nil {
			return res, errors.NewError("mapに変換できません", err)
		}

		res[data.ID().String()] = data
	}

	return res, nil
}
