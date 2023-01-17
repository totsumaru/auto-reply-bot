package app

import (
	"github.com/techstart35/auto-reply-bot/context/server/domain/model"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/block"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/rule"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// 全ての設定の更新をするブロックのリクエストです
type BlockReq struct {
	Name           string
	Keyword        []string
	Reply          []string
	MatchCondition string
	IsRandom       bool
	IsEmbed        bool
}

// URL制御のリクエストです
type URLRuleReq struct {
	IsRestrict     bool
	IsYoutubeAllow bool
	IsTwitterAllow bool
	IsGIFAllow     bool
	AllowRoleID    []string
	AllowChannelID []string
	AlertChannelID string
}

// 全ての設定を更新します
//
// IDを返します。
func (a *App) UpdateConfig(
	serverID string,
	adminRoleID string,
	blockReq []BlockReq,
	urlRuleReq URLRuleReq,
) (string, error) {
	i, err := model.NewID(serverID)
	if err != nil {
		return "", errors.NewError("idを作成できません", err)
	}

	s, err := a.Repo.FindByID(i)
	if err != nil {
		return "", errors.NewError("IDでサーバーを取得できません", err)
	}

	roleID, err := model.NewRoleID(adminRoleID)
	if err != nil {
		return "", errors.NewError("管理者のロールIDを作成できません", err)
	}

	blocks := make([]block.Block, 0)
	for _, bReq := range blockReq {
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

		bl, err := block.NewBlock(
			name,
			keyword,
			reply,
			matchCondition,
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
		for _, v := range urlRuleReq.AllowRoleID {
			alRoleID, err := model.NewRoleID(v)
			if err != nil {
				return "", errors.NewError("ロールIDを作成できません", err)
			}
			allowRoleID = append(allowRoleID, alRoleID)
		}

		allowChannelID := make([]model.ChannelID, 0)
		for _, v := range urlRuleReq.AllowChannelID {
			alChID, err := model.NewChannelID(v)
			if err != nil {
				return "", errors.NewError("チャンネルIDを作成できません", err)
			}
			allowChannelID = append(allowChannelID, alChID)
		}

		alChID := urlRuleReq.AlertChannelID
		if alChID == "" {
			alChID = "none"
		}
		alertChannelID, err := model.NewChannelID(alChID)
		if err != nil {
			return "", errors.NewError("アラートを送信するチャンネルIDを作成できません", err)
		}

		urlRule, err = rule.NewURL(
			urlRuleReq.IsRestrict,
			urlRuleReq.IsYoutubeAllow,
			urlRuleReq.IsTwitterAllow,
			urlRuleReq.IsGIFAllow,
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

	// ブロックを更新します
	if err = s.UpdateBlock(blocks); err != nil {
		return "", errors.NewError("ブロックを更新できません", err)
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
