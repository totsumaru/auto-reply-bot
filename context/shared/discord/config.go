package discord

import (
	"fmt"
)

const (
	// エラーメッセージの送信先チャンネルIDです
	ErrMsgChannelID = "1036770181118181386" // TEST SERVERの`scam-alert-log`チャンネル
)

// 埋め込みメッセージのカラーコードです
const (
	ColorBlue   = 0x0099ff
	ColorRed    = 0xff0000
	ColorOrange = 0xffa500
	ColorGreen  = 0x3cb371
	ColorPink   = 0xff69b4
	ColorBlack  = 0x000000
	ColorYellow = 0xffd700
	ColorCyan   = 0x00ffff
)

// IDをメンションの形式に変換します
func IDToMention(id string) string {
	return fmt.Sprintf("<@%s>", id)
}
