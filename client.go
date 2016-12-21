package main

import (
	"fmt"
	"github.com/LorisFriedel/go-chat/reader"
	rl "github.com/chzyer/readline"
	"io"
	"strings"
)

func main() {
	// TODO instantiate a client and listen for command

	fmt.Println("Welcome stranger ! What's your name ?") // TODO
	// get user name (while loop until not only whitespace)

	completer := rl.NewPrefixCompleter(
		rl.PcItem("!go",
			rl.PcItem("chan1"),
			rl.PcItem("chan2"), // TODO dynamique channel, or let the user type it's addresse, name etc..
		),
		rl.PcItem("!join"), // TODO same as above ? useless ?
		rl.PcItem("!bye"),  // TODO plus (optionnal) the name of the chan to exit (autocomplete here too)
		rl.PcItem("!help"), // TODO + command name for help of it OR empty for general help
		rl.PcItem("!chan",
			rl.PcItem("list"),   // TODO list all registered channel that I can connect on (need to store password)
			rl.PcItem("status"), // TODO print current status of current channel, or the given one
			rl.PcItem("create"), // TODO create new channel
			rl.PcItem("leave"),  // TODO same as !bye ? useless ?
			rl.PcItem("join"),   // TODO same as go ? useless ?
			rl.PcItem("passwd"), // TODO + new password (error if you are not the owner)
		),
		rl.PcItem("!me"),   // TODO display info about me, what channel i'm on, etc..
		rl.PcItem("!ping"), // TODO ping current channel, nothing if not in channel
	)

	rd, err := reader.Builder.
		Prefix("> ").
		PrefixColor(reader.LIGHT_CYAN).
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
		fmt.Println("PRINT LINE: ", line)
	}
}
