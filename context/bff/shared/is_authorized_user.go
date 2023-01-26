package shared

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/check"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/convert"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// 認証されているユーザーかを検証します
func IsAuthorizedUser(ctx context.Context, session *discordgo.Session, id, token string) (bool, error) {
	tmpRes, err := v1.FindByID(ctx, id)
	if err != nil {
		return false, errors.NewError("IDでサーバーを取得できません", err)
	}

	userID, err := convert.TokenToDiscordID(token)
	if err != nil {
		return false, errors.NewError("TokenをDiscordIDに変換できません", err)
	}

	guildOwnerID, err := guild.GetGuildOwnerID(session, id)
	if err != nil {
		return false, errors.NewError("ギルドのオーナーを取得できません", err)
	}

	// Dev/サーバーオーナー であればtrueを返します
	if userID == conf.TotsumaruDiscordID || userID == guildOwnerID {
		return true, nil
	}

	ok, err := check.HasRole(session, id, userID, tmpRes.AdminRoleID)
	if err != nil {
		return false, errors.NewError("ロールの所有確認に失敗しました", err)
	}

	return ok, nil
}
