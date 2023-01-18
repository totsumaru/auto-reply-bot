package server

import (
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/check"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/convert"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/initiate"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	v1 "github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"net/http"
	"os"
)

// サーバーを取得します
func Server(e *gin.Engine) {
	e.GET("/server", getServer)
}

// レスポンスです
type Res struct {
	ID          string     `json:"id"`
	AdminRoleID string     `json:"admin_role_id"`
	Block       []resBlock `json:"block"`
	// 以下はComputedです
	Token      string       `json:"token"`
	ServerName string       `json:"server_name"`
	AvatarURL  string       `json:"avatar_url"`
	Role       []resRole    `json:"role"`
	Channel    []resChannel `json:"channel"`
	Nickname   string       `json:"nickname"`
	Rule       struct {
		URL struct {
			IsRestrict     bool     `json:"is_restrict"`
			IsYoutubeAllow bool     `json:"is_youtube_allow"`
			IsTwitterAllow bool     `json:"is_twitter_allow"`
			IsGIFAllow     bool     `json:"is_gif_allow"`
			IsOpenSeaAllow bool     `json:"is_opensea_allow"`
			IsDiscordAllow bool     `json:"is_discord_allow"`
			AllowRoleID    []string `json:"allow_role_id"`
			AllowChannelID []string `json:"allow_channel_id"`
			AlertChannelID string   `json:"alert_channel_id"`
		} `json:"url"`
	} `json:"rule"`
}

// ブロックのレスポンスです
type resBlock struct {
	Name           string   `json:"name"`
	Keyword        []string `json:"keyword"`
	Reply          []string `json:"reply"`
	MatchCondition string   `json:"match_condition"`
	IsRandom       bool     `json:"is_random"`
	IsEmbed        bool     `json:"is_embed"`
}

// ロールのレスポンスです
type resRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// チャンネルのレスポンスです
type resChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// サーバーを取得します
func getServer(c *gin.Context) {
	session, err := initiate.CreateSession()
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("セッションを作成できません", err), "none")
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	ctx, tx, err := shared.CreateDBTx()
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("DBのTxを作成できません", err), "none")
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	id := c.Query("id")
	code := c.Query("code")

	var (
		token string
	)

	// クエリパラメータに指定されたサーバーです
	guildName, err := guild.GetGuildName(session, id)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("ギルド名を取得できません", err), "")
		return
	}

	// 認証されているユーザーかを検証します
	{
		tmpRes, err := v1.FindByID(ctx, id)
		if err != nil {
			message_send.SendErrMsg(session, errors.NewError("IDでサーバーを取得できません", err), guildName)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		token, err = convert.CodeToToken(code)
		if err != nil {
			// codeの不正に関してはエラー通知しません
			message_send.SendErrMsg(session, errors.NewError("codeをtokenに変換できません", err), guildName)
			c.JSON(http.StatusUnauthorized, "codeが認証されていません")
			return
		}

		userID, err := convert.TokenToDiscordID(token)
		if err != nil {
			message_send.SendErrMsg(session, errors.NewError("トークンをDiscordIDに変換できません", err), guildName)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		guildOwnerID, err := guild.GetGuildOwnerID(session, id)
		if err != nil {
			message_send.SendErrMsg(session, errors.NewError("オーナーIDを取得できません", err), guildName)
			c.JSON(http.StatusUnauthorized, "認証されていません")
			return
		}

		if !(userID == conf.TotsumaruDiscordID || userID == guildOwnerID) {
			ok, err := check.HasRole(session, id, userID, tmpRes.AdminRoleID)
			if err != nil {
				message_send.SendErrMsg(session, errors.NewError("ロールの所有確認に失敗しました", err), guildName)
				c.JSON(http.StatusUnauthorized, "認証されていません")
				return
			}

			if !ok {
				message_send.SendErrMsg(session, errors.NewError("認証に失敗しました", err), guildName)
				c.JSON(http.StatusUnauthorized, "認証されていません")
				return
			}
		}
	}

	var (
		apiRes = v1.Res{}
	)

	bffErr := (func() error {
		apiRes, err = v1.FindByID(ctx, id)
		if err != nil {
			return errors.NewError("IDでサーバーを取得できません", err)
		}

		return nil
	})()

	if bffErr != nil {
		// ロールバックを実行します
		txErr := tx.Rollback()
		if txErr != nil {
			message_send.SendErrMsg(
				session,
				errors.NewError("ロールバックに失敗しました。データに不整合が発生している可能性があります。", txErr),
				guildName,
			)
			return
		}

		message_send.SendErrMsg(session, errors.NewError("バックエンドの処理でエラーが発生しました", bffErr), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	if txErr := tx.Commit(); txErr != nil {
		message_send.SendErrMsg(session, errors.NewError("Commitに失敗しました", txErr), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	avatarURL, err := guild.GetAvatarURL(session, apiRes.ID)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("アバターURLを取得できません", err), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	allRoles, err := guild.GetAllRolesWithoutEveryone(session, apiRes.ID)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("全てのロールを取得できません", err), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	allChannels, err := guild.GetAllTextChannels(session, apiRes.ID)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("全てのチャンネルを取得できません", err), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	m, err := session.GuildMember(id, os.Getenv("DISCORD_APPLICATION_ID"))
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("botのMember情報を取得できません", err), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	res := Res{}
	res.ID = apiRes.ID
	res.AdminRoleID = apiRes.AdminRoleID
	res.Block = []resBlock{}
	res.Token = token
	res.ServerName = guildName
	res.AvatarURL = avatarURL
	res.Role = []resRole{}
	res.Channel = []resChannel{}
	res.Nickname = m.Nick
	res.Rule.URL.IsRestrict = apiRes.Rule.URL.IsRestrict
	res.Rule.URL.IsYoutubeAllow = apiRes.Rule.URL.IsYoutubeAllow
	res.Rule.URL.IsTwitterAllow = apiRes.Rule.URL.IsTwitterAllow
	res.Rule.URL.IsGIFAllow = apiRes.Rule.URL.IsGIFAllow
	res.Rule.URL.IsOpenSeaAllow = apiRes.Rule.URL.IsOpenseaAllow
	res.Rule.URL.IsDiscordAllow = apiRes.Rule.URL.IsDiscordAllow
	res.Rule.URL.AllowRoleID = apiRes.Rule.URL.AllowRoleID
	res.Rule.URL.AllowChannelID = apiRes.Rule.URL.AllowChannelID
	res.Rule.URL.AlertChannelID = apiRes.Rule.URL.AlertChannelID

	// レスポンスにブロックを追加します
	for _, v := range apiRes.Block {
		blockRes := resBlock{}
		blockRes.Name = v.Name
		blockRes.Keyword = v.Keyword
		blockRes.Reply = v.Reply
		blockRes.MatchCondition = v.MatchCondition
		blockRes.IsRandom = v.IsRandom
		blockRes.IsEmbed = v.IsEmbed

		res.Block = append(res.Block, blockRes)
	}

	// レスポンスにロールを追加します
	for roleID, roleName := range allRoles {
		tmpRole := resRole{
			ID:   roleID,
			Name: roleName,
		}

		res.Role = append(res.Role, tmpRole)
	}

	// レスポンスにチャンネルを追加します
	for channelID, channelName := range allChannels {
		tmpChannel := resChannel{
			ID:   channelID,
			Name: channelName,
		}

		res.Channel = append(res.Channel, tmpChannel)
	}

	c.JSON(http.StatusOK, res)
}
