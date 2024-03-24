package main

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

func addSpaceToEmoji(emoji string) string {
	if emoji == "" {
		return ""
	}

	emojiWidth := runewidth.StringWidth(emoji)
	spaces := strings.Repeat(" ", 3-emojiWidth)
	return fmt.Sprintf("%s%s", emoji, spaces)
}
