package convert

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/redirect"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/shared/map/seeker"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// codeからTokenを取得します
func CodeToToken(code string) (string, error) {
	values := url.Values{}
	values.Set("client_id", os.Getenv("DISCORD_CLIENT_ID"))
	values.Add("client_secret", os.Getenv("DISCORD_CLIENT_SECRET"))
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("redirect_uri", redirect.CreateRedirectURL())

	req, err := http.NewRequest(
		http.MethodPost,
		"https://discordapp.com/api/oauth2/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", errors.NewError("httpリクエストの作成に失敗しました", err)
	}

	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.NewError("httpリクエストの送信に失敗しました", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("BodyをCloseできません")
		}
	}()

	res := map[string]interface{}{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", errors.NewError("レスポンスbodyのデコードに失敗しました", err)
	}

	token := seeker.Str(res, []string{"access_token"})
	if token == "" {
		return "", errors.NewError("codeでtokenを取得できません")
	}

	return token, nil
}
