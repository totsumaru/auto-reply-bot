package rule

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	URLPrefixHTTP  = "http://"
	URLPrefixHTTPS = "https://"

	YoutubeURL    = "https://youtube.com"
	YoutubeWWWURL = "https://www.youtube.com"

	TwitterURL    = "https://twitter.com"
	TwitterWWWURL = "https://www.twitter.com"

	GIFURL    = "https://tenor.com"
	GIFWWWURL = "https://www.tenor.com"

	AllowRoleMAX    = 5
	AllowChannelMAX = 10
)

// URL制限です
//
// alertChannelIDは空の構造体が入る可能性があります。
type URL struct {
	isRestrict     bool              // URL制限をするか
	isYoutubeAllow bool              // Youtubeを許可するか
	isTwitterAllow bool              // Twitterを許可するか
	isGIFAllow     bool              // GIFを許可するか
	allowRoleID    []model.RoleID    // URL制限を受けないロールID
	allowChannelID []model.ChannelID // URL制限を受けないチャンネルID
	alertChannelID model.ChannelID   // 禁止されたURLのメッセージが送信された時にログを送信するチャンネル
}

// URL制限を作成します
func NewURL(
	isRestrict bool,
	isYoutubeAllow bool,
	isTwitterAllow bool,
	isGIFAllow bool,
	allowRoleID []model.RoleID,
	allowChannelID []model.ChannelID,
	alertChannelID model.ChannelID,
) (URL, error) {
	u := URL{}
	u.isRestrict = isRestrict
	u.isYoutubeAllow = isYoutubeAllow
	u.isTwitterAllow = isTwitterAllow
	u.isGIFAllow = isGIFAllow
	u.allowRoleID = allowRoleID
	u.allowChannelID = allowChannelID
	u.alertChannelID = alertChannelID

	if err := u.validate(); err != nil {
		return u, errors.NewError("検証に失敗しました", err)
	}

	return u, nil
}

// URL制限をするかを取得します
func (u URL) IsRestrict() bool {
	return u.isRestrict
}

// Youtubeを許可するかを取得します
func (u URL) IsYoutubeAllow() bool {
	return u.isYoutubeAllow
}

// Twitterを許可するかを取得します
func (u URL) IsTwitterAllow() bool {
	return u.isTwitterAllow
}

// GIFを許可するかを取得します
func (u URL) IsGIFAllow() bool {
	return u.isGIFAllow
}

// URL制限を受けないロールを取得します
func (u URL) AllowRoleID() []model.RoleID {
	return u.allowRoleID
}

// URL制限を受けないチャンネルを取得します
func (u URL) AllowChannelID() []model.ChannelID {
	return u.allowChannelID
}

// アラートを通知するチャンネルを取得します
func (u URL) AlertChannelID() model.ChannelID {
	return u.alertChannelID
}

// -------------------
// validation
// -------------------

// 検証します
func (u URL) validate() error {
	// ロールの数を検証します
	if len(u.allowRoleID) > AllowRoleMAX {
		return errors.NewError("ロールの数が上限を超えています")
	}

	// チャンネルの数を検証します
	if len(u.allowChannelID) > AllowChannelMAX {
		return errors.NewError("チャンネルの数が上限を超えています")
	}

	// 許可するロールIDで重複していないかを検証します
	//
	// 重複確認のmapのため、valueのboolは意味を持ちません。
	tmpRoleID := map[string]bool{}
	for _, roleID := range u.allowRoleID {
		if _, ok := tmpRoleID[roleID.String()]; ok {
			return errors.NewError("許可するロールIDが重複しています")
		}
		tmpRoleID[roleID.String()] = true
	}

	// 許可するチャンネルIDで重複していないかを検証します
	//
	// 重複確認のmapのため、valueのboolは意味を持ちません。
	tmpChannelID := map[string]bool{}
	for _, channelID := range u.allowChannelID {
		if _, ok := tmpChannelID[channelID.String()]; ok {
			return errors.NewError("許可するチャンネルIDが重複しています")
		}
		tmpChannelID[channelID.String()] = true
	}

	return nil
}

// -------------------
// JSON
// -------------------

// 構造体をJSONに変換します
func (u URL) MarshalJSON() ([]byte, error) {
	j := struct {
		IsRestrict     bool              `json:"is_restrict"`
		IsYoutubeAllow bool              `json:"is_youtube_allow"`
		IsTwitterAllow bool              `json:"is_twitter_allow"`
		IsGIFAllow     bool              `json:"is_gif_allow"`
		AllowRoleID    []model.RoleID    `json:"allow_role_id"`
		AllowChannelID []model.ChannelID `json:"allow_channel_id"`
		AlertChannelID model.ChannelID   `json:"alert_channel_id"`
	}{
		IsRestrict:     u.isRestrict,
		IsYoutubeAllow: u.isYoutubeAllow,
		IsTwitterAllow: u.isTwitterAllow,
		IsGIFAllow:     u.isGIFAllow,
		AllowRoleID:    u.allowRoleID,
		AllowChannelID: u.allowChannelID,
		AlertChannelID: u.alertChannelID,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません")
	}

	return b, nil
}

// JSONを構造体を変換します
func (u *URL) UnmarshalJSON(b []byte) error {
	j := &struct {
		IsRestrict     bool              `json:"is_restrict"`
		IsYoutubeAllow bool              `json:"is_youtube_allow"`
		IsTwitterAllow bool              `json:"is_twitter_allow"`
		IsGIFAllow     bool              `json:"is_gif_allow"`
		AllowRoleID    []model.RoleID    `json:"allow_role_id"`
		AllowChannelID []model.ChannelID `json:"allow_channel_id"`
		AlertChannelID model.ChannelID   `json:"alert_channel_id"`
	}{}

	if err := json.Unmarshal(b, &j); err != nil {
		return errors.NewError("JSONを構造体に変換できません")
	}

	u.isRestrict = j.IsRestrict
	u.isYoutubeAllow = j.IsYoutubeAllow
	u.isTwitterAllow = j.IsTwitterAllow
	u.isGIFAllow = j.IsGIFAllow
	u.allowRoleID = j.AllowRoleID
	u.allowChannelID = j.AllowChannelID
	u.alertChannelID = j.AlertChannelID

	return nil
}
