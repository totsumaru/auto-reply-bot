package block

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// ブロック名の最大文字数です
	NameMaxLen = 20
)

// 名前です
type Name struct {
	value string
}

// 名前を作成します
func NewName(v string) (Name, error) {
	n := Name{}
	n.value = v

	if err := n.validate(); err != nil {
		return Name{}, errors.NewError("検証に失敗しました", err)
	}

	return n, nil
}

// 名前の値を文字列で取得します
func (n Name) String() string {
	return n.value
}

// 名前を比較します
func (n Name) Equal(nn Name) bool {
	return n.value == nn.value
}

// 名前の値が設定されているか判別します
func (n Name) IsEmpty() bool {
	return n.value == ""
}

// 検証します
func (n Name) validate() error {
	if err := validator.New().Var(n.value, "required"); err != nil {
		return errors.NewError("値が空です", err)
	}

	if len(n.value) > NameMaxLen {
		return errors.NewError("名前の文字数が上限を超えています")
	}

	return nil
}

// 構造体をJSONに変換します
func (n Name) MarshalJSON() ([]byte, error) {
	j := struct {
		Value string `json:"value"`
	}{
		Value: n.value,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	return b, nil
}

// JSONを構造体に変換します
func (n *Name) UnmarshalJSON(b []byte) error {
	j := &struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	n.value = j.Value

	return nil
}
