package v1

import (
	"github.com/techstart35/auto-reply-bot/context/shared/map/seeker"
)

// レスポンス構造体です
type Res struct {
	ID   string
	Name string
}

// レスポンスを作成します
func CreateRes(m map[string]interface{}) (Res, error) {
	res := Res{}

	res.ID = seeker.Str(m, []string{"id", "value"})
	res.Name = seeker.Str(m, []string{"name", "value"})

	return res, nil
}
