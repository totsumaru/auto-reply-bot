package block

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// キーワードの最大文字数です
	KeywordMaxLen = 20
)

// キーワードです
type Keyword struct {
	value string
}

// キーワードを作成します
func NewKeyword(v string) (Keyword, error) {
	k := Keyword{}
	k.value = v

	if err := k.validate(); err != nil {
		return Keyword{}, errors.NewError("検証に失敗しました", err)
	}

	return k, nil
}

// キーワードの値を文字列で取得します
func (k Keyword) String() string {
	return k.value
}

// キーワードを比較します
func (k Keyword) Equal(kk Keyword) bool {
	return k.value == kk.value
}

// キーワードの値が設定されているか判別します
func (k Keyword) IsEmpty() bool {
	return k.value == ""
}

// 検証します
func (k Keyword) validate() error {
	if err := validator.New().Var(k.value, "required"); err != nil {
		return errors.NewError("値が空です", err)
	}

	if len([]rune(k.value)) > KeywordMaxLen {
		return errors.NewError("キーワードの文字数が上限を超えています")
	}

	return nil
}

// 構造体をJSONに変換します
func (k Keyword) MarshalJSON() ([]byte, error) {
	j := struct {
		Value string `json:"value"`
	}{
		Value: k.value,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	return b, nil
}

// JSONを構造体に変換します
func (k *Keyword) UnmarshalJSON(b []byte) error {
	j := &struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	k.value = j.Value

	return nil
}
