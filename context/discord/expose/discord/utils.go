package discord

import "fmt"

// IDをメンションの形式に変換します
func IDToMention(id string) string {
	return fmt.Sprintf("<@%s>", id)
}
