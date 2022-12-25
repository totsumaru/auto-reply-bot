package guild

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

// ギルド構造体を取得します
//
// 基本的には個別のメソッドを推奨します。
func GetGuild(s *discordgo.Session, guildID string) (*discordgo.Guild, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return nil, errors.NewError("ギルドを取得できません", err)
	}

	return g, nil
}

// サーバー名を取得します
func GetGuildName(s *discordgo.Session, guildID string) (string, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return "", errors.NewError("ギルドを取得できません", err)
	}

	return g.Name, nil
}

// アバターのURLを取得します
func GetAvatarURL(s *discordgo.Session, guildID string) (string, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return "", errors.NewError("ギルドを取得できません", err)
	}

	return g.IconURL(), nil
}

// ギルドの所有者(owner)を取得します
func GetGuildOwnerID(s *discordgo.Session, guildID string) (string, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return "", errors.NewError("ギルドを取得できません", err)
	}

	return g.OwnerID, nil
}

// 全てのロールを取得します
//
// ロールID:ロール名 のmapを返します。
func GetAllRoles(s *discordgo.Session, guildID string) (map[string]string, error) {
	res := map[string]string{}

	guild, err := s.Guild(guildID)
	if err != nil {
		return res, errors.NewError("ギルドを取得できません", err)
	}

	for _, role := range guild.Roles {
		if _, ok := res[role.ID]; ok {
			return res, errors.NewError("ロールが重複しています")
		}

		res[role.ID] = role.Name
	}

	return res, nil
}
