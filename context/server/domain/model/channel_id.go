package model

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// チャンネルIDです
type ChannelID struct {
	value string
}

// チャンネルIDを作成します
func NewChannelID(v string) (ChannelID, error) {
	c := ChannelID{}
	c.value = v

	if err := c.validate(); err != nil {
		return ChannelID{}, errors.NewError("検証に失敗しました", err)
	}

	return c, nil
}

// チャンネルIDの値を文字列で取得します
func (c ChannelID) String() string {
	return c.value
}

// チャンネルIDを比較します
func (c ChannelID) Equal(cc ChannelID) bool {
	return c.value == cc.value
}

// チャンネルIDの値が設定されているか判別します
func (c ChannelID) IsEmpty() bool {
	return c.value == ""
}

// 検証します
func (c ChannelID) validate() error {
	if err := validator.New().Var(c.value, "required"); err != nil {
		return errors.NewError("値が空です", err)
	}

	return nil
}

// 構造体をJSONに変換します
func (c ChannelID) MarshalJSON() ([]byte, error) {
	j := struct {
		Value string `json:"value"`
	}{
		Value: c.value,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	return b, nil
}

// JSONを構造体に変換します
func (c *ChannelID) UnmarshalJSON(b []byte) error {
	j := &struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	c.value = j.Value

	return nil
}
