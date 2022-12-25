package v1

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator/v10"
	"github.com/techstart35/auto-reply-bot/context/server/app"
	"github.com/techstart35/auto-reply-bot/context/server/gateway/di"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// ブロックのリクエストです
type BlockReq struct {
	Name       string
	Keyword    []string
	Reply      []string
	IsAllMatch bool
	IsRandom   bool
	IsEmbed    bool
}

// 設定を更新します
//
// 管理者ロールを持っている or 該当サーバーの管理者権限 のみがコールできます。
func UpdateConfig(
	s *discordgo.Session,
	ctx context.Context,
	serverID string,
	adminRoleID string,
	blockReq []BlockReq,
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

	appBlockReq := make([]app.BlockReq, 0)
	for _, v := range blockReq {
		bl := app.BlockReq{
			Name:       v.Name,
			Keyword:    v.Keyword,
			Reply:      v.Reply,
			IsAllMatch: v.IsAllMatch,
			IsRandom:   v.IsRandom,
			IsEmbed:    v.IsEmbed,
		}

		appBlockReq = append(appBlockReq, bl)
	}

	appResID, err := a.UpdateConfig(serverID, adminRoleID, appBlockReq)
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

	return res, nil
}
