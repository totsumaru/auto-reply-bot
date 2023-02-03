package server

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/bff/shared"
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
			IsOpenSeaAllow bool     `json:"is_opensea_allow"`
			IsDiscordAllow bool     `json:"is_discord_allow"`
			AllowRoleID    []string `json:"allow_role_id"`
			AllowChannelID []string `json:"allow_channel_id"`
		} `json:"url"`
	} `json:"rule"`
	// 以下はComputedです
	Token      string       `json:"token"`
	ServerName string       `json:"server_name"`
	AvatarURL  string       `json:"avatar_url"`
	Role       []resRole    `json:"role"`
	Channel    []resChannel `json:"channel"`
	Nickname   string       `json:"nickname"`
}

// ブロックのレスポンスです
type resBlock struct {
	Name             string   `json:"name"`
	Keyword          []string `json:"keyword"`
	Reply            []string `json:"reply"`
	MatchCondition   string   `json:"match_condition"`
	LimitedChannelID []string `json:"limited_channel_id"`
	IsRandom         bool     `json:"is_random"`
	IsEmbed          bool     `json:"is_embed"`
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

	// リクエストを検証します
	if id == "" || code == "" {
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
	token, err := convert.CodeToToken(code)
	if err != nil {
		// codeの不正に関してはエラー通知しません
		c.JSON(http.StatusUnauthorized, "codeが認証されていません")
		return
	}
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

	m, err := session.GuildMember(id, os.Getenv("DISCORD_APPLICATION_ID"))
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("botのMember情報を取得できません", err), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	res, err := createRes(session, apiRes, token, guildName, m.Nick)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("レスポンスを作成できません", err), guildName)
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	c.JSON(http.StatusOK, res)
}

// レスポンスを作成します
func createRes(
	session *discordgo.Session,
	apiRes v1.Res,
	token string,
	guildName string,
	nickName string,
) (Res, error) {
	res := Res{}

	avatarURL, err := guild.GetAvatarURL(session, apiRes.ID)
	if err != nil {
		return res, errors.NewError("アバターのURLを取得できません", err)
	}

	allRoles, err := guild.GetAllRolesWithoutEveryone(session, apiRes.ID)
	if err != nil {
		return res, errors.NewError("全てのロールを取得できません", err)
	}

	allChannels, err := guild.GetAllTextChannels(session, apiRes.ID)
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
	res.Rule.URL.IsOpenSeaAllow = apiRes.Rule.URL.IsOpenseaAllow
	res.Rule.URL.IsDiscordAllow = apiRes.Rule.URL.IsDiscordAllow
	res.Rule.URL.AllowRoleID = apiRes.Rule.URL.AllowRoleID
	res.Rule.URL.AllowChannelID = apiRes.Rule.URL.AllowChannelID
	// Computed
	res.Token = token
	res.ServerName = guildName
	res.AvatarURL = avatarURL
	res.Role = []resRole{}
	res.Channel = []resChannel{}
	res.Nickname = nickName

	// レスポンスにブロックを追加します
	for _, v := range apiRes.Comment.Block {
		blockRes := resBlock{}
		blockRes.Name = v.Name
		blockRes.Keyword = v.Keyword
		blockRes.Reply = v.Reply
		blockRes.MatchCondition = v.MatchCondition
		blockRes.LimitedChannelID = v.LimitedChannelID
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
