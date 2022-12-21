package discord

import (
	"encoding/json"
	"fmt"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/shared/map/seeker"
	"log"
	"net/http"
)

// tokenからDiscordIDに変換します
func TokenToDiscordID(token string) (string, error) {
	req, _ := http.NewRequest(
		http.MethodGet,
		"https://discordapp.com/api/users/@me",
		nil,
	)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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

	discordID := seeker.Str(res, []string{"id"})

	return discordID, nil
}
