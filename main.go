package main

import (
	"flag"
	"os"
	"strings"

	"errors"
	"regexp"
	"strconv"

	"fmt"
	"io"

	"bufio"

	"github.com/LorisFriedel/go-chat/console"
	"github.com/LorisFriedel/go-chat/core"
	"github.com/LorisFriedel/go-chat/format"
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
	flag.StringVar(&argNewChannels, "new", "", "Channels to be created on startup, separated by semi-colons and formatted as follow : \"(name,ip,port,passwd[,timeout])\"\n e.g. -new \"(channel1,127.0.0.1,8080,myPassword,40);(channel2,192.21.58.11,8090,qwerty123)\"")
	flag.BoolVar(&argServerMode, "servermode", false, "If specified, Go-chat will only create new channels. Username argument will be ignored.")
	flag.StringVar(&argChanToGo, "go", "", "Channel to be joined on startup, formatted as follow : \"(name,ip,port,passwd)\"\n e.g. -go \"(channel1,127.0.0.1,8080,myPassword)\"")
	flag.StringVar(&argLogLevel, "loglevel", "", "Define log level. Possible values : panic,fatal,error,warn,info,debug")
}

/*
PARSE INPUT (ENV VAR and CLI)
START corresponding mode
*/

func main() {
	// Parse arguments from everywhere
	flag.Parse()
	argsCli := parseCli()
	argsEnvVar := parsEnvVar()
	args := merge(argsEnvVar, argsCli)

	// Start the chat
	if args.serverMode {
		initServerMode(args)
	} else {
		initStandardMode(args)
	}
}

func initServerMode(args *Arguments) {
	setLogLevel(args.logLevel, log.DebugLevel)

	args.userName = ""
	fmt.Printf("All system operational.\n")

	server := core.NewClient(args.userName)
	parser := console.NewParser(prefix)
	handler := console.NewHandler(parser)

	completer := rl.NewPrefixCompleter(
		console.MakeItem(prefix, "die"),
		// makeItem(prefix, "help"), // TODO + command name for help of it OR empty for general help (DYNAMIC)
		console.MakeItem(prefix, "list"),
		// makeItem(prefix, "status"), // TODO status of the current channel (error if you are not the owner?)
		console.MakeItem(prefix, "new"),
		console.MakeItem(prefix, "close", rl.PcItemDynamic(server.ListOwnChan())),
	)

	rd, err := console.ReaderBuilder.
		Prefix("# ").PrefixColor(format.LIGHT_RED).
		HistoryFile("/tmp/go-chat").HistorySearchFold(true).
		InterruptCommand("^C").Completer(completer).Build()

	if err != nil {
		log.Panic(err)
	}
	defer rd.Close()

	// Post boot actions -----------------------------------------------------
	for _, c := range args.newChannels {
		server.CreateChan(c.name, c.address, c.port, c.password, c.timeout)
	}
	// END Post boot actions -------------------------------------------------

	cliLoop(server, rd, handler)
}

func initStandardMode(args *Arguments) {
	setLogLevel(args.logLevel, log.PanicLevel)

	if args.userName == "" {
		args.userName = getUserName()
	}

	fmt.Printf("Hi %s!\n", args.userName)

	client := core.NewClient(args.userName)
	parser := console.NewParser(prefix)
	handler := console.NewHandler(parser)

	completer := rl.NewPrefixCompleter(
		console.MakeItem(prefix, "go", rl.PcItemDynamic(client.ListKnownChan())),
		console.MakeItem(prefix, "bye"),
		console.MakeItem(prefix, "die"),
		// makeItem(prefix, "help"), // TODO + command name for help of it OR empty for general help (DYNAMIC)
		console.MakeItem(prefix, "list"),
		// makeItem(prefix, "status"), // TODO status of the current channel (error if you are not the owner)
		console.MakeItem(prefix, "new"),
		console.MakeItem(prefix, "forget", rl.PcItemDynamic(client.ListKnownChan())),
		console.MakeItem(prefix, "close", rl.PcItemDynamic(client.ListOwnChan())),
		console.MakeItem(prefix, "me"),
	)

	rd, err := console.ReaderBuilder.
		Prefix("> ").PrefixColor(format.LIGHT_CYAN).
		HistoryFile("/tmp/go-chat").HistorySearchFold(true).
		InterruptCommand("^C").Completer(completer).Build()

	if err != nil {
		log.Panic(err)
	}
	defer rd.Close()

	client.AddListener(func(msg *core.Message) {
		color := colorMap[msg.Type]
		if msg.Sender == client.Identity() {
			color = format.LIGHT_GREEN
		}

		fmt.Println(format.Msg(msg.Sender.Name, msg.Text, msg.Timestamp, color))
		rd.Refresh()
	})

	// Post boot actions -----------------------------------------------------
	for _, c := range args.newChannels {
		client.CreateChan(c.name, c.address, c.port, c.password, c.timeout)
	}

	if c := args.chanToGo; c != nil {
		client.Connect(c.name, c.address, c.port, c.password)
	}
	// END Post boot actions -------------------------------------------------

	cliLoop(client, rd, handler)
}

var colorMap map[core.TMsg]int = map[core.TMsg]int{
	core.TEXT:        format.LIGHT_BLUE,
	core.SYS_CHANNEL: format.LIGHT_RED,
	core.SYS_CLIENT:  format.RED,
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

func cliLoop(client core.IClient, rd *rl.Instance, handler console.IHandler) {
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
				fmt.Println("Go-chat is now gonna die, you murderer!")
				break
			}
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func setLogLevel(logLevelStr string, defaultLevel log.Level) {
	var logLevel log.Level
	var err error

	if logLevelStr == "" {
		logLevel = defaultLevel
	} else {
		logLevel, err = log.ParseLevel(logLevelStr)
		if err != nil {
			log.Warnln(err)
		}
	}

	log.SetLevel(logLevel)
}

///////////////////////////////////////////
//////////////// Parsing //////////////////
///////////////////////////////////////////

func parsEnvVar() *Arguments {
	return &Arguments{
		userName:    os.Getenv("GO_USERNAME"),
		newChannels: parseChannels(os.Getenv("GO_NEW")),
		serverMode:  parseBool(os.Getenv("GO_SERVER_MODE")),
		chanToGo:    parseChannel(os.Getenv("GO_GO")),
		logLevel:    os.Getenv("GO_LOG_LEVEL"),
	}
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

func parseBool(str string) bool {
	if b, err := strconv.ParseBool(str); err != nil {
		return b
	}
	return false
}
