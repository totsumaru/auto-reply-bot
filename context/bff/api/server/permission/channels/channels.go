package channels

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/initiate"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"net/http"
	"strconv"
	"strings"
)

// サーバーのチャンネルの権限を取得します
func ServerPermissionChannels(e *gin.Engine) {
	e.GET("/server/permission/channels", getServerPermissionChannels)
}

// レスポンスです
type Res struct {
	Role []resRole `json:"role"`
}

// ロールの権限のレスポンスです
type resRole struct {
	ID                string                   `json:"id"`
	Name              string                   `json:"name"`
	Color             string                   `json:"color"`
	GeneralPermission resGeneralRolePermission `json:"general_permission"`
	Channel           []resChannel             `json:"channel_permission"`
}

// 全体のロールの権限です
type resGeneralRolePermission struct {
	PermissionAdministrator       bool // 管理者
	PermissionManageServer        bool // サーバー管理
	PermissionManageChannels      bool // チャンネルの管理
	PermissionBanMembers          bool // メンバーをBAN
	PermissionKickMembers         bool // メンバーをキック
	PermissionModerateMembers     bool // メンバーをタイムアウト
	PermissionCreateInstantInvite bool // 招待を作成
}

// チャンネルのレスポンスです
type resChannel struct {
	ID         string                   `json:"id"`
	Name       string                   `json:"name"`
	Permission resChannelRolePermission `json:"permission"`
}

// チャンネルごとのロールの権限です
type resChannelRolePermission struct {
	CanViewChannel           bool `json:"can_view_channel"`             // チャンネルを見る
	CanSendMessages          bool `json:"can_send_messages"`            // メッセージを送信
	CanManageMessages        bool `json:"can_manage_messages"`          // メッセージの管理
	CanMentionEveryone       bool `json:"can_mention_everyone"`         // @everyone,@here,全てのロールにメンション
	CanEmbedLinks            bool `json:"can_embed_links"`              // 埋め込みリンク
	CanAttachFiles           bool `json:"can_attach_files"`             // ファイルを添付
	CanReadMessageHistory    bool `json:"can_read_message_history"`     // メッセージ履歴を読む
	CanUseSlashCommands      bool `json:"can_use_slash_commands"`       // アプリコマンドを使う
	CanUseExternalEmojis     bool `json:"can_use_external_emojis"`      // 外部の絵文字を使用する
	CanUseExternalStickers   bool `json:"can_use_external_stickers"`    // 外部のスタンプを使用する
	CanManageThreads         bool `json:"can_manage_threads"`           // スレッドの管理
	CanCreatePublicThreads   bool `json:"can_create_public_threads"`    // 公開スレッドの作成
	CanCreatePrivateThreads  bool `json:"can_create_private_threads"`   // プライベートスレッドの作成
	CanSendMessagesInThreads bool `json:"can_send_messages_in_threads"` // スレッドでメッセージを送信
	CanSendTTSMessages       bool `json:"can_send_tts_messages"`        // テキスト読み上げメッセージを送信する
}

// サーバーのチャンネルの権限を取得します
func getServerPermissionChannels(c *gin.Context) {
	session, err := initiate.CreateSession()
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("セッションを作成できません", err), "none")
		c.JSON(http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}

	id := c.Query("id")
	//token := c.GetHeader("Token")

	// リクエストを検証します
	//if id == "" || token == "" {
	//	c.JSON(http.StatusBadRequest, "リクエストが不正です")
	//	return
	//}

	// クエリパラメータに指定されたサーバーです
	guildName, err := guild.GetGuildName(session, id)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("ギルド名を取得できません", err), "")
		return
	}

	// 認証されているユーザーかを検証します
	//ok, err := shared.IsAuthorizedUser(ctx, session, id, token)
	//if err != nil {
	//	message_send.SendErrMsg(session, errors.NewError("認証されているかの確認に失敗しました", err), guildName)
	//	c.JSON(http.StatusUnauthorized, "認証されていません")
	//	return
	//}
	//if !ok {
	//	c.JSON(http.StatusUnauthorized, "認証されていません")
	//	return
	//}

	res, err := getAllPermission(session, id)
	if err != nil {
		message_send.SendErrMsg(session, errors.NewError("Permissionの取得に失敗しました", err), guildName)
		c.JSON(http.StatusInternalServerError, "エラーが発生しました")
		return
	}

	c.JSON(http.StatusOK, res)
}

