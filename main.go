package main

import (
	"fmt"
	"github.com/LorisFriedel/go-chat/console"
	"github.com/LorisFriedel/go-chat/core"
	rl "github.com/chzyer/readline"
	"io"
	"strings"
)

var prefix = "!"

func main() {
	fmt.Println("Welcome stranger ! What's your name ?") // TODO
	// TODO get user name (while loop until not only whitespace)

	client := core.NewClient("Loris")
	parser := console.NewCmdParser(prefix)
	handler := console.NewCmdHandler(parser)

	completer := rl.NewPrefixCompleter(
		makeItem(prefix, "go",
			rl.PcItem("chan1"),
			rl.PcItem("chan2"), // TODO dynamique channel, or let the user type it's addresse, name etc..
		),
		makeItem(prefix, "join"), // TODO same as above ? useless ?
		makeItem(prefix, "bye"),  // TODO plus (optionnal) the name of the chan to exit (autocomplete here too)
		makeItem(prefix, "help"), // TODO + command name for help of it OR empty for general help
		makeItem(prefix, "chan",
			rl.PcItem("list"),   // TODO list all registered channel that I can connect on (need to store password)
			rl.PcItem("status"), // TODO print current status of current channel, or the given one
			rl.PcItem("create"), // TODO create new channel
			rl.PcItem("leave"),  // TODO same as !bye ? useless ?
			rl.PcItem("join"),   // TODO same as go ? useless ?
			rl.PcItem("passwd"), // TODO + new password (error if you are not the owner)
		),
		makeItem(prefix, "me"),   // TODO display info about me, what channel i'm on, etc..
		makeItem(prefix, "ping"), // TODO ping current channel, nothing if not in channel
	)

	rd, err := console.ReaderBuilder.
		Prefix("> ").
		PrefixColor(console.LIGHT_CYAN).
		HistoryFile("/tmp/go-chat").
		Completer(completer).
		InterruptCommand("^C").
		HistorySearchFold(true).
		Build()

	if err != nil {
		panic(err)
	}
	defer rd.Close()

	for {
		// Read input
		line, err := rd.Readline()
		line = strings.TrimSpace(line)

		// Handle error
		if err == rl.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		// Handle input
		err = handler.Handle(client, line)
		if err != nil {
			fmt.Errorf("%v", err)
		}
	}
}

func makeItem(prefix string, name string, pc ...rl.PrefixCompleterInterface) *rl.PrefixCompleter {
	return rl.PcItem(prefix+name, pc...)
}
