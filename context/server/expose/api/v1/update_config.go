package v1

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/server/app"
	"github.com/techstart35/auto-reply-bot/context/server/gateway/di"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// リクエストです
type Req struct {
	AdminRoleID string
	Comment     struct {
		BlockReq        []BlockReq
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
		}
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

// 設定を更新します
//
// 管理者ロールを持っている or 該当サーバーの管理者権限 のみがコールできます。
func UpdateConfig(
	s *discordgo.Session,
	ctx context.Context,
	serverID string,
	req Req,
) (Res, error) {
	res := Res{}

	if err := validator.New().Var(serverID, "required"); err != nil {
		return res, errors.NewError("リクエストが不正です", err)
	}

	a, err := di.InitApp(ctx, s)
	if err != nil {
		return res, errors.NewError("アプリケーションを初期化できません", err)
	}

	q, err := di.InitQuery(ctx)
	if err != nil {
		return res, errors.NewError("クエリーを初期化できません", err)
	}

	blockReqs := make([]app.BlockReq, 0)
	for _, br := range req.Comment.BlockReq {
		bReq := app.BlockReq{
			Name:             br.Name,
			Keyword:          br.Keyword,
			Reply:            br.Reply,
			MatchCondition:   br.MatchCondition,
			LimitedChannelID: br.LimitedChannelID,
			IsRandom:         br.IsRandom,
			IsEmbed:          br.IsEmbed,
		}
		blockReqs = append(blockReqs, bReq)
	}

	appReq := app.Req{}
	appReq.AdminRoleID = req.AdminRoleID
	// Comment
	appReq.Comment.BlockReq = blockReqs
	appReq.Comment.IgnoreChannelID = req.Comment.IgnoreChannelID
	// URL-Rule
	appReq.URLRule.IsRestrict = req.Rule.URL.IsRestrict
	appReq.URLRule.IsYoutubeAllow = req.Rule.URL.IsYoutubeAllow
	appReq.URLRule.IsTwitterAllow = req.Rule.URL.IsTwitterAllow
	appReq.URLRule.IsGIFAllow = req.Rule.URL.IsGIFAllow
	appReq.URLRule.IsOpenseaAllow = req.Rule.URL.IsOpenseaAllow
	appReq.URLRule.IsDiscordAllow = req.Rule.URL.IsDiscordAllow
	appReq.URLRule.AllowRoleID = req.Rule.URL.AllowRoleID
	appReq.URLRule.AllowChannelID = req.Rule.URL.AllowChannelID

	appResID, err := a.UpdateConfig(serverID, appReq)
	if err != nil {
		return res, errors.NewError("設定を更新できません", err)
	}

	m, err := q.FindByID(appResID)
	if err != nil {
		return res, errors.NewError("IDでサーバーを取得できません", err)
	}

	res, err = CreateRes(m)
	if err != nil {
		return res, errors.NewError("レスポンスを作成できません", err)
	}

	// storeに値を保存します
	if err = updateStore(res); err != nil {
		return res, errors.NewError("storeを更新できません", err)
	}

	return res, nil
}
