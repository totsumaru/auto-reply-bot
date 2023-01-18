package config

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
)

// サーバーの設定を更新します
func ServerConfig(e *gin.Engine) {
	e.POST("/server/config", postServerConfig)
}

// リクエストBodyです
type ReqConfig struct {
	AdminRoleID string `json:"admin_role_id"`
	Block       []struct {
		Name           string   `json:"name"`
		Keyword        []string `json:"keyword"`
		Reply          []string `json:"reply"`
		MatchCondition string   `json:"match_condition"`
		IsRandom       bool     `json:"is_random"`
		IsEmbed        bool     `json:"is_embed"`
	} `json:"block"`
	Rule struct {
		URL struct {
			IsRestrict     bool     `json:"is_restrict"`
			IsYoutubeAllow bool     `json:"is_youtube_allow"`
			IsTwitterAllow bool     `json:"is_twitter_allow"`
			IsGIFAllow     bool     `json:"is_gif_allow"`
			IsOpenseaAllow bool     `json:"is_opensea_allow"`
			IsDiscordAllow bool     `json:"is_discord_allow"`
			AllowRoleID    []string `json:"allow_role_id"`
			AllowChannelID []string `json:"allow_channel_id"`
			AlertChannelID string   `json:"alert_channel_id"`
		} `json:"url"`
	} `json:"rule"`
}

// レスポンスです
type Res struct {
	ID          string     `json:"id"`
	AdminRoleID string     `json:"admin_role_id"`
	Block       []resBlock `json:"block"`
	// 以下はComputedです
	ServerName string       `json:"server_name"`
	AvatarURL  string       `json:"avatar_url"`
	Role       []resRole    `json:"role"`
	Channel    []resChannel `json:"channel"`
	Rule       struct {
		URL struct {
			IsRestrict     bool     `json:"is_restrict"`
			IsYoutubeAllow bool     `json:"is_youtube_allow"`
			IsTwitterAllow bool     `json:"is_twitter_allow"`
			IsGIFAllow     bool     `json:"is_gif_allow"`
			IsOpenseaAllow bool     `json:"is_opensea_allow"`
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

// サーバーの設定を更新します
func postServerConfig(c *gin.Context) {
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
	token := c.GetHeader("Token")

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
				message_send.SendErrMsg(session, errors.NewError("管理者ロールを持っていません"), guildName)
				c.JSON(http.StatusUnauthorized, "認証されていません")
				return
			}
		}
	}

	var (
		apiRes = v1.Res{}
	)

	bffErr := (func() error {
		req := &ReqConfig{}
		if err = c.BindJSON(req); err != nil {
			return errors.NewError("リクエストをJSONにバインドできません", err)
		}

		apiReqBlocks := make([]v1.BlockReq, 0)

		for _, rb := range req.Block {
			apiBlockReq := v1.BlockReq{}
			apiBlockReq.Name = rb.Name
			apiBlockReq.Keyword = rb.Keyword
			apiBlockReq.Reply = rb.Reply
			apiBlockReq.MatchCondition = rb.MatchCondition
			apiBlockReq.IsRandom = rb.IsRandom
			apiBlockReq.IsEmbed = rb.IsEmbed

			apiReqBlocks = append(apiReqBlocks, apiBlockReq)
		}

		apiRuleReq := v1.URLRuleReq{
			IsRestrict:     req.Rule.URL.IsRestrict,
			IsYoutubeAllow: req.Rule.URL.IsYoutubeAllow,
			IsTwitterAllow: req.Rule.URL.IsTwitterAllow,
			IsGIFAllow:     req.Rule.URL.IsGIFAllow,
			IsOpenseaAllow: req.Rule.URL.IsOpenseaAllow,
			IsDiscordAllow: req.Rule.URL.IsDiscordAllow,
			AllowRoleID:    req.Rule.URL.AllowRoleID,
			AllowChannelID: req.Rule.URL.AllowChannelID,
			AlertChannelID: req.Rule.URL.AlertChannelID,
		}

		apiRes, err = v1.UpdateConfig(session, ctx, id, req.AdminRoleID, apiReqBlocks, apiRuleReq)
		if err != nil {
			return errors.NewError("設定を更新できません", err)
		}

		return nil
	})()

	if bffErr != nil {
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
		message_send.SendErrMsg(session, errors.NewError("コミットに失敗しました", txErr), guildName)
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

	res := Res{}
	res.ID = apiRes.ID
	res.AdminRoleID = apiRes.AdminRoleID
	res.Block = []resBlock{}
	res.ServerName = guildName
	res.AvatarURL = avatarURL
	res.Role = []resRole{}
	res.Channel = []resChannel{}
	res.Rule.URL.IsRestrict = apiRes.Rule.URL.IsRestrict
	res.Rule.URL.IsYoutubeAllow = apiRes.Rule.URL.IsYoutubeAllow
	res.Rule.URL.IsTwitterAllow = apiRes.Rule.URL.IsTwitterAllow
	res.Rule.URL.IsGIFAllow = apiRes.Rule.URL.IsGIFAllow
	res.Rule.URL.IsOpenseaAllow = apiRes.Rule.URL.IsOpenseaAllow
	res.Rule.URL.IsDiscordAllow = apiRes.Rule.URL.IsDiscordAllow
	res.Rule.URL.AllowRoleID = apiRes.Rule.URL.AllowRoleID
	res.Rule.URL.AllowChannelID = apiRes.Rule.URL.AllowChannelID
	res.Rule.URL.AlertChannelID = apiRes.Rule.URL.AlertChannelID

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
