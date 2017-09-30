package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"errors"
	"regexp"
	"strconv"

	"github.com/LorisFriedel/go-chat/console"
	"github.com/LorisFriedel/go-chat/core"
	rl "github.com/chzyer/readline"
	log "github.com/sirupsen/logrus"
)

var prefix = "!"

type Arguments struct {
	userName    string
	newChannels []*ChanArgs
	serverMode  bool
	chanToGo    *ChanArgs
	logLevel    string
}

type ChanArgs struct {
	name     string
	address  string
	port     int
	password string
	timeout  int
}

var argUserName string
var argNewChannels string
var argServerMode bool
var argChanToGo string
var argLogLevel string

func init() {
	flag.StringVar(&argUserName, "username", "", "Username displayed when you send message on a channel.")
	flag.StringVar(&argNewChannels, "new", "", "Channels to be created on startup, separated by semi-colons and formatted as follow : \"(name,ip,port,passwd[,timeout])\"\n e.g. -new \"(channel1,127.0.0.1,8080,PaSsWoRd1324,40);(channel2,192.21.58.11,8090,qwerty123)\"")
	flag.BoolVar(&argServerMode, "servermode", false, "If specified, Go-chat will only create new channels. Username argument will be ignored.")
	flag.StringVar(&argChanToGo, "go", "", "Channel to be joined on startup, formatted as follow : \"(name,ip,port,passwd)\"\n e.g. -go \"(channel1,127.0.0.1,8080,PaSsWoRd1324)\"")
	flag.StringVar(&argLogLevel, "loglevel", "", "Define log level. Possible values : panic,fatal,error,warn,info,debug")
}

func main() {
	flag.Parse()
	argsCli := parseCli()
	argsEnvVar := parsEnvVar()
	args := merge(argsEnvVar, argsCli)

	setLogLevel(args.logLevel)

	// TODO isoler le client pour qu'un channel puisse lancer un client et qu'une personne se connectant en netcat
	// TODO puisse l'utiliser comme client distant
	// TODO proprifier tout ce bordel

	if args.serverMode {
		initServerMode(args)
	} else {
		initClientMode(args)
	}
}

func setLogLevel(logLevelStr string) {
	var logLevel log.Level
	var err error

	if logLevelStr == "" {
		logLevel = log.PanicLevel
	} else {
		logLevel, err = log.ParseLevel(logLevelStr)
		if err != nil {
			log.Warnln(err)
		}
	}

	log.SetLevel(logLevel)
}

func initClientMode(args *Arguments) {
	if args.userName == "" {
		args.userName = getUserName()
	}

	fmt.Printf("Hi %s!\n", args.userName)

	client := core.NewClient(args.userName)
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
		} else if msg.Type == core.SYS_CHANNEL {
			color = console.LIGHT_RED
		} else if msg.Type == core.SYS_CLIENT {
			color = console.RED
		}

		name := console.MakePromptPrefix(msg.Sender.Name, color)
		time := console.MakePromptPrefix(msg.Timestamp.Format("15:04:05"), console.LIGHT_YELLOW)
		switch msg.Type {
		case core.TEXT:
			fmt.Printf("(%s) %s: %s\n", time, name, msg.Text)
		case core.SYS_CHANNEL:
			fmt.Printf("(%s) %s: %s\n", time, name, msg.Text)
		case core.SYS_CLIENT:
			fmt.Printf("(%s) %s\n", time, msg.Text)
		}
		rd.Refresh()
	})

	for _, c := range args.newChannels {
		client.CreateChan(c.name, c.address, c.port, c.password, c.timeout)
	}

	if c := args.chanToGo; c != nil {
		client.Connect(c.name, c.address, c.port, c.password)
	}

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

func initServerMode(args *Arguments) {
	args.userName = "Go-chat"
	fmt.Printf("All system operational.\n")

	server := core.NewClient(args.userName)

	for _, c := range args.newChannels {
		server.CreateChan(c.name, c.address, c.port, c.password, c.timeout)
	}
}

func parsEnvVar() *Arguments {
	return &Arguments{
		userName:    os.Getenv("GO_USERNAME"),
		newChannels: parseChannels(os.Getenv("GO_NEW")),
		serverMode:  parseBool(os.Getenv("GO_SERVER_MODE")),
		chanToGo:    parseChannel(os.Getenv("GO_GO")),
		logLevel:    os.Getenv("GO_LOG_LEVEL"),
	}
}

func parseBool(str string) bool {
	if b, err := strconv.ParseBool(str); err != nil {
		return b
	}
	return false
}

func parseCli() *Arguments {
	return &Arguments{
		userName:    argUserName,
		newChannels: parseChannels(argNewChannels),
		serverMode:  argServerMode,
		chanToGo:    parseChannel(argChanToGo),
		logLevel:    argLogLevel,
	}
}

func parseChannels(chanArgsStrList string) []*ChanArgs {
	results := make([]*ChanArgs, 0, 1)

	if chanArgsStrList == "" {
		return results
	}

	split := strings.Split(chanArgsStrList, ";")
	if len(split) > 1 {
		for _, chanArgsStr := range split {
			if c := parseChannel(chanArgsStr); c != nil {
				results = append(results, c)
			}
		}
	} else {
		results = append(results, parseChannel(split[0]))
	}

	return results
}

func parseChannel(chanArgsStr string) *ChanArgs {
	if chanArgsStr == "" {
		return nil
	}

	re, err := regexp.Compile(`\(([^,]+),([^,]+),([^,]+),([^,]+)(?:,([0-9]+))?\)`)
	if err != nil {
		log.Panicln(err)
	}

	matches := re.FindStringSubmatch(chanArgsStr)
	if matches == nil {
		log.Panicln(errors.New("error while parsing channel arguments"))
	}

	result := &ChanArgs{}
	result.name = matches[1]
	result.address = matches[2]
	if port, err := strconv.Atoi(matches[3]); err == nil {
		result.port = port
	}
	result.password = matches[4]
	if len(matches) > 4 {
		if timeout, err := strconv.Atoi(matches[5]); err == nil {
			result.timeout = timeout
		}
	}

	return result
}

func merge(argsList ...*Arguments) *Arguments {
	result := &Arguments{}
	for _, args := range argsList {
		if args.userName != "" {
			result.userName = args.userName
		}

		result.newChannels = append(result.newChannels, args.newChannels...)

		result.serverMode = args.serverMode

		if args.chanToGo != nil {
			result.chanToGo = args.chanToGo
		}

		if args.logLevel != "" {
			result.logLevel = args.logLevel
		}
	}
	return result
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
