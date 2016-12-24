# go-chat
Simple chat application written in Go

##Available command:
If no command is specified, the written text is considered as a message to be sent on the current channel. 
 + !go name (address port)
    + Connect to a channel. If address nor port are mentioned, try to connect to a known channel. If the connection is successful, the channel is added to known channels. Be careful, known channels can be forgotten if the given name is an already known channel
 + !bye
    + Disconnect from current channel
 + !die
    + Quit chat
 + !list
    + List all known channels
 + !new name address port passwd
    + Create a new channel with given parameters. Client is automatically connected to it when created
 + !forget name
    + Forget a known channel
 + !close name
    + Close a channel (if the given name correspond to a owned channel)
 + !me
    + Display client status