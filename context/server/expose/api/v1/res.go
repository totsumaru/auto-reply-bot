package v1

import (
	"github.com/techstart35/auto-reply-bot/context/shared/map/seeker"
)

// レスポンス構造体です
type Res struct {
	ID          string
	AdminRoleID string
	Comment     struct {
		Block           []BlockRes
		IgnoreChannelID []string
	}
	Rule struct {
		URL struct {
			IsRestrict     bool
			IsYoutubeAllow bool
			IsTwitterAllow bool
			IsGIFAllow     bool
			IsOpenseaAllow bool
			IsDiscordAllow bool
			AllowRoleID    []string
			AllowChannelID []string
			//AlertChannelID string
		}
	}
}

// ブロックのレスポンスです
type BlockRes struct {
	Name             string
	Keyword          []string
	Reply            []string
	MatchCondition   string
	LimitedChannelID []string
	IsRandom         bool
	IsEmbed          bool
}

// レスポンスを作成します
func CreateRes(m map[string]interface{}) (Res, error) {
	blockRes := make([]BlockRes, 0)
	for _, bl := range seeker.Slice(m, []string{"comment", "block"}) {
		kw := make([]string, 0)
		for _, k := range seeker.Slice(bl, []string{"keyword"}) {
			kw = append(kw, seeker.Str(k, []string{"value"}))
		}

		rep := make([]string, 0)
		for _, r := range seeker.Slice(bl, []string{"reply"}) {
			rep = append(rep, seeker.Str(r, []string{"value"}))
		}

		limitedChID := make([]string, 0)
		for _, chID := range seeker.Slice(bl, []string{"limited_channel_id"}) {
			limitedChID = append(limitedChID, seeker.Str(chID, []string{"value"}))
		}

		b := BlockRes{}
		b.Name = seeker.Str(bl, []string{"name", "value"})
		b.Keyword = kw
		b.Reply = rep
		b.MatchCondition = seeker.Str(bl, []string{"match_condition", "value"})
		b.LimitedChannelID = limitedChID
		b.IsRandom = seeker.Bool(bl, []string{"is_random"})
		b.IsEmbed = seeker.Bool(bl, []string{"is_embed"})

		blockRes = append(blockRes, b)
	}

	ignoreChannelIDRes := make([]string, 0)
	for _, v := range seeker.Slice(m, []string{"comment", "ignore_channel_id"}) {
		ignoreChannelIDRes = append(ignoreChannelIDRes, seeker.Str(v, []string{"value"}))
	}

	res := Res{}
	res.ID = seeker.Str(m, []string{"id", "value"})
	res.AdminRoleID = seeker.Str(m, []string{"admin_role_id", "value"})
	// Comment
	res.Comment.Block = blockRes
	res.Comment.IgnoreChannelID = ignoreChannelIDRes

	// Rule
	{
		allowRoleID := make([]string, 0)
		for _, v := range seeker.Slice(m, []string{"rule", "url", "allow_role_id"}) {
			allowRoleID = append(allowRoleID, seeker.Str(v, []string{"value"}))
		}

		allowChannelID := make([]string, 0)
		for _, v := range seeker.Slice(m, []string{"rule", "url", "allow_channel_id"}) {
			allowChannelID = append(allowChannelID, seeker.Str(v, []string{"value"}))
		}

		res.Rule.URL.IsRestrict = seeker.Bool(m, []string{"rule", "url", "is_restrict"})
		res.Rule.URL.IsYoutubeAllow = seeker.Bool(m, []string{"rule", "url", "is_youtube_allow"})
		res.Rule.URL.IsTwitterAllow = seeker.Bool(m, []string{"rule", "url", "is_twitter_allow"})
		res.Rule.URL.IsGIFAllow = seeker.Bool(m, []string{"rule", "url", "is_gif_allow"})
		res.Rule.URL.IsOpenseaAllow = seeker.Bool(m, []string{"rule", "url", "is_opensea_allow"})
		res.Rule.URL.IsDiscordAllow = seeker.Bool(m, []string{"rule", "url", "is_discord_allow"})
		res.Rule.URL.AllowRoleID = allowRoleID
		res.Rule.URL.AllowChannelID = allowChannelID
		//res.Rule.URL.AlertChannelID = seeker.Str(m, []string{"rule", "url", "alert_channel_id", "value"})
	}

	return res, nil
}
