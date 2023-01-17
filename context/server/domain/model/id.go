package model

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// IDです
type ID struct {
	value string
}

// IDを作成します
func NewID(v string) (ID, error) {
	d := ID{}
	d.value = v

	if err := d.validate(); err != nil {
		return ID{}, errors.NewError("検証に失敗しました", err)
	}

	return d, nil
}

// IDの値を文字列で取得します
func (i ID) String() string {
	return i.value
}

// IDを比較します
func (i ID) Equal(ii ID) bool {
	return i.value == ii.value
}

// IDの値が設定されているか判別します
func (i ID) IsEmpty() bool {
	return i.value == ""
}

// 検証します
func (i ID) validate() error {
	if err := validator.New().Var(i.value, "required"); err != nil {
		return errors.NewError("値が空です", err)
	}

	return nil
}

// 構造体をJSONに変換します
func (i ID) MarshalJSON() ([]byte, error) {
	j := struct {
		Value string `json:"value"`
	}{
		Value: i.value,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	return b, nil
}

// JSONを構造体に変換します
func (i *ID) UnmarshalJSON(b []byte) error {
	j := &struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	i.value = j.Value

	return nil
}
