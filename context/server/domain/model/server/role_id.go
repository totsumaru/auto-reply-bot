package server

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// ロールIDです
type RoleID struct {
	value string
}

// ロールIDを作成します
func NewRoleID(v string) (RoleID, error) {
	r := RoleID{}
	r.value = v

	if err := r.validate(); err != nil {
		return RoleID{}, errors.NewError("検証に失敗しました", err)
	}

	return r, nil
}

// ロールIDの値を文字列で取得します
func (r RoleID) String() string {
	return r.value
}

// ロールIDを比較します
func (r RoleID) Equal(rr RoleID) bool {
	return r.value == rr.value
}

// RoleIDの値が設定されているか判別します
func (r RoleID) IsEmpty() bool {
	return r.value == ""
}

// 検証します
func (r RoleID) validate() error {
	if err := validator.New().Var(r.value, "required"); err != nil {
		return errors.NewError("値が空です", err)
	}

	return nil
}

// 構造体をJSONに変換します
func (r RoleID) MarshalJSON() ([]byte, error) {
	j := struct {
		Value string `json:"value"`
	}{
		Value: r.value,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	return b, nil
}

// JSONを構造体に変換します
func (r *RoleID) UnmarshalJSON(b []byte) error {
	j := &struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	r.value = j.Value

	return nil
}
