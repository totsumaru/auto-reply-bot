package now

import "time"

// 現在時刻を返します
//
// テストでラップするための関数です。
var Now = func() time.Time {
	return time.Now()
}

// 現在時刻のタイムスタンプを取得します
func GetNowTimeStamp() string {
	return Now().Format("2006-01-02T15:04:05+09:00")
}

// 指定した時間をタイムスタンプの形式に変換します
func ParseTimeToString(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+09:00")
}
