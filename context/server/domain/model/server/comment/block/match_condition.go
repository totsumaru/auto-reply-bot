package block

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	MatchConditionOneContain   = "one-contain"   // 1つでも含む
	MatchConditionAllContain   = "all-contain"   // 全てを含む
	MatchConditionPerfectMatch = "perfect-match" // 完全一致
)

// 一致条件です
type MatchCondition struct {
	value string
}

// 一致条件を作成します
func NewMatchCondition(v string) (MatchCondition, error) {
	m := MatchCondition{}
	m.value = v

	if err := m.validate(); err != nil {
		return MatchCondition{}, errors.NewError("検証に失敗しました", err)
	}

	return m, nil
}

// 一致条件の値を文字列で取得します
func (m MatchCondition) String() string {
	return m.value
}

// 一致条件を比較します
func (m MatchCondition) Equal(kk MatchCondition) bool {
	return m.value == kk.value
}

// 一致条件の値が設定されているか判別します
func (m MatchCondition) IsEmpty() bool {
	return m.value == ""
}

// 検証します
func (m MatchCondition) validate() error {
	if err := validator.New().Var(m.value, "required"); err != nil {
		return errors.NewError("値が空です", err)
	}

	switch m.value {
	case MatchConditionOneContain:
	case MatchConditionAllContain:
	case MatchConditionPerfectMatch:
	default:
		return errors.NewError("値が不正です")
	}

	return nil
}

// 構造体をJSONに変換します
func (m MatchCondition) MarshalJSON() ([]byte, error) {
	j := struct {
		Value string `json:"value"`
	}{
		Value: m.value,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	return b, nil
}

// JSONを構造体に変換します
func (m *MatchCondition) UnmarshalJSON(b []byte) error {
	j := &struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	m.value = j.Value

	return nil
}
