package console

import (
	"strconv"

	"github.com/chzyer/readline"
	"github.com/lann/builder"
)

type rBuilder builder.Builder

var ReaderBuilder = builder.Register(rBuilder{}, readline.Instance{}).(rBuilder)

func (b rBuilder) Prefix(prefix string) rBuilder {
	return builder.Set(b, "prefix", prefix).(rBuilder)
}

func (b rBuilder) PrefixColor(color int) rBuilder {
	return builder.Set(b, "prefixColor", color).(rBuilder)
}

func (b rBuilder) InterruptCommand(command string) rBuilder {
	return builder.Set(b, "interruptCommand", command).(rBuilder)
}

func (b rBuilder) HistoryFile(path string) rBuilder {
	return builder.Set(b, "historyFile", path).(rBuilder)
}

func (b rBuilder) Completer(completer readline.AutoCompleter) rBuilder {
	return builder.Set(b, "completer", completer).(rBuilder)
}

func (b rBuilder) HistorySearchFold(search bool) rBuilder {
	return builder.Set(b, "historySearchFold", search).(rBuilder)
}

func (b rBuilder) Build() (*readline.Instance, error) {
	return newReader(b)
}

func newReader(rb rBuilder) (*readline.Instance, error) {
	var (
		prefix            string                 = "> "
		prefixColor       int                    = WHITE
		historyFile       string                 = "/tmp/go-reader.tmp"
		completer         readline.AutoCompleter = nil
		interruptCommand  string                 = "^C"
		historySearchFold bool                   = true
	)

	if val, set := builder.Get(rb, "prefix"); set {
		prefix = val.(string)
	}

	if val, set := builder.Get(rb, "prefixColor"); set {
		prefixColor = val.(int)
	}

	if val, set := builder.Get(rb, "historyFile"); set {
		historyFile = val.(string)
	}

	if val, set := builder.Get(rb, "completer"); set {
		completer = val.(readline.AutoCompleter)
	}

	if val, set := builder.Get(rb, "interruptCommand"); set {
		interruptCommand = val.(string)
	}

	if val, set := builder.Get(rb, "historySearchFold"); set {
		historySearchFold = val.(bool)
	}

	return readline.NewEx(&readline.Config{
		Prompt:            MakePromptPrefix(prefix, prefixColor),
		HistoryFile:       historyFile,
		AutoComplete:      completer,
		InterruptPrompt:   interruptCommand,
		EOFPrompt:         "exit",
		HistorySearchFold: historySearchFold,
	})
}

func MakePromptPrefix(prefix string, prefixColor int) string {
	return "\033[" + strconv.Itoa(prefixColor) + "m" + prefix + "\033[0m"
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

func MakeItem(prefix string, name string, pc ...readline.PrefixCompleterInterface) *readline.PrefixCompleter {
	return readline.PcItem(prefix+name, pc...)
}
