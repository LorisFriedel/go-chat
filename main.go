package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/LorisFriedel/go-chat/console"
	"github.com/LorisFriedel/go-chat/core"
	rl "github.com/chzyer/readline"
	"io"
	"log"
	"os"
	"strings"
)

var prefix = "!"

func main() {
	flag.Parse() // glog need that

	userName := getUserName()
	fmt.Printf("Hi %s !\n", userName)

	client := core.NewClient(userName)
	parser := console.NewParser(prefix)
	handler := console.NewHandler(parser)

	completer := rl.NewPrefixCompleter(
		makeItem(prefix, "go",
			rl.PcItem("chan1"),
			rl.PcItem("chan2"), // TODO DYNAMIC channel, or let the user type it's addresse, name etc..
		),
		makeItem(prefix, "bye"),    // TODO plus (optionnal) the name of the chan to exit (DYNAMIC)
		makeItem(prefix, "help"),   // TODO + command name for help of it OR empty for general help (DYNAMIC)
		makeItem(prefix, "list"),   // TODO list all registered channel that I can connect on (need to store password)
		makeItem(prefix, "status"), // TODO print current status of current channel, or the given one
		makeItem(prefix, "new"),    // TODO create new channel
		makeItem(prefix, "forget"), // TODO delete known channel (DYNAMIC)
		makeItem(prefix, "delete"), // TODO delete own channel (DYNAMIC)
		makeItem(prefix, "passwd"), // TODO + new password (error if you are not the owner)
		makeItem(prefix, "me"),     // TODO display info about me, what channel i'm on, etc..
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

	client.AddListener(func(msg core.Message) {
		color := console.LIGHT_BLUE
		if msg.Sender == client.Identity() {
			color = console.LIGHT_GREEN
		}

		name := console.MakePromptPrefix(msg.Sender.Name, color)
		time := console.MakePromptPrefix(msg.Timestamp.Format("15:04:05"), console.LIGHT_YELLOW)
		fmt.Printf("(%s) %s: %s\n", time, name, msg.Text)
		rd.Refresh()
	})

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
		if len(line) == 0 {
			continue
		}

		err = handler.Handle(client, line)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func makeItem(prefix string, name string, pc ...rl.PrefixCompleterInterface) *rl.PrefixCompleter {
	return rl.PcItem(prefix+name, pc...)
}

func getUserName() string {
	fmt.Println("Welcome stranger ! What's your name ?")
	fmt.Print("> ")
	name, err := bufio.NewReader(os.Stdin).ReadString('\n')

	if err != nil {
		log.Fatal("Oh noooo! Invalid user name, bye bye :(")
	}

	return strings.Trim(name, "\n")
}