package v1

import (
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/block"
	"github.com/techstart35/auto-reply-bot/context/shared/map/seeker"
)

// レスポンス構造体です
type Res struct {
	ID          string
	AdminRoleID string
	Block       []BlockRes
}

// ブロックのレスポンスです
type BlockRes struct {
	Name           string
	Keyword        []string
	Reply          []string
	MatchCondition string
	IsRandom       bool
	IsEmbed        bool
}

// レスポンスを作成します
func CreateRes(m map[string]interface{}) (Res, error) {
	blockRes := make([]BlockRes, 0)

	for _, bl := range seeker.Slice(m, []string{"block"}) {
		kw := make([]string, 0)
		for _, k := range seeker.Slice(bl, []string{"keyword"}) {
			kw = append(kw, seeker.Str(k, []string{"value"}))
		}

		rep := make([]string, 0)
		for _, r := range seeker.Slice(bl, []string{"reply"}) {
			rep = append(rep, seeker.Str(r, []string{"value"}))
		}

		// TODO: DBが全て更新されたら削除(修正)
		matchCondition := seeker.Str(bl, []string{"match_condition", "value"})
		if matchCondition == "" {
			if seeker.Bool(bl, []string{"is_all_match"}) {
				matchCondition = block.MatchConditionAllContain
			} else {
				matchCondition = block.MatchConditionOneContain
			}
		}

		b := BlockRes{}
		b.Name = seeker.Str(bl, []string{"name", "value"})
		b.Keyword = kw
		b.Reply = rep
		b.MatchCondition = matchCondition
		b.IsRandom = seeker.Bool(bl, []string{"is_random"})
		b.IsEmbed = seeker.Bool(bl, []string{"is_embed"})

		blockRes = append(blockRes, b)
	}

	res := Res{}

	res.ID = seeker.Str(m, []string{"id", "value"})
	res.AdminRoleID = seeker.Str(m, []string{"admin_role_id", "value"})
	res.Block = blockRes

	return res, nil
}
