package user

// ドメインモデルのデータベースから読み込み、データベースへの書き込み処理を提供します

// テーブルの設定です
type Config struct {
	TableName  string
	ColumnName *ColumnName
}

// テーブルのカラム名です
type ColumnName struct {
	ID      string
	Content string
	Created string
	Updated string
}

// 結果の行です
type Row struct {
	ID      string
	Content string
	Created string
	Updated string
}

var c = &Config{
	TableName: "users",
	ColumnName: &ColumnName{
		ID:      "id",
		Content: "content",
		Created: "created",
		Updated: "updated",
	},
}

// GetConfig テーブルの設定を返します
func GetConfig() *Config {
	return c
}

var Conf = GetConfig()
