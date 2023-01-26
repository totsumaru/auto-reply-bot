package config

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
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
type ReqBody struct {
	AdminRoleID string `json:"admin_role_id"`
	Comment     struct {
		Block           []blockReq `json:"block"`
		IgnoreChannelID []string   `json:"ignore_channel_id"`
	} `json:"comment"`
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
		} `json:"url"`
	} `json:"rule"`
}

// ブロックのリクエストです
type blockReq struct {
	Name           string   `json:"name"`
	Keyword        []string `json:"keyword"`
	Reply          []string `json:"reply"`
	MatchCondition string   `json:"match_condition"`
	IsRandom       bool     `json:"is_random"`
	IsEmbed        bool     `json:"is_embed"`
}

// レスポンスです
type Res struct {
	ID          string `json:"id"`
	AdminRoleID string `json:"admin_role_id"`
	Comment     struct {
		Block           []resBlock `json:"block"`
		IgnoreChannelID []string   `json:"ignore_channel_id"`
	} `json:"comment"`
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
		} `json:"url"`
	} `json:"rule"`
	// 以下はComputedです
	ServerName string       `json:"server_name"`
	AvatarURL  string       `json:"avatar_url"`
	Role       []resRole    `json:"role"`
	Channel    []resChannel `json:"channel"`
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

	// リクエストを検証します
	if id == "" || token == "" {
		c.JSON(http.StatusBadRequest, "リクエストが不正です")
		return
	}

	// クエリパラメータに指定されたサーバーです
	guildName, err := guild.GetGuildName(session, id)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("ギルド名を取得できません", err), "")
		return
	}

	// 認証されているユーザーかを検証します
	ok, err := shared.IsAuthorizedUser(ctx, session, id, token)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("認証されているかの確認に失敗しました", err), guildName)
		c.JSON(http.StatusUnauthorized, "認証されていません")
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, "認証されていません")
		return
	}

	apiRes := v1.Res{}

	bffErr := (func() error {
		reqBody := &ReqBody{}
		if err = c.BindJSON(reqBody); err != nil {
			return errors.NewError("リクエストをJSONにバインドできません", err)
		}

		// BodyのリクエストからAPIのリクエストを作成します
		apiReq, err := castReqBodyToAPIReq(reqBody)
		if err != nil {
			return errors.NewError("BodyからAPIのリクエストを作成できません", err)
		}

		// APIをコールします
		apiRes, err = v1.UpdateConfig(session, ctx, id, apiReq)
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

	// レスポンスを作成します
	res, err := createRes(session, apiRes, guildName)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("レスポンスを作成できません", err), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	c.JSON(http.StatusOK, res)
}

// ServerコンテキストへのAPIリクエストを作成します
func castReqBodyToAPIReq(reqBody *ReqBody) (v1.Req, error) {
	// ブロックのリクエストを作成します
	apiReqBlocks := make([]v1.BlockReq, 0)
	for _, rb := range reqBody.Comment.Block {
		apiBlockReq := v1.BlockReq{}
		apiBlockReq.Name = rb.Name
		apiBlockReq.Keyword = rb.Keyword
		apiBlockReq.Reply = rb.Reply
		apiBlockReq.MatchCondition = rb.MatchCondition
		apiBlockReq.IsRandom = rb.IsRandom
		apiBlockReq.IsEmbed = rb.IsEmbed

		apiReqBlocks = append(apiReqBlocks, apiBlockReq)
	}

	apiReq := v1.Req{}
	apiReq.AdminRoleID = reqBody.AdminRoleID
	// Comment
	apiReq.Comment.BlockReq = apiReqBlocks
	apiReq.Comment.IgnoreChannelID = reqBody.Comment.IgnoreChannelID
	// Rule
	apiReq.Rule.URL.IsRestrict = reqBody.Rule.URL.IsRestrict
	apiReq.Rule.URL.IsYoutubeAllow = reqBody.Rule.URL.IsYoutubeAllow
	apiReq.Rule.URL.IsTwitterAllow = reqBody.Rule.URL.IsTwitterAllow
	apiReq.Rule.URL.IsGIFAllow = reqBody.Rule.URL.IsGIFAllow
	apiReq.Rule.URL.IsOpenseaAllow = reqBody.Rule.URL.IsOpenseaAllow
	apiReq.Rule.URL.IsDiscordAllow = reqBody.Rule.URL.IsDiscordAllow
	apiReq.Rule.URL.AllowRoleID = reqBody.Rule.URL.AllowRoleID
	apiReq.Rule.URL.AllowChannelID = reqBody.Rule.URL.AllowChannelID

	return apiReq, nil
}

// レスポンスを作成します
func createRes(s *discordgo.Session, apiRes v1.Res, guildName string) (Res, error) {
	res := Res{}

	avatarURL, err := guild.GetAvatarURL(s, apiRes.ID)
	if err != nil {
		return res, errors.NewError("アバターのURLを取得できません", err)
	}

	allRoles, err := guild.GetAllRolesWithoutEveryone(s, apiRes.ID)
	if err != nil {
		return res, errors.NewError("全てのロールを取得できません", err)
	}

	allChannels, err := guild.GetAllTextChannels(s, apiRes.ID)
	if err != nil {
		return res, errors.NewError("全てのチャンネルを取得できません", err)
	}

	res.ID = apiRes.ID
	res.AdminRoleID = apiRes.AdminRoleID
	// Comment
	res.Comment.Block = []resBlock{}
	res.Comment.IgnoreChannelID = apiRes.Comment.IgnoreChannelID
	// Rule
	res.Rule.URL.IsRestrict = apiRes.Rule.URL.IsRestrict
	res.Rule.URL.IsYoutubeAllow = apiRes.Rule.URL.IsYoutubeAllow
	res.Rule.URL.IsTwitterAllow = apiRes.Rule.URL.IsTwitterAllow
	res.Rule.URL.IsGIFAllow = apiRes.Rule.URL.IsGIFAllow
	res.Rule.URL.IsOpenseaAllow = apiRes.Rule.URL.IsOpenseaAllow
	res.Rule.URL.IsDiscordAllow = apiRes.Rule.URL.IsDiscordAllow
	res.Rule.URL.AllowRoleID = apiRes.Rule.URL.AllowRoleID
	res.Rule.URL.AllowChannelID = apiRes.Rule.URL.AllowChannelID
	// Computed
	res.ServerName = guildName
	res.AvatarURL = avatarURL
	res.Role = []resRole{}
	res.Channel = []resChannel{}

	// レスポンスにブロックを追加します
	for _, v := range apiRes.Comment.Block {
		blockRes := resBlock{}
		blockRes.Name = v.Name
		blockRes.Keyword = v.Keyword
		blockRes.Reply = v.Reply
		blockRes.MatchCondition = v.MatchCondition
		blockRes.IsRandom = v.IsRandom
		blockRes.IsEmbed = v.IsEmbed

		res.Comment.Block = append(res.Comment.Block, blockRes)
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

	return res, nil
}
