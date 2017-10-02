package format

import (
	"fmt"
	"strconv"

	"time"
)

func MakePromptPrefix(prefix string, prefixColor int) string {
	return "\033[" + strconv.Itoa(prefixColor) + "m" + prefix + "\033[0m"
}

func Msg(senderName, text string, timestamp time.Time, color int) string {
	name := MakePromptPrefix(senderName, color)
	time := MakePromptPrefix(timestamp.Format("15:04:05"), LIGHT_YELLOW)
	return fmt.Sprintf("(%s) %s: %s", time, name, text)
}

const (
	BLACK        = 30
	RED          = 31
	GREEN        = 32
	BLUE         = 34
	LIGHT_GRAY   = 37
	LIGHT_RED    = 91
	LIGHT_GREEN  = 92
	LIGHT_YELLOW = 93
	LIGHT_BLUE   = 94
	LIGHT_CYAN   = 96
	WHITE        = 97
)
