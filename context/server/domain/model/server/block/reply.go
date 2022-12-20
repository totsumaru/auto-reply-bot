package block

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// 返信の最大文字数です
	ReplyMaxLen = 500
)

// 返信です
type Reply struct {
	value string
}

// 返信を作成します
func NewReply(v string) (Reply, error) {
	r := Reply{}
	r.value = v

	if err := r.validate(); err != nil {
		return Reply{}, errors.NewError("検証に失敗しました", err)
	}

	return r, nil
}

// 返信の値を文字列で取得します
func (r Reply) String() string {
	return r.value
}

// 返信を比較します
func (r Reply) Equal(rr Reply) bool {
	return r.value == rr.value
}

// 返信の値が設定されているか判別します
func (r Reply) IsEmpty() bool {
	return r.value == ""
}

// 検証します
func (r Reply) validate() error {
	if err := validator.New().Var(r.value, "required"); err != nil {
		return errors.NewError("値が空です", err)
	}

	if len(r.value) > ReplyMaxLen {
		return errors.NewError("返信の文字数が上限を超えています")
	}

	return nil
}

// 構造体をJSONに変換します
func (r Reply) MarshalJSON() ([]byte, error) {
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
func (r *Reply) UnmarshalJSON(b []byte) error {
	j := &struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	r.value = j.Value

	return nil
}
