package shared

import (
	"fmt"
	"os"
)

// DiscordログインのリダイレクトURLを作成します
func CreateDiscordLoginRedirectURL(serverID string) string {
	return fmt.Sprintf(
		"%s?id=%s",
		os.Getenv("FE_ROOT_URL"),
		serverID,
	)
}
