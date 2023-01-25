package comment

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/comment/block"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// ブロックの上限数です
	BlockMaxAmount = 30
)

// コメントです
type Comment struct {
	block []block.Block
}

// コメントを作成します
func NewComment(block []block.Block) (Comment, error) {
	c := Comment{}
	c.block = block

	if err := c.validate(); err != nil {
		return c, errors.NewError("検証に失敗しました", err)
	}

	return c, nil
}

// ブロックを取得します
func (c Comment) Block() []block.Block {
	return c.block
}

// 検証します
func (c Comment) validate() error {
	if len(c.block) > BlockMaxAmount {
		return errors.NewError("ブロックの数が上限を超えています")
	}

	return nil
}

// -------------------
// JSON
// -------------------

// 構造体をJSONに変換します
func (c Comment) MarshalJSON() ([]byte, error) {
	j := struct {
		Block []block.Block `json:"block"`
	}{
		Block: c.block,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません")
	}

	return b, nil
}

// JSONを構造体を変換します
func (c *Comment) UnmarshalJSON(b []byte) error {
	j := &struct {
		Block []block.Block `json:"block"`
	}{}

	if err := json.Unmarshal(b, &j); err != nil {
		return errors.NewError("JSONを構造体に変換できません")
	}

	c.block = j.Block

	return nil
}
