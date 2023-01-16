package rule

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// ルールです
type Rule struct {
	url URL
}

// ルールを作成します
func NewRule(url URL) (*Rule, error) {
	s := &Rule{}
	s.url = url

	if err := s.validate(); err != nil {
		return s, errors.NewError("検証に失敗しました", err)
	}

	return s, nil
}

// URLを取得します
func (u *Rule) URL() URL {
	return u.url
}

// 検証します
func (u *Rule) validate() error {
	return nil
}

// -------------------
// JSON
// -------------------

// 構造体をJSONに変換します
func (u *Rule) MarshalJSON() ([]byte, error) {
	j := struct {
		URL URL `json:"url"`
	}{
		URL: u.url,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません")
	}

	return b, nil
}

// JSONを構造体を変換します
func (u *Rule) UnmarshalJSON(b []byte) error {
	j := &struct {
		URL URL `json:"url"`
	}{}

	if err := json.Unmarshal(b, &j); err != nil {
		return errors.NewError("JSONを構造体に変換できません")
	}

	u.url = j.URL

	return nil
}
