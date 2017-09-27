# Go-chat
Simple chat application written in Go.

## Available command:
If no command is specified, the written text is considered as a message to be sent on the current channel.
 + !go name \[address port password\]
    + Connect to a channel. If address nor port are mentioned, try to connect to a known channel. If the connection is successful, the channel is added to known channels. Be careful, known channels can be forgotten if the given name is an already known channel
 + !bye
    + Disconnect from current channel
 + !die
    + Quit chat
 + !list
    + List all known channels
 + !new name address port password
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
Bob start the chat and create a new channel :

~ ./run.sh

~ Welcome stranger ! What's your name ?

~ > Bob

~ Hi Bob!

~ !new MyAwesomeChannel 127.0.0.1 8080 myawesomepassword

~ (16:19:55) Bob: Now connected to MyAwesomeChannel (127.0.0.1:8080)

~ > (16:19:55) 127.0.0.1:8080: Bob joined the channel.


Then, John join the channel created by Bob :

~ > !go MyAwesomeChannel 127.0.0.1 8080 myawesomepassword

~ (16:20:10) John: Now connected to chan (127.0.0.1:8080)

~ > (16:20:10) 127.0.0.1:8080: John joined the channel.

~ > (16:20:13) Bob: Hello John! How are you?

## Build

To build Go-chat executable, simply run the 'build.sh' script.
You need Go (Golang) installed OR Docker installed and running.
However, even without Go installed, you need the GOPATH environment variable to be set.

## Build Docker image

To build the Go-chat Docker image, run the script 'img_build.sh'.
Warning: Go-chat Docker image is not working due to input reading issues.
Please feel free to fix it and send a pull request or wait until I fix it.

## Run

To run Go-chat, execute the 'run.sh' script.
Requirements here are the same as for the build section given that the exe has to be built if it's the first time you want to run Go-chat.
