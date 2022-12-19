package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const (
	AlphaNumLetters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphaLetters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LCAlphaNumLetters = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// ランダムなstring型の値を返します
func RandStr(arg ...interface{}) string {
	length := 32
	letters := AlphaNumLetters

	switch len(arg) {
	case 0:
	case 1:
		a0, ok := arg[0].(int)
		if !ok {
			panic("引数が不正")
		}

		length = a0
	case 2:
		a0, ok := arg[0].(int)
		if !ok {
			panic("引数が不正")
		}

		length = a0

		a1, ok := arg[1].(string)
		if !ok {
			panic("引数が不正")
		}

		letters = a1
	}

	var (
		result string
	)

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic("ランダムな文字列を作成できなかった")
	}

	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}

	return result
}

// ランダムなbool型の値を返します
func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		return true
	}

	return false
}

// ランダムなint型の値を返します
//
// `arg`で上限値を指定できます。
func RandInt(arg ...interface{}) int {
	rand.Seed(time.Now().UnixNano())

	min := 0
	max := 0

	var ok bool
	switch len(arg) {
	case 0:
		return rand.Int()
	case 1:
		min, ok = arg[0].(int)
		if min == 0 {
			panic("最小値は0以上を指定してください")
		}

		if !ok {
			panic("引数をintに変換できません")
		}

		return rand.Intn(min)
	case 2:
		min, ok = arg[0].(int)
		if min == 0 {
			panic("最小値は0以上を指定してください")
		}
		if !ok {
			panic("引数をintに変換できません")
		}

		max, ok = arg[1].(int)
		if !ok {
			panic("引数をintに変換できません")
		}
		if min >= max {
			panic("最大値は最小値より大きい値を指定してください")
		}
		if !ok {
			panic("引数をintに変換できません")
		}
	}

	return rand.Intn(max-min) + min
}

// `s`のスライスの要素をランダムで返します
func RandStrSlice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	return s[rand.Intn(len(s))]
}

// 文字列からtime.Time構造体を作成します
func StrToTime(t *testing.T, v string) time.Time {
	l, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal("予期しないエラーが発生した", err)
	}
	if v == "" {
		return time.Time{}
	}
	ti, err := time.ParseInLocation("2006-01-02 15:04:05", v, l)
	if err != nil {
		t.Fatal("予期しないエラーが発生した", err)
	}
	return ti
}

// メッセージにケースをつけます
func AddCaseInfo(m string, n int) string {
	return fmt.Sprintf(m+" (case: %d)", n)
}
