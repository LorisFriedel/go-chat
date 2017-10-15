# GoChat

![GoChat logo](./images/go-chat-small.png)

Another simple IRC application written in Go.

And, like the others, it is unique.

But why it is unique Jean-Pierre?

Because it support command line input, environment variable input and in-app input, because you can chat without running the client at all, because you can run in full server mode in case you wanna create 147894816545 channels of pure awesomeness, because you can build and run it with Docker like you've always wanted to do, and above all because ITS LOGO ROCKS.

## Available command:
If no command is specified, the written text is considered as a message to be sent on the current channel.
 + !go name \[address port \[password\]\]
    + Connect to a channel. If address nor port are mentioned, try to connect to a known channel. If the connection is successful, the channel is added to known channels. Be careful, known channels can be forgotten if the given name is an already known channel
 + !bye
    + Disconnect from current channel
 + !die
    + Quit chat
 + !list
    + List all known channels
 + !new name address port \[password\]
    + Create a new channel with given parameters. Client is automatically connected to it when created
 + !forget name
    + Forget a known channel
 + !close name
    + Close a channel (if the given name correspond to a owned channel)
 + !me
    + Display client status

Coming soon :

 + !help \[command\]
    + Display all available commands if no command name is specified, otherwise display help for that particular command.

And channel auto-save, to let you reconnect to known channels even after exiting Go-chat!

## Usage example

### Standard method

Bob start the Go-chat program and create a new channel :

~ ./run.sh

~ Welcome stranger ! What's your name ?

~ > McFly

~ Hi McFly!

~ !new MyAwesomeChannel 104.20.90.7 8080 myAwesomePassword

~ (16:19:55) McFly: Now connected to MyAwesomeChannel (104.20.90.7:8080)

~ > (16:19:55) 104.20.90.7:8080: McFly joined the channel.


Then, Carlito join the channel created by McFly :

~ > !go MyAwesomeChannel 104.20.90.7 8080 myAwesomePassword

~ (16:20:10) Carlito: Now connected to chan (104.20.90.7:8080)

~ > (16:20:10) 104.20.90.7:8080: Carlito joined the channel.

~ > (16:20:13) McFly: Hello John! How are you?

~ > (16:20:18) Carlito: WOW AMAZING THIS IS WORKING! SO SMOOTH!

And so on..

### Barbarian method

If Alicia want to join the chat but can't download the client for mysterious reasons, she can simply do as follow:

~ nc 104.20.90.7 8080

~ Welcome stranger ! What's your name ?

~ Alicia

~ Hi Alicia!

~ (16:21:32) Alicia: Now connected to chan (104.20.90.7:8080)

~ > (16:21:32) 104.20.90.7:8080: Alicia joined the channel.

And from here she can use Go-chat like she had the executable client. M.A.G.I.C.

## Build

To build Go-chat executable, simply run the 'build.sh' script.
You need Go (Golang) installed OR Docker installed and running.
However, even without Go installed, you need the GOPATH environment variable to be set.

## Build Docker image

To build the Go-chat Docker image, run the script 'img_build.sh'.

## Run

To run Go-chat, execute the 'run.sh' script.
Requirements here are the same as for the build section given that the exe has to be built if it's the first time you want to run Go-chat.

## Run with Docker

To run Go-chat Docker image, execute :

~ sudo docker run -it lorisfriedel/go-chat

## Execute command at start

If you want to create a new channel when running Go-chat without using in-app command, or join an existing channel, well.. you can!
In fact you can do whatever you want and however you want, like :
 + Docker + environment variable
 + Docker + command line arguments
 + Docker + in-app command
 + Docker-compose + environment variable + command line arguments
 + Run script + command line arguments
 + Run script + environment variable
 + Run script + in-app command
 + Executable + command line arguments
 + Executable + environment variable
 + Executable + in-app command
 + Eat a cookie and connect using the good old netcat without the intention to execute any command because you're a BADASS. Now that's BADASS!
 + This list is too long but obviously you don't care because, like the old chinese proverb sayin': tl,dr.

## Available options
This is a list of all available options that can be used to customize the way you run your Go-chat. For command line arguments (to pass when running Go-chat, with or without the script run.sh) or environment variables, the input format is the same. However, the command line arguments will always have priority over environments variables.
 + -username (cli args) | GO_USERNAME (ENV VAR)
    + Username displayed when you send message on a channel (empty string by default, asked when starting the chat)
 + -new (cli args) | GO_NEW (ENV VAR)
    + Channels to be created on startup (none by default), separated by semi-colons and formatted as follow (timeout in second): \"(name,ip,port,passwd\[,timeout\])\"
    + Example: ~ ./go-chat -new \"(channel1,127.0.0.1,8080,myPassword,40);(channel2,192.21.58.11,8090,qwerty123)\"
 + -servermode (cli args) | GO_SERVER_MODE (ENV VAR)
    + If specified, Go-chat will only create new channels. Username and Go arguments will be ignored. (false by default)
 + -go (cli args) | GO_GO (ENV VAR)
    + Channel to be joined on startup, formatted as follow : \"(name,ip,port,passwd)\"
    + Example: ~ ./go-chat -go \"(channel1,127.0.0.1,8080,myPassword)\"
 + -loglevel (cli args) | GO_LOG_LEVEL (ENV VAR)
    + Define log level. Possible values : panic,fatal,error,warn,info,debug