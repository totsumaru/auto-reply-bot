package user

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// ユーザーです
type User struct {
	id   ID
	name Name
}

// ユーザーを作成します
func NewUser(d ID) (*User, error) {
	u := &User{}
	u.id = d

	if err := u.validate(); err != nil {
		return u, errors.NewError("検証に失敗しました", err)
	}

	return u, nil
}

// -------------------
// setter
// -------------------

// 名前を更新します
func (u *User) UpdateName(n Name) error {
	u.name = n

	if err := u.validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}

// -------------------
// getter
// -------------------

// IDを取得します
func (u *User) ID() ID {
	return u.id
}

// 名前を取得します
func (u *User) Name() Name {
	return u.name
}

// -------------------
// validation
// -------------------

// 検証します
func (u *User) validate() error {
	return nil
}

// -------------------
// JSON
// -------------------

// 構造体をJSONに変換します
func (u *User) MarshalJSON() ([]byte, error) {
	j := struct {
		ID   ID   `json:"id"`
		Name Name `json:"name"`
	}{
		ID:   u.id,
		Name: u.name,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません")
	}

	return b, nil
}

// JSONを構造体を変換します
func (u *User) UnmarshalJSON(b []byte) error {
	j := &struct {
		ID   ID   `json:"id"`
		Name Name `json:"name"`
	}{}

	if err := json.Unmarshal(b, &j); err != nil {
		return errors.NewError("JSONを構造体に変換できません")
	}

	u.id = j.ID
	u.name = j.Name

	return nil
}