// 全てのチャンネル情報を取得します
func getAllPermission(s *discordgo.Session, guildID string) (Res, error) {
	res := Res{
		[]resRole{},
	}

	// 全てのチャンネルを取得します
	allTextChannels, err := guild.GetAllTextChannelsByRow(s, guildID)
	if err != nil {
		return res, errors.NewError("全てのテキストチャンネルを取得できません", err)
	}

	// 全てのロールを取得します
	// @everyoneも含みます
	allRoles, err := guild.GetAllRolesByRow(s, guildID)
	if err != nil {
		return res, errors.NewError("全てのロールを取得できません", err)
	}

	// 1ロールずつ処理します
	for _, rawRole := range allRoles {
		rr := resRole{}

		// ID,Name,Colorを設定します
		{
			rr.ID = rawRole.ID
			rr.Name = rawRole.Name
			// Colorを0x...から#...に置換します
			rr.Color = strings.Replace(
				strconv.Itoa(rawRole.Color), "0x", "#", -1,
			)
		}

		// 全体の権限を設定します
		rr.GeneralPermission.PermissionAdministrator = hasPermission(rawRole, discordgo.PermissionAdministrator)
		rr.GeneralPermission.PermissionManageServer = hasPermission(rawRole, discordgo.PermissionManageServer)
		rr.GeneralPermission.PermissionManageChannels = hasPermission(rawRole, discordgo.PermissionManageChannels)
		rr.GeneralPermission.PermissionBanMembers = hasPermission(rawRole, discordgo.PermissionBanMembers)
		rr.GeneralPermission.PermissionKickMembers = hasPermission(rawRole, discordgo.PermissionKickMembers)
		rr.GeneralPermission.PermissionModerateMembers = hasPermission(rawRole, discordgo.PermissionModerateMembers)
		rr.GeneralPermission.PermissionCreateInstantInvite = hasPermission(rawRole, discordgo.PermissionCreateInstantInvite)

		// チャンネルに関する権限のデフォルトを作成します
		defaultChannelRolePermission := resChannelRolePermission{
			CanViewChannel:           hasPermission(rawRole, discordgo.PermissionViewChannel),
			CanSendMessages:          hasPermission(rawRole, discordgo.PermissionSendMessages),
			CanManageMessages:        hasPermission(rawRole, discordgo.PermissionManageMessages),
			CanMentionEveryone:       hasPermission(rawRole, discordgo.PermissionMentionEveryone),
			CanEmbedLinks:            hasPermission(rawRole, discordgo.PermissionEmbedLinks),
			CanAttachFiles:           hasPermission(rawRole, discordgo.PermissionAttachFiles),
			CanReadMessageHistory:    hasPermission(rawRole, discordgo.PermissionReadMessageHistory),
			CanUseSlashCommands:      hasPermission(rawRole, discordgo.PermissionUseSlashCommands),
			CanUseExternalEmojis:     hasPermission(rawRole, discordgo.PermissionUseExternalEmojis),
			CanUseExternalStickers:   hasPermission(rawRole, discordgo.PermissionUseExternalStickers),
			CanManageThreads:         hasPermission(rawRole, discordgo.PermissionManageThreads),
			CanCreatePublicThreads:   hasPermission(rawRole, discordgo.PermissionCreatePublicThreads),
			CanCreatePrivateThreads:  hasPermission(rawRole, discordgo.PermissionCreatePrivateThreads),
			CanSendMessagesInThreads: hasPermission(rawRole, discordgo.PermissionSendMessagesInThreads),
			CanSendTTSMessages:       hasPermission(rawRole, discordgo.PermissionSendTTSMessages),
		}

		// チャンネルを1つずつ処理します
		for _, channel := range allTextChannels {
			// 権限はまずはデフォルトを設定します
			rc := resChannel{
				ID:         channel.ID,
				Name:       channel.Name,
				Permission: defaultChannelRolePermission,
			}

			// チャンネルで独自に設定されているロールがあれば、権限を上書きします
			for _, permission := range channel.PermissionOverwrites {
				// 上書きされているロールIDと処理中のロールIDが一致した場合は上書きします
				if permission.ID == rawRole.ID {
					overwritePermission := resChannelRolePermission{
						CanViewChannel:           hasPermissionByBasePermission(permission.Allow, discordgo.PermissionViewChannel),
						CanSendMessages:          hasPermissionByBasePermission(permission.Allow, discordgo.PermissionSendMessages),
						CanManageMessages:        hasPermissionByBasePermission(permission.Allow, discordgo.PermissionManageMessages),
						CanMentionEveryone:       hasPermissionByBasePermission(permission.Allow, discordgo.PermissionMentionEveryone),
						CanEmbedLinks:            hasPermissionByBasePermission(permission.Allow, discordgo.PermissionEmbedLinks),
						CanAttachFiles:           hasPermissionByBasePermission(permission.Allow, discordgo.PermissionAttachFiles),
						CanReadMessageHistory:    hasPermissionByBasePermission(permission.Allow, discordgo.PermissionReadMessageHistory),
						CanUseSlashCommands:      hasPermissionByBasePermission(permission.Allow, discordgo.PermissionUseSlashCommands),
						CanUseExternalEmojis:     hasPermissionByBasePermission(permission.Allow, discordgo.PermissionUseExternalEmojis),
						CanUseExternalStickers:   hasPermissionByBasePermission(permission.Allow, discordgo.PermissionUseExternalStickers),
						CanManageThreads:         hasPermissionByBasePermission(permission.Allow, discordgo.PermissionManageThreads),
						CanCreatePublicThreads:   hasPermissionByBasePermission(permission.Allow, discordgo.PermissionCreatePublicThreads),
						CanCreatePrivateThreads:  hasPermissionByBasePermission(permission.Allow, discordgo.PermissionCreatePrivateThreads),
						CanSendMessagesInThreads: hasPermissionByBasePermission(permission.Allow, discordgo.PermissionSendMessagesInThreads),
						CanSendTTSMessages:       hasPermissionByBasePermission(permission.Allow, discordgo.PermissionSendTTSMessages),
					}

					rc.Permission = overwritePermission
				}
			}

			rr.Channel = append(rr.Channel, rc)
		}

		// レスポンスに追加します
		res.Role = append(res.Role, rr)
	}

	return res, nil
}

// ロールが指定の権限を持っているかを確認します
func hasPermission(targetRole *discordgo.Role, targetPermission int64) bool {
	return targetRole.Permissions&targetPermission == targetPermission
}

// ロールが指定のチャンネルで指定の権限を持っているかを確認します
func hasPermissionByBasePermission(basePermission int64, targetPermission int64) bool {
	return basePermission&targetPermission == targetPermission
}
