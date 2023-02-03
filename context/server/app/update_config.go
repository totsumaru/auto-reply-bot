package app

import (
	"github.com/techstart35/auto-reply-bot/context/server/domain/model"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/comment"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/comment/block"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/rule"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// リクエストです
type Req struct {
	AdminRoleID string
	Comment     struct {
		BlockReq        []BlockReq
		IgnoreChannelID []string
	}
	URLRule struct {
		IsRestrict     bool
		IsYoutubeAllow bool
		IsTwitterAllow bool
		IsGIFAllow     bool
		IsOpenseaAllow bool
		IsDiscordAllow bool
		AllowRoleID    []string
		AllowChannelID []string
	}
}

// ブロックのリクエストです
type BlockReq struct {
	Name             string
	Keyword          []string
	Reply            []string
	MatchCondition   string
	LimitedChannelID []string
	IsRandom         bool
	IsEmbed          bool
}

// 全ての設定を更新します
//
// IDを返します。
func (a *App) UpdateConfig(serverID string, req Req) (string, error) {
	i, err := model.NewID(serverID)
	if err != nil {
		return "", errors.NewError("idを作成できません", err)
	}

	s, err := a.Repo.FindByID(i)
	if err != nil {
		return "", errors.NewError("IDでサーバーを取得できません", err)
	}

	roleID, err := model.NewRoleID(req.AdminRoleID)
	if err != nil {
		return "", errors.NewError("管理者のロールIDを作成できません", err)
	}

	blocks := make([]block.Block, 0)
	for _, bReq := range req.Comment.BlockReq {
		// 名前
		name, err := block.NewName(bReq.Name)
		if err != nil {
			return "", errors.NewError("ブロック名を作成できません", err)
		}

		// キーワード
		keyword := make([]block.Keyword, 0)
		{
			for _, kw := range bReq.Keyword {
				k, err := block.NewKeyword(kw)
				if err != nil {
					return "", errors.NewError("キーワードを作成できません", err)
				}

				keyword = append(keyword, k)
			}
		}

		// 返信
		reply := make([]block.Reply, 0)
		{
			for _, rep := range bReq.Reply {
				r, err := block.NewReply(rep)
				if err != nil {
					return "", errors.NewError("返信を作成できません", err)
				}

				reply = append(reply, r)
			}
		}

		// 一致条件
		matchCondition, err := block.NewMatchCondition(bReq.MatchCondition)
		if err != nil {
			return "", errors.NewError("一致条件を作成できません", err)
		}

		// 限定起動するチャンネルID
		limitedChID := make([]model.ChannelID, 0)
		for _, ch := range bReq.LimitedChannelID {
			chID, err := model.NewChannelID(ch)
			if err != nil {
				return "", errors.NewError("チャンネルIDを作成できません", err)
			}

			limitedChID = append(limitedChID, chID)
		}

		bl, err := block.NewBlock(
			name,
			keyword,
			reply,
			matchCondition,
			limitedChID,
			bReq.IsRandom,
			bReq.IsEmbed,
		)
		if err != nil {
			return "", errors.NewError("ブロックを作成できません", err)
		}

		blocks = append(blocks, bl)
	}

	// URLのルールを作成します
	urlRule := rule.URL{}
	{
		allowRoleID := make([]model.RoleID, 0)
		for _, v := range req.URLRule.AllowRoleID {
			alRoleID, err := model.NewRoleID(v)
			if err != nil {
				return "", errors.NewError("ロールIDを作成できません", err)
			}
			allowRoleID = append(allowRoleID, alRoleID)
		}

		allowChannelID := make([]model.ChannelID, 0)
		for _, v := range req.URLRule.AllowChannelID {
			alChID, err := model.NewChannelID(v)
			if err != nil {
				return "", errors.NewError("チャンネルIDを作成できません", err)
			}
			allowChannelID = append(allowChannelID, alChID)
		}

		// FEから通知チャンネルを受け取らないようにしたため、
		// 空のチャンネルIDを生成します。
		alertChannelID := model.ChannelID{}

		urlRule, err = rule.NewURL(
			req.URLRule.IsRestrict,
			req.URLRule.IsYoutubeAllow,
			req.URLRule.IsTwitterAllow,
			req.URLRule.IsGIFAllow,
			req.URLRule.IsOpenseaAllow,
			req.URLRule.IsDiscordAllow,
			allowRoleID,
			allowChannelID,
			alertChannelID,
		)
		if err != nil {
			return "", errors.NewError("URLのルールを作成できません", err)
		}
	}

	// ルールを作成します
	r, err := rule.NewRule(urlRule)
	if err != nil {
		return "", errors.NewError("ルールを作成できません", err)
	}

	// ロールIDを更新します
	if err = s.UpdateAdminRoleID(roleID); err != nil {
		return "", errors.NewError("管理者のロールIDを更新できません", err)
	}

	// コメントを無視するチャンネルを作成します
	ignoreCh := make([]model.ChannelID, 0)
	for _, chID := range req.Comment.IgnoreChannelID {
		cID, err := model.NewChannelID(chID)
		if err != nil {
			return "", errors.NewError("チャンネルIDを作成できません", err)
		}

		ignoreCh = append(ignoreCh, cID)
	}

	// コメントを作成します
	c, err := comment.NewComment(blocks, ignoreCh)
	if err != nil {
		return "", errors.NewError("コメントを作成できません", err)
	}

	// コメントを更新します
	if err = s.UpdateComment(c); err != nil {
		return "", errors.NewError("コメントを更新できません", err)
	}

	// ルールを更新します
	if err = s.UpdateRule(r); err != nil {
		return "", errors.NewError("ルールを更新できません", err)
	}

	if err = a.Repo.Update(s); err != nil {
		return "", errors.NewError("サーバーを更新できません", err)
	}

	return s.ID().String(), nil
}
