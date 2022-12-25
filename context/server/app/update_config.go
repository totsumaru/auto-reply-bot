package app

import (
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/block"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// 全ての設定の更新をするブロックのリクエストです
type BlockReq struct {
	Name       string
	Keyword    []string
	Reply      []string
	IsAllMatch bool
	IsRandom   bool
	IsEmbed    bool
}

// 全ての設定を更新します
//
// IDを返します。
func (a *App) UpdateConfig(
	serverID string,
	adminRoleID string,
	blockReq []BlockReq,
) (string, error) {
	i, err := server.NewID(serverID)
	if err != nil {
		return "", errors.NewError("idを作成できません", err)
	}

	s, err := a.Repo.FindByID(i)
	if err != nil {
		return "", errors.NewError("IDでサーバーを取得できません", err)
	}

	roleID, err := server.NewRoleID(adminRoleID)
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

		bl, err := block.NewBlock(
			name,
			keyword,
			reply,
			bReq.IsAllMatch,
			bReq.IsRandom,
			bReq.IsEmbed,
		)
		if err != nil {
			return "", errors.NewError("ブロックを作成できません", err)
		}

		blocks = append(blocks, bl)
	}

	// ロールIDを更新します
	if err = s.UpdateAdminRoleID(roleID); err != nil {
		return "", errors.NewError("管理者のロールIDを更新できません", err)
	}

	// ブロックを更新します
	if err = s.UpdateBlock(blocks); err != nil {
		return "", errors.NewError("ブロックを更新できません", err)
	}

	if err = a.Repo.Update(s); err != nil {
		return "", errors.NewError("サーバーを更新できません", err)
	}

	return s.ID().String(), nil
}
