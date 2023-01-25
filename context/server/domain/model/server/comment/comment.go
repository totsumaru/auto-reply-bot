package comment

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/comment/block"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// ブロックの上限数です
	BlockMaxAmount = 30

	// コメント機能を実行しないチャンネルの上限数です
	IgnoreChannelIDMaxAmount = 10
)

// コメントです
type Comment struct {
	block           []block.Block     // ブロック
	ignoreChannelID []model.ChannelID // コメント機能を実行しないチャンネル
}

// コメントを作成します
func NewComment(block []block.Block, ignoreChannel []model.ChannelID) (Comment, error) {
	c := Comment{}
	c.block = block
	c.ignoreChannelID = ignoreChannel

	if err := c.validate(); err != nil {
		return c, errors.NewError("検証に失敗しました", err)
	}

	return c, nil
}

// ブロックを取得します
func (c Comment) Block() []block.Block {
	return c.block
}

// コメントを実行しないチャンネルを取得します
func (c Comment) IgnoreChannel() []model.ChannelID {
	return c.ignoreChannelID
}

// 検証します
func (c Comment) validate() error {
	if len(c.block) > BlockMaxAmount {
		return errors.NewError("ブロックの数が上限を超えています")
	}

	if len(c.ignoreChannelID) > IgnoreChannelIDMaxAmount {
		return errors.NewError("コメントを実行しないチャンネルが上限を超えています")
	}

	return nil
}

// -------------------
// JSON
// -------------------

// 構造体をJSONに変換します
func (c Comment) MarshalJSON() ([]byte, error) {
	j := struct {
		Block         []block.Block     `json:"block"`
		IgnoreChannel []model.ChannelID `json:"ignore_channel"`
	}{
		Block:         c.block,
		IgnoreChannel: c.ignoreChannelID,
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
		Block         []block.Block     `json:"block"`
		IgnoreChannel []model.ChannelID `json:"ignore_channel"`
	}{}

	if err := json.Unmarshal(b, &j); err != nil {
		return errors.NewError("JSONを構造体に変換できません")
	}

	c.block = j.Block
	c.ignoreChannelID = j.IgnoreChannel

	return nil
}
