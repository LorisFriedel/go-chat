package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/LorisFriedel/go-chat/console"
	"github.com/LorisFriedel/go-chat/core"
	rl "github.com/chzyer/readline"
	log "github.com/sirupsen/logrus"
)

var prefix = "!"

// TODO add flag to start a channel from command line (with parameters)

func main() {
	flag.Parse() // glog need that
	log.SetLevel(log.PanicLevel)

	userName := getUserName()
	fmt.Printf("Hi %s!\n", userName)

	client := core.NewClient(userName)
	parser := console.NewParser(prefix)
	handler := console.NewHandler(parser)

	completer := rl.NewPrefixCompleter(
		makeItem(prefix, "go", rl.PcItemDynamic(client.ListKnownChan())),
		makeItem(prefix, "bye"),
		makeItem(prefix, "die"),
		// makeItem(prefix, "help"), // TODO + command name for help of it OR empty for general help (DYNAMIC)
		makeItem(prefix, "list"),
		// makeItem(prefix, "status"), // TODO status of the current channel (error if you are not the owner)
		makeItem(prefix, "new"),
		makeItem(prefix, "forget", rl.PcItemDynamic(client.ListKnownChan())),
		makeItem(prefix, "close", rl.PcItemDynamic(client.ListOwnChan())),
		makeItem(prefix, "me"),
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
		rd.Refresh() // TODO fix display issue
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
			if err == core.ClientSuicide {
				fmt.Println("Bye bye !")
				break
			}
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func makeItem(prefix string, name string, pc ...rl.PrefixCompleterInterface) *rl.PrefixCompleter {
	return rl.PcItem(prefix+name, pc...)
}

func getUserName() string {
	fmt.Println("Welcome stranger! What's your name?")
	var name string
	var err error

	done := false
	for !done {
		fmt.Print("> ")
		name, err = bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil || len(strings.TrimSpace(name)) == 0 {
			fmt.Println("Hmmm, tell me your name again?")
		} else {
			done = true
		}
	}

	return strings.TrimSpace(name)
}
