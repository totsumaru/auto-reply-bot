package db

import (
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"os"
)

// 環境変数の構造体です
type Config struct {
	EnvKeyDBHost         string
	EnvKeyDBPort         string
	EnvKeyDBName         string
	EnvKeyDBUserName     string
	EnvKeyDBUserPassword string
}

// 環境変数からデータベースの設定用の構造体を新規作成します
func NewConf() (c *Config, err error) {
	c = &Config{}

	c.EnvKeyDBHost = os.Getenv("DB_HOST")
	if c.EnvKeyDBHost == "" {
		return nil, errors.NewError("DB_HOSTの環境変数が設定されていません")
	}

	c.EnvKeyDBPort = os.Getenv("DB_PORT")
	if c.EnvKeyDBPort == "" {
		return nil, errors.NewError("DB_PORTの環境変数が設定されていません")
	}

	c.EnvKeyDBName = os.Getenv("DB_NAME")
	if c.EnvKeyDBName == "" {
		return nil, errors.NewError("DB_NAMEの環境変数が設定されていません")
	}

	c.EnvKeyDBUserName = os.Getenv("DB_USER_NAME")
	if c.EnvKeyDBUserName == "" {
		return nil, errors.NewError("DB_USER_NAMEの環境変数が設定されていません")
	}

	c.EnvKeyDBUserPassword = os.Getenv("DB_USER_PASSWORD")
	if c.EnvKeyDBUserPassword == "" {
		return nil, errors.NewError("DB_USER_PASSWORDの環境変数が設定されていません")
	}

	return
}
