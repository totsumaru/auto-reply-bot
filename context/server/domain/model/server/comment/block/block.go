package block

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// 1つのブロックに設定できるキーワードの上限数です
	KeyWordMaxAmount = 10

	// 1つのブロックに設定できる返信の上限数です
	ReplyMaxAmount = 10

	// 1つのブロックに設定できる限定発動チャンネルの上限数です
	LimitedChannelIDMaxAmount = 3
)

// ブロックです
type Block struct {
	name             Name              // 表示名
	keyword          []Keyword         // キーワード
	reply            []Reply           // 返信
	matchCondition   MatchCondition    // 一致条件(完全一致,部分一致など)
	limitedChannelID []model.ChannelID // 自動返信はこのチャンネルでのみ発動(任意)
	isRandom         bool              // 返信のランダムフラグ
	isEmbed          bool              // 埋め込みフラグ
}

// ブロックを作成します
func NewBlock(
	n Name,
	kw []Keyword,
	r []Reply,
	matchCondition MatchCondition,
	limitedChannelID []model.ChannelID,
	isRandom bool,
	isEmbed bool,
) (Block, error) {
	b := Block{}
	b.name = n
	b.keyword = kw
	b.reply = r
	b.matchCondition = matchCondition
	b.limitedChannelID = limitedChannelID
	b.isRandom = isRandom
	b.isEmbed = isEmbed

	if err := b.validate(); err != nil {
		return Block{}, errors.NewError("検証に失敗しました", err)
	}

	return b, nil
}

// 名前を取得します
func (b Block) Name() Name {
	return b.name
}

// キーワードを取得します
func (b Block) Keyword() []Keyword {
	return b.keyword
}

// 返信を取得します
func (b Block) Reply() []Reply {
	return b.reply
}

// 一致条件を取得します
func (b Block) MatchCondition() MatchCondition {
	return b.matchCondition
}

// ここでのみ起動するチャンネルIDを取得します
func (b Block) LimitedChannelID() []model.ChannelID {
	return b.limitedChannelID
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

	if len(b.limitedChannelID) > LimitedChannelIDMaxAmount {
		return errors.NewError("限定起動のチャンネル数が不正です")
	}

	// 同じキーワードが設定できないように制限
	// ※同じ返信は設定できます
	tmpKeyword := map[string]bool{}
	for _, kw := range b.keyword {
		if _, ok := tmpKeyword[kw.value]; ok {
			return errors.NewError("キーワードが重複しています")
		}

		tmpKeyword[kw.value] = true
	}

	// ランダムではない時は返信は必ず1つ
	if !b.isRandom {
		// ※同じ返信は設定できます
		if len(b.reply) > 1 {
			return errors.NewError("ランダムではない場合は返信は1つしか記述できません")
		}
	}

	// 完全一致の場合はキーワードは必ず1つ
	if b.matchCondition.String() == MatchConditionPerfectMatch {
		if len(b.keyword) > 1 {
			return errors.NewError("完全一致の場合はキーワードは1つしか設定できません")
		}
	}

	return nil
}

// 構造体をJSONに変換します
func (b Block) MarshalJSON() ([]byte, error) {
	j := struct {
		Name             Name              `json:"name"`
		Keyword          []Keyword         `json:"keyword"`
		Reply            []Reply           `json:"reply"`
		MatchCondition   MatchCondition    `json:"match_condition"`
		LimitedChannelID []model.ChannelID `json:"limited_channel_id"`
		IsRandom         bool              `json:"is_random"`
		IsEmbed          bool              `json:"is_embed"`
	}{
		Name:             b.name,
		Keyword:          b.keyword,
		Reply:            b.reply,
		MatchCondition:   b.matchCondition,
		LimitedChannelID: b.limitedChannelID,
		IsRandom:         b.isRandom,
		IsEmbed:          b.isEmbed,
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
		Name             Name              `json:"name"`
		Keyword          []Keyword         `json:"keyword"`
		Reply            []Reply           `json:"reply"`
		MatchCondition   MatchCondition    `json:"match_condition"`
		LimitedChannelID []model.ChannelID `json:"limited_channel_id"`
		IsRandom         bool              `json:"is_random"`
		IsEmbed          bool              `json:"is_embed"`
	}{}

	if err := json.Unmarshal(bb, j); err != nil {
		return errors.NewError("JSONを構造体に変換できません", err)
	}

	b.name = j.Name
	b.keyword = j.Keyword
	b.reply = j.Reply
	b.matchCondition = j.MatchCondition
	b.limitedChannelID = j.LimitedChannelID
	b.isRandom = j.IsRandom
	b.isEmbed = j.IsEmbed

	return nil
}
