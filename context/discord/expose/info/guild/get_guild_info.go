package guild

import (
	"fmt"
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
		return "", errors.NewError(fmt.Sprintf("ギルドを取得できません[ID: %s]", guildID), err)
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
//
// @everyoneは除外します。
func GetAllRolesWithoutEveryone(s *discordgo.Session, guildID string) (map[string]string, error) {
	res := map[string]string{}

	guild, err := s.Guild(guildID)
	if err != nil {
		return res, errors.NewError("ギルドを取得できません", err)
	}

	for _, role := range guild.Roles {
		if _, ok := res[role.ID]; ok {
			return res, errors.NewError("ロールが重複しています")
		}

		if role.ID != guildID {
			res[role.ID] = role.Name
		}
	}

	return res, nil
}

// 全てのテキストチャンネルを取得します
//
// チャンネルID:チャンネル名 のmapを返します。
func GetAllTextChannels(s *discordgo.Session, guildID string) (map[string]string, error) {
	res := map[string]string{}

	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return res, errors.NewError("チャンネル一覧を取得できません", err)
	}

	for _, channel := range channels {
		if _, ok := res[channel.ID]; ok {
			return res, errors.NewError("チャンネルが重複しています")
		}

		// テキストチャンネルのみ追加します
		if channel.Type == discordgo.ChannelTypeGuildText {
			res[channel.ID] = channel.Name
		}
	}

	return res, nil
}
