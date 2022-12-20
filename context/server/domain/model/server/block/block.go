package block

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// 1つのブロックに設定できるキーワードの上限数です
	KeyWordMaxAmount = 5

	// 1つのブロックに設定できる返信の上限数です
	ReplyMaxAmount = 5
)

// ブロックです
type Block struct {
	keyword    []Keyword
	reply      []Reply
	isAllMatch bool // キーワードの完全一致フラグ(true=and, false=or)
	isRandom   bool // 返信のランダムフラグ
	isEmbed    bool // 埋め込みフラグ
}

// ブロックを作成します
func NewBlock(
	kw []Keyword,
	r []Reply,
	isAllMention bool,
	isRandom bool,
	isEmbed bool,
) (Block, error) {
	b := Block{}
	b.keyword = kw
	b.reply = r
	b.isAllMatch = isAllMention
	b.isRandom = isRandom
	b.isEmbed = isEmbed

	if err := b.validate(); err != nil {
		return Block{}, errors.NewError("検証に失敗しました", err)
	}

	return b, nil
}

// キーワードを取得します
func (b Block) Keyword() []Keyword {
	return b.keyword
}

// 返信を取得します
func (b Block) Reply() []Reply {
	return b.reply
}

// 完全一致フラグを取得します
func (b Block) IsAllMatch() bool {
	return b.isAllMatch
}

// ランダムフラグを取得します
func (b Block) IsRandom() bool {
	return b.isRandom
}

// 埋め込みフラグを取得します
func (b Block) IsEmbed() bool {
	return b.isEmbed
}

// 検証します
func (b Block) validate() error {
	if len(b.keyword) > KeyWordMaxAmount || len(b.keyword) == 0 {
		return errors.NewError("キーワドの数が不正です")
	}

	if len(b.reply) > ReplyMaxAmount || len(b.reply) == 0 {
		return errors.NewError("返信の数が不正です")
	}

	// ランダムではない時は返信は必ず1つ
	if !b.isRandom {
		if len(b.reply) > 1 {
			return errors.NewError("ランダムではない場合は返信は1つしか記述できません")
		}
	}

	return nil
}

// 構造体をJSONに変換します
func (b Block) MarshalJSON() ([]byte, error) {
	j := struct {
		Keyword    []Keyword `json:"keyword"`
		Reply      []Reply   `json:"reply"`
		IsAllMatch bool      `json:"is_all_match"`
		IsRandom   bool      `json:"is_random"`
		IsEmbed    bool      `json:"is_embed"`
	}{
		Keyword:    b.keyword,
		Reply:      b.reply,
		IsAllMatch: b.isAllMatch,
		IsRandom:   b.isRandom,
		IsEmbed:    b.isEmbed,
	}

	bb, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	return bb, nil
}

// JSONを構造体に変換します
func (b *Block) UnmarshalJSON(bb []byte) error {
	j := &struct {
		Keyword    []Keyword `json:"keyword"`
		Reply      []Reply   `json:"reply"`
		IsAllMatch bool      `json:"is_all_match"`
		IsRandom   bool      `json:"is_random"`
		IsEmbed    bool      `json:"is_embed"`
	}{}

	if err := json.Unmarshal(bb, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	b.keyword = j.Keyword
	b.reply = j.Reply
	b.isAllMatch = j.IsAllMatch
	b.isRandom = j.IsRandom
	b.isEmbed = j.IsEmbed

	return nil
}
